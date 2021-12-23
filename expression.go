package calc_engnieer

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"unicode"
)

type Expression struct {
	tokens []*Token
}

// NewExpressFromStr 从字符串中创建表达式
func NewExpressFromStr(str string) (*Expression, error) {
	if len(str) <= 0 {
		return nil, IllegalExpressionStr
	}
	s := newStream(str)
	tokenList := make([]*Token, 0)
	for s.canRead() {
		token, has, err := splitToken(s)
		if err != nil {
			return nil, err
		}
		if !has {
			break
		}
		if token != nil {
			tokenList = append(tokenList, token)
		}
	}
	e := &Expression{}
	e.tokens = tokenList
	// 检查 括号
	err := blankCheck(e)
	if err != nil {
		return nil, err
	}
	// 合并数值
	mergeNumber(e)
	return e, nil
}

func (e *Expression) Execute(paramMap map[string]interface{}) (interface{}, error) {
	ts := newTokenStream(e.tokens)
	rts, err := execute(ts, paramMap)
	if err != nil {
		return nil, err
	}
	if rts.canRead() {
		return rts.read().Value, nil
	}
	return nil, nil
}

// Execute 执行计算
// 遇到左括号进栈
// 遇到右括号出栈
func execute(ts *tokenStream, paramMap map[string]interface{}) (*tokenStream, error) {
	tempTokenList := make([]*Token, 0)
	for ts.canRead() {
		token := ts.read()
		// 遇到 (
		if token.TokenType == ClauseOpen {
			ts.back(1)
			// 读出下一个完整的()
			tempStream, start, end := ts.realBlankArea()
			nts, err := execute(tempStream, paramMap)
			if err != nil {
				return nil, err
			}
			if len(tempTokenList) > 0 {
				lastToken := tempTokenList[len(tempTokenList)-1]
				if lastToken.TokenType == Function {
					ts.back(1)
					// 如果这是一个公式的()
					formulaResult, formulaErr := calcFormula(lastToken, nts)
					if formulaErr != nil {
						return nil, formulaErr
					}
					ts.replace([]*Token{formulaResult}, start-1, end)
					tempTokenList = tempTokenList[0 : len(tempTokenList)-1]
				}
				continue
			}
			ts.replace(nts.tokenList, start, end)
			continue
		}
		// 遇到变量
		if token.TokenType == Variable {
			val := paramMap[token.Value.(string)]
			if val == nil {
				token.TokenType = Null
				continue
			}
			tp := reflect.TypeOf(val)
			tpKind := tp.Kind()
			if tpKind == reflect.Int {
				token.TokenType = Number
				val = float64(val.(int))
			} else if tpKind == reflect.Float64 {
				token.TokenType = Number
			} else {
				token.TokenType = String
			}
			token.Value = val
			tempTokenList = append(tempTokenList, token)
			continue
		}
		// 遇到字符串
		if token.TokenType == String {
			tempTokenList = append(tempTokenList, token)
		}
		// 遇到比较符
		if token.TokenType == Compare {
			tempTokenList = append(tempTokenList, token)
		}
		// 遇到运算符
		if token.TokenType == Operator {
			tempTokenList = append(tempTokenList, token)
		}
		// 遇到分隔符
		if token.TokenType == Separator {
			tempTokenList = append(tempTokenList, token)
		}
		if token.TokenType == Number {
			tempTokenList = append(tempTokenList, token)
		}
		if token.TokenType == Function {
			tempTokenList = append(tempTokenList, token)
		}
		if token.TokenType == Bool {
			tempTokenList = append(tempTokenList, token)
		}
	}
	newTempList := make([]*Token, 0)
	resultList := make([]*Token, 0)
	for _, item := range tempTokenList {
		if item.TokenType != Separator {
			newTempList = append(newTempList, item)
		} else {
			resultToken, err := calcTokens(newTempList)
			if err != nil {
				return nil, err
			}
			resultList = append(resultList, resultToken)
			resultList = append(resultList, item)
		}
	}
	if len(newTempList) > 0 {
		resultToken, err := calcTokens(newTempList)
		if err != nil {
			return nil, err
		}
		resultList = append(resultList, resultToken)
	}
	return newTokenStream(resultList), nil
}

