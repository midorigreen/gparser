package main

import (
	"bufio"
	"errors"
	"fmt"
	"go/ast"
	"go/constant"
	"go/parser"
	"io"
	"os"
)

var idents = []*ast.Ident{}

// copy from "http://qiita.com/tenntenn/items/a312d2c5381e36cf4cd3"
func main() {
	repl(os.Stdin)
}

func parse(str string) (string, error) {
	expr, err := parser.ParseExpr(str)
	if err != nil {
		return "", err
	}
	print(expr)

	v, err := evalExpr(expr)
	if err != nil {
		return "", err
	}
	return v.String(), nil
}

func evalBinary(expr *ast.BinaryExpr) (constant.Value, error) {
	x, err := evalExpr(expr.X)
	if err != nil {
		return constant.MakeUnknown(), errors.New("left operand faild")
	}
	y, err := evalExpr(expr.Y)
	if err != nil {
		return constant.MakeUnknown(), errors.New("right operand faild")
	}
	return constant.BinaryOp(x, expr.Op, y), nil
}

func evalUnary(expr *ast.UnaryExpr) (constant.Value, error) {
	x, err := evalExpr(expr.X)
	if err != nil {
		return constant.MakeUnknown(), err
	}

	return constant.UnaryOp(expr.Op, x, 0), nil
}

func evalIdent(e *ast.Ident) (constant.Value, error) {
	var u ast.Ident = ast.Ident{}
	var using *ast.Ident = &u
	for _, v := range idents {
		if v.Name == e.Name {
			using = v
			break
		}
	}

	if using.Name == "" {
		using = e
	}

	switch e.Obj.Kind {
	case ast.Var:
		fmt.Println(e)
	}

	return constant.MakeUnknown(), errors.New("Error: Idents")
}

func evalExpr(expr ast.Expr) (constant.Value, error) {
	switch e := expr.(type) {
	case *ast.ParenExpr:
		return evalExpr(e.X)
	case *ast.BinaryExpr:
		return evalBinary(e)
	case *ast.UnaryExpr:
		return evalUnary(e)
	case *ast.Ident:
		return evalIdent(e)
	case *ast.BasicLit:
		return constant.MakeFromLiteral(e.Value, e.Kind, 0), nil
	}
	return constant.MakeUnknown(), errors.New("unknown node")
}

func repl(r io.Reader) {
	s := bufio.NewScanner(r)
	for {
		fmt.Print(">")
		if !s.Scan() {
			break
		}

		l := s.Text()
		switch {
		case l == "exit":
			return
		default:
			r, err := parse(l)
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			fmt.Println(r)
		}
	}

	if err := s.Err(); err != nil {
		fmt.Println("Error:", err)
	}
}

// print is ast structure printing
func print(expr ast.Expr) {
	ast.Inspect(expr, func(n ast.Node) bool {
		if n != nil {
			fmt.Printf("%[1]v(%[1]T)\n", n)
		}
		return true
	})
}
