package task02

import "chislennie-metodi/internal/taskmeta"

// Meta возвращает описание второго задания.
func Meta() taskmeta.Meta {
	return taskmeta.Meta{
		ID:          2,
		Title:       "Задание 2",
		Description: "Решить СЛАУ методом прогонки и методом Зейделя",
		Status:      "ready",
	}
}
