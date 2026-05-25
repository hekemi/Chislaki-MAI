package task10

import "chislennie-metodi/internal/taskmeta"

// Meta возвращает описание десятого задания.
func Meta() taskmeta.Meta {
	return taskmeta.Meta{
		ID:          10,
		Title:       "Задание 10",
		Description: "Заготовка для десятого модуля численных методов.",
		Status:      "placeholder",
	}
}
