package task05

import "chislennie-metodi/internal/taskmeta"

// Meta возвращает описание пятого задания.
func Meta() taskmeta.Meta {
	return taskmeta.Meta{
		ID:          5,
		Title:       "Задание 5",
		Description: "Заготовка для пятого модуля численных методов.",
		Status:      "placeholder",
	}
}
