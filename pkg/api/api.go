package api

import "net/http"

func Init() {
	http.HandleFunc("/api/signin", SigninHandler)
	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", AuthMiddleware(TaskHandler))
	http.HandleFunc("/api/tasks", AuthMiddleware(TasksHandler))
	http.HandleFunc("/api/task/done", AuthMiddleware(TaskDoneHandler))
}
