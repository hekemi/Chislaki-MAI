package task10

import (
	"chislennie-metodi/internal/tasks/task04"
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"
)

type SolveRequest struct {
	Equation string  `json:"equation"`
	A        float64 `json:"a"`
	B        float64 `json:"b"`
	Y0       float64 `json:"y0"`
	H        float64 `json:"h"`
}

type RKStep struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type AdamsRow struct {
	X     float64 `json:"x"`
	YPred float64 `json:"y_pred"`
	YCorr float64 `json:"y_corr"`
	Delta float64 `json:"delta"`
}

type SolveResponse struct {
	RK4   []RKStep   `json:"rk4"`
	Adams []AdamsRow `json:"adams"`
}

func makeFunc(eqStr string) (func(float64, float64) (float64, error), error) {
	orig := eqStr
	eq := task04.ParsedEquation{Expression: eqStr}
	return func(x, y float64) (float64, error) {
		// Replace standalone 'y' occurrences with numeric literal, but avoid identifiers like 'exp'
		var b strings.Builder
		runes := []rune(orig)
		for i, r := range runes {
			if r == 'y' {
				prevIsLetter := false
				nextIsLetter := false
				if i-1 >= 0 {
					p := runes[i-1]
					if (p >= 'A' && p <= 'Z') || (p >= 'a' && p <= 'z') || p == '_' {
						prevIsLetter = true
					}
				}
				if i+1 < len(runes) {
					n := runes[i+1]
					if (n >= 'A' && n <= 'Z') || (n >= 'a' && n <= 'z') || n == '_' {
						nextIsLetter = true
					}
				}
				if !prevIsLetter && !nextIsLetter {
					b.WriteString("(" + strconv.FormatFloat(y, 'g', -1, 64) + ")")
					continue
				}
			}
			b.WriteRune(r)
		}
		eq.Expression = b.String()
		return eq.EvaluateFunc(x)
	}, nil
}

func rk4Solve(f func(float64, float64) (float64, error), a, b, y0 float64, h float64) ([]RKStep, error) {
	if h <= 0 {
		h = 0.1
	}
	n := int(math.Ceil((b - a) / h))
	if n <= 0 {
		n = 1
	}
	h = (b - a) / float64(n)

	steps := make([]RKStep, 0, n+1)
	x := a
	y := y0
	steps = append(steps, RKStep{X: x, Y: y})
	for i := 0; i < n; i++ {
		k1, err := f(x, y)
		if err != nil {
			return nil, err
		}
		k2, err := f(x+h/2, y+h*k1/2)
		if err != nil {
			return nil, err
		}
		k3, err := f(x+h/2, y+h*k2/2)
		if err != nil {
			return nil, err
		}
		k4, err := f(x+h, y+h*k3)
		if err != nil {
			return nil, err
		}
		y = y + (h/6.0)*(k1+2*k2+2*k3+k4)
		x = x + h
		steps = append(steps, RKStep{X: x, Y: y})
	}
	return steps, nil
}

func adamsSolve(f func(float64, float64) (float64, error), a, b, y0 float64, h float64) ([]AdamsRow, error) {
	// compute n and use RK4 to get first 3 steps
	if h <= 0 {
		h = 0.1
	}
	n := int(math.Ceil((b - a) / h))
	if n <= 0 {
		n = 1
	}
	h = (b - a) / float64(n)

	// initial RK4 steps
	rk, err := rk4Solve(f, a, a+3*h, y0, h)
	if err != nil {
		return nil, err
	}
	// ensure we have 4 points (0..3)
	if len(rk) < 4 {
		return nil, nil
	}

	// store f evaluations
	xs := make([]float64, 0, n+1)
	ys := make([]float64, 0, n+1)
	fs := make([]float64, 0, n+1)
	// push initial rk points
	for i := 0; i < len(rk); i++ {
		xs = append(xs, rk[i].X)
		ys = append(ys, rk[i].Y)
		fv, ferr := f(rk[i].X, rk[i].Y)
		if ferr != nil {
			return nil, ferr
		}
		fs = append(fs, fv)
	}

	// continue steps
	rows := make([]AdamsRow, 0, n)
	// add initial points as rows with pred==corr==y
	for i := 0; i < len(xs); i++ {
		rows = append(rows, AdamsRow{X: xs[i], YPred: ys[i], YCorr: ys[i], Delta: 0})
	}

	for i := len(xs); i <= n; i++ {
		xnext := a + float64(i)*h
		// predictor AB4 using last 4 f's (fs[len-1]..fs[len-4])
		m := len(fs)
		if m < 4 {
			// shouldn't happen
			break
		}
		f_n := fs[m-1]
		f_n1 := fs[m-2]
		f_n2 := fs[m-3]
		f_n3 := fs[m-4]
		y_n := ys[m-1]
		ypred := y_n + (h/24.0)*(55*f_n-59*f_n1+37*f_n2-9*f_n3)

		// evaluate f at predicted
		f_pred, ferr := f(xnext, ypred)
		if ferr != nil {
			return nil, ferr
		}
		// corrector (Adams-Moulton 4-step)
		ycorr := y_n + (h/24.0)*(9*f_pred+19*f_n-5*f_n1+f_n2)
		delta := math.Abs(ycorr - ypred)

		rows = append(rows, AdamsRow{X: xnext, YPred: ypred, YCorr: ycorr, Delta: delta})

		// append to arrays for next iter
		xs = append(xs, xnext)
		ys = append(ys, ycorr)
		fv2, ferr2 := f(xnext, ycorr)
		if ferr2 != nil {
			return nil, ferr2
		}
		fs = append(fs, fv2)
	}

	return rows, nil
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
	if req.H <= 0 {
		req.H = 0.1
	}
	if req.B <= req.A {
		http.Error(w, "invalid interval", http.StatusBadRequest)
		return
	}

	f, _ := makeFunc(req.Equation)

	rk4, err := rk4Solve(f, req.A, req.B, req.Y0, req.H)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	adams, err := adamsSolve(f, req.A, req.B, req.Y0, req.H)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(SolveResponse{RK4: rk4, Adams: adams})
}
