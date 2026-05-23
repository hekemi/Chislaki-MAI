package task01

import (
	"fmt"
	"math"
	"strings"

	"chislennie-metodi/internal/taskmeta"
)

// SolveRequest описывает входные данные для решения СЛАУ.
//
// Поля специально сделаны универсальными, чтобы можно было решать не только
// один конкретный вариант из задания, но и произвольную квадратную СЛАУ.
type SolveRequest struct {
	Matrix        [][]float64 `json:"matrix"`
	Vector        []float64   `json:"vector"`
	Epsilon       float64     `json:"epsilon"`
	MaxIterations int         `json:"maxIterations"`
	InitialGuess  []float64   `json:"initialGuess"`
}

// IterationRow хранит одну строку таблицы для метода простой итерации.
// Здесь есть не только значения приближения, но и критерии остановки,
// чтобы их можно было показать и в таблице, и на графике.
type IterationRow struct {
	Iteration int       `json:"iteration"`
	Values    []float64 `json:"values"`
	Delta     float64   `json:"delta"`
	Residual  float64   `json:"residual"`
}

// GaussianResult хранит результат метода Гаусса.
type GaussianResult struct {
	Solution []float64      `json:"solution,omitempty"`
	Error    string         `json:"error,omitempty"`
	Steps    []GaussianStep `json:"steps,omitempty"`
}

// GaussianStep хранит один шаг преобразования расширенной матрицы.
// Такие шаги показываются на фронтенде как последовательность действий
// с матрицей при прямом ходе метода Гаусса.
type GaussianStep struct {
	Title  string      `json:"title"`
	Matrix [][]float64 `json:"matrix"`
}

// IterativeResult хранит результат метода простой итерации.
type IterativeResult struct {
	Solution   []float64      `json:"solution,omitempty"`
	Converged  bool           `json:"converged"`
	Iterations []IterationRow `json:"iterations,omitempty"`
	Error      string         `json:"error,omitempty"`
	FinalDelta float64        `json:"finalDelta,omitempty"`
	FinalResid float64        `json:"finalResidual,omitempty"`
}

// SolveResponse возвращает данные сразу для двух методов:
// метода Гаусса и метода простой итерации.
type SolveResponse struct {
	Title     string          `json:"title"`
	Input     SolveRequest    `json:"input"`
	Gaussian  GaussianResult  `json:"gaussian"`
	Iterative IterativeResult `json:"iterative"`
	Notes     []string        `json:"notes,omitempty"`
}

// DefaultRequest возвращает систему из варианта 13, которую можно сразу
// показать в интерфейсе как готовый пример для расчета.
func DefaultRequest() SolveRequest {
	return SolveRequest{
		Matrix: [][]float64{
			{2.82, 0.43, -0.57},
			{-0.35, 1.12, -0.48},
			{0.48, 0.23, 2.37},
		},
		Vector:        []float64{0.48, 0.52, 1.44},
		Epsilon:       0.01,
		MaxIterations: 100,
		InitialGuess:  []float64{0, 0, 0},
	}
}

// Solve решает СЛАУ обоими методами и собирает готовый ответ для фронтенда.
func Solve(req SolveRequest) SolveResponse {
	req = normalizeRequest(req)

	gaussianSolution, gaussianSteps, gaussianErr := solveGaussian(req.Matrix, req.Vector)
	iterativeSolution, iterations, converged, iterativeErr, finalDelta, finalResidual := solveJacobi(req.Matrix, req.Vector, req.Epsilon, req.MaxIterations, req.InitialGuess)

	notes := make([]string, 0, 2)
	if gaussianErr != nil {
		notes = append(notes, gaussianErr.Error())
	}
	if iterativeErr != nil {
		notes = append(notes, iterativeErr.Error())
	}

	resp := SolveResponse{
		Title:     "Задание 1: решение СЛАУ методом Гаусса и методом простой итерации",
		Input:     req,
		Gaussian:  GaussianResult{},
		Iterative: IterativeResult{},
		Notes:     notes,
	}

	if gaussianErr != nil {
		resp.Gaussian.Error = gaussianErr.Error()
	} else {
		resp.Gaussian.Solution = gaussianSolution
		resp.Gaussian.Steps = gaussianSteps
	}

	resp.Iterative = IterativeResult{
		Solution:   iterativeSolution,
		Converged:  converged,
		Iterations: iterations,
		FinalDelta: finalDelta,
		FinalResid: finalResidual,
	}
	if iterativeErr != nil {
		resp.Iterative.Error = iterativeErr.Error()
	}

	return resp
}

// normalizeRequest приводит входные данные к безопасному виду.
// Это нужно, чтобы фронтенд мог прислать неполный JSON, а сервер сам
// подставил значения по умолчанию.
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

