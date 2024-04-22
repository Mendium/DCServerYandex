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

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// секретный ключ для jwt токена
var secretKey = []byte("bobr_kurwa")

// адрес бд
const dsn = "docker_test_exo:1111@tcp(db:3306)/docker_test"

// авторизация
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// проверяем метод запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// парсим JSON из тела запроса
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Ошибка при чтении JSON", http.StatusBadRequest)
		return
	}

	// проверяем, существует ли пользователь в базе данных
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

	// проверяем совпадение пароля
	if user.Password != dbPassword {
		http.Error(w, "Неправильные логин или пароль", http.StatusUnauthorized)
		return
	}

	// создание нового jwt токена на основе логина пользователя (срок действия - 1 час)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"login": user.Login,
		"exp":   time.Now().Add(time.Hour * 1).Unix(),
	})

	//подписывание токена секрет. ключом
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		http.Error(w, "Ошибка при создании токена", http.StatusInternalServerError)
		return
	}

	// Отправляем успешный ответ клиенту с JWT токеном
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Вход выполнен успешно. JWT токен: %s", tokenString)
}

// регистрация нового пользователя
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	//проверка метода запроса
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// парсинг json из body запроса
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Ошибка при чтении JSON", http.StatusBadRequest)
		return
	}

	// соединение с бд
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		http.Error(w, "Ошибка при подключении к базе данных", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	existsU, _ := userExists(db, user.Login)
	if existsU {
		http.Error(w, "Такой логин уже существует.", http.StatusUnauthorized)
		return
	}

	// запрос для добавления пользователя в бд
	query := "INSERT INTO Users (Login, Password) VALUES (?, ?)"
	_, err = db.Exec(query, user.Login, user.Password)
	if err != nil {
		http.Error(w, "Ошибка при добавлении пользователя в базу данных", http.StatusInternalServerError)
		fmt.Println(err)
		return
	}

	//отправляем успешный ответ клиенту
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Регистрация завершена успешно")
}

func userExists(db *sql.DB, login string) (bool, error) {
	query := "SELECT COUNT(*) FROM Users WHERE Login = ?"

	var count int
	err := db.QueryRow(query, login).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
