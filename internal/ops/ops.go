/*
 * Copyright (C) 2024 by Jason Figge
 */

package ops

// ****** OpType Enum *********************************************************

type OpType uint8

const ( // values pinned
	OpTypeEmpty OpType = 1 << iota
	OpTypeOperand
	OpTypeOperator
	OpTypeParentheses
)

func (oc OpType) String() string {
	return [...]string{"Operator", "Parentheses", "Operand"}[oc-1]
}

func (oc OpType) EnumIndex() int {
	return int(oc)
}

// ****** Op Structure ********************************************************

type Op struct {
	opType OpType
	source any
}

func (o *Op) Type() OpType {
	return o.opType
}

func (o *Op) Source() any {
	return o.source
}

func (o *Op) Operator() *Operator {
	return o.source.(*Operator)
}

func ParseOp(token string, topType OpType) (*Op, error) {
	if len(token) == 1 {
		if operator, ok := OperatorFromSymbol(token[0], topType); ok {
			return &Op{opType: OpTypeOperator, source: operator}, nil
		}
		if parenthesis, ok := ParenthesisFromSymbol(token[0]); ok {
			return &Op{opType: OpTypeParentheses, source: parenthesis}, nil
		}
	}
	operand, err := OperandFromToken(token)
	if err != nil {
		return nil, err
	}

	return &Op{opType: OpTypeOperand, source: operand}, nil
}

// ****** Op Interfaces *******************************************************

type OpOperator interface {
	Operands() int
	Precedence() int
	Solver() func([]any) (any, error)
}

type OpParentheses interface {
	IsOpen() bool
}

type OpOperand interface {
	IsFloat() bool
	Float() float64
	Int() int64
}
