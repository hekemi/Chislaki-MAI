package task03

import "chislennie-metodi/internal/taskmeta"

// Meta возвращает описание третьего задания.
func Meta() taskmeta.Meta {
	return taskmeta.Meta{
		ID:          3,
		Title:       "Задание 3",
		Description: "Заготовка для третьего модуля численных методов.",
		Status:      "placeholder",
	}
}
