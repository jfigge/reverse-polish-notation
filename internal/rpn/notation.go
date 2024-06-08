/*
 * Copyright (C) 2024 by Jason Figge
 */

package rpn

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	ErrInvalidExpression = fmt.Errorf("invalid expression")
	ErrInvalidSyntax     = fmt.Errorf("invalid syntax")
)

var (
	tokenizer = regexp.MustCompile(
		fmt.Sprintf(`^\s*(%s|0|[1-9][0-9]*(?:\.[0-9]*)?)(.*)$`, OperatorRegEx()),
	)
)

const ( // values pinned
	TokenEmpty OpType = 1 << iota
	TokenOperand
	TokenOperator
	TokenParentheses
)

type OpType uint8
type Op interface {
	Type() OpType
	String() string
}
type Notation []Op

func (rpn Notation) String() string {
	result := ""
	for _, token := range rpn {
		result = fmt.Sprintf("%s%s", result, token.String())
	}
	return result
}

func Parse(exp string) Notation {
	notation, opStack := Notation{}, make([]*Operator, 0)
	for i, lastOpType := 0, TokenEmpty; strings.TrimSpace(exp) != ""; i++ {
		token := nextToken(&exp, lastOpType)
		lastOpType = token.Type()
		if token.Type() == TokenOperand {
			notation = append(notation, token.(*Operand))
		} else if op := token.(*Operator); !op.Exclude() {
			if op.String() == "(" {
				opStack = append(opStack, op)
				lastOpType = TokenEmpty
			} else if op.String() == ")" {
				notation = decantStack(notation, &opStack, func(i int) bool { return true })
				if len(opStack) == 0 {
					panic(fmt.Errorf("%w: Too many close parenthesis", ErrInvalidExpression))
				}
				opStack = opStack[:len(opStack)-1]
			} else {
				notation = decantStack(notation, &opStack, func(i int) bool { return opStack[i].Presedence() > op.Presedence() })
				opStack = append(opStack, op)
			}
		}
	}
	notation = decantStack(notation, &opStack, func(i int) bool { return true })
	if len(opStack) > 0 && opStack[len(opStack)-1].String() == "(" {
		panic(fmt.Errorf("%w: Unclosed parenthesis", ErrInvalidExpression))
	}
	return notation
}

func decantStack(notation Notation, opStack *[]*Operator, f func(i int) bool) Notation {
	for i := len(*opStack) - 1; i >= 0 && f(i) && (*opStack)[i].String() != "("; i-- {
		notation = append(notation, (*opStack)[i])
		*opStack = (*opStack)[:i]
	}
	return notation
}

func nextToken(exp *string, lastToken OpType) Op {
	parts := tokenizer.FindStringSubmatch(*exp)
	if len(parts) != 3 || len(parts[1]) == 0 {
		panic(fmt.Errorf("%w: no valid token found", ErrInvalidSyntax))
	}
	*exp = parts[2]
	if operator, ok := OperatorFromToken(parts[1][0], lastToken); ok {
		return operator
	}
	return OperandFromToken(parts[1])
}

func (rpn Notation) Solve() *Operand {
	operandStack := make([]*Operand, 0)
	for _, token := range rpn {
		switch op := any(token).(type) {
		case *Operator:
			length := uint8(len(operandStack))
			if length >= op.Operands() {
				answer := op.Solve(operandStack[length-op.Operands() : length])
				operandStack = operandStack[:length-op.Operands()]
				operandStack = append(operandStack, answer)
			} else {
				panic(fmt.Errorf("%w: insufficient operands %d != %d", ErrInvalidExpression, op.Operands(), length))
			}
		case *Operand:
			operandStack = append(operandStack, op)
		}

	}
	if len(operandStack) != 1 {
		panic(fmt.Errorf("%w: not all operands consumed", ErrInvalidExpression))
	}
	return operandStack[0]
}
