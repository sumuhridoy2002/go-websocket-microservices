package main

import (
	"log"
	"net/http"
	"sync"

	"my-ecommerce-app/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Image       string  `json:"image"`
	Stock       int     `json:"stock"`
}

var mu sync.Mutex
var mockProducts = []Product{
	{ID: "1", Name: "Laptop", Price: 45000, Description: "High performance laptop", Image: "https://images.unsplash.com/photo-1517336714731-489689fd1ca8?w=500", Stock: 10},
	{ID: "2", Name: "Smartphone", Price: 25000, Description: "Latest 5G smartphone", Image: "https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?w=500", Stock: 15},
	{ID: "3", Name: "Headphones", Price: 3000, Description: "Noise cancelling headphones", Image: "https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=500", Stock: 25},
	{ID: "4", Name: "Smartwatch", Price: 4500, Description: "Fitness tracking smartwatch", Image: "https://images.unsplash.com/photo-1523275335684-37898b6baf30?w=500", Stock: 8},
}

// WebSocket clients and broadcast channel
var clients = make(map[*websocket.Conn]bool)
var broadcast = make(chan []Product)
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleConnections(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	mu.Lock()
	clients[conn] = true
	mu.Unlock()

	// By WebSocket, send the current product list to the connected client
	conn.WriteJSON(mockProducts)

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			mu.Lock()
			delete(clients, conn)
			mu.Unlock()
			break
		}
	}
}

func handleMessages() {
	for {
		products := <-broadcast
		mu.Lock()
		for client := range clients {
			err := client.WriteJSON(products)
			if err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		mu.Unlock()
	}
}

func main() {
	r := gin.Default()

	go handleMessages()

	r.GET("/ws", handleConnections)

	r.GET("/api/v1/products", func(c *gin.Context) {
		mu.Lock()
		defer mu.Unlock()
		response.Success(c, http.StatusOK, "Products fetched successfully", mockProducts)
	})

	r.POST("/api/v1/products/update-stock", func(c *gin.Context) {
		var req struct {
			ProductID string `json:"product_id"`
			Quantity  int    `json:"quantity"`
		}

		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid payload", err.Error())
			return
		}

		mu.Lock()
		for i, p := range mockProducts {
			if p.ID == req.ProductID {
				if mockProducts[i].Stock < req.Quantity {
					mu.Unlock()
					response.Error(c, http.StatusBadRequest, "Out of stock!", nil)
					return
				}
				mockProducts[i].Stock -= req.Quantity
				break
			}
		}
		updatedProducts := mockProducts
		mu.Unlock()

		// By WebSocket, notify all connected clients about the stock update
		broadcast <- updatedProducts

		response.Success(c, http.StatusOK, "Stock updated successfully", nil)
	})

	log.Println("🛍️ Product Service (with WebSocket) running on port :9002")
	r.Run(":9002")
}
