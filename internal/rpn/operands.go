/*
 * Copyright (C) 2024 by Jason Figge
 */

package rpn

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInvalidOperand = fmt.Errorf("invalid operand")
)

type Operand struct {
	tokenType OpType
	i64       int64
	f64       float64
}

func (o *Operand) IsFloat() bool {
	return o.f64 != 0
}

func (o *Operand) String() string {
	return fmt.Sprintf("%v", o.Value())
}

func (o *Operand) Value() any {
	if o.f64 != 0 {
		return o.f64
	}
	return o.i64
}

func (o *Operand) Type() OpType {
	return o.tokenType
}

func (o *Operand) ToFloat() {
	if o.i64 != 0 {
		o.f64 = float64(o.i64)
		o.i64 = 0
	}
}

func OperandFromToken(token string) *Operand {
	var err error
	op := &Operand{tokenType: TokenOperand}
	if strings.Contains(token, ".") && !strings.HasSuffix(token, ".") {
		op.f64, err = strconv.ParseFloat(token, 64)
	} else {
		op.i64, err = strconv.ParseInt(token, 10, 64)
	}
	if err != nil {
		panic(fmt.Errorf("%w: cannot parse %q", ErrInvalidOperand, token))
	}
	return op
}
