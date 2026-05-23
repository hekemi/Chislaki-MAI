package task07

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
)

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type RegressionLine struct {
	A        float64 `json:"a"`
	B        float64 `json:"b"`
	Equation string  `json:"equation"`
	RMSE     float64 `json:"rmse"`
	SSE      float64 `json:"sse"`
}

type QuadraticCurve struct {
	A        float64 `json:"a"`
	B        float64 `json:"b"`
	C        float64 `json:"c"`
	Equation string  `json:"equation"`
	RMSE     float64 `json:"rmse"`
	SSE      float64 `json:"sse"`
}

type ObservationRow struct {
	Index        int     `json:"index"`
	X            float64 `json:"x"`
	Y            float64 `json:"y"`
	Linear       float64 `json:"linear"`
	Quadratic    float64 `json:"quadratic"`
	LinearError  float64 `json:"linear_error"`
	QuadraticErr float64 `json:"quadratic_error"`
}

type SolveRequest struct {
	Points []Point `json:"points"`
}

type SolveResponse struct {
	Points    []Point          `json:"points"`
	Linear    RegressionLine   `json:"linear"`
	Quadratic QuadraticCurve   `json:"quadratic"`
	Rows      []ObservationRow `json:"rows"`
}

var presetPoints = []Point{
	{X: 0.00, Y: 1.0},
	{X: 0.12, Y: 1.2},
	{X: 0.19, Y: 1.6},
	{X: 0.35, Y: 2.6},
	{X: 0.40, Y: 1.8},
	{X: 0.45, Y: 2.7},
	{X: 0.62, Y: 3.5},
	{X: 0.71, Y: 4.4},
	{X: 0.84, Y: 4.5},
	{X: 0.91, Y: 5.2},
	{X: 1.00, Y: 6.3},
}

func clamp(v float64) float64 {
	if math.Abs(v) < 1e-12 {
		return 0
	}
	return v
}

func linearLeastSquares(points []Point) RegressionLine {
	n := float64(len(points))
	var sumX, sumY, sumXX, sumXY float64
	for _, p := range points {
		sumX += p.X
		sumY += p.Y
		sumXX += p.X * p.X
		sumXY += p.X * p.Y
	}
	den := n*sumXX - sumX*sumX
	b := 0.0
	if den != 0 {
		b = (n*sumXY - sumX*sumY) / den
	}
	a := (sumY - b*sumX) / n
	sse := 0.0
	for _, p := range points {
		e := p.Y - (a + b*p.X)
		sse += e * e
	}
	rmse := math.Sqrt(sse / n)
	return RegressionLine{
		A:        clamp(a),
		B:        clamp(b),
		Equation: fmt.Sprintf("y = %.6f + %.6fx", a, b),
		SSE:      sse,
		RMSE:     rmse,
	}
}

func quadraticLeastSquares(points []Point) QuadraticCurve {
	var s0, s1, s2, s3, s4 float64
	var t0, t1, t2 float64
	for _, p := range points {
		x := p.X
		y := p.Y
		x2 := x * x
		s0 += 1
		s1 += x
		s2 += x2
		s3 += x2 * x
		s4 += x2 * x2
		t0 += y
		t1 += x * y
		t2 += x2 * y
	}

	m := [3][4]float64{
		{s0, s1, s2, t0},
		{s1, s2, s3, t1},
		{s2, s3, s4, t2},
	}

	for i := 0; i < 3; i++ {
		pivot := i
		for j := i + 1; j < 3; j++ {
			if math.Abs(m[j][i]) > math.Abs(m[pivot][i]) {
				pivot = j
			}
		}
		m[i], m[pivot] = m[pivot], m[i]
		div := m[i][i]
		if div == 0 {
			return QuadraticCurve{}
		}
		for k := i; k < 4; k++ {
			m[i][k] /= div
		}
		for j := 0; j < 3; j++ {
			if j == i {
				continue
			}
			ratio := m[j][i]
			for k := i; k < 4; k++ {
				m[j][k] -= ratio * m[i][k]
			}
		}
	}

	a, b, c := m[0][3], m[1][3], m[2][3]
	sse := 0.0
	for _, p := range points {
		e := p.Y - (a + b*p.X + c*p.X*p.X)
		sse += e * e
	}
	rmse := math.Sqrt(sse / float64(len(points)))
	return QuadraticCurve{
		A:        clamp(a),
		B:        clamp(b),
		C:        clamp(c),
		Equation: fmt.Sprintf("y = %.6f + %.6fx + %.6fx²", a, b, c),
		SSE:      sse,
		RMSE:     rmse,
	}
}

func evaluateLinear(line RegressionLine, x float64) float64 {
	return line.A + line.B*x
}

func evaluateQuadratic(curve QuadraticCurve, x float64) float64 {
	return curve.A + curve.B*x + curve.C*x*x
}

func solve(points []Point) SolveResponse {
	linear := linearLeastSquares(points)
	quadratic := quadraticLeastSquares(points)
	rows := make([]ObservationRow, 0, len(points))
	for i, p := range points {
		linearValue := evaluateLinear(linear, p.X)
		quadraticValue := evaluateQuadratic(quadratic, p.X)
		rows = append(rows, ObservationRow{
			Index:        i + 1,
			X:            p.X,
			Y:            p.Y,
			Linear:       linearValue,
			Quadratic:    quadraticValue,
			LinearError:  p.Y - linearValue,
			QuadraticErr: p.Y - quadraticValue,
		})
	}
	return SolveResponse{
		Points:    points,
		Linear:    linear,
		Quadratic: quadratic,
		Rows:      rows,
	}
}

func parsePointsJSON(body []byte, req *SolveRequest) error {
	return json.Unmarshal(body, req)
}

func defaultPoints() []Point {
	points := make([]Point, len(presetPoints))
	copy(points, presetPoints)
	return points
}

func SolveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "unable to read request body", http.StatusBadRequest)
		return
	}

	points := defaultPoints()
	if len(body) > 0 {
		var req SolveRequest
		if err := parsePointsJSON(body, &req); err == nil && len(req.Points) > 0 {
			points = req.Points
		}
	}

	sort.Slice(points, func(i, j int) bool { return points[i].X < points[j].X })
	resp := solve(points)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
