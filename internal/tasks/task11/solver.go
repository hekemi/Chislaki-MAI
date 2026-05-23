package task11

import (
	"chislennie-metodi/internal/tasks/task04"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type SolveRequest struct {
	PExpr string  `json:"p_expr"`
	QExpr string  `json:"q_expr"`
	RHS   string  `json:"rhs_expr"`
	A     float64 `json:"a"`
	B     float64 `json:"b"`
	H     float64 `json:"h"`

	LeftAlpha  float64 `json:"left_alpha"`
	LeftBeta   float64 `json:"left_beta"`
	LeftGamma  float64 `json:"left_gamma"`
	RightAlpha float64 `json:"right_alpha"`
	RightBeta  float64 `json:"right_beta"`
	RightGamma float64 `json:"right_gamma"`
}

type Point struct {
	Index int     `json:"index"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
}

type SweepRow struct {
	Index int     `json:"index"`
	X     float64 `json:"x"`
	A     float64 `json:"a"`
	B     float64 `json:"b"`
	C     float64 `json:"c"`
	D     float64 `json:"d"`
	Alpha float64 `json:"alpha"`
	Beta  float64 `json:"beta"`
	Y     float64 `json:"y"`
}

type SolveResponse struct {
	Points   []Point    `json:"points"`
	Sweep    []SweepRow `json:"sweep"`
	StepUsed float64    `json:"step_used"`
}

func compileExpr(expr string) func(float64) (float64, error) {
	expr = strings.TrimSpace(expr)
	if v, err := strconv.ParseFloat(expr, 64); err == nil {
		return func(float64) (float64, error) {
			return v, nil
		}
	}
	if strings.HasPrefix(expr, "-") || strings.HasPrefix(expr, "+") {
		// Парсер task04 не поддерживает унарный знак в начале выражения, нормализуем к бинарному виду.
		expr = "0" + expr
	}

	eq := task04.ParsedEquation{Expression: expr}
	return func(x float64) (float64, error) {
		return eq.EvaluateFunc(x)
	}
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

	// Значения по умолчанию для варианта из задания.
	if req.PExpr == "" {
		req.PExpr = "2/x"
	}
	if req.QExpr == "" {
		req.QExpr = "-3"
	}
	if req.RHS == "" {
		req.RHS = "2"
	}
	if req.A == 0 && req.B == 0 {
		req.A = 0.8
		req.B = 1.1
	}
	if req.H <= 0 {
		req.H = 0.1
	}
	if req.LeftAlpha == 0 && req.LeftBeta == 0 && req.LeftGamma == 0 && req.RightAlpha == 0 && req.RightBeta == 0 && req.RightGamma == 0 {
		req.LeftAlpha = 0
		req.LeftBeta = 1
		req.LeftGamma = 1.5
		req.RightAlpha = 2
		req.RightBeta = 1
		req.RightGamma = 3
	}

	if req.B <= req.A {
		http.Error(w, "invalid interval: b must be greater than a", http.StatusBadRequest)
		return
	}

	n := int(math.Round((req.B - req.A) / req.H))
	if n < 2 {
		n = 2
	}
	h := (req.B - req.A) / float64(n)

	p := compileExpr(req.PExpr)
	q := compileExpr(req.QExpr)
	f := compileExpr(req.RHS)

	xs := make([]float64, n+1)
	for i := 0; i <= n; i++ {
		xs[i] = req.A + float64(i)*h
	}

	aCoef := make([]float64, n+1)
	bCoef := make([]float64, n+1)
	cCoef := make([]float64, n+1)
	dCoef := make([]float64, n+1)

	for i := 1; i <= n-1; i++ {
		x := xs[i]
		pv, err := p(x)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid p(x) at x=%g: %v", x, err), http.StatusBadRequest)
			return
		}
		qv, err := q(x)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid q(x) at x=%g: %v", x, err), http.StatusBadRequest)
			return
		}
		fv, err := f(x)
		if err != nil {
			http.Error(w, fmt.Sprintf("invalid rhs(x) at x=%g: %v", x, err), http.StatusBadRequest)
			return
		}

		aCoef[i] = 1.0/(h*h) - pv/(2*h)
		bCoef[i] = -2.0/(h*h) + qv
		cCoef[i] = 1.0/(h*h) + pv/(2*h)
		dCoef[i] = fv
	}

	// Левая граница: alpha1*y(a) + beta1*y'(a) = gamma1,
	// y'(a) ≈ (-3y0 + 4y1 - y2)/(2h), затем исключаем y2 через уравнение в i=1.
	L0 := req.LeftAlpha - 3.0*req.LeftBeta/(2*h)
	L1 := 2.0 * req.LeftBeta / h
	L2 := -req.LeftBeta / (2 * h)
	A1 := aCoef[1]
	B1 := bCoef[1]
	C1 := cCoef[1]
	D1 := dCoef[1]
	if math.Abs(C1) < 1e-14 {
		http.Error(w, "degenerate first interior equation", http.StatusBadRequest)
		return
	}
	aCoef[0] = 0
	bCoef[0] = L0 - L2*A1/C1
	cCoef[0] = L1 - L2*B1/C1
	dCoef[0] = req.LeftGamma - L2*D1/C1

	// Правая граница: alpha2*y(b) + beta2*y'(b) = gamma2,
	// y'(b) ≈ (3yN - 4yN-1 + yN-2)/(2h), затем исключаем yN-2 через i=N-1.
	R0 := req.RightBeta / (2 * h)
	R1 := -2.0 * req.RightBeta / h
	R2 := req.RightAlpha + 3.0*req.RightBeta/(2*h)
	AN1 := aCoef[n-1]
	BN1 := bCoef[n-1]
	CN1 := cCoef[n-1]
	DN1 := dCoef[n-1]
	if math.Abs(AN1) < 1e-14 {
		http.Error(w, "degenerate last interior equation", http.StatusBadRequest)
		return
	}
	aCoef[n] = R1 - R0*BN1/AN1
	bCoef[n] = R2 - R0*CN1/AN1
	cCoef[n] = 0
	dCoef[n] = req.RightGamma - R0*DN1/AN1

	alpha := make([]float64, n+1)
	beta := make([]float64, n+1)
	y := make([]float64, n+1)

	den := bCoef[0]
	if math.Abs(den) < 1e-14 {
		http.Error(w, "zero denominator at sweep start", http.StatusBadRequest)
		return
	}
	alpha[0] = -cCoef[0] / den
	beta[0] = dCoef[0] / den

	for i := 1; i <= n; i++ {
		den = bCoef[i] + aCoef[i]*alpha[i-1]
		if math.Abs(den) < 1e-14 {
			http.Error(w, fmt.Sprintf("zero denominator at sweep row %d", i), http.StatusBadRequest)
			return
		}
		alpha[i] = -cCoef[i] / den
		beta[i] = (dCoef[i] - aCoef[i]*beta[i-1]) / den
	}

	y[n] = beta[n]
	for i := n - 1; i >= 0; i-- {
		y[i] = alpha[i]*y[i+1] + beta[i]
	}

	points := make([]Point, 0, n+1)
	rows := make([]SweepRow, 0, n+1)
	for i := 0; i <= n; i++ {
		points = append(points, Point{Index: i, X: xs[i], Y: y[i]})
		rows = append(rows, SweepRow{
			Index: i,
			X:     xs[i],
			A:     aCoef[i],
			B:     bCoef[i],
			C:     cCoef[i],
			D:     dCoef[i],
			Alpha: alpha[i],
			Beta:  beta[i],
			Y:     y[i],
		})
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(SolveResponse{Points: points, Sweep: rows, StepUsed: h})
}
