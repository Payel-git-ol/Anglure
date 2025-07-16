@echo off
docker-compose down
docker-compose up -d --build
docker ps
echo "Приложение запущено! Откройте http://localhost:8080"