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

import "math"

// MathOperators is default set for math expressions
var MathOperators = []*Operator{
	{Op: "+", Assoc: LeftAssoc, Priority: 10, Fn: func(a float64, b float64) (float64, error) { return a + b, nil }},
	{Op: "-", Assoc: LeftAssoc, Priority: 10, Fn: func(a float64, b float64) (float64, error) { return a - b, nil }},
	{Op: "*", Assoc: LeftAssoc, Priority: 20, Fn: func(a float64, b float64) (float64, error) { return a * b, nil }},
	{Op: "/", Assoc: LeftAssoc, Priority: 20, Fn: func(a float64, b float64) (float64, error) { return a / b, nil }},
	{Op: "^", Assoc: RightAssoc, Priority: 30, Fn: func(a, b float64) (float64, error) { return math.Pow(a, b), nil }},
}

// LogicOperators is default set for logic expressions
var LogicOperators = []*Operator{
	{Op: "==", Assoc: LeftAssoc, Priority: 0, Fn: func(a float64, b float64) (float64, error) {
		if a == b {
			return 1, nil
		}
		return 0, nil
	}},
	{Op: "!=", Assoc: LeftAssoc, Priority: 0, Fn: func(a float64, b float64) (float64, error) {
		if a != b {
			return 1, nil
		}
		return 0, nil
	}},
	{Op: ">", Assoc: LeftAssoc, Priority: 0, Fn: func(a float64, b float64) (float64, error) {
		if a > b {
			return 1, nil
		}
		return 0, nil
	}},
	{Op: "<", Assoc: LeftAssoc, Priority: 0, Fn: func(a float64, b float64) (float64, error) {
		if a < b {
			return 1, nil
		}
		return 0, nil
	}},
	{Op: ">=", Assoc: LeftAssoc, Priority: 0, Fn: func(a float64, b float64) (float64, error) {
		if a >= b {
			return 1, nil
		}
		return 0, nil
	}},
	{Op: "<=", Assoc: LeftAssoc, Priority: 0, Fn: func(a float64, b float64) (float64, error) {
		if a <= b {
			return 1, nil
		}
		return 0, nil
	}},
}
