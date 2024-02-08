package main

import (
	"regexp"
	"strconv"
	"strings"
)

func (e Expression) Evaluate() string {
	switch e.Type {
	case "Literal":
		return evaluateLiteral(e)
	case "BinaryExpression":
		return evaluateBinary(e)
	case "UnaryExpression":
		return evaluateUnary(e)
	case "LogicalExpression":
		return evaluateLogical(e)
	case "ConditionalExpression":
		return evaluateConditional(e)
	}
	return ""
}

func evaluateLiteral(e Expression) string {
	return newValue(e.Raw)
}

func evaluateBinary(e Expression) string {
	left := e.Left.Evaluate()
	right := e.Right.Evaluate()
	if hasError(left) {
		return left
	} else if hasError(right) {
		return right
	}
	if valueIsBoolean(left) || valueIsBoolean(right) {
		return newError("invalid binary type(s)")
	}

	leftVal := getNumberFromValue(left)
	rightVal := getNumberFromValue(right)
	if isArithmeticOperator(e.Operator) {
		// Arithmetic evaluation
		return doMath(leftVal, rightVal, e.Operator)
	} else {
		// Relational evaluation
		return doComparison(leftVal, rightVal, e.Operator)
	}
}

func evaluateUnary(e Expression) string {
	arg := e.Argument.Evaluate()
	if hasError(arg) {
		return arg
	}
	if !valueIsBoolean(arg) {
		return newError("invalid unary type")
	}
	b := getBoolFromValue(arg)
	return newValue(boolAsString(!b))
}

func evaluateLogical(e Expression) string {
	left := e.Left.Evaluate()
	right := e.Right.Evaluate()
	if hasError(left) {
		return left
	} else if hasError(right) {
		return right
	}
	if !valueIsBoolean(left) || !valueIsBoolean(right) {
		return newError("invalid logical type(s)")
	}

	leftVal := getBoolFromValue(left)
	rightVal := getBoolFromValue(right)
	var result bool
	switch e.Operator {
	case "||":
		result = leftVal || rightVal
	case "&&":
		result = leftVal && rightVal
	}
	return newValue(boolAsString(result))
}

func evaluateConditional(e Expression) string {
	test := e.Test.Evaluate()
	if hasError(test) {
		return test
	}
	if !valueIsBoolean(test) {
		return newError("invalid conditional type(s)")
	}

	if getBoolFromValue(test) {
		return e.Consequent.Evaluate()
	} else {
		return e.Alternate.Evaluate()
	}
}

func doMath(left int, right int, op string) string {
	var answer int
	switch op {
	case "+":
		answer = left + right
	case "-":
		answer = left - right
	case "*":
		answer = left * right
	case "/":
		if right == 0 {
			return newError("divide by zero")
		} else {
			answer = left / right
		}
	}
	return newValue(strconv.Itoa(answer))
}

func doComparison(left int, right int, op string) string {
	var result bool
	switch op {
	case "==":
		result = left == right
	case "<":
		result = left < right
	case ">":
		result = left > right
	case "<=":
		result = left <= right
	case ">=":
		result = left >= right
	}
	return newValue(boolAsString(result))
}

func isArithmeticOperator(op string) bool {
	return op == "+" || op == "-" || op == "*" || op == "/"
}

func stringIsWholeNumber(str string) bool {
	_, err := strconv.Atoi(str)
	return err == nil
}

func hasError(str string) bool {
	return strings.Contains(str, "error")
}

func valueIsBoolean(str string) bool {
	return strings.Contains(str, "boolean")
}

func getNumberFromValue(str string) int {
	re := regexp.MustCompile(`\((value \(number (\d+)\))\)`)

	match := re.FindStringSubmatch(str)

	numberStr := match[2]
	number, _ := strconv.Atoi(numberStr)

	return number
}

func getBoolFromValue(str string) bool {
	return strings.Contains(str, "true")
}

func newError(str string) string {
	return "(error (\"banana: " + str + "\"))"
}

func newValue(str string) string {
	var result string
	if str == "true" || str == "false" {
		result = "boolean " + str
	} else {
		if !stringIsWholeNumber(str) {
			return newError("not a whole number")
		}
		result = "number " + str
	}
	return "(value (" + result + "))"
}

func boolAsString(b bool) string {
	if b {
		return "true"
	} else {
		return "false"
	}
}
