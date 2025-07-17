package main

import (
	"go_final_project/pkg/db"
	"go_final_project/pkg/server"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	defer db.DB.Close()

	// Получаем путь к БД из переменной окружения
	godotenv.Load()
	dbFile := os.Getenv("TODO_DBFILE")

	// Если путь к файлу БД отсутствует - присваеваем переменной нужное нам значение
	if dbFile == "" {
		dbFile = "scheduler.db"
	}

	// Инициализируем БД
	var err error
	if err = db.Init(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации БД: %v", err)
	}

	// Запускаем сервер
	server.Run()
}
