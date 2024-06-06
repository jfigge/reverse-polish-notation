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

func ParseToken(token string, lastToken TokenType) (Token, error) {
	if len(token) == 1 {
		if operator, ok := OperatorFromToken(token[0], lastToken); ok {
			return operator, nil
		}
		if parenthesis, ok := ParenthesisFromToken(token[0]); ok {
			return parenthesis, nil
		}
	}
	operand, err := OperandFromToken(token)
	if err != nil {
		return nil, err
	}

	return operand, nil
}

type Token interface {
	Type() TokenType
	String() string
}
