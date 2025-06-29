package api

import (
	"go_final_project/pkg/db"
	"net/http"
	"strconv"
)

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := db.DeleteTask(id); err != nil {
		writeJSONError(w, "ошибка при удалении задачи", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, struct{}{}, http.StatusOK)
}
