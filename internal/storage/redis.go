// internal/storage/redis.go
package storage

import (
	"context"
	"fmt"
	"time"

	"s3-at-home/internal/models"

	"github.com/go-redis/redis/v8"
)

type Store interface {
	CreateBucket(bucket *models.Bucket) error
	DeleteBucket(bucketName string) error
	ListBuckets() ([]string, error)
	PutObject(bucketName string, object *models.Object) error
	GetObject(bucketName, objectKey string) (*models.Object, error)
	DeleteObject(bucketName, objectKey string) error
	ListObjects(bucketName string) ([]string, error)
	BucketExists(bucketName string) (bool, error)
	ObjectExists(bucketName, objectKey string) (bool, error)
}

type RedisStore struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStore(addr string) (*RedisStore, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})

	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	return &RedisStore{
		client: client,
		ctx:    ctx,
	}, nil
}

func (s *RedisStore) CreateBucket(bucket *models.Bucket) error {
	key := fmt.Sprintf("bucket:%s", bucket.Name)
	exists, err := s.client.Exists(s.ctx, key).Result()
	if err != nil {
		return fmt.Errorf("error checking bucket existence: %w", err)
	}
	if exists == 1 {
		return fmt.Errorf("bucket already exists")
	}

	err = s.client.HSet(s.ctx, key, "creation_date", bucket.CreationDate.Format(time.RFC3339)).Err()
	if err != nil {
		return fmt.Errorf("error creating bucket: %w", err)
	}

	return nil
}

func (s *RedisStore) DeleteBucket(bucketName string) error {
	key := fmt.Sprintf("bucket:%s", bucketName)
	exists, err := s.client.Exists(s.ctx, key).Result()
	if err != nil {
		return fmt.Errorf("error checking bucket existence: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("bucket not found")
	}

	// Delete all objects in the bucket
	objects, err := s.ListObjects(bucketName)
	if err != nil {
		return fmt.Errorf("error listing objects: %w", err)
	}
	for _, objectKey := range objects {
		if err := s.DeleteObject(bucketName, objectKey); err != nil {
			return fmt.Errorf("error deleting object %s: %w", objectKey, err)
		}
	}

	// Delete bucket metadata and contents set
	_, err = s.client.Del(s.ctx, key, fmt.Sprintf("bucket:%s:contents", bucketName)).Result()
	if err != nil {
		return fmt.Errorf("error deleting bucket: %w", err)
	}

	return nil
}

func (s *RedisStore) ListBuckets() ([]string, error) {
	var buckets []string
	iter := s.client.Scan(s.ctx, 0, "bucket:*", 0).Iterator()
	for iter.Next(s.ctx) {
		key := iter.Val()
		if key[len(key)-9:] != ":contents" {
			buckets = append(buckets, key[7:])
		}
	}
	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("error listing buckets: %w", err)
	}
	return buckets, nil
}

