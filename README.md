# DCServerYandex
 # Yandex Lyceum - Final Project - Distributed Computing Server
 ![Logo](https://github.com/Mendium/DCServerYa/blob/main/orig.png)

## Требования: Docker, Go

## Запуск: 
 1. Запуск сервисов (огранайзер и оркестратор) и базы данных в Docker-контейнерах:
    
     ```bash
     docker-compose up --build
    ```
 2. Инициализация базы данных (перед этим переходим в директорию db_init):
    
    ```bash
     go run main.go
    ```
## Готово! Если у вас есть Docker Desktop, вы можете увидеть запущенные контейнеры:
![Logo](https://github.com/Mendium/DCServerYa/blob/main/orig.png)  
