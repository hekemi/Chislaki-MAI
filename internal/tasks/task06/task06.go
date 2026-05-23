package task06

import "chislennie-metodi/internal/taskmeta"

// Meta возвращает описание шестого задания.
func Meta() taskmeta.Meta {
	return taskmeta.Meta{
		ID:          6,
		Title:       "Задание 6",
		Description: "Построение натурального кубического сплайна и таблицы коэффициентов.",
		Status:      "ready",
	}
}
