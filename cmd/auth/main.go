package main

import (
	"log"
	"net/http"

	"my-ecommerce-app/pkg/response"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func main() {
	r := gin.Default()

	r.POST("/api/v1/auth/login", func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			response.Error(c, http.StatusBadRequest, "Invalid request payload", err.Error())
			return
		}

		// Simple authentication logic (for demonstration purposes)
		if req.Email == "user@ecommerce.com" && req.Password == "123456" {
			response.Success(c, http.StatusOK, "Authentication successful", gin.H{
				"token": "ecommerce-jwt-secret-token-xyz",
				"user": gin.H{
					"name":  "Ecommerce Customer",
					"email": req.Email,
				},
			})
			return
		}

		response.Unauthorized(c, "Invalid email or password")
	})

	log.Println("🔑 Auth Service running on port :9001")
	r.Run(":9001")
}
