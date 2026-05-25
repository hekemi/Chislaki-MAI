package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"

	"chislennie-metodi/internal/tasks"
	"chislennie-metodi/internal/tasks/task01"
	"chislennie-metodi/internal/tasks/task02"
	"chislennie-metodi/internal/tasks/task04"
	"chislennie-metodi/internal/tasks/task05"
	"chislennie-metodi/internal/tasks/task06"
	"chislennie-metodi/internal/tasks/task07"
	"chislennie-metodi/internal/tasks/task08"
	"chislennie-metodi/internal/tasks/task09"
	"chislennie-metodi/internal/tasks/task10"
	"chislennie-metodi/internal/tasks/task11"
)

//go:embed web
var embeddedWeb embed.FS

var (
	webRoot        fs.FS
	frontendHandler http.Handler
)

func init() {
	sub, err := fs.Sub(embeddedWeb, "web")
	if err != nil {
		log.Fatal(err)
	}
	webRoot = sub
	frontendHandler = http.FileServer(http.FS(webRoot))
}

// main запускает HTTP-сервер, который обрабатывает:
// - статический фронтенд из папки web,http://127.0.0.1:8080
// - API со списком заданий,
// - API для получения одного задания по номеру.
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/tasks", handleTasksList)
	mux.HandleFunc("/api/tasks/", handleTaskByID)
	mux.HandleFunc("/api/task1/solve", handleTask1Solve)
	mux.HandleFunc("/api/task2/solve", handleTask2Solve)
	mux.HandleFunc("/api/task4/solve", handleTask4Solve)
	mux.HandleFunc("/api/task5/solve", handleTask5Solve)
	mux.HandleFunc("/api/task6/solve", handleTask6Solve)
	mux.HandleFunc("/api/task7/solve", handleTask7Solve)
	mux.HandleFunc("/api/task8/solve", handleTask8Solve)
	mux.HandleFunc("/api/task9/solve", handleTask9Solve)
	mux.HandleFunc("/api/task10/solve", handleTask10Solve)
	mux.HandleFunc("/api/task11/solve", handleTask11Solve)
	mux.HandleFunc("/", handleFrontend)

	addr := "127.0.0.1:18080"
	log.Printf("server started at http://%s", addr)
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
func handleTaskByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/api/tasks/")
	if path == "" {
		writeJSON(w, http.StatusOK, tasks.All())
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

// handleTask1Solve принимает СЛАУ из фронтенда и возвращает решения
// методом Гаусса и методом простой итерации.
func handleTask1Solve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req task01.SolveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, task01.Solve(req))
}

// handleTask2Solve принимает СЛАУ из фронтенда и возвращает решения
// методом прогонки и методом Зейделя.
func handleTask2Solve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req task02.SolveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, task02.Solve(req))
}

// handleTask4Solve принимает нелинейное уравнение из фронтенда и возвращает решения
// методом бисекции, простой итерации и методом Ньютона.
func handleTask4Solve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req task04.SolveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}

	writeJSON(w, http.StatusOK, task04.Solve(req))
}

// handleFrontend отдает встроенные фронтенд-файлы из exe.
func handleFrontend(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		serveEmbeddedIndex(w)
		return
	}

	// Файлы со "точкой" считаем статикой: /app.js, /styles.css, /images/...
	if strings.Contains(path.Base(r.URL.Path), ".") {
		frontendHandler.ServeHTTP(w, r)
		return
	}

	// Для роутов SPA возвращаем index.html
	serveEmbeddedIndex(w)
}

func serveEmbeddedIndex(w http.ResponseWriter) {
	data, err := fs.ReadFile(webRoot, "index.html")
	if err != nil {
		http.Error(w, "index.html not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(data)
}

// handleTask5Solve принимает данные для интерполяции и возвращает результаты
func handleTask5Solve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Передаем управление в пакет task05
	task05.SolveHandler(w, r)
}

// handleTask6Solve принимает данные для сплайна и возвращает результаты.
func handleTask6Solve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	task06.SolveHandler(w, r)
}

// handleTask7Solve принимает данные для МНК и возвращает результаты.
func handleTask7Solve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	task07.SolveHandler(w, r)
}

// handleTask8Solve принимает данные для задания 8 (производные через сплайн)
func handleTask8Solve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	task08.SolveHandler(w, r)
}

// handleTask9Solve принимает данные для задания 9 (интеграл Симпсона)
func handleTask9Solve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	task09.SolveHandler(w, r)
}

// handleTask10Solve принимает данные для задания 10 (ODE RK4 + Adams)
func handleTask10Solve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	task10.SolveHandler(w, r)
}

// handleTask11Solve принимает данные для задания 11 (краевая задача, метод прогонки)
func handleTask11Solve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	task11.SolveHandler(w, r)
}

// writeJSON централизует отправку JSON-ответов и установку заголовков.
func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil {
		fmt.Println("json encode error:", err)
	}
}
