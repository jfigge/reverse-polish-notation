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

var (
	tokenizer = regexp.MustCompile(
		fmt.Sprintf(`^\s*(%s|%s|0|[1-9][0-9]*(?:\.[0-9]*)?)(.*)$`, ops.OperatorRegEx(), ops.ParenthesisRegEx()),
	)
)

type Notation []ops.Token

func (rpn Notation) String() string {
	result := ""
	for _, token := range rpn {
		result = fmt.Sprintf("%s%s", result, token.String())
	}
	return result
}

func Parse(exp string) (Notation, error) {
	n, _, e := parse(exp)
	return n, e
}
func parse(exp string) (Notation, string, error) {
	notation := Notation{}
	operatorStack := make([]*ops.Operator, 0)
	lastOpType := ops.TokenEmpty
	var token ops.Token
	var i int
	var err error
forLoop:
	for exp != "" {
		i++
		token, exp, err = nextToken(exp, lastOpType)
		if err != nil {
			return nil, "", err
		}
		lastOpType = token.Type()
		switch op := any(token).(type) {
		case *ops.Operator:
			if op.Exclude() {
				break
			}
			length := len(operatorStack) - 1
			for length >= 0 && operatorStack[length].Precedence() > op.Precedence() {
				notation = append(notation, operatorStack[length])
				operatorStack = operatorStack[:length]
				length--
			}
			operatorStack = append(operatorStack, op)
		case *ops.Operand:
			notation = append(notation, op)
		case *ops.Parenthesis:
			if op.IsStart() {
				var subNotation Notation
				subNotation, exp, err = parse(exp)
				if err != nil {
					return nil, "", err
				}
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

func nextToken(exp string, lastToken ops.TokenType) (ops.Token, string, error) {
	parts := tokenizer.FindStringSubmatch(exp)
	if len(parts) != 3 {
		return nil, exp, fmt.Errorf("%w: no valid token found", ErrInvalidSyntax)
	}
	op, err := ops.ParseToken(parts[1], lastToken)
	if err != nil {
		return nil, exp, err
	}
	return op, parts[2], nil
}

func (rpn Notation) Solve() (*ops.Operand, error) {
	operandStack := make([]*ops.Operand, 0)
	for _, token := range rpn {
		switch op := any(token).(type) {
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
			return nil, fmt.Errorf("%w: Unexpected value at %q", ErrInvalidExpression, token)
		}

	}
	if len(operandStack) != 1 {
		return nil, fmt.Errorf("%w: Not all operands consumed", ErrInvalidExpression)
	}
	return operandStack[0], nil
}
