package server

import (
	"fmt"
	"go_final_project/pkg/api"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func Run() error {
	api.Init()

	// Используем переменную окружения для определения порта
	godotenv.Load()
	port := os.Getenv("TODO_PORT")

	// Если значение порта не установлено - присваеваем переменной нужное нам значение
	if port == "" {
		port = "7540"
	}

	// Запускаем сервер
	http.Handle("/", http.FileServer(http.Dir("web")))
	return http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
