# GoMathExecutor [![GoDoc](https://godoc.org/github.com/neonxp/GoMathExecutor?status.svg)](https://godoc.org/github.com/neonxp/GoMathExecutor) [![Build Status](https://travis-ci.org/neonxp/GoMathExecutor.svg?branch=master)](https://travis-ci.org/neonxp/GoMathExecutor) [![codecov](https://codecov.io/gh/neonxp/GoMathExecutor/branch/master/graph/badge.svg)](https://codecov.io/gh/neonxp/GoMathExecutor)

Package GoMathExecutor provides simple expression executor.

## Installation

`go get github.com/neonxp/GoMathExecutor`

## Usage

```
calc := executor.NewCalc()
calc.AddOperators(executor.MathOperators) // Loads default MathOperators (see: defaults.go)
calc.Prepare("2+2*2") // Prepare expression
calc.Execute(nil) // == 6, nil
calc.Prepare("x * (y+z)") // Prepare another expression with variables
calc.Execute(map[string]float64{
	"x": 3,
	"y": 2,
	"z": 1,
}) // == 9, nil
```
