package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"math/big"
	"os"
	"strings"
)

var (
	pt = fmt.Printf
)

func main() {
	// error handling
	var err error
	defer func() {
		if err != nil {
			log.Fatal(err)
		}
	}()
	defer ct(&err)

	// get lines
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		// parse
		expr, err := parser.ParseExpr(scanner.Text())
		ce(err, "parse: %s", line)

		// print
		value, err := eval(expr)
		ce(err, "evaluate: %s", line)
		if value.IsInt() {
			pt("%s\n", value.RatString())
		} else {
			pt("%s\n", strings.TrimRight(value.FloatString(1024), "0"))
		}
	}
}

func eval(expr ast.Expr) (*big.Rat, error) {
	switch expr := expr.(type) {
	// literal
	case *ast.BasicLit:
		v := new(big.Rat)
		v, ok := v.SetString(expr.Value)
		if !ok {
			return nil, me(nil, "invalid literal: %s", expr.Value)
		}
		return v, nil
	// binary expression
	case *ast.BinaryExpr:
		a, err := eval(expr.X)
		ce(err, "evaluate left operant")
		b, err := eval(expr.Y)
		ce(err, "evaluate right operant")
		v := new(big.Rat)
		switch expr.Op {
		// add
		case token.ADD:
			v.Add(a, b)
			return v, nil
		// sub
		case token.SUB:
			v.Sub(a, b)
			return v, nil
		// mul
		case token.MUL:
			v.Mul(a, b)
			return v, nil
		// quo
		case token.QUO:
			v.Quo(a, b)
			return v, nil
		default:
			return nil, me(nil, "unknown operator: %v", expr.Op)
		}
	// quotes
	case *ast.ParenExpr:
		return eval(expr.X)
	default:
		return nil, me(nil, "unknown expression type: %T", expr)
	}
	return nil, nil
}
