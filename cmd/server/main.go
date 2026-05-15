package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"chislennie-metodi/internal/tasks"
)

// main запускает HTTP-сервер, который одновременно обслуживает:
// - статический фронтенд из папки web,
// - API со списком заданий,
// - API для получения одного задания по номеру.
//
// Такой каркас удобен для модульной разработки: каждое задание можно
// держать в отдельном пакете и расширять независимо от остальных.
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/tasks", handleTasksList)
	mux.HandleFunc("/api/tasks/", handleTaskByID)
	mux.HandleFunc("/", handleFrontend)

	addr := ":8080"
	log.Printf("server started at http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}

// handleTasksList возвращает все задания, которые будут показаны в левой панели.
func handleTasksList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	writeJSON(w, http.StatusOK, tasks.All())
}

// handleTaskByID возвращает одно задание по номеру.
// Этот endpoint пригодится позже, когда у каждого модуля появится собственная логика.
func handleTaskByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	if path == "" {
		http.Error(w, "task id is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(strings.Split(path, "/")[0])
	if err != nil {
		http.Error(w, "task id must be a number", http.StatusBadRequest)
		return
	}

	task, ok := tasks.ByID(id)
	if !ok {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	writeJSON(w, http.StatusOK, task)
}

// handleFrontend отдает фронтенд-файлы.
// Если запрошен конкретный файл из папки web, он отдается напрямую.
// Если файл не найден, возвращается index.html, чтобы приложение открывалось
// как единая оболочка и в будущем могло поддерживать маршруты по hash.
func handleFrontend(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		requested := filepath.Join("web", filepath.FromSlash(strings.TrimPrefix(r.URL.Path, "/")))
		if info, err := os.Stat(requested); err == nil && !info.IsDir() {
			http.ServeFile(w, r, requested)
			return
		}
	}

	http.ServeFile(w, r, filepath.Join("web", "index.html"))
}

// writeJSON централизует отправку JSON-ответов и установку заголовков.
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		fmt.Println("json encode error:", err)
	}
}
