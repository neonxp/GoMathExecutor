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
	"fmt"
	"strconv"
)

type tokenizer struct {
	str           string
	numberBuffer  string
	strBuffer     string
	allowNegative bool
	tkns          []*token
	operators     map[string]*Operator
}

func newTokenizer(str string, operators map[string]*Operator) *tokenizer {
	return &tokenizer{str: str, numberBuffer: "", strBuffer: "", allowNegative: true, tkns: []*token{}, operators: operators}
}

func (t *tokenizer) emptyNumberBufferAsLiteral() error {
	if t.numberBuffer != "" {
		f, err := strconv.ParseFloat(t.numberBuffer, 64)
		if err != nil {
			return fmt.Errorf("invalid number %s", t.numberBuffer)
		}
		t.tkns = append(t.tkns, newToken(literalType, "", f))
	}
	t.numberBuffer = ""
	return nil
}

func (t *tokenizer) emptyStrBufferAsVariable() {
	if t.strBuffer != "" {
		t.tkns = append(t.tkns, newToken(variableType, t.strBuffer, 0))
		t.strBuffer = ""
	}
}

func (t *tokenizer) tokenize() error {
	for _, ch := range t.str {
		if ch == ' ' {
			continue
		}
		ch := byte(ch)
		switch true {
		case isAlpha(ch):
			if t.numberBuffer != "" {
				if err := t.emptyNumberBufferAsLiteral(); err != nil {
					return err
				}
				t.tkns = append(t.tkns, newToken(operatorType, "*", 0))
				t.numberBuffer = ""
			}
			t.allowNegative = false
			t.strBuffer += string(ch)
		case isNumber(ch):
			t.numberBuffer += string(ch)
			t.allowNegative = false
		case isDot(ch):
			t.numberBuffer += string(ch)
			t.allowNegative = false
		case isLP(ch):
			if t.strBuffer != "" {
				t.tkns = append(t.tkns, newToken(functionType, t.strBuffer, 0))
				t.strBuffer = ""
			} else if t.numberBuffer != "" {
				if err := t.emptyNumberBufferAsLiteral(); err != nil {
					return err
				}
				t.tkns = append(t.tkns, newToken(operatorType, "*", 0))
				t.numberBuffer = ""
			}
			t.allowNegative = true
			t.tkns = append(t.tkns, newToken(leftParenthesisType, "", 0))
		case isRP(ch):
			if err := t.emptyNumberBufferAsLiteral(); err != nil {
				return err
			}
			t.emptyStrBufferAsVariable()
			t.allowNegative = false
			t.tkns = append(t.tkns, newToken(rightParenthesisType, "", 0))
		case isComma(ch):
			if err := t.emptyNumberBufferAsLiteral(); err != nil {
				return err
			}
			t.emptyStrBufferAsVariable()
			t.tkns = append(t.tkns, newToken(funcSep, "", 0))
			t.allowNegative = true
		default:
			if t.allowNegative && ch == '-' {
				t.numberBuffer += "-"
				t.allowNegative = false
				continue
			}
			if err := t.emptyNumberBufferAsLiteral(); err != nil {
				return err
			}
			t.emptyStrBufferAsVariable()
			if len(t.tkns) > 0 && t.tkns[len(t.tkns)-1].Type == operatorType {
				t.tkns[len(t.tkns)-1].SValue += string(ch)
			} else {
				t.tkns = append(t.tkns, newToken(operatorType, string(ch), 0))
			}
			t.allowNegative = true
		}
	}
	if err := t.emptyNumberBufferAsLiteral(); err != nil {
		return err
	}
	t.emptyStrBufferAsVariable()
	return nil
}

func (t *tokenizer) toRPN() ([]*token, error) {
	var tkns []*token
	var stack tokenStack
	for _, tkn := range t.tkns {
		switch tkn.Type {
		case literalType:
			tkns = append(tkns, tkn)
		case variableType:
			tkns = append(tkns, tkn)
		case functionType:
			stack.Push(tkn)
		case funcSep:
			for stack.Head().Type != leftParenthesisType {
				if stack.Head().Type == eof {
					return nil, ErrInvalidExpression
				}
				tkns = append(tkns, stack.Pop())
			}
		case operatorType:
			leftOp, ok := t.operators[tkn.SValue]
			if !ok {
				return nil, fmt.Errorf("unknown operator: %s", tkn.SValue)
			}
			for {
				if stack.Head().Type == operatorType {
					rightOp, ok := t.operators[stack.Head().SValue]
					if !ok {
						return nil, fmt.Errorf("unknown operator: %s", stack.Head().SValue)
					}
					if leftOp.Priority < rightOp.Priority || (leftOp.Priority == rightOp.Priority && leftOp.Assoc == RightAssoc) {
						tkns = append(tkns, stack.Pop())
						continue
					}
				}
				break
			}
			stack.Push(tkn)
		case leftParenthesisType:
			stack.Push(tkn)
		case rightParenthesisType:
			for stack.Head().Type != leftParenthesisType {
				if stack.Head().Type == eof {
					return nil, ErrInvalidParenthesis
				}
				tkns = append(tkns, stack.Pop())
			}
			stack.Pop()
			if stack.Head().Type == functionType {
				tkns = append(tkns, stack.Pop())
			}
		}
	}
	for stack.Head().Type != eof {
		if stack.Head().Type == leftParenthesisType {
			return nil, ErrInvalidParenthesis
		}
		tkns = append(tkns, stack.Pop())
	}
	return tkns, nil
}

type tokenStack struct {
	ts []*token
}

func (ts *tokenStack) Push(t *token) {
	ts.ts = append(ts.ts, t)
}

func (ts *tokenStack) Pop() *token {
	if len(ts.ts) == 0 {
		return &token{Type: eof}
	}
	var head *token
	head, ts.ts = ts.ts[len(ts.ts)-1], ts.ts[:len(ts.ts)-1]
	return head
}

func (ts *tokenStack) Head() *token {
	if len(ts.ts) == 0 {
		return &token{Type: eof}
	}
	return ts.ts[len(ts.ts)-1]
}

func isComma(ch byte) bool {
	return ch == ','
}

func isDot(ch byte) bool {
	return ch == '.'
}

func isNumber(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isAlpha(ch byte) bool {
	return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' || ch == '_'
}

func isLP(ch byte) bool {
	return ch == '('
}

func isRP(ch byte) bool {
	return ch == ')'
}
