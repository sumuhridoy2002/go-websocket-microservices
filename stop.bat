@echo off
echo Stopping all microservices, Air, and Next.js frontend...

:: ১. Go এবং Air প্রসেসগুলো ফোর্সফুলি বন্ধ করা
taskkill /f /im go.exe >nul 2>&1
taskkill /f /im air.exe >nul 2>&1

:: ২. সব Node.js প্রসেস বন্ধ করা
taskkill /f /im node.exe >nul 2>&1

:: ৩. নির্দিষ্ট পোর্টগুলোতে (3000, 8080, 9001, 9002, 9003) কোনো প্রসেস আটকে থাকলে তা কিল করা
for %%p in (3000 8080 9001 9002 9003) do (
    for /f "tokens=5" %%a in ('netstat -ano ^| findstr :%%p') do (
        taskkill /f /pid %%a >nul 2>&1
    )
)

echo All services and Next.js stopped successfully!
pause