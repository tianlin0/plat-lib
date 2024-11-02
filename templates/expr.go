package templates

import (
	"fmt"
	"github.com/Knetic/govaluate"
	"github.com/tianlin0/plat-lib/conv"
	"strings"
)

// RuleExpr 字符串规则引擎，也是模版的一种
func RuleExpr(exprStr string, data interface{}) (bool, error) {
	controlMap := make(map[string]interface{})
	_ = conv.Unmarshal(data, &controlMap)
	if len(controlMap) > 0 {
		// interface{} int string int64 float64 四种类型
		result, err := ruleExprMap(exprStr, controlMap)
		if err == nil {
			return result, nil
		}
	}

	// 如果有错误，则表示不支持格式 map[string]interface{}
	controlStr := conv.String(data)
	result, err := ruleExprStr(exprStr, controlStr)
	if err == nil {
		return result, nil
	}

	return false, err
}

// ExpressEvaluate 备注
func ExpressEvaluate(exprStr string, controlMap map[string]interface{},
	functions ...map[string]govaluate.ExpressionFunction) (interface{}, error) {
	if len(controlMap) == 0 || controlMap == nil {
		controlMap = nil
	}

	if len(functions) == 0 {
		expression, err := govaluate.NewEvaluableExpression(exprStr)
		if err != nil {
			return nil, err
		}
		return expression.Evaluate(controlMap)
	}

	expression, err := govaluate.NewEvaluableExpressionWithFunctions(exprStr, functions[0])
	if err != nil {
		return nil, err
	}
	return expression.Evaluate(controlMap)
}

// ExpressEvaluateFromToken 格式
func ExpressEvaluateFromToken(when string, controlMap ...map[string]interface{}) (interface{}, error) {
	var whenMap map[string]interface{}
	if len(controlMap) == 0 || controlMap[0] == nil {
		whenMap = nil
	} else {
		whenMap = controlMap[0]
	}

	expression, err := govaluate.NewEvaluableExpression(when)
	if err != nil {
		if strings.Contains(err.Error(), "Invalid token") {
			return nil, fmt.Errorf("invalid 'when' expression '%s': %v (hint: try wrapping the affected "+
				"expression in quotes (\"))", when, err)
		}
		return nil, fmt.Errorf("invalid 'when' expression '%s': %v", when, err)
	}
	tokens := expression.Tokens()
	for i, tok := range tokens {
		switch tok.Kind {
		case govaluate.VARIABLE:
			tok.Kind = govaluate.STRING
		default:
			continue
		}
		tokens[i] = tok
	}
	expression, err = govaluate.NewEvaluableExpressionFromTokens(tokens)
	if err != nil {
		return nil, fmt.Errorf("failed to parse 'when' expression '%s': %v", when, err)
	}
	result, err := expression.Evaluate(whenMap)
	if err != nil {
		return nil, fmt.Errorf("failed to evaluate 'when' expresion '%s': %v", when, err)
	}
	return result, nil
}
