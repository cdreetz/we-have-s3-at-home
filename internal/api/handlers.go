package api

import (
	"io"
	"net/http"

	"s3-at-home/internal/models"
	"s3-at-home/internal/storage"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	store storage.Store
}

func NewHandlers(store storage.Store) *Handlers {
	return &Handlers{store: store}
}

func (h *Handlers) ListAllBuckets(c *gin.Context) {
	buckets, err := h.store.ListBuckets()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, buckets)
}

func (h *Handlers) CreateNewBucket(c *gin.Context) {
	bucketName := c.Param("bucket")
	bucket := models.NewBucket(bucketName)
	if err := h.store.CreateBucket(bucket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusCreated)
}

func (h *Handlers) RemoveBucket(c *gin.Context) {
	bucketName := c.Param("bucket")
	if err := h.store.DeleteBucket(bucketName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handlers) ListBucketObject(c *gin.Context) {
	bucketName := c.Param("bucket")
	objects, err := h.store.ListObjects(bucketName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, objects)
}

func (h *Handlers) UploadObject(c *gin.Context) {
	bucketName := c.Param("bucket")
	objectKey := c.Param("object")
	contentType := c.Request.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	data, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
		return
	}

	object := models.NewObject(objectKey, contentType, data)
	if err := h.store.PutObject(bucketName, object); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusCreated)
}

func (h *Handlers) DownloadObject(c *gin.Context) {
	bucketName := c.Param("bucket")
	objectKey := c.Param("object")

	object, err := h.store.GetObject(bucketName, objectKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Object not found"})
		return
	}

	c.Data(http.StatusOK, object.ContentType, object.Data)
}

func (h *Handlers) RemoveObject(c *gin.Context) {
	bucketName := c.Param("bucket")
	objectKey := c.Param("object")

	if err := h.store.DeleteObject(bucketName, objectKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
