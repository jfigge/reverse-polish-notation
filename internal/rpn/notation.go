/*
 * Copyright (C) 2024 by Jason Figge
 */

package rpn

import (
	"fmt"
	"regexp"
	"strings"

	"us.figge.rpn/internal/ops"
)

var (
	ErrInvalidExpression = fmt.Errorf("invalid expression")
	ErrInvalidSyntax     = fmt.Errorf("invalid syntax")
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

func Parse(exp string) Notation {
	return parse(&exp, 0)
}
func parse(exp *string, subExpression int) Notation {
	notation, operatorStack := Notation{}, make([]*ops.Operator, 0)
	subExpressionStart := subExpression
	for i, lastOpType := 0, ops.TokenEmpty; strings.TrimSpace(*exp) != "" && subExpressionStart == subExpression; i++ {
		token := nextToken(exp, lastOpType)
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
				notation = append(notation, parse(exp, 1)...)
			} else {
				subExpression -= 1
			}
		}
	}
	if subExpression > 0 {
		panic(fmt.Errorf("%w: Unclosed parenthesis", ErrInvalidExpression))
	} else if subExpression < 0 {
		panic(fmt.Errorf("%w: Too many close parenthesis", ErrInvalidExpression))
	}
	for len(operatorStack) > 0 {
		notation = append(notation, operatorStack[len(operatorStack)-1])
		operatorStack = operatorStack[:len(operatorStack)-1]
	}
	return notation
}

func nextToken(exp *string, lastToken ops.TokenType) ops.Token {
	parts := tokenizer.FindStringSubmatch(*exp)
	if len(parts) != 3 {
		panic(fmt.Errorf("%w: no valid token found", ErrInvalidSyntax))
	}
	*exp = parts[2]
	return ops.ParseToken(parts[1], lastToken)
}

func (rpn Notation) Solve() *ops.Operand {
	operandStack := make([]*ops.Operand, 0)
	for _, token := range rpn {
		switch op := any(token).(type) {
		case *ops.Operator:
			length := uint8(len(operandStack))
			if length >= op.Operands() {
				answer := op.Solve(operandStack[length-op.Operands() : length])
				operandStack = operandStack[:length-op.Operands()]
				operandStack = append(operandStack, answer)
			} else {
				panic(fmt.Errorf("%w: insufficient operands %d != %d", ErrInvalidExpression, op.Operands(), length))
			}
		case *ops.Operand:
			operandStack = append(operandStack, op)
		}

	}
	if len(operandStack) != 1 {
		panic(fmt.Errorf("%w: not all operands consumed", ErrInvalidExpression))
	}
	return operandStack[0]
}
