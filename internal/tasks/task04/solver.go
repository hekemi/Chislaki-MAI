package task04

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
)

// ParsedEquation содержит информацию о разобранном уравнении
type ParsedEquation struct {
	Expression string // исходное выражение
}

// EvaluateFunc вычисляет f(x)
func (eq *ParsedEquation) EvaluateFunc(x float64) (float64, error) {
	expr := eq.Expression

	// Поддержка записи e^x и e^(...) — приводим к exp(...)
	expr = strings.ReplaceAll(expr, "e^(", "exp(")
	expr = strings.ReplaceAll(expr, "e^x", "exp(x)")

	// Если уравнение задано через '=', приводим к форме left-right
	if strings.Contains(expr, "=") {
		parts := strings.SplitN(expr, "=", 2)
		expr = fmt.Sprintf("(%s)-(%s)", parts[0], parts[1])
	}

	// Заменяем отдельную константу e на её числовое значение
	reE := regexp.MustCompile(`\be\b`)
	expr = reE.ReplaceAllString(expr, "("+strconv.FormatFloat(math.E, 'g', -1, 64)+")")

	// Заменяем переменную x, включая случаи типа 3x, но не внутри идентификаторов (например exp)
	var b strings.Builder
	runes := []rune(expr)
	for i, r := range runes {
		if r == 'x' {
			prevIsLetter := false
			if i-1 >= 0 {
				prev := runes[i-1]
				if (prev >= 'A' && prev <= 'Z') || (prev >= 'a' && prev <= 'z') || prev == '_' {
					prevIsLetter = true
				}
			}

			if !prevIsLetter {
				b.WriteString("(" + strconv.FormatFloat(x, 'g', -1, 64) + ")")
				continue
			}
		}
		b.WriteRune(r)
	}
	expr = b.String()

	// Вставляем явный оператор умножения там, где было неявное: число|')' перед '('
	expr = insertImplicitMultiplication(expr)

	return evaluateExpression(expr)
}

// EvaluateDerivative вычисляет f'(x) численно
func (eq *ParsedEquation) EvaluateDerivative(x float64) (float64, error) {
	h := 1e-8
	f1, err := eq.EvaluateFunc(x + h)
	if err != nil {
		return 0, err
	}
	f2, err := eq.EvaluateFunc(x - h)
	if err != nil {
		return 0, err
	}
	return (f1 - f2) / (2 * h), nil
}

// Iteration содержит информацию об одной итерации
type Iteration struct {
	Iteration int     `json:"iteration"`
	X         float64 `json:"x"`
	FX        float64 `json:"fx"`
	Delta     float64 `json:"delta"`
	Residual  float64 `json:"residual"`
	ErrorMsg  string  `json:"error_msg,omitempty"`
}

// BisectionResult результат метода бисекции
type BisectionResult struct {
	Solution      float64     `json:"solution"`
	Error         string      `json:"error"`
	Converged     bool        `json:"converged"`
	Iterations    []Iteration `json:"iterations"`
	FinalDelta    float64     `json:"final_delta"`
	FinalResidual float64     `json:"final_residual"`
}

// SimpleIterationResult результат метода простой итерации
type SimpleIterationResult struct {
	Solution      float64     `json:"solution"`
	Error         string      `json:"error"`
	Converged     bool        `json:"converged"`
	Iterations    []Iteration `json:"iterations"`
	FinalDelta    float64     `json:"final_delta"`
	FinalResidual float64     `json:"final_residual"`
}

// NewtonResult результат метода Ньютона
type NewtonResult struct {
	Solution      float64     `json:"solution"`
	Error         string      `json:"error"`
	Converged     bool        `json:"converged"`
	Iterations    []Iteration `json:"iterations"`
	FinalDelta    float64     `json:"final_delta"`
	FinalResidual float64     `json:"final_residual"`
}

// SolveRequest запрос на решение уравнения
type SolveRequest struct {
	Equation      string  `json:"equation"`
	A             float64 `json:"a"`
	B             float64 `json:"b"`
	Epsilon       float64 `json:"epsilon"`
	MaxIterations int     `json:"max_iterations"`
	X0            float64 `json:"x0"` // начальное приближение для простой итерации и Ньютона
}

// SolveResponse ответ с решениями всеми тремя методами
type SolveResponse struct {
	Bisection       BisectionResult       `json:"bisection"`
	SimpleIteration SimpleIterationResult `json:"simple_iteration"`
	Newton          NewtonResult          `json:"newton"`
}

// evaluateExpression вычисляет математическое выражение
func evaluateExpression(expr string) (float64, error) {
	// Простая подстановка и вычисление для базовых операций
	// Поддерживаем: +, -, *, /, exp(), sin(), cos(), sqrt()

	expr = strings.TrimSpace(expr)

	// Заменяем встроенные функции на их значения (упрощённо)
	for {
		prevExpr := expr
		expr = evaluateBuiltins(expr)
		if expr == prevExpr {
			break
		}
	}

	result, err := parseAndEvaluate(expr)
	return result, err
}

