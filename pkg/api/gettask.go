package api

import (
	"database/sql"
	"go_final_project/pkg/db"
	"net/http"
	"strconv"
)

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		writeJSONError(w, "не указан идентификатор", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeJSONError(w, "неверный формат идентификатора", http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, "задача не найдена", http.StatusNotFound)
		} else {
			writeJSONError(w, "ошибка при получении задачи", http.StatusInternalServerError)
		}
		return
	}

	writeJSONResponse(w, task, http.StatusOK)
}
