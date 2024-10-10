// internal/api/routes.go

package api

import (
	"s3-at-home/internal/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(store storage.Store) *gin.Engine {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	r.Use(cors.New(config))

	handlers := NewHandlers(store)

	r.GET("/", handlers.ListAllBuckets)
	r.PUT("/:bucket", handlers.CreateNewBucket)
	r.DELETE("/:bucket", handlers.RemoveBucket)
	r.GET("/:bucket", handlers.ListBucketObject)
	r.PUT("/:bucket/:object", handlers.UploadObject)
	r.GET("/:bucket/:object", handlers.DownloadObject)
	r.DELETE("/:bucket/:object", handlers.RemoveObject)

	return r
}