// solveGaussian реализует метод Гаусса с частичным выбором главного элемента.
func solveGaussian(matrix [][]float64, vector []float64) ([]float64, []GaussianStep, error) {
	aug, err := buildAugmentedMatrix(matrix, vector)
	if err != nil {
		return nil, nil, err
	}

	n := len(aug)
	steps := []GaussianStep{{
		Title:  "Исходная расширенная матрица",
		Matrix: cloneMatrix(aug),
	}}

	for col := 0; col < n; col++ {
		pivotRow := col
		pivotValue := math.Abs(aug[col][col])
		for row := col + 1; row < n; row++ {
			if value := math.Abs(aug[row][col]); value > pivotValue {
				pivotValue = value
				pivotRow = row
			}
		}

		if pivotValue < 1e-12 {
			return nil, steps, fmt.Errorf("метод Гаусса не может продолжить работу: матрица вырождена")
		}

		if pivotRow != col {
			aug[col], aug[pivotRow] = aug[pivotRow], aug[col]
			steps = append(steps, GaussianStep{
				Title:  fmt.Sprintf("Перестановка строк %d и %d", col+1, pivotRow+1),
				Matrix: cloneMatrix(aug),
			})
		}

		for row := col + 1; row < n; row++ {
			factor := aug[row][col] / aug[col][col]
			for k := col; k <= n; k++ {
				aug[row][k] -= factor * aug[col][k]
			}
			steps = append(steps, GaussianStep{
				Title:  fmt.Sprintf("R%d <- R%d - %.4f*R%d", row+1, row+1, factor, col+1),
				Matrix: cloneMatrix(aug),
			})
		}
	}

	solution := make([]float64, n)
	for row := n - 1; row >= 0; row-- {
		sum := aug[row][n]
		for col := row + 1; col < n; col++ {
			sum -= aug[row][col] * solution[col]
		}
		if math.Abs(aug[row][row]) < 1e-12 {
			return nil, steps, fmt.Errorf("метод Гаусса не может продолжить работу: нулевой ведущий элемент")
		}
		solution[row] = sum / aug[row][row]
	}

	steps = append(steps, GaussianStep{
		Title:  "Треугольный вид после прямого хода",
		Matrix: cloneMatrix(aug),
	})

	return solution, steps, nil
}

// solveJacobi реализует метод простой итерации для СЛАУ.
//
// Внутри метода мы записываем все приближения. Именно они потом идут в таблицу,
// а значения residual можно использовать для построения графика сходимости.
func solveJacobi(matrix [][]float64, vector []float64, epsilon float64, maxIterations int, initialGuess []float64) ([]float64, []IterationRow, bool, error, float64, float64) {
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
		next := make([]float64, n)
		for i := 0; i < n; i++ {
			if math.Abs(matrix[i][i]) < 1e-12 {
				return nil, iterations, false, fmt.Errorf("метод простой итерации не может работать: нулевой диагональный элемент в строке %d", i+1), lastDelta, lastResidual
			}

			sum := vector[i]
			for j := 0; j < n; j++ {
				if j == i {
					continue
				}
				sum -= matrix[i][j] * current[j]
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

	return current, iterations, false, fmt.Errorf("метод простой итерации не сошелся за %d итераций", maxIterations), lastDelta, lastResidual
}

// buildAugmentedMatrix создает расширенную матрицу [A|b].
func buildAugmentedMatrix(matrix [][]float64, vector []float64) ([][]float64, error) {
	if err := validateSystem(matrix, vector); err != nil {
		return nil, err
	}

	n := len(matrix)
	aug := make([][]float64, n)
	for i := 0; i < n; i++ {
		aug[i] = make([]float64, n+1)
		copy(aug[i], matrix[i])
		aug[i][n] = vector[i]
	}
	return aug, nil
}

// validateSystem проверяет, что матрица является квадратной и согласована с вектором.
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

// cloneVector создает независимую копию вектора.
func cloneVector(src []float64) []float64 {
	if src == nil {
		return nil
	}
	dup := make([]float64, len(src))
	copy(dup, src)
	return dup
}

// cloneMatrix создает глубокую копию матрицы, чтобы шаги Гаусса
// сохраняли состояние матрицы на каждом этапе преобразования.
func cloneMatrix(src [][]float64) [][]float64 {
	if src == nil {
		return nil
	}
	dup := make([][]float64, len(src))
	for i := range src {
		dup[i] = make([]float64, len(src[i]))
		copy(dup[i], src[i])
	}
	return dup
}

// maxAbsDiff возвращает максимальное абсолютное различие между двумя векторами.
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

// residualNorm оценивает невязку ||Ax - b||_inf для текущего приближения.
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

// FormatEquation создает строковое представление одного уравнения.
// Эта функция полезна для красивого вывода примера на фронтенде.
func FormatEquation(row []float64, rhs float64) string {
	parts := make([]string, 0, len(row)+1)
	for i, coefficient := range row {
		parts = append(parts, fmt.Sprintf("%+.2fx%d", coefficient, i+1))
	}
	return strings.Join(parts, " ") + fmt.Sprintf(" = %.2f", rhs)
}

// RecommendedMeta возвращает текст, который можно показать в карточке задания.
func RecommendedMeta() taskmeta.Meta {
	return taskmeta.Meta{
		ID:          1,
		Title:       "Задание 1",
		Description: "Решение СЛАУ методом Гаусса и методом простой итерации",
		Status:      "ready",
	}
}
