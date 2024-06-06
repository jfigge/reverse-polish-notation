/*
 * Copyright (C) 2024 by Jason Figge
 */

package ops

import (
	"fmt"
	"strings"
)

var (
	ErrInvalidOperandCount = fmt.Errorf("invalid number of operands")
	ErrOperandTypeMismatch = fmt.Errorf("operand type mistmatch")
)

var ( // https://www.tutorialspoint.com/go/go_operators_precedence.htm
	opMap     = map[byte][]*Operator{}
	operators = [...]Operator{
		{precedence: 11, symbol: '-', operands: 2, solver: subtract, qualifiers: OpTypeOperand | OpTypeParentheses},
		{precedence: 12, symbol: '+', operands: 2, solver: add, qualifiers: OpTypeOperand | OpTypeParentheses},
		{precedence: 21, symbol: '*', operands: 2, solver: multiply},
		{precedence: 22, symbol: '%', operands: 2, solver: mod},
		{precedence: 23, symbol: '/', operands: 2, solver: divide},
		{precedence: 31, symbol: '+', operands: 1, qualifiers: OpTypeEmpty | OpTypeOperator},
		{precedence: 32, symbol: '-', operands: 1, solver: negative, qualifiers: OpTypeEmpty | OpTypeOperator},
	}
	opRegEx string
)

type Operator struct {
	precedence uint8
	symbol     byte
	operands   uint8
	solver     func([]*Operand) (*Operand, error)
	qualifiers OpType
}

func init() {
	parts := make([]string, 0, len(operators))
	for i := 0; i < len(operators); i++ {
		op := operators[i]
		ops, exists := opMap[op.symbol]
		if !exists {
			ops = []*Operator{}
			var escape string
			if op.symbol == '+' || op.symbol == '*' {
				escape = `\`
			}
			parts = append(parts, escape+string(op.symbol))
		}
		opMap[op.symbol] = append(ops, &op)
	}
	opRegEx = strings.Join(parts, "|")
}

func OperatorFromSymbol(symbol byte, topOpType OpType) (*Operator, bool) {
	ops, ok := opMap[symbol]
	if !ok {
		return nil, false
	} else if ops[0].qualifiers == 0 {
		return ops[0], true
	}
	for _, op := range ops {
		if op.qualifiers&topOpType == topOpType {
			return op, true
		}
	}
	return nil, false
}

func OperatorRegEx() string {
	return opRegEx
}

func (o *Operator) Operands() uint8 {
	return o.operands
}

func (o *Operator) Precedence() uint8 {
	return o.precedence
}

func (o *Operator) Symbol() byte {
	return o.symbol
}

func (o *Operator) Exclude() bool {
	return o.solver == nil
}

func (o *Operator) Solve(args []*Operand) (*Operand, error) {
	return o.solver(args)
}

// ****** Operations **********************************************************

func validateOne(operands []*Operand) error {
	if len(operands) != 1 {
		return fmt.Errorf("%w: expected 1, received %d", ErrInvalidOperandCount, len(operands))
	}
	return nil
}
func validateTwo(operands []*Operand) error {
	if len(operands) != 2 {
		return fmt.Errorf("%w: expected 2, received %d", ErrInvalidOperandCount, len(operands))
	} else if operands[0].IsFloat() != operands[1].IsFloat() {
		return fmt.Errorf("%w: Cannot mix int with float", ErrOperandTypeMismatch)
	}
	return nil
}

func subtract(operands []*Operand) (*Operand, error) {
	if err := validateTwo(operands); err != nil {
		return nil, err
	}
	if operands[0].IsFloat() {
		return &Operand{f64: operands[0].f64 - operands[1].f64}, nil
	}
	return &Operand{i64: operands[0].i64 - operands[1].i64}, nil
}

func add(operands []*Operand) (*Operand, error) {
	if err := validateTwo(operands); err != nil {
		return nil, err
	}
	if operands[0].IsFloat() {
		return &Operand{f64: operands[0].f64 + operands[1].f64}, nil
	}
	return &Operand{i64: operands[0].i64 + operands[1].i64}, nil
}

func multiply(operands []*Operand) (*Operand, error) {
	if err := validateTwo(operands); err != nil {
		return nil, err
	}
	if operands[0].IsFloat() {
		return &Operand{f64: operands[0].f64 * operands[1].f64}, nil
	}
	return &Operand{i64: operands[0].i64 * operands[1].i64}, nil
}

func divide(operands []*Operand) (*Operand, error) {
	if err := validateTwo(operands); err != nil {
		return nil, err
	}
	if operands[0].IsFloat() {
		return &Operand{f64: operands[0].f64 / operands[1].f64}, nil
	}
	return &Operand{i64: operands[0].i64 / operands[1].i64}, nil
}

func mod(operands []*Operand) (*Operand, error) {
	if err := validateTwo(operands); err != nil {
		return nil, err
	}
	if operands[0].IsFloat() {
		return nil, fmt.Errorf("%w: cannot perform modulus operation with floats", ErrInvalidOperand)
	}
	return &Operand{i64: operands[0].i64 % operands[1].i64}, nil
}

func negative(operands []*Operand) (*Operand, error) {
	if err := validateOne(operands); err != nil {
		return nil, err
	}
	if operands[0].IsFloat() {
		return &Operand{f64: -operands[0].f64}, nil
	}
	return &Operand{i64: -operands[0].i64}, nil
}
