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

import "testing"

func TestCalc(t *testing.T) {
	funcs := []*Function{
		NewFunction("negative", func(args ...float64) (f float64, err error) {
			return -args[0], nil
		}, 1),
		NewFunction("sum", func(args ...float64) (f float64, err error) {
			return args[0] + args[1], nil
		}, 2),
	}
	operators := []*Operator{
		NewOperator("==", 1, LeftAssoc, func(a, b float64) (float64, error) {
			if a == b {
				return 1, nil
			}
			return 0, nil
		}),
	}
	tests := []struct {
		name       string
		expression string
		expected   float64
		vars       map[string]float64
		funcs      []*Function
		operators  []*Operator
	}{
		{"simple", "((15/(7-(1+1)))*-3)-(-2+(1+1))", ((15.0 / (7.0 - (1.0 + 1.0))) * -3.0) - (-2.0 + (1.0 + 1.0)), nil, nil, nil},
		{"variables", "a+b*c", 14.0, map[string]float64{"a": 2.0, "b": 3.0, "c": 4.0}, nil, nil},
		{"functions 1 arg", "negative(10)", -10.0, nil, funcs, nil},
		{"functions 2 arg", "negative(sum(10, 20)+20)", -50.0, nil, funcs, nil},
		{"custom operator", "10 == 10", 1, nil, nil, operators},
		{"custom operator 2", "10 == 12", 0, nil, nil, operators},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := NewCalc()
			c.AddOperators(MathOperators)
			if test.funcs != nil {
				c.AddFunctions(test.funcs)
			}
			if test.operators != nil {
				c.AddOperators(test.operators)
			}
			if err := c.Prepare(test.expression); err != nil {
				t.Error(err)
			}
			actual, err := c.Execute(test.vars)
			if err != nil {
				t.Error(err)
			}
			if actual != test.expected {
				t.Errorf("Expected %f, actual %f", test.expected, actual)
			}
		})
	}
}

func TestCalc2(t *testing.T) {
	c := NewCalc()
	c.AddOperators(MathOperators)
	if err := c.Prepare("((15/(7-(1+1)))*-3)-(-2+(1+1))"); err != nil {
		t.Error(err)
	}
	actual, err := c.Execute(nil)
	if err != nil {
		t.Error(err)
	}
	expected := ((15.0 / (7.0 - (1.0 + 1.0))) * -3.0) - (-2.0 + (1.0 + 1.0))
	if actual != expected {
		t.Errorf("Expected %f, actual %f", expected, actual)
	}
}
