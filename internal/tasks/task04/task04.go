package task04

import "chislennie-metodi/internal/taskmeta"

// Meta возвращает описание четвертого задания.
func Meta() taskmeta.Meta {
	return taskmeta.Meta{
		ID:          4,
		Title:       "Задание 4",
		Description: "Заготовка для четвертого модуля численных методов.",
		Status:      "placeholder",
	}
}
