package task07

import "chislennie-metodi/internal/taskmeta"

// Meta возвращает описание седьмого задания.
func Meta() taskmeta.Meta {
	return taskmeta.Meta{
		ID:          7,
		Title:       "Задание 7",
		Description: "Аппроксимация линейным и квадратичным многочленом по заданной таблице.",
		Status:      "ready",
	}
}
