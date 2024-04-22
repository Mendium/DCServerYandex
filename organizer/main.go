package main

import (
	"fmt"
	"github.com/mendium/orchestrator-c/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/register", handlers.RegisterHandler)
	http.HandleFunc("/login", handlers.LoginHandler)
	http.HandleFunc("/tasks", handlers.TasksHandler)

	fmt.Println("Сервер запущен на порту :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
