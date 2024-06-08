/*
 * Copyright (C) 2024 by Jason Figge
 */

package rpn

import (
	"fmt"
	"strings"
)

var (
	ErrInvalidOperation = fmt.Errorf("invalid operation")
)

var (
	opMap     = map[byte][]*Operator{}
	operators = [...]Operator{
		{tokenType: TokenParentheses, precedence: 0, symbol: '('},
		{tokenType: TokenParentheses, precedence: 0, symbol: ')'},
		{tokenType: TokenOperator, precedence: 11, symbol: '-', operands: 2, solver: subtract, qualifiers: TokenOperand | TokenParentheses},
		{tokenType: TokenOperator, precedence: 12, symbol: '+', operands: 2, solver: add, qualifiers: TokenOperand | TokenParentheses},
		{tokenType: TokenOperator, precedence: 21, symbol: '*', operands: 2, solver: multiply},
		{tokenType: TokenOperator, precedence: 22, symbol: '%', operands: 2, solver: mod},
		{tokenType: TokenOperator, precedence: 23, symbol: '/', operands: 2, solver: divide},
		{tokenType: TokenOperator, precedence: 31, symbol: '+', operands: 1, qualifiers: TokenEmpty | TokenOperator},
		{tokenType: TokenOperator, precedence: 32, symbol: '-', operands: 1, solver: negative, qualifiers: TokenEmpty | TokenOperator},
	}
	opRegEx string
)

type Operator struct {
	tokenType  OpType
	precedence uint8
	symbol     byte
	operands   uint8
	solver     func([]*Operand) *Operand
	qualifiers OpType
}

func OperatorFromToken(symbol byte, lastToken OpType) (*Operator, bool) {
	ops, ok := opMap[symbol]
	if !ok {
		return nil, false
	} else if ops[0].qualifiers == 0 {
		return ops[0], true
	}
	for _, op := range ops {
		if op.qualifiers&lastToken == lastToken {
			return op, true
		}
	}
	return nil, false
}

func OperatorRegEx() string {
	if opRegEx == "" {
		parts := make([]string, 0, len(operators))
		for i := 0; i < len(operators); i++ {
			op := operators[i]
			ops, exists := opMap[op.symbol]
			if !exists {
				ops = []*Operator{}
				var escape string
				if op.symbol == '+' || op.symbol == '*' || op.symbol == '(' || op.symbol == ')' {
					escape = `\`
				}
				parts = append(parts, escape+string(op.symbol))
			}
			opMap[op.symbol] = append(ops, &op)
		}
		opRegEx = strings.Join(parts, "|")
	}
	return opRegEx
}

func (o *Operator) Operands() uint8 {
	return o.operands
}

func (o *Operator) Presedence() uint8 {
	return o.precedence
}

func (o *Operator) String() string {
	return string(o.symbol)
}

func (o *Operator) Exclude() bool {
	return o.solver == nil && o.Operands() > 0
}

func (o *Operator) Solve(args []*Operand) *Operand {
	return o.solver(args)
}

func (o *Operator) Type() OpType {
	return o.tokenType
}

// ****** Operations **********************************************************

func validateOne(operands []*Operand) {
	if len(operands) != 1 {
		panic(fmt.Errorf("%w: invalid number of operands 1 != %d", ErrInvalidOperation, len(operands)))
	}
}
func validateTwo(operands []*Operand) {
	if len(operands) != 2 {
		panic(fmt.Errorf("%w: invalid number of operands 2 != %d", ErrInvalidOperation, len(operands)))
	}
}
func matchTypes(operands []*Operand) {
	validateTwo(operands)
	if operands[0].IsFloat() && !operands[1].IsFloat() {
		operands[1].ToFloat()
	} else if !operands[0].IsFloat() && operands[1].IsFloat() {
		operands[0].ToFloat()
	}
}

func subtract(operands []*Operand) *Operand {
	matchTypes(operands)
	if operands[0].IsFloat() {
		return &Operand{f64: operands[0].f64 - operands[1].f64}
	}
	return &Operand{i64: operands[0].i64 - operands[1].i64}
}

func add(operands []*Operand) *Operand {
	matchTypes(operands)
	if operands[0].IsFloat() {
		return &Operand{f64: operands[0].f64 + operands[1].f64}
	}
	return &Operand{i64: operands[0].i64 + operands[1].i64}
}

func multiply(operands []*Operand) *Operand {
	matchTypes(operands)
	if operands[0].IsFloat() {
		return &Operand{f64: operands[0].f64 * operands[1].f64}
	}
	return &Operand{i64: operands[0].i64 * operands[1].i64}
}

func divide(operands []*Operand) *Operand {
	matchTypes(operands)
	if operands[0].IsFloat() {
		return &Operand{f64: operands[0].f64 / operands[1].f64}
	}
	return &Operand{f64: float64(operands[0].i64) / float64(operands[1].i64)}
}

func mod(operands []*Operand) *Operand {
	validateTwo(operands)
	if operands[0].IsFloat() || operands[1].IsFloat() {
		panic(fmt.Errorf("%w: cannot perform modulus operation with floats", ErrInvalidOperand))
	}
	return &Operand{i64: operands[0].i64 % operands[1].i64}
}

func negative(operands []*Operand) *Operand {
	validateOne(operands)
	if operands[0].IsFloat() {
		return &Operand{f64: -operands[0].f64}
	}
	return &Operand{i64: -operands[0].i64}
}
