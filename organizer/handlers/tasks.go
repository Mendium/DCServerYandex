package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	pb "github.com/mendium/orchestrator-c/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
	"strconv"
)

// GET
type UserRequestGet struct {
	Token  string `json:"token"`
	TaskID int    `json:"task_id"`
}

// POST
type UserRequestPost struct {
	Token      string `json:"token"`
	Expression string `json:"expression"`
}

// Структура для хранения информации о задаче
type Task struct {
	Expression string `json:"expression"`
	Status     string `json:"status"`
	Answer     string `json:"answer"`
}

// JWTHandler обрабатывает запросы для проверки JWT токена и извлечения логина пользователя
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodPost:
		postTaskHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)

	}
}

// Функция для получения информации о задаче из базы данных
func getTask(taskID int, login string) *Task {
	// Подключение к базе данных
	db, err := sql.Open("mysql", "docker_test_exo:1111@tcp(localhost:3306)/docker_test")
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return nil
	}
	defer db.Close()

	// Запрос к базе данных
	var task Task
	err = db.QueryRow("SELECT expression, status, answer FROM Tasks WHERE task_id = ? AND login = ?", taskID, login).Scan(&task.Expression, &task.Status, &task.Answer)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Задача не найдена")
		} else {
			fmt.Println("Ошибка при выполнении запроса:", err)
		}
		return nil
	}

	return &task
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Парсим JSON из тела запроса
	var userReq UserRequestGet
	err := json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		http.Error(w, "Ошибка при чтении JSON", http.StatusBadRequest)
		return
	}

	// Парсинг и проверка токена
	token, err := jwt.Parse(userReq.Token, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи токена
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Неправильный метод подписи токена")
		}
		return secretKey, nil // Возвращаем секретный ключ для проверки подписи токена
	})
	if err != nil {
		http.Error(w, "Неверный токен", http.StatusUnauthorized)
		return
	}

	// Проверяем, валиден ли токен
	if !token.Valid {
		http.Error(w, "Токен недействителен", http.StatusUnauthorized)
		return
	}

	// Извлекаем логин пользователя из токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Неправильный формат токена", http.StatusInternalServerError)
		return
	}
	login, ok := claims["login"].(string)
	if !ok {
		http.Error(w, "Неправильный формат утверждения 'login'", http.StatusInternalServerError)
		return
	}

	// Получаем информацию о задаче из базы данных
	task := getTask(userReq.TaskID, login)
	if task == nil {
		http.Error(w, "Задача не найдена или не принадлежит текущему пользователю", http.StatusNotFound)
		return
	}

	// Проверяем статус задачи и возвращаем ответ
	switch task.Status {
	case "pending":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Задача с ID %d находится в процессе вычисления", userReq.TaskID)
	case "ready":
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Выражение посчиталось! Ответ на задачу с ID %d: %s", userReq.TaskID, task.Answer)
	default:
		http.Error(w, "Неподдерживаемый статус задачи", http.StatusInternalServerError)
	}
}

func postTaskHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "docker_test_exo:1111@tcp(localhost:3306)/docker_test")
	if err != nil {
		fmt.Println("Ошибка при подключении к базе данных:", err)
		return
	}
	defer db.Close()

	var userReq UserRequestPost
	err = json.NewDecoder(r.Body).Decode(&userReq)
	if err != nil {
		http.Error(w, "Ошибка при чтении JSON", http.StatusBadRequest)
		return
	}

	// Парсинг и проверка токена
	token, err := jwt.Parse(userReq.Token, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи токена
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Неправильный метод подписи токена")
		}
		return secretKey, nil // Возвращаем секретный ключ для проверки подписи токена
	})
	if err != nil {
		http.Error(w, "Неверный токен", http.StatusUnauthorized)
		return
	}

	// Проверяем, валиден ли токен
	if !token.Valid {
		http.Error(w, "Токен недействителен", http.StatusUnauthorized)
		return
	}

	// Извлекаем логин пользователя из токена
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		http.Error(w, "Неправильный формат токена", http.StatusInternalServerError)
		return
	}
	login, ok := claims["login"].(string)
	if !ok {
		http.Error(w, "Неправильный формат утверждения 'login'", http.StatusInternalServerError)
		return
	}
	query := "INSERT INTO Tasks (expression, status, answer, login) VALUES (?, ?, ?, ?)"
	result, err := db.Exec(query, userReq.Expression, "pending", 0, login)
	if err != nil {
		http.Error(w, "Ошибка при добавлении выражения в базу данных", http.StatusInternalServerError)
		return
	}
	taskID, err := result.LastInsertId()
	fmt.Fprintf(w, "Ваша задача добавлена в очередь. Ее идентификатор: "+strconv.Itoa(int(taskID)))

	host := "localhost"
	port := "5000"

	addr := fmt.Sprintf("%s:%s", host, port) // используем адрес сервера
	// установим соединение
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Println("could not connect to grpc server: ", err)
		os.Exit(1)
	}
	// закроем соединение, когда выйдем из функции
	defer conn.Close()
	grpcClient := pb.NewOrchestratorServiceClient(conn)

	_, err = grpcClient.Orchestrate(context.TODO(), &pb.Expression{
		Expression: userReq.Expression,
		TaskId:     int32(taskID),
	})
	if err != nil {
		log.Println("Не удалось передать выражение оркестратору: ", err)
	}

}
