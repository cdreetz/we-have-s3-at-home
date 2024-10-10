package config

import "os"

type Config struct {
  RedisAddr  string
  ServerAddr string
}

func Load() (*Config, error) {
  return &Config{
    RedisAddr:  getEnv("REDIS_ADDR", "localhost:6379"),
    ServerAddr: getEnv("SERVER_ADDR", ":8080"),
  }, nil
}

func getEnv(key, defaultValue string) string {
  if value, exists := os.LookupEnv(key); exists {
    return value
  }
  return defaultValue
}
