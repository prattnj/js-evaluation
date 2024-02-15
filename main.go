package main

import (
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	var program Program

	// For development
	/*content, err := ioutil.ReadFile("test.json")
	if err != nil {
		return
	}*/

	err := json.Unmarshal([]byte(os.Args[1]), &program)
	if err != nil {
		fmt.Println(err)
		return
	}

	// fmt.Println(program) // Parser
	fmt.Println(program.Evaluate()) // Expressions
}
