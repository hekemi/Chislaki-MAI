package main

import (
	"chislennie-metodi/internal/tasks/task04"
	"fmt"
)

func main() {
	eq := &task04.ParsedEquation{Expression: "3x - e^x = 0"}
	v, err := eq.EvaluateFunc(0.5)
	fmt.Println("value:", v, "err:", err)
}
