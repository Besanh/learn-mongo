package server

import (
	responsetime "mongo-fundamential/middleware"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	serviceName = "mongo-fundamential"
	version     = "v1.0.0"
)

func NewHTTPServer() *gin.Engine {
	// For Production
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.MaxMultipartMemory = 100 << 20
	engine.Use(CORSMiddleware())
	engine.Use(allowOptionsMethod())
	engine.Use(responsetime.Handler)
	engine.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"service": serviceName,
			"version": version,
			"time":    time.Now().Unix(),
		})
	})
	return engine
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-Tenant-Id, X-Tenant-Uuid, X-Tenant-Name, X-User-Id, X-User-Level, X-User-Name, Token, Auth-Url, Source, Validator-Header")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, PATCH, DELETE, OPTIONS")
		c.Next()
	}
}

func allowOptionsMethod() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	}
}
