package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // Драйвер для MySQL
)

func main() {
	// Строка подключения к вашей базе данных MySQL
	connString := "docker_test_exo:1111@tcp(localhost:3306)/docker_test"

	// Устанавливаем соединение с базой данных
	db, err := sql.Open("mysql", connString)
	if err != nil {
		log.Fatalf("Ошибка при подключении к базе данных: %v\n", err)
	}
	defer db.Close()

	query := `
		CREATE TABLE IF NOT EXISTS Users (
			ID INT AUTO_INCREMENT PRIMARY KEY,
			Login VARCHAR(255) NOT NULL,
			Password VARCHAR(255) NOT NULL
		)
	`
	query_2 := `
		CREATE TABLE IF NOT EXISTS Tasks (
			task_id INT AUTO_INCREMENT PRIMARY KEY,
			expression VARCHAR(255),
			status VARCHAR(255),
			answer VARCHAR(255),
			login VARCHAR(255)
		)
	`

	// Выполнение SQL запроса
	_, err = db.Exec(query)
	if err != nil {
		fmt.Println("Ошибка при создании таблицы Users:", err)
		return
	}

	_, err = db.Exec(query_2)
	if err != nil {
		fmt.Println("Ошибка при создании таблицы Tasks:", err)
		return
	}

	// Проверяем подключение к базе данных
	err = db.Ping()
	if err != nil {
		log.Fatalf("Ошибка при проверке подключения к базе данных: %v\n", err)
	}

	fmt.Println("Успешно подключились к базе данных!")

	fmt.Println("Таблица Users успешно создана")
	fmt.Println("Таблица Tasks успешно создана")
	// Теперь вы можете выполнить запросы к вашей базе данных здесь

}
