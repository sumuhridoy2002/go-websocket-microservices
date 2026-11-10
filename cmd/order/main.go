package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"my-ecommerce-app/pkg/response"

	"github.com/gin-gonic/gin"
)

func internalGuard() gin.HandlerFunc {
	expectedSecret := os.Getenv("INTERNAL_GATEWAY_SECRET")
	if expectedSecret == "" {
		expectedSecret = "ecommerce-super-secret-key"
	}

	return func(c *gin.Context) {
		if c.GetHeader("X-Internal-Gateway-Secret") != expectedSecret {
			response.Forbidden(c, "Direct microservice access is restricted! Use API Gateway.")
			c.Abort()
			return
		}
		c.Next()
	}
}

func main() {
	r := gin.Default()

	orders := r.Group("/api/v1/orders")
	orders.Use(internalGuard())
	{
		orders.POST("", func(c *gin.Context) {
			var body struct {
				ProductID string  `json:"product_id"`
				Item      string  `json:"item"`
				Price     float64 `json:"price"`
			}

			if err := c.ShouldBindJSON(&body); err != nil {
				response.Error(c, http.StatusBadRequest, "Invalid order payload", err.Error())
				return
			}

			// By WebSocket, notify all connected clients about the stock update
			stockReqBody, _ := json.Marshal(map[string]interface{}{
				"product_id": body.ProductID,
				"quantity":   1,
			})

			productServiceURL := "http://localhost:9002/api/v1/products/update-stock"
			httpReq, _ := http.NewRequest("POST", productServiceURL, bytes.NewBuffer(stockReqBody))
			httpReq.Header.Set("Content-Type", "application/json")

			client := &http.Client{Timeout: 3 * time.Second}
			resp, err := client.Do(httpReq)
			if err != nil || resp.StatusCode != http.StatusOK {
				response.Error(c, http.StatusInternalServerError, "Failed to update product stock", nil)
				return
			}

			// Random order ID generation
			rng := rand.New(rand.NewSource(time.Now().UnixNano()))
			orderID := fmt.Sprintf("ORD-%d", rng.Intn(90000)+10000)

			response.Success(c, http.StatusCreated, "Order placed & stock updated successfully", gin.H{
				"order_id":    orderID,
				"item":        body.Item,
				"total_price": body.Price,
				"status":      "PROCESSING",
			})
		})
	}

	log.Println("📦 Order Service running on port :9003")
	r.Run(":9003")
}
