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
		fmt.Sprintf(`^\s*(%s|0|[1-9][0-9]*(?:\.[0-9]*)?)(.*)$`, ops.OperatorRegEx()),
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
	notation, opStack := Notation{}, make([]*ops.Operator, 0)
	for i, lastOpType := 0, ops.TokenEmpty; strings.TrimSpace(exp) != ""; i++ {
		token := nextToken(&exp, lastOpType)
		lastOpType = token.Type()
		if token.Type() == ops.TokenOperand {
			notation = append(notation, token.(*ops.Operand))
		} else if op := token.(*ops.Operator); !op.Exclude() {
			if op.String() == "(" {
				opStack = append(opStack, op)
				lastOpType = ops.TokenEmpty
			} else if op.String() == ")" {
				notation = decantStack(notation, &opStack, func(i int) bool {
					return i >= 0 && (opStack)[i].String() != "("
				})
				opStack = opStack[:len(opStack)-1]
			} else {
				notation = decantStack(notation, &opStack, func(i int) bool {
					return i >= 0 && opStack[i].Precedence() > op.Precedence()
				})
				opStack = append(opStack, op)
			}
		}
	}
	return decantStack(notation, &opStack, func(i int) bool { return i >= 0 })
}

func decantStack(notation Notation, operatorStack *[]*ops.Operator, f func(i int) bool) Notation {
	for i := len(*operatorStack) - 1; f(i); i-- {
		notation = append(notation, (*operatorStack)[i])
		*operatorStack = (*operatorStack)[:i]
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
