# DCServerYandex
 # Yandex Lyceum - Final Project - Distributed Computing Server
 ![Logo](https://github.com/Mendium/DCServerYa/blob/main/orig.png)

## Требования: Docker (желательно Desktop), Go, Postman

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
![Ex](docs/1355.png)


# Синтаксис запросов (на примере Postman):
## Регистрация нового пользователя в БД (/register):
### Метод POST; Body выбираем raw:
```bash
{
    "login": "testuser123",
    "password": "testpassword123"
}
```
![Ex](docs/5352.png)
## Вход пользователя в систему и получение JWT-токена сроком на час (/login):
### Метод POST
```
