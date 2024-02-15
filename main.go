package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	devMode := false

	var program Program

	if devMode {
		content, err := os.ReadFile("test8.json")
		if err != nil {
			return
		}
		err = json.Unmarshal(content, &program)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		err := json.Unmarshal([]byte(os.Args[1]), &program)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// fmt.Println(program) // Parser
	fmt.Println(program.Evaluate()) // Expressions
}