func (s *RedisStore) PutObject(bucketName string, object *models.Object) error {
	bucketKey := fmt.Sprintf("bucket:%s", bucketName)
	exists, err := s.client.Exists(s.ctx, bucketKey).Result()
	if err != nil {
		return fmt.Errorf("error checking bucket existence: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("bucket not found")
	}

	metadataKey := fmt.Sprintf("object:%s:%s:metadata", bucketName, object.Key)
	dataKey := fmt.Sprintf("object:%s:%s:data", bucketName, object.Key)

	metadata := map[string]interface{}{
		"content_type":  object.ContentType,
		"size":          len(object.Data),
		"creation_date": object.CreationDate.Format(time.RFC3339),
	}

	err = s.client.HSet(s.ctx, metadataKey, metadata).Err()
	if err != nil {
		return fmt.Errorf("error storing object metadata: %w", err)
	}

	err = s.client.Set(s.ctx, dataKey, object.Data, 0).Err()
	if err != nil {
		return fmt.Errorf("error storing object data: %w", err)
	}

	err = s.client.SAdd(s.ctx, fmt.Sprintf("bucket:%s:contents", bucketName), object.Key).Err()
	if err != nil {
		return fmt.Errorf("error adding object to bucket contents: %w", err)
	}

	return nil
}

func (s *RedisStore) GetObject(bucketName, objectKey string) (*models.Object, error) {
	metadataKey := fmt.Sprintf("object:%s:%s:metadata", bucketName, objectKey)
	dataKey := fmt.Sprintf("object:%s:%s:data", bucketName, objectKey)

	metadata, err := s.client.HGetAll(s.ctx, metadataKey).Result()
	if err != nil {
		return nil, fmt.Errorf("error getting object metadata: %w", err)
	}
	if len(metadata) == 0 {
		return nil, fmt.Errorf("object not found")
	}

	data, err := s.client.Get(s.ctx, dataKey).Bytes()
	if err != nil {
		return nil, fmt.Errorf("error getting object data: %w", err)
	}

	creationDate, err := time.Parse(time.RFC3339, metadata["creation_date"])
	if err != nil {
		return nil, fmt.Errorf("error parsing creation date: %w", err)
	}

	return &models.Object{
		Key:          objectKey,
		ContentType:  metadata["content_type"],
		Data:         data,
		CreationDate: creationDate,
	}, nil
}

func (s *RedisStore) DeleteObject(bucketName, objectKey string) error {
	bucketKey := fmt.Sprintf("bucket:%s", bucketName)
	exists, err := s.client.Exists(s.ctx, bucketKey).Result()
	if err != nil {
		return fmt.Errorf("error checking bucket existence: %w", err)
	}
	if exists == 0 {
		return fmt.Errorf("bucket not found")
	}

	metadataKey := fmt.Sprintf("object:%s:%s:metadata", bucketName, objectKey)
	dataKey := fmt.Sprintf("object:%s:%s:data", bucketName, objectKey)

	_, err = s.client.Del(s.ctx, metadataKey, dataKey).Result()
	if err != nil {
		return fmt.Errorf("error deleting object: %w", err)
	}

	err = s.client.SRem(s.ctx, fmt.Sprintf("bucket:%s:contents", bucketName), objectKey).Err()
	if err != nil {
		return fmt.Errorf("error removing object from bucket contents: %w", err)
	}

	return nil
}

func (s *RedisStore) ListObjects(bucketName string) ([]string, error) {
	bucketKey := fmt.Sprintf("bucket:%s", bucketName)
	exists, err := s.client.Exists(s.ctx, bucketKey).Result()
	if err != nil {
		return nil, fmt.Errorf("error checking bucket existence: %w", err)
	}
	if exists == 0 {
		return nil, fmt.Errorf("bucket not found")
	}

	objects, err := s.client.SMembers(s.ctx, fmt.Sprintf("bucket:%s:contents", bucketName)).Result()
	if err != nil {
		return nil, fmt.Errorf("error listing objects: %w", err)
	}

	return objects, nil
}

// Optional: Add a method to check if a bucket exists
func (s *RedisStore) BucketExists(bucketName string) (bool, error) {
	bucketKey := fmt.Sprintf("bucket:%s", bucketName)
	exists, err := s.client.Exists(s.ctx, bucketKey).Result()
	if err != nil {
		return false, fmt.Errorf("error checking bucket existence: %w", err)
	}
	return exists == 1, nil
}

// Optional: Add a method to check if an object exists
func (s *RedisStore) ObjectExists(bucketName, objectKey string) (bool, error) {
	metadataKey := fmt.Sprintf("object:%s:%s:metadata", bucketName, objectKey)
	exists, err := s.client.Exists(s.ctx, metadataKey).Result()
	if err != nil {
		return false, fmt.Errorf("error checking object existence: %w", err)
	}
	return exists == 1, nil
}
