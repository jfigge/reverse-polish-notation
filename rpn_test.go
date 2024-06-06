/*
 * Copyright (C) 2024 by Jason Figge
 */

package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"us.figge.rpn/internal/rpn"
)

func Test_rpn(t *testing.T) {
	tests := map[string]struct {
		exp string
		rpn string
		ans any
	}{
		"basic-decimal": {
			exp: "1+2*4-3",
			rpn: "124*+3-",
			ans: int64(6),
		},
		"basic-float": {
			exp: "1.0+2.5*4.0-3.2",
			rpn: "12.54*+3.2-",
			ans: 7.8,
		},
		"parentheses": {
			exp: "(1+2)",
			rpn: "12+",
			ans: int64(3),
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
		"unary": {
			exp: "2*-3",
			rpn: "23-*",
			ans: int64(-6),
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
			exp: "-4",
			rpn: "4-",
			ans: int64(-4),
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
	}
	for name, test := range tests {
		t.Run(name, func(tt *testing.T) {
			notation, err := rpn.Parse(test.exp)
			if err != nil {
				assert.FailNowf(tt, "Unexpected parse error", "%v", err)
			}
			assert.Equal(tt, test.rpn, notation.String())
			fmt.Printf("%s\n", notation.String())
			ans, err := notation.Solve()
			if err != nil {
				assert.FailNowf(tt, "Unexpected solve error", "%v", err)
			}
			assert.Equal(tt, test.ans, ans.Value())
		})
	}
}

func TestAdhoc(t *testing.T) {
	notation, err := rpn.Parse("1+2-3")
	if err != nil {
		fmt.Printf("%v\n", err)
		t.FailNow()
	}
	fmt.Printf("%s\n", notation.String())
	answer, err := notation.Solve()
	if err != nil {
		fmt.Printf("%v\n", err)
		t.FailNow()
	}
	fmt.Printf("%d\n", answer.Value())
}
