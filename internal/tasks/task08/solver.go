package task08

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type SolveRequest struct {
	CSV            string    `json:"csv"`
	EvaluatePoints []float64 `json:"evaluate_points"`
	BoundaryLeft   *float64  `json:"boundary_left_second,omitempty"`
	BoundaryRight  *float64  `json:"boundary_right_second,omitempty"`
}

type Node struct {
	Index int     `json:"index"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
}

type EvalRow struct {
	X  float64 `json:"x"`
	S  float64 `json:"s"`
	S1 float64 `json:"s1"`
	S2 float64 `json:"s2"`
}

type SolveResponse struct {
	Nodes       []Node    `json:"nodes"`
	Evaluations []EvalRow `json:"evaluations"`
}

func parseCSV(s string) ([][2]float64, error) {
	var nodes [][2]float64
	// Support both comma-separated and whitespace-separated rows.
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// normalize commas to spaces, then split by whitespace
		ln := strings.ReplaceAll(line, ",", " ")
		fields := strings.Fields(ln)
		if len(fields) < 2 {
			continue
		}
		x, err := strconv.ParseFloat(strings.TrimSpace(fields[0]), 64)
		if err != nil {
			return nil, err
		}
		y, err := strconv.ParseFloat(strings.TrimSpace(fields[1]), 64)
		if err != nil {
			return nil, err
		}
		nodes = append(nodes, [2]float64{x, y})
	}
	if len(nodes) < 2 {
		return nil, fmt.Errorf("нужно минимум 2 узла")
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i][0] < nodes[j][0] })
	for i := 1; i < len(nodes); i++ {
		if nodes[i][0] == nodes[i-1][0] {
			return nil, fmt.Errorf("x-координаты должны быть различны")
		}
	}
	return nodes, nil
}

// buildSpline (similar to task06) builds spline segments and returns c array too
func buildSpline(nodes [][2]float64, leftC *float64, rightC *float64) ([][7]float64, error) {
	// each segment: x0,x1,a,b,c,d, (we'll pack as array)
	n := len(nodes)
	h := make([]float64, n-1)
	for i := 0; i < n-1; i++ {
		h[i] = nodes[i+1][0] - nodes[i][0]
		if h[i] <= 0 {
			return nil, fmt.Errorf("узлы должны быть упорядочены по возрастанию")
		}
	}

	rhs := make([]float64, n)
	for i := 1; i < n-1; i++ {
		rhs[i] = 3 * ((nodes[i+1][1]-nodes[i][1])/h[i] - (nodes[i][1]-nodes[i-1][1])/h[i-1])
	}

	leftVal := 0.0
	rightVal := 0.0
	if leftC != nil {
		leftVal = *leftC
	}
	if rightC != nil {
		rightVal = *rightC
	}

	m := n - 2
	c := make([]float64, n)
	if m > 0 {
		aTri := make([]float64, m)
		bTri := make([]float64, m)
		cTri := make([]float64, m)
		dVec := make([]float64, m)
		for i := 0; i < m; i++ {
			idx := i + 1
			bTri[i] = 2 * (h[idx-1] + h[idx])
			dVec[i] = rhs[idx]
			if i > 0 {
				aTri[i] = h[idx-1]
			}
			if i < m-1 {
				cTri[i] = h[idx]
			}
		}
		dVec[0] -= h[0] * leftVal
		dVec[m-1] -= h[n-2] * rightVal
		for i := 1; i < m; i++ {
			w := aTri[i] / bTri[i-1]
			bTri[i] = bTri[i] - w*cTri[i-1]
			dVec[i] = dVec[i] - w*dVec[i-1]
		}
		if bTri[m-1] == 0 {
			return nil, fmt.Errorf("не удалось построить сплайн: вырожденная система")
		}
		x := make([]float64, m)
		x[m-1] = dVec[m-1] / bTri[m-1]
		for i := m - 2; i >= 0; i-- {
			x[i] = (dVec[i] - cTri[i]*x[i+1]) / bTri[i]
		}
		for i := 0; i < m; i++ {
			c[i+1] = x[i]
		}
	}
	c[0] = leftVal
	c[n-1] = rightVal

	segments := make([][7]float64, 0, n-1)
	for j := 0; j < n-1; j++ {
		a := nodes[j][1]
		b := (nodes[j+1][1]-nodes[j][1])/h[j] - h[j]*(c[j+1]+2*c[j])/3
		d := (c[j+1] - c[j]) / (3 * h[j])
		segments = append(segments, [7]float64{nodes[j][0], nodes[j+1][0], a, b, c[j], d, 0})
	}
	return segments, nil
}

func findSegment(segments [][7]float64, x float64) int {
	if len(segments) == 0 {
		return -1
	}
	for i, s := range segments {
		x0 := s[0]
		x1 := s[1]
		if x >= x0 && x <= x1 {
			return i
		}
	}
	if x < segments[0][0] {
		return 0
	}
	return len(segments) - 1
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

	nodesRaw, err := parseCSV(strings.TrimSpace(req.CSV))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	segments, err := buildSpline(nodesRaw, req.BoundaryLeft, req.BoundaryRight)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := SolveResponse{
		Nodes:       make([]Node, 0, len(nodesRaw)),
		Evaluations: make([]EvalRow, 0),
	}
	for i, n := range nodesRaw {
		resp.Nodes = append(resp.Nodes, Node{Index: i + 1, X: n[0], Y: n[1]})
	}

	pts := req.EvaluatePoints
	if len(pts) == 0 {
		// default: try two points near middle
		mid := (nodesRaw[0][0] + nodesRaw[len(nodesRaw)-1][0]) / 2
		pts = []float64{mid}
	}

	for _, x := range pts {
		i := findSegment(segments, x)
		if i < 0 {
			continue
		}
		s := segments[i]
		x0 := s[0]
		// a:=s[2]
		b := s[3]
		c := s[4]
		d := s[5]
		dx := x - x0
		sval := s[2] + b*dx + c*dx*dx + d*dx*dx*dx
		s1 := b + 2*c*dx + 3*d*dx*dx
		s2 := 2*c + 6*d*dx
		resp.Evaluations = append(resp.Evaluations, EvalRow{X: x, S: sval, S1: s1, S2: s2})
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
