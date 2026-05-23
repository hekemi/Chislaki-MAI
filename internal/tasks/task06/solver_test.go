package task06

import (
	"math"
	"testing"
)

func TestBuildSplineNatural_InterpolatesNodesAndNaturalBoundaries(t *testing.T) {
	nodes := [][2]float64{
		{1.375, 5.04192},
		{1.379, 5.17744},
		{1.383, 5.32016},
		{1.387, 5.47069},
		{1.391, 5.62968},
		{1.395, 5.79788},
	}

	segments, left, right, err := buildSplineWithBoundaries(nodes, nil, nil)
	if err != nil {
		t.Fatalf("buildSplineWithBoundaries returned error: %v", err)
	}

	if math.Abs(left) > 1e-12 || math.Abs(right) > 1e-12 {
		t.Fatalf("expected natural boundaries c0=cn=0, got left=%g right=%g", left, right)
	}

	for i := range nodes {
		x := nodes[i][0]
		want := nodes[i][1]
		got := evaluateSpline(segments, x)
		if math.Abs(got-want) > 1e-8 {
			t.Fatalf("node interpolation mismatch at x=%g: got=%g want=%g", x, got, want)
		}
	}
}

func TestBuildSplineWithBoundarySeconds_RespectsGivenValues(t *testing.T) {
	// y=x^3, y''(x)=6x, so at x=0 -> 0 and at x=2 -> 12
	nodes := [][2]float64{{0, 0}, {1, 1}, {2, 8}}
	left := 0.0
	right := 12.0

	segments, gotLeft, gotRight, err := buildSplineWithBoundaries(nodes, &left, &right)
	if err != nil {
		t.Fatalf("buildSplineWithBoundaries returned error: %v", err)
	}

	if math.Abs(gotLeft-left) > 1e-12 || math.Abs(gotRight-right) > 1e-12 {
		t.Fatalf("unexpected boundary values: got left=%g right=%g", gotLeft, gotRight)
	}

	for i := range nodes {
		x := nodes[i][0]
		want := nodes[i][1]
		got := evaluateSpline(segments, x)
		if math.Abs(got-want) > 1e-9 {
			t.Fatalf("node interpolation mismatch at x=%g: got=%g want=%g", x, got, want)
		}
	}

}
