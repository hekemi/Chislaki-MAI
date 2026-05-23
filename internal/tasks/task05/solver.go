package task05

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// Request и Response структуры
type SolveRequest struct {
	CSV   string  `json:"csv"`    // строки вида x,y на каждой строке
	XStar float64 `json:"x_star"` // точка, в которой искать погрешность
}

type Iteration struct {
	Index int     `json:"index"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
}

type SolveResponse struct {
	Nodes       []Iteration `json:"nodes"`
	Lagrange    float64     `json:"lagrange"`
	Newton      float64     `json:"newton"`
	ErrorLagr   float64     `json:"err_lagrange"`
	ErrorNewton float64     `json:"err_newton"`
}

// parseCSV парсит CSV строку в массив узлов
func parseCSV(s string) ([][2]float64, error) {
	r := csv.NewReader(strings.NewReader(s))
	r.TrimLeadingSpace = true
	var out [][2]float64
	for {
		rec, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(rec) < 2 {
			continue
		}
		x, err := strconv.ParseFloat(strings.TrimSpace(rec[0]), 64)
		if err != nil {
			return nil, err
		}
		y, err := strconv.ParseFloat(strings.TrimSpace(rec[1]), 64)
		if err != nil {
			return nil, err
		}
		out = append(out, [2]float64{x, y})
	}
	return out, nil
}

// lagrangeEvaluate вычисляет значение многочлена Лагранжа в x
func lagrangeEvaluate(nodes [][2]float64, x float64) float64 {
	n := len(nodes)
	res := 0.0
	for i := 0; i < n; i++ {
		xi := nodes[i][0]
		yi := nodes[i][1]
		li := 1.0
		for j := 0; j < n; j++ {
			if i == j {
				continue
			}
			xj := nodes[j][0]
			li *= (x - xj) / (xi - xj)
		}
		res += yi * li
	}
	return res
}

// newtonEvaluate вычисляет значение по формуле Ньютона (деленое разностное)
func newtonEvaluate(nodes [][2]float64, x float64) float64 {
	coeffs := newtonCoefficients(nodes)
	return newtonEvaluateWithCoefficients(nodes, coeffs, x)
}

func newtonCoefficients(nodes [][2]float64) []float64 {
	n := len(nodes)
	coeffs := make([]float64, n)
	for i := 0; i < n; i++ {
		coeffs[i] = nodes[i][1]
	}
	for k := 1; k < n; k++ {
		for i := n - 1; i >= k; i-- {
			coeffs[i] = (coeffs[i] - coeffs[i-1]) / (nodes[i][0] - nodes[i-k][0])
		}
	}
	return coeffs
}

func newtonEvaluateWithCoefficients(nodes [][2]float64, coeffs []float64, x float64) float64 {
	n := len(nodes)
	if n == 0 {
		return 0
	}
	res := coeffs[n-1]
	for i := n - 2; i >= 0; i-- {
		res = res*(x-nodes[i][0]) + coeffs[i]
	}
	return res
}

func estimateFromHighestNewtonTerm(nodes [][2]float64, coeffs []float64, x float64) float64 {
	n := len(nodes)
	if n == 0 {
		return 0
	}
	term := coeffs[n-1]
	for i := 0; i < n-1; i++ {
		term *= x - nodes[i][0]
	}
	return math.Abs(term)
}

func SolveHandler(w http.ResponseWriter, r *http.Request) {
	// простой парсинг body (form-urlencoded или raw)
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// read body
	buf := new(strings.Builder)
	_, _ = io.Copy(buf, r.Body)
	body := buf.String()
	// ожидаем формат: csv=...&x_star=...
	vals := map[string]string{}
	for _, part := range strings.Split(body, "&") {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			v, err := url.QueryUnescape(kv[1])
			if err != nil {
				v = kv[1]
			}
			vals[kv[0]] = v
		}
	}
	csvStr, ok := vals["csv"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("missing csv"))
		return
	}
	xStar := 0.0
	if s, ok := vals["x_star"]; ok {
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			xStar = v
		}
	}

	nodes, err := parseCSV(csvStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	resp := SolveResponse{}
	coeffs := newtonCoefficients(nodes)
	for i, n := range nodes {
		resp.Nodes = append(resp.Nodes, Iteration{Index: i + 1, X: n[0], Y: n[1]})
	}

	resp.Lagrange = lagrangeEvaluate(nodes, xStar)
	resp.Newton = newtonEvaluateWithCoefficients(nodes, coeffs, xStar)

	// Оценка погрешности: модуль последнего члена формы Ньютона.
	resp.ErrorLagr = estimateFromHighestNewtonTerm(nodes, coeffs, xStar)
	resp.ErrorNewton = resp.ErrorLagr

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
