package models

import "time"

type Object struct {
  Key          string
  ContentType  string
  Data         []byte
  CreationDate time.Time
}

func NewObject(key, contentType string, data []byte) *Object {
  return &Object{
    Key:          key,
    ContentType:  contentType,
    Data:         data,
    CreationDate: time.Now(),
  }
}
