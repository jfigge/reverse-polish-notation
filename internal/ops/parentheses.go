/*
 * Copyright (C) 2024 by Jason Figge
 */

package ops

import (
	"strings"
)

var (
	parenthesisMap = map[byte]Parenthesis{}
	parentheses    = [...]Parenthesis{
		{symbol: '('},
		{symbol: ')'},
	}
	parenthesisRegEx string
)

type Parenthesis struct {
	symbol byte
}

func init() {
	parts := make([]string, len(parentheses))
	for i, parenthesis := range parentheses {
		parenthesisMap[parenthesis.symbol] = parenthesis
		parts[i] += "\\" + string(parenthesis.symbol)
	}
	parenthesisRegEx = strings.Join(parts, "|")
}

func ParenthesisFromSymbol(symbol byte) (*Parenthesis, bool) {
	parenthesis, ok := parenthesisMap[symbol]
	return &parenthesis, ok
}

func ParenthesisRegEx() string {
	return parenthesisRegEx
}

func (p *Parenthesis) Symbol() byte {
	return p.symbol
}

func (p *Parenthesis) IsStart() bool {
	return p.symbol == '('
}