// blankCheck 括号检查
func blankCheck(e *Expression) error {
	index := 0
	for _, t := range e.tokens {
		if t.TokenType == ClauseOpen {
			index++
		}
		if t.TokenType == ClauseClose {
			index--
		}
		if index < 0 {
			return HasIllegalClause
		}
	}
	if index != 0 {
		return HasIllegalClause
	}
	return nil
}

// mergeNumber 合并数值
func mergeNumber(e *Expression) {
	i := 0
	for index, t := range e.tokens {
		if t.TokenType == Operator {
			i++
			if fmt.Sprintf("%v", t.Value) == fmt.Sprintf("%c", Sub) && i == 2 {
				if index+1 < len(e.tokens) {
					t2 := e.tokens[index+1]
					if t2.TokenType == Number {
						t2.Value = t2.Value.(float64) * -1
						e.tokens = append(e.tokens[0:index], e.tokens[index+1:]...)
						mergeNumber(e)
					}
					continue
				}
				break
			}
		} else {
			i = 0
		}
	}
}

// splitToken 切分token
func splitToken(s *stream) (*Token, bool, error) {
	t := &Token{}
	for s.canRead() {
		c := s.read()
		// 如果是空串 就跳过
		if unicode.IsSpace(c) {
			continue
		}
		// 如果是数值类型的
		if isNumericStrict(c) {
			valStr := splitTokenTilFalse(s, isNumeric)
			floatNumber, transError := strconv.ParseFloat(valStr, 64)
			if transError != nil {
				return t, false, errors.New(fmt.Sprintf("%s can't trans into float", valStr))
			}
			t.TokenType = Number
			t.Value = floatNumber
			break
		}
		// 如果是变量类型
		if c == LeftSquare {
			valStr := ""
			hasEnd := false
			for s.canRead() {
				c2 := s.read()
				if c2 == RightSquare {
					hasEnd = true
					break
				}
				valStr = fmt.Sprintf("%s%c", valStr, c2)
			}
			if !hasEnd {
				return t, false, errors.New(fmt.Sprintf("variable %s don't have end board", valStr[1:]))
			}
			t.TokenType = Variable
			t.Value = valStr
			break
		}
		// 如果是字符串
		if c == Quote {
			valStr := fmt.Sprintf("%c", c)
			hasEnd := false
			for s.canRead() {
				c2 := s.read()
				valStr = fmt.Sprintf("%s%c", valStr, c2)
				if c2 == Quote {
					hasEnd = true
					break
				}
			}
			if !hasEnd {
				return t, false, errors.New(fmt.Sprintf("string %s don't have end board", valStr[1:]))
			}
			t.TokenType = String
			t.Value = valStr
			break
		}
		// 如果是运算符
		if isOperator(c) {
			t.TokenType = Operator
			t.Value = fmt.Sprintf("%c", c)
			break
		}
		// 如果是比较符号
		if isCompare(c) {
			valStr := splitTokenTilFalse(s, isCompare)
			t.TokenType = Compare
			t.Value = valStr
			break
		}
		// 如果是其他字符 按照公式处理
		if isLetter(c) {
			valStr := splitTokenTilFalse(s, isLetter)
			t.TokenType = Function
			t.Value = valStr
			break
		}
		// 左括号
		if c == LeftBlank {
			t.TokenType = ClauseOpen
			t.Value = fmt.Sprintf("%c", c)
			break
		}
		// 右括号
		if c == RightBlank {
			t.TokenType = ClauseClose
			t.Value = fmt.Sprintf("%c", c)
			break
		}
		// 分隔符号
		if c == Comma {
			t.TokenType = Separator
			t.Value = fmt.Sprintf("%c", c)
			break
		}
		return t, false, nil
	}
	return t, true, nil
}

// splitTokenTilFalse 读取token 直到不再符合条件
func splitTokenTilFalse(s *stream, cdtFunc func(c rune) bool) string {
	s.back(1)
	str := ""
	willBack := false
	for s.canRead() {
		c := s.read()
		if cdtFunc(c) == true {
			str = fmt.Sprintf("%s%c", str, c)
		} else {
			willBack = true
			break
		}
	}
	if willBack {
		s.back(1)
	}
	return str
}
