package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

func main() {
	devMode := false

	var program Program
	var data string

	// Read piped input from os.Stdin
	scanner := bufio.NewScanner(os.Stdin)

	// Read input line by line
	for scanner.Scan() {
		input := scanner.Text()
		data += input
	}

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
		for {
			err := json.Unmarshal([]byte(data), &program)
			if err == nil {
				break
			} else {
				// necessary for strange file encodings
				data = data[1:]
			}
		}
	}

	// fmt.Println(program) // Parser
	fmt.Println(program.Evaluate()) // Expressions
}
