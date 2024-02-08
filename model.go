package main

type Program struct {
	Body []ExpressionStatement `json:"body"`
}

func (p Program) String() string {
	return p.Body[0].Expression.String()
}

type ExpressionStatement struct {
	Expression Expression `json:"expression"`
}

type Expression struct {
	Type       string      `json:"type"`
	Left       *Expression `json:"left"`
	Operator   string      `json:"operator"`
	Right      *Expression `json:"right"`
	Test       *Expression `json:"test"`
	Consequent *Expression `json:"consequent"`
	Alternate  *Expression `json:"alternate"`
	Raw        string      `json:"raw"`
	Argument   *Expression `json:"argument"`
}

func (e Expression) String() string {
	var result string
	switch e.Type {
	case "Literal":
		if e.Raw == "true" || e.Raw == "false" {
			result = "boolean " + e.Raw
		} else {
			result = "number " + e.Raw
		}
	case "BinaryExpression":
		if isArithmeticOperator(e.Operator) {
			result = "arithmetic "
		} else {
			result = "relational "
		}
		result += e.Operator + " " + e.Left.String() + " " + e.Right.String()
	case "UnaryExpression":
		result = "unary " + e.Operator + " " + e.Argument.String()
	case "LogicalExpression":
		result = "logical " + e.Operator + " " + e.Left.String() + " " + e.Right.String()
	case "ConditionalExpression":
		result = "conditional " + e.Test.String() + " " + e.Consequent.String() + " " + e.Alternate.String()
	default:
		return ""
	}
	return "(" + result + ")"
}
