/*
 * Copyright (C) 2024 by Jason Figge
 */

package ops

type TokenType uint8

const ( // values pinned
	TokenEmpty TokenType = 1 << iota
	TokenOperand
	TokenOperator
	TokenParentheses
)

func ParseToken(token string, lastToken TokenType) Token {
	if len(token) == 1 {
		if operator, ok := OperatorFromToken(token[0], lastToken); ok {
			return operator
		}
		if parenthesis, ok := ParenthesisFromToken(token[0]); ok {
			return parenthesis
		}
	}
	return OperandFromToken(token)
}

type Token interface {
	Type() TokenType
	String() string
}
