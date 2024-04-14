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
	return newFinalValue("function")
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
		return evaluateConditional(e)
	case "FunctionExpression":
		return newFunctionValue(evaluateFunction(e))
	case "CallExpression":
		return evaluateCall(e)
	case "AssignmentExpression":
		return newStringValue(evaluateAssignment(e))
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
	leftNumeric, err := getNumberFromValue(left)
	if err != "" {
		return newError("invalid binary type(s)")
	}
	rightNumeric, err := getNumberFromValue(right)
	if err != "" {
		return newError("invalid binary type(s)")
	}

	if isArithmeticOperator(e.Operator) {
		// Arithmetic evaluation
		return doMath(leftNumeric, rightNumeric, e.Operator)
	} else {
		// Relational evaluation
		return doComparison(leftNumeric, rightNumeric, e.Operator)
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

func evaluateConditional(e Expression) Value {
	e.Test.Scope = e.Scope
	test := e.Test.Evaluate().StringValue
	if hasError(test) {
		return newStringValue(test)
	}
	if !valueIsBoolean(test) {
		return newStringValue(newError("invalid conditional type(s)"))
	}

	if getBoolFromValue(test) {
		e.Consequent.Scope = e.Scope
		return e.Consequent.Evaluate()
	} else {
		e.Alternate.Scope = e.Scope
		return e.Alternate.Evaluate()
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
			return newStringValue(newError("not a function"))
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

func evaluateAssignment(e Expression) string {
	// check if left identifier exists
	_, exists := getIdentifierValue(e.Left.Name, e.Scope)
	if !exists {
		return newError("unbound identifier")
	}
	// evaluate right side
	e.Right.Scope = e.Scope
	value := e.Right.Evaluate()
	if hasError(value.StringValue) {
		return value.StringValue
	}
	// assign identifier to value
	//e.Scope.Variables[e.Left.Name] = value
	setIdentifierValue(e.Left.Name, value, e.Scope)
	return newFinalValue("void")
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
	return handleBody(f.Body.Body, bottomScope)
}

// for the whole program, function bodies, and loop bodies
func handleBody(body []BlockChild, scope *Scope) Value {
	for i, statement := range body {
		if statement.Type == "VariableDeclaration" {
			err := handleDeclarations(statement.Declarations, scope)
			if err != nil {
				return newStringValue(err.Error())
			}
		} else if statement.Type == "ForStatement" || statement.Type == "WhileStatement" {
			result := handleLoop(statement, scope)
			if result.StringValue != "" { // either there is an error or a return statement in the loop
				if result.StringValue == "f" { // double-check this in handleForLoop
					return newFunctionValue(result.FunctionValue)
				}
				return result
			}
		} else if statement.Type == "ReturnStatement" {
			statement.Argument.Scope = scope
			return statement.Argument.Evaluate()
		} else { // any ExpressionStatement
			statement.Expression.Scope = scope
			value := statement.Expression.Evaluate()
			if hasError(value.StringValue) {
				return value
			}
			if scope.Parent == nil && i == len(body)-1 {
				// this is a top-level program, return the value of the final expression
				return value
			}
		}
	}
	return Value{}
}

// Helper methods
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

func handleLoop(loop BlockChild, scope *Scope) Value {
	loopScope := Scope{
		Variables: map[string]Value{},
		Parent:    scope,
	}
	if len(loop.Declarations) == 0 {
		return handleWhileLoop(loop, &loopScope)
	} else {
		return handleForLoop(loop, &loopScope)
	}
}

func handleForLoop(loop BlockChild, loopScope *Scope) Value {
	err := handleDeclarations(loop.Declarations, loopScope)
	if err != nil {
		return newStringValue(err.Error())
	}
	return Value{}
}

func handleWhileLoop(loop BlockChild, loopScope *Scope) Value {
	return Value{}
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

func getNumberFromValue(str string) (int, string) {
	re := regexp.MustCompile(`\((value \(number (-?\d+)\))\)`)

	match := re.FindStringSubmatch(str)
	if len(match) < 3 {
		return 0, "error"
	}

	numberStr := match[2]
	number, _ := strconv.Atoi(numberStr)

	return number, ""
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
	} else if str == "void" {
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

func setIdentifierValue(id string, val Value, scope *Scope) {
	for {
		if scope == nil {
			return
		}
		_, ok := scope.Variables[id]
		if ok {
			scope.Variables[id] = val
			return
		} else {
			scope = scope.Parent
		}
	}
}
