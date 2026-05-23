package task06

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type SolveRequest struct {
	CSV                 string    `json:"csv"`
	XStar               float64   `json:"x_star"`
	BoundaryLeftSecond  *float64  `json:"boundary_left_second,omitempty"`
	BoundaryRightSecond *float64  `json:"boundary_right_second,omitempty"`
	EvaluatePoints      []float64 `json:"evaluate_points,omitempty"`
}

type Node struct {
	Index int     `json:"index"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
}

type Segment struct {
	Index int     `json:"index"`
	X0    float64 `json:"x0"`
	X1    float64 `json:"x1"`
	A     float64 `json:"a"`
	B     float64 `json:"b"`
	C     float64 `json:"c"`
	D     float64 `json:"d"`
	Value float64 `json:"value"`
}

type SolveResponse struct {
	Nodes         []Node    `json:"nodes"`
	Segments      []Segment `json:"segments"`
	SplineValue   float64   `json:"spline_value"`
	BoundaryLeft  float64   `json:"boundary_left_second"`
	BoundaryRight float64   `json:"boundary_right_second"`
	Evaluations   []Segment `json:"evaluations"`
}

func parseCSV(s string) ([][2]float64, error) {
	r := csv.NewReader(strings.NewReader(s))
	r.TrimLeadingSpace = true
	var nodes [][2]float64
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

func buildSplineWithBoundaries(nodes [][2]float64, leftC *float64, rightC *float64) ([]Segment, float64, float64, error) {
	n := len(nodes)
	h := make([]float64, n-1)
	for i := 0; i < n-1; i++ {
		h[i] = nodes[i+1][0] - nodes[i][0]
		if h[i] <= 0 {
			return nil, 0, 0, fmt.Errorf("узлы должны быть упорядочены по возрастанию")
		}
	}

	// Build tridiagonal system for c[0..n-1]:
	// For interior i=1..n-2: h[i-1]*c[i-1] + 2*(h[i-1]+h[i])*c[i] + h[i]*c[i+1] = rhs[i]
	// where rhs[i] = 3*((a[i+1]-a[i])/h[i] - (a[i]-a[i-1])/h[i-1])

	// Prepare RHS for interior nodes
	rhs := make([]float64, n)
	for i := 1; i < n-1; i++ {
		rhs[i] = 3 * ((nodes[i+1][1]-nodes[i][1])/h[i] - (nodes[i][1]-nodes[i-1][1])/h[i-1])
	}

	// If boundaries not provided, assume natural (0)
	leftVal := 0.0
	rightVal := 0.0
	if leftC != nil {
		leftVal = *leftC
	}
	if rightC != nil {
		rightVal = *rightC
	}

	// Solve for interior c[1..n-2] using Thomas algorithm on (n-2)x(n-2) system
	m := n - 2
	c := make([]float64, n)
	if m > 0 {
		aTri := make([]float64, m) // lower diag
		bTri := make([]float64, m) // main diag
		cTri := make([]float64, m) // upper diag
		dVec := make([]float64, m) // rhs

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

		// adjust rhs for known boundary c0 and cn
		if m > 0 {
			dVec[0] -= h[0] * leftVal
			dVec[m-1] -= h[n-2] * rightVal
		}

		// Thomas algorithm
		for i := 1; i < m; i++ {
			w := aTri[i] / bTri[i-1]
			bTri[i] = bTri[i] - w*cTri[i-1]
			dVec[i] = dVec[i] - w*dVec[i-1]
		}

		// back substitution
		x := make([]float64, m)
		if bTri[m-1] == 0 {
			return nil, 0, 0, fmt.Errorf("не удалось построить сплайн: вырожденная система")
		}
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

	// compute a,b,d
	a := make([]float64, n-1)
	b := make([]float64, n-1)
	d := make([]float64, n-1)
	for j := 0; j < n-1; j++ {
		a[j] = nodes[j][1]
		b[j] = (nodes[j+1][1]-nodes[j][1])/h[j] - h[j]*(c[j+1]+2*c[j])/3
		d[j] = (c[j+1] - c[j]) / (3 * h[j])
	}

	segments := make([]Segment, 0, n-1)
	for i := 0; i < n-1; i++ {
		segments = append(segments, Segment{
			Index: i + 1,
			X0:    nodes[i][0],
			X1:    nodes[i+1][0],
			A:     a[i],
			B:     b[i],
			C:     c[i],
			D:     d[i],
		})
	}

	return segments, c[0], c[n-1], nil
}

func evaluateSpline(segments []Segment, x float64) float64 {
	if len(segments) == 0 {
		return 0
	}
	segment := segments[len(segments)-1]
	for _, candidate := range segments {
		if x >= candidate.X0 && x <= candidate.X1 {
			segment = candidate
			break
		}
	}
	if x < segments[0].X0 {
		segment = segments[0]
	} else if x > segments[len(segments)-1].X1 {
		segment = segments[len(segments)-1]
	}
	dx := x - segment.X0
	return segment.A + segment.B*dx + segment.C*dx*dx + segment.D*dx*dx*dx
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

	var req SolveRequest
	if len(body) > 0 && body[0] == '{' {
		if err := json.Unmarshal(body, &req); err != nil {
			http.Error(w, "invalid json body", http.StatusBadRequest)
			return
		}
	} else {
		vals := map[string]string{}
		for _, part := range strings.Split(string(body), "&") {
			kv := strings.SplitN(part, "=", 2)
			if len(kv) == 2 {
				v, err := url.QueryUnescape(kv[1])
				if err != nil {
					v = kv[1]
				}
				vals[kv[0]] = v
			}
		}

		req.CSV = vals["csv"]
		if s, ok := vals["x_star"]; ok {
			if v, err := strconv.ParseFloat(s, 64); err == nil {
				req.XStar = v
			}
		}
		if s, ok := vals["boundary_left_second"]; ok {
			if v, err := strconv.ParseFloat(s, 64); err == nil {
				req.BoundaryLeftSecond = &v
			}
		}
		if s, ok := vals["boundary_right_second"]; ok {
			if v, err := strconv.ParseFloat(s, 64); err == nil {
				req.BoundaryRightSecond = &v
			}
		}
		if s, ok := vals["evaluate_points"]; ok {
			// comma separated
			parts := strings.Split(s, ",")
			for _, p := range parts {
				if v, err := strconv.ParseFloat(strings.TrimSpace(p), 64); err == nil {
					req.EvaluatePoints = append(req.EvaluatePoints, v)
				}
			}
		}
	}

	csvStr := strings.TrimSpace(req.CSV)
	if csvStr == "" {
		http.Error(w, "missing csv", http.StatusBadRequest)
		return
	}
	xStar := req.XStar

	nodesRaw, err := parseCSV(csvStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	segments, leftSecond, rightSecond, err := buildSplineWithBoundaries(nodesRaw, req.BoundaryLeftSecond, req.BoundaryRightSecond)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := SolveResponse{
		Nodes:         make([]Node, 0, len(nodesRaw)),
		Segments:      make([]Segment, 0, len(segments)),
		SplineValue:   evaluateSpline(segments, xStar),
		BoundaryLeft:  leftSecond,
		BoundaryRight: rightSecond,
		Evaluations:   make([]Segment, 0),
	}
	for i, node := range nodesRaw {
		resp.Nodes = append(resp.Nodes, Node{Index: i + 1, X: node[0], Y: node[1]})
	}
	resp.Segments = append(resp.Segments, segments...)

	// prepare evaluation points: if EvaluatePoints provided use them, else use single xStar
	evalPoints := req.EvaluatePoints
	if len(evalPoints) == 0 {
		evalPoints = []float64{xStar}
	}

	for _, x := range evalPoints {
		val := evaluateSpline(segments, x)
		// reuse Segment struct to report point: Index=0, X0=x, Value=val
		resp.Evaluations = append(resp.Evaluations, Segment{Index: 0, X0: x, Value: val})
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
}