// evaluateBuiltins заменяет встроенные функции на их значения
func evaluateBuiltins(expr string) string {
	// exp(...)
	expr = replaceFunction(expr, "exp", func(x float64) float64 { return math.Exp(x) })
	// sin(...)
	expr = replaceFunction(expr, "sin", func(x float64) float64 { return math.Sin(x) })
	// cos(...)
	expr = replaceFunction(expr, "cos", func(x float64) float64 { return math.Cos(x) })
	// sqrt(...)
	expr = replaceFunction(expr, "sqrt", func(x float64) float64 { return math.Sqrt(x) })
	return expr
}

// replaceFunction заменяет функцию вида func(x) на её числовое значение
func replaceFunction(expr, funcName string, fn func(float64) float64) string {
	for {
		idx := strings.Index(expr, funcName+"(")
		if idx == -1 {
			break
		}

		// Найдём закрывающую скобку
		depth := 0
		endIdx := -1
		for i := idx + len(funcName) + 1; i < len(expr); i++ {
			if expr[i] == '(' {
				depth++
			} else if expr[i] == ')' {
				if depth == 0 {
					endIdx = i
					break
				}
				depth--
			}
		}

		if endIdx == -1 {
			break
		}

		argExpr := expr[idx+len(funcName)+1 : endIdx]
		val, err := parseAndEvaluate(argExpr)
		if err != nil {
			break
		}

		result := fn(val)
		expr = expr[:idx] + fmt.Sprintf("%f", result) + expr[endIdx+1:]
	}
	return expr
}

// parseAndEvaluate парсит и вычисляет основное выражение
func parseAndEvaluate(expr string) (float64, error) {
	expr = strings.TrimSpace(expr)

	// Убираем лишние скобки
	for strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
		expr = expr[1 : len(expr)-1]
		expr = strings.TrimSpace(expr)
	}

	// Парсим выражение слева направо
	return parseAdditive(expr)
}

func parseAdditive(expr string) (float64, error) {
	// Разбираем выражение слева направо, учитывая скобки
	start := 0
	depth := 0
	var parts []float64
	var ops []rune

	for i, ch := range expr {
		if ch == '(' {
			depth++
		} else if ch == ')' {
			depth--
		} else if depth == 0 && (ch == '+' || ch == '-') {
			partStr := strings.TrimSpace(expr[start:i])
			val, err := parseMultiplicative(partStr)
			if err != nil {
				return 0, err
			}
			parts = append(parts, val)
			ops = append(ops, ch)
			start = i + 1
		}
	}

	lastPart := strings.TrimSpace(expr[start:])
	v, err := parseMultiplicative(lastPart)
	if err != nil {
		return 0, err
	}
	parts = append(parts, v)

	result := parts[0]
	for i, op := range ops {
		if op == '+' {
			result += parts[i+1]
		} else {
			result -= parts[i+1]
		}
	}

	return result, nil
}

func parseMultiplicative(expr string) (float64, error) {
	// Разделяем по * и / на верхнем уровне, затем используем степень/primary
	start := 0
	depth := 0
	var parts []float64
	var ops []rune

	for i, ch := range expr {
		if ch == '(' {
			depth++
		} else if ch == ')' {
			depth--
		} else if depth == 0 && (ch == '*' || ch == '/') {
			partStr := strings.TrimSpace(expr[start:i])
			val, err := parsePower(partStr)
			if err != nil {
				return 0, err
			}
			parts = append(parts, val)
			ops = append(ops, ch)
			start = i + 1
		}
	}

	lastPart := strings.TrimSpace(expr[start:])
	v, err := parsePower(lastPart)
	if err != nil {
		return 0, err
	}
	parts = append(parts, v)

	result := parts[0]
	for i, op := range ops {
		if op == '*' {
			result *= parts[i+1]
		} else {
			result /= parts[i+1]
		}
	}

	return result, nil
}

// parsePower обрабатывает оператор возведения в степень '^' (право-ассоциативно)
func parsePower(expr string) (float64, error) {
	expr = strings.TrimSpace(expr)
	depth := 0
	last := -1
	for i := 0; i < len(expr); i++ {
		ch := expr[i]
		if ch == '(' {
			depth++
		} else if ch == ')' {
			depth--
		} else if depth == 0 && ch == '^' {
			last = i
		}
	}

	if last == -1 {
		return parsePrimary(expr)
	}

	leftStr := strings.TrimSpace(expr[:last])
	rightStr := strings.TrimSpace(expr[last+1:])

	leftVal, err := parsePrimary(leftStr)
	if err != nil {
		return 0, err
	}
	rightVal, err := parsePower(rightStr)
	if err != nil {
		return 0, err
	}

	return math.Pow(leftVal, rightVal), nil
}

// insertImplicitMultiplication вставляет '*' между числом или ')' и '('
func insertImplicitMultiplication(s string) string {
	var b strings.Builder
	var prev rune
	for i, ch := range s {
		if i > 0 {
			if (prev >= '0' && prev <= '9' || prev == ')') && ch == '(' {
				b.WriteRune('*')
			}
		}
		b.WriteRune(ch)
		prev = ch
	}
	return b.String()
}

