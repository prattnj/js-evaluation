package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	var program Program

	err := json.Unmarshal([]byte(os.Args[1]), &program)
	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Println(program) // Parser
	fmt.Println(program.Body[0].Expression.Evaluate()) // Expressions
}
