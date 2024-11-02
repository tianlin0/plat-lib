package internal

import "strings"

// SubStr 截取字符串，支持多字节字符
// start：起始下标，负数从从尾部开始，最后一个为-1
// length：截取长度，负数表示截取到末尾
func SubStr(str string, start int, length int) (result string) {
	s := []rune(str)
	total := len(s)
	if total == 0 {
		return
	}
	// 允许从尾部开始计算
	if start < 0 {
		start = total + start
		if start < 0 {
			return
		}
	}
	if start > total {
		return
	}
	// 到末尾
	if length < 0 {
		length = total
	}

	end := start + length
	if end > total {
		result = string(s[start:])
	} else {
		result = string(s[start:end])
	}
	return
}

// SnakeString XxYy to xx_yy , XxYY to xx_yy
func SnakeString(s string) string {
	//先将连续的大写转换为第一个大写，后面小写，如果后面大写还接小写，则将最近的小写转大写
	//userIDC => user_idc
	//userADDR rmb => user_addr_rmb
	//userAbc => user_abc

	num := len(s)
	newData := make([]string, 0, len(s))
	for i := 0; i < num; i++ {
		d := s[i]
		newData = append(newData, string(d))
		if d >= 'A' && d <= 'Z' {
			n := 0
			i = i + 1
			for ; i < num; i++ {
				d2 := s[i]
				if d2 >= 'A' && d2 <= 'Z' {
					n = n + 1
					newData = append(newData, strings.ToLower(string(d2)))
				} else {
					//表示是多个大写后的第一个小写
					if n > 0 {
						newData = append(newData, strings.ToUpper(string(d2)))
					} else {
						newData = append(newData, string(d2))
					}
					break
				}
			}
		}
	}
	s = strings.Join(newData, "")

	data := make([]byte, 0, len(s)*2)
	j := false
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}

	snakeStringTemp := strings.ToLower(string(data[:]))
	return snakeStringTemp
}

// PascalString xx_yy to XxYy
func PascalString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}
