package task02

import (
	"fmt"
	"math"

	"chislennie-metodi/internal/taskmeta"
)

// SolveRequest описывает входные данные для второго задания.
type SolveRequest struct {
	Matrix        [][]float64 `json:"matrix"`
	Vector        []float64   `json:"vector"`
	Epsilon       float64     `json:"epsilon"`
	MaxIterations int         `json:"maxIterations"`
	InitialGuess  []float64   `json:"initialGuess"`
}

// IterationRow хранит одну итерацию метода Зейделя.
type IterationRow struct {
	Iteration int       `json:"iteration"`
	Values    []float64 `json:"values"`
	Delta     float64   `json:"delta"`
	Residual  float64   `json:"residual"`
}

// SweepResult хранит результат метода прогонки.
type SweepResult struct {
	Solution      []float64      `json:"solution,omitempty"`
	Converged     bool           `json:"converged"`
	Iterations    []IterationRow `json:"iterations,omitempty"`
	FinalDelta    float64        `json:"finalDelta,omitempty"`
	FinalResidual float64        `json:"finalResidual,omitempty"`
	Error         string         `json:"error,omitempty"`
}

// SeidelResult хранит результат метода Зейделя.
type SeidelResult struct {
	Solution   []float64      `json:"solution,omitempty"`
	Converged  bool           `json:"converged"`
	Iterations []IterationRow `json:"iterations,omitempty"`
	Error      string         `json:"error,omitempty"`
	FinalDelta float64        `json:"finalDelta,omitempty"`
	FinalResid float64        `json:"finalResidual,omitempty"`
}

// SolveResponse возвращает решение двумя методами.
type SolveResponse struct {
	Title  string       `json:"title"`
	Input  SolveRequest `json:"input"`
	Sweep  SweepResult  `json:"sweep"`
	Seidel SeidelResult `json:"seidel"`
	Notes  []string     `json:"notes,omitempty"`
}

// DefaultRequest возвращает исходную систему из задания 2.
func DefaultRequest() SolveRequest {
	return SolveRequest{
		Matrix: [][]float64{
			{8, 2, 0, 0},
			{-3, 9, -2, 0},
			{0, 1, 10, 1},
			{0, 0, 1, 6},
		},
		Vector:        []float64{15, 5.5, 15, 9.5},
		Epsilon:       0.01,
		MaxIterations: 100,
		InitialGuess:  []float64{0, 0, 0, 0},
	}
}

// Solve решает систему методом прогонки и методом Зейделя.
func Solve(req SolveRequest) SolveResponse {
	req = normalizeRequest(req)

	sweepSolution, sweepErr := solveSweep(req.Matrix, req.Vector)
	seidelSolution, seidelIterations, converged, seidelErr, finalDelta, finalResidual := solveSeidel(req.Matrix, req.Vector, req.Epsilon, req.MaxIterations, req.InitialGuess)

	notes := make([]string, 0, 2)
	if sweepErr != nil {
		notes = append(notes, sweepErr.Error())
	}
	if seidelErr != nil {
		notes = append(notes, seidelErr.Error())
	}

	resp := SolveResponse{
		Title: "Задание 2: метод прогонки и метод Зейделя",
		Input: req,
		Sweep: SweepResult{},
		Seidel: SeidelResult{
			Converged: converged,
		},
		Notes: notes,
	}

	if sweepErr != nil {
		resp.Sweep.Error = sweepErr.Error()
	} else {
		resp.Sweep.Solution = sweepSolution
		resp.Sweep.Converged = true
		// Для метода прогонки — прямой метод, итераций нет, но вернём одну запись
		// с корректной невязкой, чтобы фронтенд мог отобразить значение и график.
		finalRes := residualNorm(req.Matrix, sweepSolution, req.Vector)
		resp.Sweep.FinalDelta = 0
		resp.Sweep.FinalResidual = finalRes
		resp.Sweep.Iterations = []IterationRow{{
			Iteration: 1,
			Values:    cloneVector(sweepSolution),
			Delta:     0,
			Residual:  finalRes,
		}}
	}

	resp.Seidel.Solution = seidelSolution
	resp.Seidel.Iterations = seidelIterations
	resp.Seidel.FinalDelta = finalDelta
	resp.Seidel.FinalResid = finalResidual
	if seidelErr != nil {
		resp.Seidel.Error = seidelErr.Error()
	}

	return resp
}

func normalizeRequest(req SolveRequest) SolveRequest {
	if len(req.Matrix) == 0 || len(req.Vector) == 0 {
		req = DefaultRequest()
	}

	n := len(req.Matrix)
	if req.Epsilon <= 0 {
		req.Epsilon = 0.01
	}
	if req.MaxIterations <= 0 {
		req.MaxIterations = 100
	}
	if len(req.InitialGuess) != n {
		req.InitialGuess = make([]float64, n)
	}

	return req
}

