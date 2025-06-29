package api

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

const DateFormat = "20060102"

func NextDateHandler(w http.ResponseWriter, r *http.Request) {
	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	// Расчитываем параметр now, если он не был оперделен
	var now time.Time
	if nowStr == "" {
		now = time.Now().UTC().Truncate(24 * time.Hour)
	} else {
		var err error
		now, err = time.Parse(DateFormat, nowStr)
		if err != nil {
			writeJSONError(w, "неверно указан параметр now", http.StatusBadRequest)
			return
		}
	}

	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(""))
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
}

// NextDate вычисляет следующую дату задачи в зависимости от правила повторения
func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	date, err := time.Parse(DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("невозможная начальная дата: %v", err)
	}

	if repeat == "" {
		return "", nil
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", errors.New("неверный формат правила повторения")
	}

	switch parts[0] {
	case "d":
		return handleDaily(now, date, parts)
	case "y":
		return handleYearly(now, date)
	case "w":
		return handleWeekly(now, date, parts)
	case "m":
		return handleMonthly(now, date, parts)
	default:
		return "", errors.New("неподдерживаемый формат правила повторения")
	}
}

// handleDaily — обработка ежедневного повторения d N
func handleDaily(now, date time.Time, parts []string) (string, error) {
	var next time.Time

	if len(parts) < 2 {
		return "", errors.New("неверный формат правила d")
	}

	interval, err := strconv.Atoi(parts[1])
	if err != nil || interval < 1 || interval > 400 {
		return "", errors.New("интервал должен быть между 1 и 400")
	}

	if date.Equal(now) && afterNow(date, now) {
		next = date
		return next.Format(DateFormat), nil
	} else {
		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				return date.Format(DateFormat), nil
			}
		}
	}
}

// handleYearly — обработка годового повторения
func handleYearly(now, date time.Time) (string, error) {
	next := date.AddDate(1, 0, 0)
	for !next.After(now) {
		next = next.AddDate(1, 0, 0)
	}
	return next.Format(DateFormat), nil
}

// handleWeekly — обработка недельного повторения
func handleWeekly(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("неверный формат правила w")
	}

	dayStrs := strings.Split(parts[1], ",")
	days := make(map[int]bool)
	for _, s := range dayStrs {
		d, err := strconv.Atoi(s)
		if err != nil || d < 1 || d > 7 {
			return "", errors.New("день должен быть между 1 и 7")
		}
		days[d] = true
	}

	// Начинаем поиск с завтрашнего дня
	date = date.AddDate(0, 0, 1)
	for {
		weekday := int(date.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		if days[weekday] && afterNow(date, now) {
			return date.Format(DateFormat), nil
		}
		date = date.AddDate(0, 0, 1)
	}
}

// handleMonthly — обработка месячного повторения
func handleMonthly(now, date time.Time, parts []string) (string, error) {
	if len(parts) < 2 {
		return "", errors.New("неверный формат правила m")
	}

	// Парсим дни месяца
	dayStrs := strings.Split(parts[1], ",")
	days := make([]int, 0, len(dayStrs))
	for _, s := range dayStrs {
		d, err := strconv.Atoi(s)
		if err != nil || d == 0 || d < -2 || d > 31 {
			return "", errors.New("день должен быть в промежутке между -2 и 31, исключая 0")
		}
		days = append(days, d)
	}

	// Парсим месяцы (если указаны)
	months := []int{}
	if len(parts) > 2 {
		monthStrs := strings.Split(parts[2], ",")
		for _, s := range monthStrs {
			m, err := strconv.Atoi(s)
			if err != nil || m < 1 || m > 12 {
				return "", errors.New("месяц должен быть в промежутке между 1 и 12")
			}
			months = append(months, m)
		}
	}

	// Начинаем поиск с завтрашнего дня
	date = date.AddDate(0, 0, 1)
	for range 36 * 31 { // до 3 лет вперед
		// Проверяем месяц
		if len(months) > 0 {
			found := false
			currentMonth := int(date.Month())

			if slices.Contains(months, currentMonth) {
				found = true
			}

			if !found {
				date = date.AddDate(0, 0, 1)
				continue
			}
		}

		// Проверяем день
		dayMatch := false
		currentDay := date.Day()
		lastDay := time.Date(date.Year(), date.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()
		for _, d := range days {
			switch {
			case d > 0:
				if currentDay == d {
					dayMatch = true
				}
			case d == -1:
				if currentDay == lastDay {
					dayMatch = true
				}
			case d == -2:
				if currentDay == lastDay-1 {
					dayMatch = true
				}
			}
		}

		if dayMatch && afterNow(date, now) {
			return date.Format(DateFormat), nil
		}
		date = date.AddDate(0, 0, 1)
	}

	return "", errors.New("не найдена подходящая дата для m")
}
