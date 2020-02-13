// Copyright (c) 2020 Alexander Kiryukhin <a.kiryukhin@mail.ru>

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package executor

import (
	"errors"
	"fmt"
)

// Calc calculates expressions
type Calc struct {
	preparedTokens []*token
	functions      map[string]*Function
	operators      map[string]*Operator
}

// NewCalc instantinates new calculator
func NewCalc() *Calc {
	c := &Calc{
		functions: map[string]*Function{},
		operators: map[string]*Operator{},
	}
	return c
}

// Prepare expression before execution
func (c *Calc) Prepare(expression string) error {
	t := newTokenizer(expression, c.operators)
	if err := t.tokenize(); err != nil {
		return err
	}
	tkns, err := t.toRPN()
	if err != nil {
		return err
	}
	c.preparedTokens = tkns
	return nil
}

// Execute prepared expression with variables at `vars` argument
func (c *Calc) Execute(vars map[string]float64) (float64, error) {
	if len(c.preparedTokens) == 0 {
		return 0, errors.New("must prepare expression")
	}
	if vars == nil {
		vars = map[string]float64{}
	}
	var stack []float64
	for _, tkn := range c.preparedTokens {
		switch tkn.Type {
		case literalType:
			stack = append(stack, tkn.FValue)
		case operatorType:
			sz := len(stack)
			if sz < 2 {
				return 0, errors.New("empty stack")
			}
			var args []float64
			args, stack = stack[sz-2:], stack[:sz-2]

			if op, ok := c.operators[tkn.SValue]; ok {
				res, err := op.Fn(args[0], args[1])
				if err != nil {
					return 0, err
				}
				stack = append(stack, res)
			} else {
				return 0, fmt.Errorf("unknown operator '%s'", tkn.SValue)
			}
		case functionType:
			fn, exists := c.functions[tkn.SValue]
			if !exists {
				return 0, fmt.Errorf("unknown function '%s'", tkn.SValue)
			}
			sz := len(stack)
			if sz < fn.Places {
				return 0, errors.New("not enough args")
			}
			var args []float64
			args, stack = stack[sz-fn.Places:], stack[:sz-fn.Places]
			res, err := fn.Fn(args...)
			if err != nil {
				return 0, err
			}
			stack = append(stack, res)
		case variableType:
			res, exists := vars[tkn.SValue]
			if !exists {
				return 0, fmt.Errorf("unknown variable '%s'", tkn.SValue)
			}
			stack = append(stack, res)
		default:
			return 0, fmt.Errorf("unknown token %d, %s, %f", tkn.Type, tkn.SValue, tkn.FValue)
		}
	}
	if len(stack) != 1 {
		return 0, errors.New("invalid expression")
	}
	return stack[0], nil
}

// AddFunction adds custom function
func (c *Calc) AddFunction(cf *Function) {
	c.functions[cf.Name] = cf
}

// AddOperator adds custom operator
func (c *Calc) AddOperator(op *Operator) {
	c.operators[op.Op] = op
}

// AddFunctions xadds many custom functions
func (c *Calc) AddFunctions(funcs []*Function) {
	for _, fn := range funcs {
		c.AddFunction(fn)
	}
}

// AddOperators adds many custom operators
func (c *Calc) AddOperators(operators []*Operator) {
	for _, op := range operators {
		c.AddOperator(op)
	}
}
