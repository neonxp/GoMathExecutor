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
	"testing"
)

func TestTokenize(t *testing.T) {
	operators := map[string]*Operator{
		"+": {Op: "+", Assoc: LeftAssoc, Priority: 1, Fn: func(a float64, b float64) (float64, error) { return a + b, nil }},
		"-": {Op: "-", Assoc: LeftAssoc, Priority: 1, Fn: func(a float64, b float64) (float64, error) { return a - b, nil }},
		"*": {Op: "*", Assoc: LeftAssoc, Priority: 2, Fn: func(a float64, b float64) (float64, error) { return a * b, nil }},
		"/": {Op: "/", Assoc: LeftAssoc, Priority: 2, Fn: func(a float64, b float64) (float64, error) { return a / b, nil }},
	}
	tk := newTokenizer("((15/(7-(1+1)))*-3)-(-2+(1+1))", operators)
	if err := tk.tokenize(); err != nil {
		t.Error(err)
	}
	tkns, err := tk.toRPN()
	if err != nil {
		t.Error(err)
	}
	expected := []token{
		{literalType, "", 15},
		{literalType, "", 7},
		{literalType, "", 1},
		{literalType, "", 1},
		{operatorType, "+", 0},
		{operatorType, "-", 0},
		{operatorType, "/", 0},
		{literalType, "", -3},
		{operatorType, "*", 0},
		{literalType, "", -2},
		{literalType, "", 1},
		{literalType, "", 1},
		{operatorType, "+", 0},
		{operatorType, "+", 0},
		{operatorType, "-", 0},
	}
	if len(tkns) != len(expected) {
		t.Errorf("Expected len = %d, got %d", len(expected), len(tkns))
	}
	for i, tkn := range tkns {
		if tkn.Type != expected[i].Type {
			t.Errorf("Expected type %d, got %d at pos %d", expected[i].Type, tkn.Type, i)
		}
		if tkn.SValue != expected[i].SValue {
			t.Errorf("Expected %s, got %s at pos %d", expected[i].SValue, tkn.SValue, i)
		}
		if tkn.FValue != expected[i].FValue {
			t.Errorf("Expected %f, got %f at pos %d", expected[i].FValue, tkn.FValue, i)
		}
	}
	tk = newTokenizer("a**b==10", operators)
	if err := tk.tokenize(); err != nil {
		t.Error(err)
	}
	expected = []token{
		{variableType, "a", 0},
		{operatorType, "**", 0},
		{variableType, "b", 0},
		{operatorType, "==", 0},
		{literalType, "", 10},
	}
	if len(tk.tkns) != len(expected) {
		t.Errorf("Expected len = %d, got %d", len(expected), len(tkns))
	}
	for i, tkn := range tk.tkns {
		if tkn.Type != expected[i].Type {
			t.Errorf("Expected type %d, got %d at pos %d", expected[i].Type, tkn.Type, i)
		}
		if tkn.SValue != expected[i].SValue {
			t.Errorf("Expected %s, got %s at pos %d", expected[i].SValue, tkn.SValue, i)
		}
		if tkn.FValue != expected[i].FValue {
			t.Errorf("Expected %f, got %f at pos %d", expected[i].FValue, tkn.FValue, i)
		}
	}
}