func parsePrimary(expr string) (float64, error) {
	expr = strings.TrimSpace(expr)

	if strings.HasPrefix(expr, "(") && strings.HasSuffix(expr, ")") {
		return parseAndEvaluate(expr[1 : len(expr)-1])
	}

	val := 0.0
	_, err := fmt.Sscanf(expr, "%f", &val)
	return val, err
}

func splitByTopLevel(expr string, ops ...rune) []string {
	var result []string
	var current strings.Builder
	depth := 0

	for _, ch := range expr {
		if ch == '(' {
			depth++
		} else if ch == ')' {
			depth--
		} else if depth == 0 {
			for _, op := range ops {
				if ch == op {
					result = append(result, current.String())
					current.Reset()
					continue
				}
			}
		}
		current.WriteRune(ch)
	}

	if current.Len() > 0 {
		result = append(result, current.String())
	}

	return result
}

// SolveBisection решает уравнение методом бисекции
func SolveBisection(eq *ParsedEquation, a, b, epsilon float64, maxIter int) BisectionResult {
	result := BisectionResult{
		Iterations: []Iteration{},
	}

	fa, err := eq.EvaluateFunc(a)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	fb, err := eq.EvaluateFunc(b)
	if err != nil {
		result.Error = err.Error()
		return result
	}

	if fa*fb > 0 {
		result.Error = "f(a) и f(b) должны иметь разные знаки"
		return result
	}

	prevX := a
	for iter := 0; iter < maxIter; iter++ {
		c := (a + b) / 2
		fc, _ := eq.EvaluateFunc(c)

		delta := math.Abs(c - prevX)
		residual := math.Abs(fc)

		result.Iterations = append(result.Iterations, Iteration{
			Iteration: iter + 1,
			X:         c,
			FX:        fc,
			Delta:     delta,
			Residual:  residual,
		})

		result.FinalDelta = delta
		result.FinalResidual = residual

		if residual < epsilon || delta < epsilon {
			result.Solution = c
			result.Converged = true
			return result
		}

		if fa*fc < 0 {
			b = c
			fb = fc
		} else {
			a = c
			fa = fc
		}

		prevX = c
	}

	result.Solution = (a + b) / 2
	return result
}

// SolveSimpleIteration решает уравнение методом простой итерации
func SolveSimpleIteration(eq *ParsedEquation, x0, epsilon float64, maxIter int) SimpleIterationResult {
	result := SimpleIterationResult{
		Iterations: []Iteration{},
	}

	// Для простой итерации нужно привести f(x)=0 к виду x=g(x)
	// g(x) = x - 0.1*f(x) (масштабированное приближение)

	x := x0
	for iter := 0; iter < maxIter; iter++ {
		fx, err := eq.EvaluateFunc(x)
		if err != nil {
			result.Error = err.Error()
			return result
		}

		// Используем x_next = x - 0.1*f(x) как итерационное уравнение
		xNext := x - 0.1*fx

		delta := math.Abs(xNext - x)
		residual := math.Abs(fx)

		result.Iterations = append(result.Iterations, Iteration{
			Iteration: iter + 1,
			X:         xNext,
			FX:        fx,
			Delta:     delta,
			Residual:  residual,
		})

		result.FinalDelta = delta
		result.FinalResidual = residual

		if residual < epsilon || delta < epsilon {
			result.Solution = xNext
			result.Converged = true
			return result
		}

		x = xNext
	}

	result.Solution = x
	return result
}

// SolveNewton решает уравнение методом Ньютона
func SolveNewton(eq *ParsedEquation, x0, epsilon float64, maxIter int) NewtonResult {
	result := NewtonResult{
		Iterations: []Iteration{},
	}

	x := x0
	for iter := 0; iter < maxIter; iter++ {
		fx, err := eq.EvaluateFunc(x)
		if err != nil {
			result.Error = err.Error()
			return result
		}

		fpx, err := eq.EvaluateDerivative(x)
		if err != nil {
			result.Error = err.Error()
			return result
		}

		if math.Abs(fpx) < 1e-10 {
			result.Error = "производная близка к нулю"
			return result
		}

		xNext := x - fx/fpx

		delta := math.Abs(xNext - x)
		residual := math.Abs(fx)

		result.Iterations = append(result.Iterations, Iteration{
			Iteration: iter + 1,
			X:         xNext,
			FX:        fx,
			Delta:     delta,
			Residual:  residual,
		})

		result.FinalDelta = delta
		result.FinalResidual = residual

		if residual < epsilon || delta < epsilon {
			result.Solution = xNext
			result.Converged = true
			return result
		}

		x = xNext
	}

	result.Solution = x
	return result
}

// Solve решает уравнение всеми тремя методами
func Solve(req SolveRequest) SolveResponse {
	eq := &ParsedEquation{Expression: req.Equation}

	return SolveResponse{
		Bisection:       SolveBisection(eq, req.A, req.B, req.Epsilon, req.MaxIterations),
		SimpleIteration: SolveSimpleIteration(eq, req.X0, req.Epsilon, req.MaxIterations),
		Newton:          SolveNewton(eq, req.X0, req.Epsilon, req.MaxIterations),
	}
}
