package validators

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ValidateDate проверяет корректность даты в формате ДД.ММ.ГГГГ
func ValidateDate(dateStr string) (time.Time, bool) {
	// Регулярное выражение для проверки формата
	re := regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}$`)
	if !re.MatchString(dateStr) {
		return time.Time{}, false
	}

	// Парсинг даты
	date, err := time.Parse("02.01.2006", dateStr)
	if err != nil {
		return time.Time{}, false
	}

	// Проверка, что дата не в будущем
	if date.After(time.Now()) {
		return time.Time{}, false
	}

	return date, true
}

// ValidateSalary проверяет корректность ввода зарплаты
func ValidateSalary(salaryStr string) (float64, bool) {
	// Удаляем пробелы и заменяем запятую на точку
	salaryStr = strings.ReplaceAll(salaryStr, " ", "")
	salaryStr = strings.ReplaceAll(salaryStr, ",", ".")

	salary, err := strconv.ParseFloat(salaryStr, 64)
	if err != nil {
		return 0, false
	}

	if salary < 0 {
		return 0, false
	}

	return salary, true
}

// ValidateDays проверяет корректность количества дней
func ValidateDays(daysStr string) (int, bool) {
	days, err := strconv.Atoi(daysStr)
	if err != nil {
		return 0, false
	}

	if days < 0 || days > 730 {
		return 0, false
	}

	return days, true
}
