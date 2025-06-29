package api

import (
	"encoding/json"
	"go_final_project/pkg/db"
	"net/http"
	"time"
)

func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSONError(w, "ошибка десериализации JSON", http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJSONError(w, "не указан заголовок задания", http.StatusBadRequest)
		return
	}

	var now time.Time
	var err error
	nowStr := r.URL.Query().Get("now")
	if nowStr == "" {
		now = time.Now().UTC().Truncate(24 * time.Hour)
	} else {
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			writeJSONError(w, "неверный формат now", http.StatusBadRequest)
			return
		}
	}

	if task.Date == "" || task.Date == now.Format(DateFormat) {
		task.Date = now.Format(DateFormat)
	} else if _, err := time.Parse(DateFormat, task.Date); err != nil {
		writeJSONError(w, "неверный формат даты", http.StatusBadRequest)
		return
	}

	parsedDate, _ := time.Parse(DateFormat, task.Date)
	var next string

	if task.Repeat != "" {
		if afterNow(now, parsedDate) {
			next, err = NextDate(now, task.Date, task.Repeat)
			if err != nil {
				writeJSONError(w, "неверный формат правила повторения", http.StatusBadRequest)
				return
			}
		} else {
			next = parsedDate.Format(DateFormat)
		}
		task.Date = next
	} else {
		if parsedDate.Before(now) {
			task.Date = now.Format(DateFormat)
		} else {
			task.Date = parsedDate.Format(DateFormat)
		}
	}

	id, err := db.AddTask(db.DB, &task)
	if err != nil {
		writeJSONError(w, "ошибка при добавлении задачи", http.StatusInternalServerError)
		return
	}

	writeJSONResponse(w, map[string]int64{"id": id}, http.StatusOK)
}
