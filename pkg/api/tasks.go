package api

import (
	//"fmt"
	"go_final_project/pkg/db"
	"net/http"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// Обработчик Get-запросов
func TasksHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSONError(w, "неверный метод", http.StatusMethodNotAllowed)
		return
	}

	search := r.URL.Query().Get("search")
	tasks, err := db.Tasks(50, search)
	if err != nil {
		writeJSONError(w, "ошибка сервера", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, TasksResp{Tasks: tasks}, http.StatusOK)
}
