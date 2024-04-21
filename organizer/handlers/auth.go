package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
)

// User представляет данные о пользователе для аутентификации
type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// Секретный ключ для подписи JWT токена
var secretKey = []byte("bobr_kurwa")

// Строка подключения к базе данных MySQL
const dsn = "docker_test_exo:1111@tcp(localhost:3306)/docker_test"

// LoginHandler обрабатывает запросы на аутентификацию пользователя
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Парсим JSON из тела запроса
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Ошибка при чтении JSON", http.StatusBadRequest)
		return
	}

	// Проверяем, существует ли пользователь в базе данных
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		http.Error(w, "Ошибка при подключении к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var dbLogin, dbPassword string
	err = db.QueryRow("SELECT Login, Password FROM Users WHERE Login = ?", user.Login).Scan(&dbLogin, &dbPassword)
	if err != nil {
		http.Error(w, "Неправильные логин или пароль", http.StatusUnauthorized)
		return
	}

	// Проверяем совпадение пароля
	if user.Password != dbPassword {
		http.Error(w, "Неправильные логин или пароль", http.StatusUnauthorized)
		return
	}

	// Создаем новый JWT токен с именем пользователя и сроком действия 1 час
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": user.Login,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	})

	// Подписываем токен с секретным ключом
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Ошибка при создании токена", http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ клиенту с JWT токеном
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Вход выполнен успешно. JWT токен: %s", tokenString)
}

// RegisterHandler обрабатывает запросы на регистрацию нового пользователя
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Парсим JSON из тела запроса
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Ошибка при чтении JSON", http.StatusBadRequest)
		return
	}

	// Устанавливаем соединение с базой данных
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		http.Error(w, "Ошибка при подключении к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// SQL запрос для добавления пользователя в базу данных
	query := "INSERT INTO Users (Login, Password) VALUES (?, ?)"
	_, err = db.Exec(query, user.Login, user.Password)
	if err != nil {
		http.Error(w, "Ошибка при добавлении пользователя в базу данных", http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ клиенту
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Регистрация завершена успешно")
}
