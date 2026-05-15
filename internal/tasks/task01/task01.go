package task01

import "chislennie-metodi/internal/taskmeta"

// Meta возвращает описание первого задания.
// Пока это заглушка, но именно здесь позже можно разместить свой модуль,
// функции расчета и структуру данных для первого численного метода.
func Meta() taskmeta.Meta {
	return taskmeta.Meta{
		ID:          1,
		Title:       "Задание 1",
		Description: "Подключи сюда первый численный метод и его форму ввода.",
		Status:      "placeholder",
	}
}
