package task04

import "chislennie-metodi/internal/taskmeta"

// Meta возвращает описание четвертого задания.
func Meta() taskmeta.Meta {
	return taskmeta.Meta{
		ID:          4,
		Title:       "Задание 4",
		Description: "Решить нелинейное уравнение методами бисекции, простой итерации и Ньютона",
		Status:      "ready",
	}
}
