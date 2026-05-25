package task11

import "chislennie-metodi/internal/taskmeta"

// Meta возвращает описание одиннадцатого задания.
func Meta() taskmeta.Meta {
	return taskmeta.Meta{
		ID:          11,
		Title:       "Задание 11",
		Description: "Краевая задача для ОДУ 2-го порядка: метод прогонки, O(h^2).",
		Status:      "ready",
	}
}
