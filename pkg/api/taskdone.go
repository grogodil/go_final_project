package api

import (
	"database/sql"
	"go_final_project/pkg/db"
	"net/http"
	"strconv"
	"time"
)

func TaskDoneHandler(w http.ResponseWriter, r *http.Request) {
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
			writeJSONError(w, "ошибка БД", http.StatusInternalServerError)
		}
		return
	}

	now := time.Now().UTC().Truncate(24 * time.Hour)

	if task.Repeat == "" {
		// Разовая задача - удаляем
		if err := db.DeleteTask(id); err != nil {
			writeJSONError(w, "ошибка при удалении задачи", http.StatusInternalServerError)
			return
		}
	} else {
		// Повторяющаяся задача - пересчитываем дату
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			writeJSONError(w, "невозможное правило повторения", http.StatusBadRequest)
			return
		}

		if err := db.UpdateDate(next, id); err != nil {
			writeJSONError(w, "ошибка обновления задачи", http.StatusInternalServerError)
			return
		}
	}

	writeJSONResponse(w, struct{}{}, http.StatusOK)
}
