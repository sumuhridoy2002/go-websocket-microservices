# ⚡ Go WebSocket Microservices

A modern, high-performance distributed architecture constructed with **Go (Golang)** and **Next.js**, showcasing real-time data synchronization via **WebSockets**, **API Gateway pattern**, and **inter-service communication**.

## 🎯 Architecture Summary
* **API Gateway (`:8080`)**: Entry point routing traffic to downstream microservices with security headers.
* **Auth Service (`:9001`)**: Issues JWT tokens and handles client authorization.
* **Product Service (`:9002`)**: Manages product state & hosts the central WebSocket broadcast hub.
* **Order Service (`:9003`)**: Processes orders and communicates with Product Service for inventory updates.
* **Next.js Frontend (`:3000`)**: Dynamic UI with auto-reconnecting WebSocket client for live state updates.