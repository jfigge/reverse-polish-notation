/*
 * Copyright (C) 2024 by Jason Figge
 */

package rpn

import (
	"fmt"
	"regexp"

	"us.figge.rpn/internal/ops"
)

var (
	ErrInsufficientOperands = fmt.Errorf("insufficient operands")
	ErrInvalidExpression    = fmt.Errorf("invalid expression")
	ErrInvalidSyntax        = fmt.Errorf("invalid syntax")
)

type Notation []*ops.Op

func (rpn Notation) String() string {
	result := ""
	for _, o := range rpn {
		switch op := o.Source().(type) {
		case *ops.Operand:
			result = fmt.Sprintf("%s%s", result, op.String())
		case *ops.Operator:
			result = fmt.Sprintf("%s%s", result, string(op.Symbol()))
		default:
		}
	}
	return result
}

var (
	tokenizer = regexp.MustCompile(
		fmt.Sprintf(`^\s*(%s|%s|0|[1-9][0-9]*(?:\.[0-9]*)?)(.*)$`, ops.OperatorRegEx(), ops.ParenthesisRegEx()),
	)
)

func Parse(exp string) (Notation, error) {
	n, _, e := parse(exp)
	return n, e
}
func parse(exp string) (Notation, string, error) {
	lastOpType := ops.OpTypeEmpty
	operatorStack := make([]*ops.Op, 0)
	notation := Notation{}
	var o *ops.Op
	var i int
	var err error
forLoop:
	for exp != "" {
		i++
		o, exp, err = nextToken(exp, lastOpType)
		if err != nil {
			return nil, "", err
		}
		switch op := o.Source().(type) {
		case *ops.Operator:
			if op.Exclude() {
				break
			}
			lastOpType = ops.OpTypeOperator
			length := len(operatorStack) - 1
			for length >= 0 && operatorStack[length].Operator().Precedence() > op.Precedence() {
				notation = append(notation, operatorStack[length])
				operatorStack = operatorStack[:length]
				length--
			}
			operatorStack = append(operatorStack, o)
		case *ops.Operand:
			lastOpType = ops.OpTypeOperand
			notation = append(notation, o)
		case *ops.Parenthesis:
			if op.IsStart() {
				var subNotation Notation
				subNotation, exp, err = parse(exp)
				if err != nil {
					return nil, "", err
				}
				lastOpType = ops.OpTypeOperand
				notation = append(notation, subNotation...)
			} else {
				break forLoop
			}
		default:
			return nil, "", fmt.Errorf("%w: Invalid token at %d", ErrInvalidExpression, i)
		}
	}
	for len(operatorStack) > 0 {
		notation = append(notation, operatorStack[len(operatorStack)-1])
		operatorStack = operatorStack[:len(operatorStack)-1]
	}
	return notation, exp, nil
}

func nextToken(exp string, topType ops.OpType) (*ops.Op, string, error) {
	parts := tokenizer.FindStringSubmatch(exp)
	if len(parts) != 3 {
		return nil, exp, fmt.Errorf("%w: no valid token found", ErrInvalidSyntax)
	}
	op, err := ops.ParseOp(parts[1], topType)
	if err != nil {
		return nil, exp, err
	}
	return op, parts[2], nil
}

func (rpn Notation) Solve() (*ops.Operand, error) {
	operandStack := make([]*ops.Operand, 0)
	for _, a := range rpn {
		switch op := a.Source().(type) {
		case *ops.Operator:
			length := uint8(len(operandStack))
			if length >= op.Operands() {
				answer, err := op.Solve(operandStack[length-op.Operands() : length])
				if err != nil {
					return nil, err
				}
				operandStack = operandStack[:length-op.Operands()]
				operandStack = append(operandStack, answer)
			} else {
				return nil, fmt.Errorf("%w: expected %d, received: %d", ErrInsufficientOperands, length, op.Operands())
			}
		case *ops.Operand:
			operandStack = append(operandStack, op)
		default:
			return nil, fmt.Errorf("%w: Unexpected value at %q", ErrInvalidExpression, a)
		}

	}
	if len(operandStack) != 1 {
		return nil, fmt.Errorf("%w: Not all operands consumed", ErrInvalidExpression)
	}
	return operandStack[0], nil
}
