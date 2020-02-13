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

// Package GoMathExecutor provides simple expression executor.
//
// Usage:
//
// ```
// calc := executor.NewCalc()
// calc.AddOperators(executor.MathOperators) // Loads default MathOperators (see: defaults.go)
// calc.Prepare("2+2*2") // Prepare expression
// calc.Execute(nil) // == 6, nil
// calc.Prepare("x * (y+z)") // Prepare another expression with variables
// calc.Execute(map[string]float64{
// 	"x": 3,
// 	"y": 2,
// 	"z": 1,
// }) // == 9, nil
// ```

package executor // import "github.com/neonxp/GoMathExecutor"