// solveSweep реализует метод прогонки (алгоритм Томаса).
func solveSweep(matrix [][]float64, vector []float64) ([]float64, error) {
	if err := validateSystem(matrix, vector); err != nil {
		return nil, err
	}
	if err := validateTridiagonal(matrix); err != nil {
		return nil, err
	}

	n := len(matrix)
	a := make([]float64, n)
	b := make([]float64, n)
	c := make([]float64, n)
	d := cloneVector(vector)

	for i := 0; i < n; i++ {
		b[i] = matrix[i][i]
		if i > 0 {
			a[i] = matrix[i][i-1]
		}
		if i < n-1 {
			c[i] = matrix[i][i+1]
		}
	}

	alpha := make([]float64, n)
	beta := make([]float64, n)

	if math.Abs(b[0]) < 1e-12 {
		return nil, fmt.Errorf("метод прогонки не может работать: нулевой главный элемент в первой строке")
	}
	alpha[0] = -c[0] / b[0]
	beta[0] = d[0] / b[0]

	for i := 1; i < n; i++ {
		denominator := b[i] + a[i]*alpha[i-1]
		if math.Abs(denominator) < 1e-12 {
			return nil, fmt.Errorf("метод прогонки не может продолжить работу: нулевой знаменатель в строке %d", i+1)
		}

		if i < n-1 {
			alpha[i] = -c[i] / denominator
		}
		beta[i] = (d[i] - a[i]*beta[i-1]) / denominator
	}

	x := make([]float64, n)
	x[n-1] = beta[n-1]
	for i := n - 2; i >= 0; i-- {
		x[i] = alpha[i]*x[i+1] + beta[i]
	}

	return x, nil
}

// solveSeidel реализует итерационный метод Зейделя.
func solveSeidel(matrix [][]float64, vector []float64, epsilon float64, maxIterations int, initialGuess []float64) ([]float64, []IterationRow, bool, error, float64, float64) {
	if err := validateSystem(matrix, vector); err != nil {
		return nil, nil, false, err, 0, 0
	}

	n := len(matrix)
	current := cloneVector(initialGuess)
	if len(current) != n {
		current = make([]float64, n)
	}

	iterations := make([]IterationRow, 0, maxIterations)
	var lastDelta float64
	var lastResidual float64

	for iter := 1; iter <= maxIterations; iter++ {
		next := cloneVector(current)

		for i := 0; i < n; i++ {
			if math.Abs(matrix[i][i]) < 1e-12 {
				return nil, iterations, false, fmt.Errorf("метод Зейделя не может работать: нулевой диагональный элемент в строке %d", i+1), lastDelta, lastResidual
			}

			sum := vector[i]
			for j := 0; j < n; j++ {
				if j == i {
					continue
				}
				if j < i {
					sum -= matrix[i][j] * next[j]
				} else {
					sum -= matrix[i][j] * current[j]
				}
			}
			next[i] = sum / matrix[i][i]
		}

		lastDelta = maxAbsDiff(next, current)
		lastResidual = residualNorm(matrix, next, vector)
		iterations = append(iterations, IterationRow{
			Iteration: iter,
			Values:    cloneVector(next),
			Delta:     lastDelta,
			Residual:  lastResidual,
		})

		current = next
		if lastDelta < epsilon || lastResidual < epsilon {
			return current, iterations, true, nil, lastDelta, lastResidual
		}
	}

	return current, iterations, false, fmt.Errorf("метод Зейделя не сошелся за %d итераций", maxIterations), lastDelta, lastResidual
}

func validateSystem(matrix [][]float64, vector []float64) error {
	if len(matrix) == 0 {
		return fmt.Errorf("матрица системы не должна быть пустой")
	}
	if len(matrix) != len(vector) {
		return fmt.Errorf("размер матрицы и вектора не совпадает")
	}

	n := len(matrix)
	for i := 0; i < n; i++ {
		if len(matrix[i]) != n {
			return fmt.Errorf("матрица должна быть квадратной")
		}
	}

	return nil
}

func validateTridiagonal(matrix [][]float64) error {
	for i := range matrix {
		for j := range matrix[i] {
			if math.Abs(float64(i-j)) > 1 && math.Abs(matrix[i][j]) > 1e-12 {
				return fmt.Errorf("метод прогонки применим только к трехдиагональной матрице")
			}
		}
	}

	return nil
}

func cloneVector(src []float64) []float64 {
	if src == nil {
		return nil
	}
	dup := make([]float64, len(src))
	copy(dup, src)
	return dup
}

func maxAbsDiff(a, b []float64) float64 {
	limit := len(a)
	if len(b) < limit {
		limit = len(b)
	}
	maxValue := 0.0
	for i := 0; i < limit; i++ {
		if value := math.Abs(a[i] - b[i]); value > maxValue {
			maxValue = value
		}
	}
	return maxValue
}

func residualNorm(matrix [][]float64, x []float64, vector []float64) float64 {
	maxValue := 0.0
	for i := range matrix {
		sum := 0.0
		for j := range matrix[i] {
			sum += matrix[i][j] * x[j]
		}
		if value := math.Abs(sum - vector[i]); value > maxValue {
			maxValue = value
		}
	}
	return maxValue
}

// RecommendedMeta возвращает метаинформацию второго задания.
func RecommendedMeta() taskmeta.Meta {
	return taskmeta.Meta{
		ID:          2,
		Title:       "Задание 2",
		Description: "Решить СЛАУ методом прогонки и методом Зейделя",
		Status:      "ready",
	}
}
