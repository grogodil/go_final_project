package db

import (
	"database/sql"
	"errors"
	"time"
)

type Task struct {
	ID      int64  `json:"id,string"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(db *sql.DB, task *Task) (int64, error) {
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

func Tasks(limit int, search string) ([]*Task, error) {
	tasks := make([]*Task, 0)
	query := "SELECT id, date, title, comment, repeat FROM scheduler"
	args := []any{}

	if search != "" {
		if t, err := time.Parse("02.01.2006", search); err == nil {
			searchDate := t.UTC().Truncate(24 * time.Hour).Format("20060102")
			query += " WHERE date = ?"
			args = append(args, searchDate)
		} else {
			searchTerm := "%" + search + "%"
			query += " WHERE title LIKE ? OR comment LIKE ?"
			args = append(args, searchTerm, searchTerm)
		}
	}

	query += " ORDER BY date ASC LIMIT ?"
	args = append(args, limit)

	rows, err := DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat); err != nil {
			return nil, err
		}
		tasks = append(tasks, &t)
	}

	return tasks, nil
}

func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date=?, title=?, comment=?, repeat=? WHERE id=?`
	res, err := DB.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("задание не найдено")
	}

	return nil
}

func GetTask(id int64) (*Task, error) {
	var t Task
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`
	err := DB.QueryRow(query, id).Scan(&t.ID, &t.Date, &t.Title, &t.Comment, &t.Repeat)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &t, nil
}

func DeleteTask(id int64) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	result, err := DB.Exec(query, id)
	if err != nil {
		return err
	}

	// Проверяем, что задача была удалена
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func UpdateDate(next string, id int64) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	result, err := DB.Exec(query, next, id)
	if err != nil {
		return err
	}

	// Проверяем, что задача была обновлена
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
