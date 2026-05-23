package task04

import "testing"

func TestSolvePresetEquationHasIterations(t *testing.T) {
	eq := &ParsedEquation{Expression: "3*x - exp(x)"}

	bis := SolveBisection(eq, 0, 1, 0.0001, 100)
	if bis.Error != "" {
		t.Fatalf("bisection error: %s", bis.Error)
	}
	if len(bis.Iterations) == 0 {
		t.Fatal("bisection returned no iterations")
	}

	newton := SolveNewton(eq, 1, 0.0001, 100)
	if newton.Error != "" {
		t.Fatalf("newton error: %s", newton.Error)
	}
	if len(newton.Iterations) == 0 {
		t.Fatal("newton returned no iterations")
	}
}
