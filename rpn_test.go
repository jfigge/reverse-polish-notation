/*
 * Copyright (C) 2024 by Jason Figge
 */

package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"us.figge.rpn/internal/ops"
	"us.figge.rpn/internal/rpn"
)

func Test_rpn(t *testing.T) {
	tests := map[string]struct {
		exp string
		rpn string
		ans any
		err string
	}{
		"basic-decimal": {
			exp: "1+2*4-3",
			rpn: "124*+3-",
			ans: int64(6),
		},
		"basic-float": {
			exp: "1.0+2.5*4.0-3.2/1.0",
			rpn: "12.54*+3.21/-",
			ans: 7.8,
		},
		"parentheses": {
			exp: "(4/2)",
			rpn: "42/",
			ans: 2.0,
		},
		"parentheses-1": {
			exp: "(1+2)*4-3",
			rpn: "12+4*3-",
			ans: int64(9),
		},
		"parentheses-2": {
			exp: "1+2*(4-3)",
			rpn: "1243-*+",
			ans: int64(3),
		},
		"parentheses-3": {
			exp: "1+((4+3)-2)*2-2",
			rpn: "143+2-2*+2-",
			ans: int64(9),
		},
		"parentheses-4": {
			exp: "(1+2)*(1--2)",
			rpn: "12+12--*",
			ans: int64(9),
		},
		"parentheses-5": {
			exp: "(((1+2)))*(1--2)",
			rpn: "12+12--*",
			ans: int64(9),
		},
		"parentheses-6": {
			exp: "(4+3)-2",
			rpn: "43+2-",
			ans: int64(5),
		},
		"unary": {
			exp: "2*-3",
			rpn: "23-*",
			ans: int64(-6),
		},
		"space": {
			exp: " 1 - 2 +  3   *  ( 4  / 5 ) + + 6 - - 7  ",
			rpn: "12345/*6++7---",
			ans: -16.4,
		},
		"unary-1": {
			exp: "-2*-3",
			rpn: "2-3-*",
			ans: int64(6),
		},
		"unary-2": {
			exp: "(+2--2)*-3",
			rpn: "22--3-*",
			ans: int64(-12),
		},
		"unary-3": {
			exp: "-4.1",
			rpn: "4.1-",
			ans: -4.1,
		},
		"3": {
			exp: "3",
			rpn: "3",
			ans: int64(3),
		},
		"mod": {
			exp: "6%(3-1)",
			rpn: "631-%",
			ans: int64(0),
		},
		"cannot parse": {
			exp: "-99999999999999999999",
			err: "invalid operand: cannot parse \"99999999999999999999\"",
		},
		"cannot mix int with float": {
			exp: "2%1.4",
			rpn: "21.4%",
			err: "invalid operand: cannot perform modulus operation with floats",
		},
		"cannot mix float with int": {
			exp: "1.4%2",
			rpn: "1.42%",
			err: "invalid operand: cannot perform modulus operation with floats",
		},
		"cannot perform modulus operation with floats": {
			exp: "1.4%2.3",
			rpn: "1.42.3%",
			err: "invalid operand: cannot perform modulus operation with floats",
		},
		"Unclosed parenthesis": {
			exp: "(4-1",
			err: "invalid expression: Unclosed parenthesis",
		},
		"Too many close parenthesis": {
			exp: "4-1)",
			err: "invalid expression: Too many close parenthesis",
		},
		"No valid token found": {
			exp: "2=1",
			err: "invalid syntax: no valid token found",
		},
		"Not all operands consumed": {
			exp: "1 2-1",
			rpn: "121-",
			err: "invalid expression: not all operands consumed",
		},
		"insufficient operands": {
			exp: "1*",
			rpn: "1*",
			err: "invalid expression: insufficient operands 2 != 1",
		},
		"type-mismatch-plus": {
			exp: "1.0+2",
			rpn: "12+",
			ans: 3.0,
		},
		"type-mismatch-minus": {
			exp: "2-1.0",
			rpn: "21-",
			ans: 1.0,
		},
	}
	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {

			defer func() {
				if e := recover(); e != nil {
					txt := fmt.Sprintf("%v", e)
					assert.Equal(tt, test.err, txt)
				} else {
					assert.Equal(tt, "", test.err)
				}
			}()
			notation := rpn.Parse(test.exp)
			assert.Equal(tt, test.rpn, notation.String())
			fmt.Printf("%s\n", notation.String())
			ans := notation.Solve()
			assert.Equal(tt, test.ans, ans.Value())
		})
	}
}

func Test_Unreachable(t *testing.T) {
	t.Run("invalid operator qualifier", func(tt *testing.T) {
		op, ok := ops.OperatorFromToken('-', ops.TokenType(16))
		assert.Nil(tt, op)
		assert.False(tt, ok)
	})

	t.Run("validate 1 != 2", func(tt *testing.T) {
		defer func() {
			if e := recover(); e != nil {
				txt := fmt.Sprintf("%v", e)
				assert.Equal(tt, "invalid operation: invalid number of operands 1 != 2", txt)
			} else {
				tt.Errorf("failed to catch error")
			}
		}()

		op, ok := ops.OperatorFromToken('-', ops.TokenEmpty)
		assert.NotNil(tt, op)
		assert.True(tt, ok)

		op1 := ops.OperandFromToken("123")
		op2 := ops.OperandFromToken("456")
		_ = op.Solve([]*ops.Operand{op1, op2})
	})

	t.Run("validate 2 != 1", func(tt *testing.T) {
		defer func() {
			if e := recover(); e != nil {
				txt := fmt.Sprintf("%v", e)
				assert.Equal(tt, "invalid operation: invalid number of operands 2 != 1", txt)
			} else {
				tt.Errorf("failed to catch error")
			}
		}()

		op, ok := ops.OperatorFromToken('*', ops.TokenEmpty)
		assert.NotNil(tt, op)
		assert.True(tt, ok)

		op1 := ops.OperandFromToken("123")
		_ = op.Solve([]*ops.Operand{op1})
	})
}

func TestAdhoc(t *testing.T) {
	notation := rpn.Parse("1+2-3")
	fmt.Printf("%s\n", notation.String())
	answer := notation.Solve()
	fmt.Printf("%d\n", answer.Value())
}
