package templates

import (
	"fmt"
	"github.com/lqiz/expr"
	"github.com/tianlin0/plat-lib/cond"
	"github.com/tianlin0/plat-lib/conv"
	"github.com/tidwall/gjson"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"strings"
)

func ruleExprMap(exprStr string, controlMap map[string]interface{}) (bool, error) {
	engine, err := expr.NewEngine(exprStr)
	if err != nil {
		return false, err
	}
	result, err := engine.RunRule(controlMap)
	if err != nil {
		return false, err
	}
	return result, nil
}

// 支持操作：>, <, ==, !=, &&, ||, >=, <=
// 支持括号嵌套
func ruleExprStr(exprStr string, json string) (bool, error) {
	exprAst, errParse := parser.ParseExpr(exprStr)
	if errParse != nil {
		return false, errParse
	}

	ret := runExpr(exprAst, json)
	if retBool, ok := ret.(bool); ok {
		return retBool, nil
	}

	return false, fmt.Errorf(exprStr + " error " + json)
}

// 执行表达式
// 支持操作：>, <, ==, !=, &&, ||, >=, <=
// 支持括号嵌套
func runExpr(expr ast.Expr, json string) interface{} {
	// 二元表达式
	if binaryExpr, ok := expr.(*ast.BinaryExpr); ok {
		opStr := strings.TrimSpace(binaryExpr.Op.String())
		x := runExpr(binaryExpr.X, json)
		y := runExpr(binaryExpr.Y, json)
		if cond.IsNil(x) || cond.IsNil(y) {
			return nil
		}
		{ //类型需要转换成y的类型 TODO

		}

		if opStr == "&&" || opStr == "||" {
			if opStr == "&&" {
				return x.(bool) && y.(bool)
			}
			return x.(bool) || y.(bool)
		}

		if opStr == "==" {
			return conv.String(x) == conv.String(y)
		}
		if opStr == "!=" {
			return conv.String(x) != conv.String(y)
		}
		xInt, _ := conv.Int64(x)
		yInt, _ := conv.Int64(y)

		if opStr == ">" || opStr == ">=" {
			if opStr == ">" {
				return xInt > yInt
			}
			return xInt >= yInt
		}
		if opStr == "<" || opStr == "<=" {
			if opStr == "<" {
				return xInt < yInt
			}
			return xInt <= yInt
		}
	}
	// 基本类型值
	if basicLit, ok := expr.(*ast.BasicLit); ok {
		switch basicLit.Kind {
		case token.INT:
			v, _ := strconv.Atoi(basicLit.Value)
			return v
		case token.FLOAT:
			v, _ := strconv.ParseFloat(basicLit.Value, 64)
			return v
		default:
			v := conv.String(basicLit.Value)
			return v
		}
	}
	// 标识符
	if ident, ok := expr.(*ast.Ident); ok {
		r := gjson.Get(json, ident.Name)
		if r.Exists() {
			return r.Raw
		}
		return ident.Name
	}
	// 括号表达式
	if parenExpr, ok := expr.(*ast.ParenExpr); ok {
		return runExpr(parenExpr.X, json)
	}

	if selectorExpr, ok := expr.(*ast.SelectorExpr); ok {
		list := make([]string, 0)
		getSelectorList(selectorExpr, &list)
		if len(list) > 0 {
			keyName := strings.Join(list, ".")
			ident := new(ast.Ident)
			ident.Name = keyName
			return runExpr(ident, json)
		}
	}

	return nil
}

// 根据.运算返回整个数组
func getSelectorList(selectorExpr *ast.SelectorExpr, list *[]string) {
	if one, ok := selectorExpr.X.(*ast.SelectorExpr); ok {
		getSelectorList(one, list)
		*list = append(*list, selectorExpr.Sel.Name)
		return
	}
	if one, ok := selectorExpr.X.(*ast.Ident); ok {
		*list = append(*list, one.Name)
		*list = append(*list, selectorExpr.Sel.Name)
		return
	}
}
