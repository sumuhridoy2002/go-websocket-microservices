@echo off
echo Starting Ecommerce Microservices with Auto-Restart (Air)...

:: 1. Auth Service (Port 9001)
cd cmd/auth
start cmd /k "air"
cd ../..

:: 2. Product Service (Port 9002)
cd cmd/product
start cmd /k "air"
cd ../..

:: 3. Order Service (Port 9003)
cd cmd/order
start cmd /k "air"
cd ../..

:: 4. API Gateway (Port 8080)
cd cmd/gateway
start cmd /k "air"
cd ../..

:: 5. Frontend (Next.js - অলরেডি নিজস্ব হট-রিলোড আছে)
cd frontend
start cmd /k "npm run dev"

echo All microservices and frontend are running with live reload!
pause