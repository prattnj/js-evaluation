package main

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

func (p Program) Evaluate() string {
	topScope := &Scope{
		Variables: map[string]Value{},
		Parent:    nil,
	}
	result := handleBody(p.Body, topScope)
	if result.StringValue != "" {
		return result.StringValue
	}
	return evaluateValue(result)
}

func (e Expression) Evaluate() Value {
	switch e.Type {
	case "Identifier":
		val, str := evaluateIdentifier(e)
		if str != "" {
			return newStringValue(str)
		} else {
			return val
		}
	case "Literal":
		return newStringValue(evaluateLiteral(e))
	case "BinaryExpression":
		return newStringValue(evaluateBinary(e))
	case "UnaryExpression":
		return newStringValue(evaluateUnary(e))
	case "LogicalExpression":
		return newStringValue(evaluateLogical(e))
	case "ConditionalExpression":
		return newStringValue(evaluateConditional(e))
	case "FunctionExpression":
		return newFunctionValue(evaluateFunction(e))
	case "CallExpression":
		return evaluateCall(e)
	}
	return Value{}
}

func evaluateIdentifier(e Expression) (Value, string) {
	expr, ok := getIdentifierValue(e.Name, e.Scope)
	if !ok {
		return Value{}, newError("unbound identifier")
	}
	return expr, ""
}

func evaluateLiteral(e Expression) string {
	return newFinalValue(e.Raw)
}

func evaluateBinary(e Expression) string {
	e.Left.Scope = e.Scope
	left := e.Left.Evaluate().StringValue
	if hasError(left) {
		return left
	}
	e.Right.Scope = e.Scope
	right := e.Right.Evaluate().StringValue
	if hasError(right) {
		return right
	}
	if valueIsBoolean(left) || valueIsBoolean(right) || left == "" || right == "" {
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
	e.Argument.Scope = e.Scope
	arg := e.Argument.Evaluate().StringValue
	if hasError(arg) {
		return arg
	}
	if !valueIsBoolean(arg) {
		return newError("invalid unary type")
	}
	b := getBoolFromValue(arg)
	return newFinalValue(boolAsString(!b))
}

func evaluateLogical(e Expression) string {
	e.Left.Scope = e.Scope
	left := e.Left.Evaluate().StringValue
	if hasError(left) {
		return left
	}
	e.Right.Scope = e.Scope
	right := e.Right.Evaluate().StringValue
	if hasError(right) {
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
	return newFinalValue(boolAsString(result))
}

func evaluateConditional(e Expression) string {
	e.Test.Scope = e.Scope
	test := e.Test.Evaluate().StringValue
	if hasError(test) {
		return test
	}
	if !valueIsBoolean(test) {
		return newError("invalid conditional type(s)")
	}

	if getBoolFromValue(test) {
		e.Consequent.Scope = e.Scope
		return e.Consequent.Evaluate().StringValue
	} else {
		e.Alternate.Scope = e.Scope
		return e.Alternate.Evaluate().StringValue
	}
}

func evaluateFunction(e Expression) Function {
	var params []Parameter
	for _, param := range e.Params {
		params = append(params, Parameter{
			Name:  param.Name,
			Value: Value{},
		})
	}
	for _, line := range e.Body.Body {
		line.Argument.Scope = &Scope{
			Variables: map[string]Value{},
			Parent:    e.Scope,
		}
	}
	return Function{
		Parameters: params,
		Body:       *e.Body,
		Scope:      e.Scope,
	}
}

func evaluateCall(e Expression) Value {
	// e.Callee, e.Arguments
	// evaluate arguments FIRST
	var args []Value
	for _, arg := range e.Arguments {
		arg.Scope = e.Scope
		result := arg.Evaluate()
		if hasError(result.StringValue) {
			return result
		} else {
			args = append(args, result)
		}
	}
	var f Function
	if e.Callee.Name == "" {
		// recursive calling
		e.Callee.Scope = e.Scope
		result := e.Callee.Evaluate()
		if result.StringValue != "" {
			return result
		}
		f = result.FunctionValue
	} else {
		val, ok := getIdentifierValue(e.Callee.Name, e.Scope)
		if !ok {
			return newStringValue(newError("unbound identifier"))
		}
		if val.StringValue != "" {
			return newStringValue(newError("not a function"))
		}
		f = val.FunctionValue
	}
	for i, arg := range args {
		f.Parameters[i].Value = arg
	}
	return f.ExecuteFunction()
}

func (f Function) ExecuteFunction() Value {
	// look through bottom-most scope, then parameters, then f.scope...
	paramScope := Scope{
		Variables: map[string]Value{},
		Parent:    f.Scope,
	}
	for _, param := range f.Parameters {
		paramScope.Variables[param.Name] = param.Value
	}
	bottomScope := &Scope{
		Variables: map[string]Value{},
		Parent:    &paramScope,
	}
	for _, statement := range f.Body.Body {
		if statement.Type == "VariableDeclaration" {
			err := handleDeclarations(statement.Declarations, bottomScope)
			if err != nil {
				return newStringValue(err.Error())
			}
		} else { // the final return statement
			statement.Argument.Scope = bottomScope
			return statement.Argument.Evaluate()
		}
	}
	return Value{}
}

// Helper methods
func handleBody(body []ProgramChild, scope *Scope) Value { // for the whole program
	for _, child := range body {
		if child.Type == "VariableDeclaration" {
			err := handleDeclarations(child.Declarations, scope)
			if err != nil {
				return newStringValue(err.Error())
			}
		} else {
			child.Expression.Scope = scope
			return child.Expression.Evaluate()
		}
	}
	return Value{}
}

func handleDeclarations(declarations []Expression, scope *Scope) error {
	for _, decl := range declarations {
		name := decl.Id.Name
		if decl.Init.Type == "FunctionExpression" {
			decl.Init.Scope = scope
			fun := evaluateFunction(*decl.Init)
			scope.Variables[name] = newFunctionValue(fun)
		} else {
			decl.Init.Scope = scope
			expr := decl.Init.Evaluate()
			if hasError(expr.StringValue) {
				return errors.New(expr.StringValue)
			} else {
				scope.Variables[name] = expr
			}
		}
	}
	return nil
}

func evaluateValue(val Value) string {
	if val.StringValue != "" {
		return newFinalValue(val.StringValue)
	} else {
		return newFinalValue("function")
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
	return newFinalValue(strconv.Itoa(answer))
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
	return newFinalValue(boolAsString(result))
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
	re := regexp.MustCompile(`\((value \(number (-?\d+)\))\)`)

	match := re.FindStringSubmatch(str)

	numberStr := match[2]
	number, _ := strconv.Atoi(numberStr)

	return number
}

func getBoolFromValue(str string) bool {
	return strings.Contains(str, "true")
}

func newError(str string) string {
	return "(error \"" + str + " banana\")"
}

func newStringValue(str string) Value {
	return Value{
		StringValue:   str,
		FunctionValue: Function{},
	}
}

func newFunctionValue(fun Function) Value {
	return Value{
		StringValue:   "",
		FunctionValue: fun,
	}
}

func newFinalValue(str string) string {
	var result string
	if str == "true" || str == "false" {
		result = "boolean " + str
	} else if str == "function" {
		result = str
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

func getIdentifierValue(id string, scope *Scope) (Value, bool) {
	for {
		if scope == nil {
			return Value{}, false
		}
		expr, ok := scope.Variables[id]
		if ok {
			return expr, true
		} else {
			scope = scope.Parent
		}
	}
}
