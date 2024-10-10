package main

import (
  "log"

  "s3-at-home/config"
  "s3-at-home/internal/api"
  "s3-at-home/internal/storage"
)

func main() {
  cfg, err := config.Load()
  if err != nil {
    log.Fatalf("Failed to load configuration: %v", err)
  }

  store, err := storage.NewRedisStore(cfg.RedisAddr)
  if err != nil {
    log.Fatalf("Failed to create Redis store: %v", err)
  }

  router := api.SetupRouter(store)

  log.Printf("Starting server on %s", cfg.ServerAddr)
  if err := router.Run(cfg.ServerAddr); err != nil {
    log.Fatalf("Failed to start server: %v", err)
  }
}
