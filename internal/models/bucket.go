package models

import "time"

type Bucket struct {
  Name         string
  CreationDate time.Time
}

func NewBucket(name string) *Bucket {
  return &Bucket{
    Name:         name,
    CreationDate: time.Now(),
  }
}
