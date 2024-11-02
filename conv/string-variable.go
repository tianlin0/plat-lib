package conv

import (
	"github.com/tianlin0/plat-lib/internal"
	"strings"
)

// ChangeVariableName 将驼峰与小写互转
func ChangeVariableName(varName string, toType ...string) string {
	if varName == "" {
		return ""
	}
	typeList := []string{"snake", "camel", "pascal"}
	if len(toType) == 1 {
		typeName := toType[0]
		if typeName == typeList[0] {
			return internal.SnakeString(varName)
		}
		if typeName == typeList[1] {
			return camelString(varName)
		}
		if typeName == typeList[2] {
			return internal.PascalString(varName)
		}
		return varName
	}

	//检测是否都为小写
	isLower := true
	for i := 0; i < len(varName); i++ {
		c := varName[i]
		if isASCIIUpper(c) {
			isLower = false
			break
		}
	}

	if !isLower {
		return internal.SnakeString(varName)
	}
	return internal.PascalString(varName)
}

func isASCIIUpper(c byte) bool {
	return 'A' <= c && c <= 'Z'
}

func camelString(s string) string {
	newStr := internal.PascalString(s)
	if newStr == "" {
		return ""
	}
	first := string(newStr[0])
	firstName := strings.ToLower(first)
	return firstName + newStr[1:]
}
