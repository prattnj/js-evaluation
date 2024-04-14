package main

type Program struct {
	Body []BlockChild `json:"body"`
}

type BlockChild struct {
	Type         string         `json:"type"`
	Argument     Expression     `json:"argument"`
	Expression   Expression     `json:"expression"`
	Declarations []Expression   `json:"declarations"`
	Test         Expression     `json:"test"`   // while, for
	Body         BlockStatement `json:"body"`   // while, for
	Init         Expression     `json:"init"`   // for
	Update       Expression     `json:"update"` // for
}

type Expression struct {
	Type       string          `json:"type"` // all of them
	Scope      *Scope          // all of them
	Left       *Expression     `json:"left"`       // binary, logical, assignment
	Operator   string          `json:"operator"`   // binary, logical, assignment
	Right      *Expression     `json:"right"`      // binary, logical, assignment
	Test       *Expression     `json:"test"`       // conditional
	Consequent *Expression     `json:"consequent"` // conditional
	Alternate  *Expression     `json:"alternate"`  // conditional
	Raw        string          `json:"raw"`        // literal
	Argument   *Expression     `json:"argument"`   // unary
	Name       string          `json:"name"`       // identifier
	Params     []Expression    `json:"params"`     // function
	Body       *BlockStatement `json:"body"`       // function
	Callee     *Expression     `json:"callee"`     // call
	Arguments  []Expression    `json:"arguments"`  // call
	Id         *Expression     `json:"id"`         // variable declarator
	Init       *Expression     `json:"init"`       // variable declarator
}

type Function struct {
	Parameters []Parameter
	Body       Program
	Scope      *Scope // the scope that the function lies within. If x is declared in this function, this scope doesn't have it
}

type Value struct {
	StringValue   string
	FunctionValue Function
}

type Parameter struct {
	Name  string
	Value Value
}

type Scope struct {
	Variables map[string]Value
	Parent    *Scope
}
