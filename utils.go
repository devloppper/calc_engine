package calc_engnieer

import "unicode"

func isNumericStrict(character rune) bool {
	// 严格判断 开头的必须正式
	return unicode.IsDigit(character)
}

func isNumeric(character rune) bool {
	return unicode.IsDigit(character) || character == '.'
}

// isOperator 判断是否为运算符
func isOperator(c rune) bool {
	return c == Add || c == Sub || c == Mul || c == Div
}

// isCompare 判断是否为比较符
func isCompare(c rune) bool {
	return c == Gt || c == Lt || c == Eq || c == Not
}

// isLetter 判断是否是字符
func isLetter(c rune) bool {
	return unicode.IsLetter(c)
}
