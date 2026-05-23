package tasks

import (
	"chislennie-metodi/internal/taskmeta"
	"sort"

	"chislennie-metodi/internal/tasks/task01"
	"chislennie-metodi/internal/tasks/task02"
	"chislennie-metodi/internal/tasks/task03"
	"chislennie-metodi/internal/tasks/task04"
	"chislennie-metodi/internal/tasks/task05"
	"chislennie-metodi/internal/tasks/task06"
	"chislennie-metodi/internal/tasks/task07"
	"chislennie-metodi/internal/tasks/task08"
	"chislennie-metodi/internal/tasks/task09"
	"chislennie-metodi/internal/tasks/task10"
	"chislennie-metodi/internal/tasks/task11"
)

// allTaskProviders хранит ссылки на пакеты заданий.
// Это дает желаемую модульность: каждое задание живет в своем пакете,
// но фронтенд получает единый список из одного места.
var allTaskProviders = []func() taskmeta.Meta{
	task01.Meta,
	task02.Meta,
	task03.Meta,
	task04.Meta,
	task05.Meta,
	task06.Meta,
	task07.Meta,
	task08.Meta,
	task09.Meta,
	task10.Meta,
	task11.Meta,
}

func All() []taskmeta.Meta {
	result := make([]taskmeta.Meta, 0, len(allTaskProviders))
	for _, provider := range allTaskProviders {
		result = append(result, provider())
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].ID < result[j].ID
	})

	return result
}

// ByID находит одно задание по номеру.
func ByID(id int) (taskmeta.Meta, bool) {
	for _, task := range All() {
		if task.ID == id {
			return task, true
		}
	}

	return taskmeta.Meta{}, false
}
