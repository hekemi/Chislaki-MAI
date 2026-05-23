package main

import (
	"chislennie-metodi/internal/tasks/task02"
	"encoding/json"
	"fmt"
)

func main() {
	req := task02.DefaultRequest()
	resp := task02.Solve(req)
	b, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(b))
}
