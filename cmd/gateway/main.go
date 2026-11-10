package main

import (
	"log"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func proxyHandler(target string) gin.HandlerFunc {
	remote, err := url.Parse(target)
	if err != nil {
		log.Fatalf("Failed to parse target URL: %v", err)
	}

	proxy := httputil.NewSingleHostReverseProxy(remote)

	return func(c *gin.Context) {
		secretKey := os.Getenv("INTERNAL_GATEWAY_SECRET")
		if secretKey == "" {
			secretKey = "ecommerce-super-secret-key"
		}

		// Set the internal secret header for inter-service communication
		c.Request.Header.Set("X-Internal-Gateway-Secret", secretKey)
		c.Request.Host = remote.Host

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://127.0.0.1:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	authURL := "http://localhost:9001"
	productURL := "http://localhost:9002"
	orderURL := "http://localhost:9003"

	// API Gateway routes
	api := r.Group("/api/v1")
	{
		api.Any("/auth/*path", proxyHandler(authURL))
		api.Any("/auth", proxyHandler(authURL))

		api.Any("/products/*path", proxyHandler(productURL))
		api.Any("/products", proxyHandler(productURL))

		api.Any("/orders/*path", proxyHandler(orderURL))
		api.Any("/orders", proxyHandler(orderURL))
	}

	log.Println("🚀 API Gateway running on port :8080")
	r.Run(":8080")
}
