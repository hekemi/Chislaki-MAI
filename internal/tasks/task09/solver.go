package task09

import (
	"chislennie-metodi/internal/tasks/task04"
	"encoding/json"
	"math"
	"net/http"
)

type SolveRequest struct {
	Equation string  `json:"equation"`
	A        float64 `json:"a"`
	B        float64 `json:"b"`
	Epsilon  float64 `json:"epsilon"`
	MaxIter  int     `json:"max_iterations"`
}

type IterRow struct {
	Iteration int     `json:"iteration"`
	X         float64 `json:"x"`
	Delta     float64 `json:"delta"`
	Residual  float64 `json:"residual"`
}

type SolveResponse struct {
	Iterations []IterRow `json:"iterations"`
	Result     float64   `json:"result"`
}

func compositeSimpson(f func(float64) (float64, error), a, b float64, n int) (float64, error) {
	if n%2 == 1 {
		n++
	}
	h := (b - a) / float64(n)
	sum := 0.0
	for i := 0; i <= n; i++ {
		x := a + float64(i)*h
		fx, err := f(x)
		if err != nil {
			return 0, err
		}
		coef := 1.0
		if i == 0 || i == n {
			coef = 1.0
		} else if i%2 == 1 {
			coef = 4.0
		} else {
			coef = 2.0
		}
		sum += coef * fx
	}
	return (h / 3.0) * sum, nil
}

func SolveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var req SolveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	if req.MaxIter <= 0 {
		req.MaxIter = 10
	}
	if req.Epsilon <= 0 {
		req.Epsilon = 1e-6
	}
	// use task04 parser
	eq := task04.ParsedEquation{Expression: req.Equation}
	f := func(x float64) (float64, error) {
		return eq.EvaluateFunc(x)
	}

	iterations := make([]IterRow, 0)
	n := 2
	prev, err := compositeSimpson(f, req.A, req.B, n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	iterations = append(iterations, IterRow{Iteration: n, X: prev, Delta: 0, Residual: 0})

	for k := 0; k < req.MaxIter; k++ {
		n = n * 2
		cur, err := compositeSimpson(f, req.A, req.B, n)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		delta := math.Abs(cur - prev)
		residual := delta / 15.0 // Simpson error estimate
		iterations = append(iterations, IterRow{Iteration: n, X: cur, Delta: delta, Residual: residual})
		if residual < req.Epsilon {
			// done
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(SolveResponse{Iterations: iterations, Result: cur})
			return
		}
		prev = cur
	}
	// if reached here, return last
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(SolveResponse{Iterations: iterations, Result: prev})
}
