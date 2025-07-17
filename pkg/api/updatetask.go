package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"go_final_project/pkg/db"
	"net/http"
	"time"
)

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	if task.ID == 0 {
		writeJSONError(w, "не указан идентификатор", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSONError(w, "не указан заголовок задания", http.StatusBadRequest)
		return
	}

	// Получаем текущее время в UTC
	now := time.Now()
	// Обработка даты
	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	} else if _, err := time.Parse(DateFormat, task.Date); err != nil {
		writeJSONError(w, "неверный формат даты", http.StatusBadRequest)
		return
	}

	taskDate, _ := time.Parse(DateFormat, task.Date)
	var next string
	var err error

	// Для повторяющихся задач пересчитываем дату
	if task.Repeat != "" {
		if afterNow(taskDate, now) || taskDate.Equal(now) {
			next, err = NextDate(taskDate, task.Date, task.Repeat)
			if err != nil {
				writeJSONError(w, "неверный формат правила повторения", http.StatusBadRequest)
				return
			}
		} else {
			next, err = NextDate(now, task.Date, task.Repeat)
			if err != nil {
				writeJSONError(w, "неверный формат правила повторения", http.StatusBadRequest)
				return
			}
		}
		task.Date = next
	} else {
		// Для одноразовых задач корректируем просроченные
		if parsedDate, _ := time.Parse(DateFormat, task.Date); parsedDate.Before(now) {
			task.Date = now.Format(DateFormat)
		}
	}

	if err := db.UpdateTask(&task); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSONError(w, "задачи не найдено", http.StatusNotFound)
		} else {
			writeJSONError(w, "ошибка при обновлении задачи", http.StatusInternalServerError)
		}
		return
	}

	writeJSONResponse(w, struct{}{}, http.StatusOK)
}
