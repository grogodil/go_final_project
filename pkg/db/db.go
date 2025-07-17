package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
	CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(255) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(50) NOT NULL DEFAULT ""
	);

	CREATE INDEX idx_date ON scheduler (date);
	`

var DB *sql.DB

func Init(dbFile string) error {
	var err error

	// Открываем БД
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия БД: %w", err)
	}

	// Проверяем существование БД файла
	// Если БД не существует - создаем ее
	_, err = os.Stat(dbFile)
	if os.IsNotExist(err) {
		if _, err = DB.Exec(schema); err != nil {
			return fmt.Errorf("ошибка создания БД: %w", err)
		}
	}

	return nil
}
