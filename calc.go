package calc_engnieer

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
)

// calcTokens 计算token数组
// 这一步只有四则运算和逻辑运算
func calcTokens(tokenList []*Token) (*Token, error) {
	var lastToken *Token = nil
	oper := false
	comp := false
	for index, token := range tokenList {
		if token.TokenType == Operator {
			oper = true
			if lastToken == nil {
				return nil, errors.New(fmt.Sprintf("args missing ? %s", TokenTypeDict[token.TokenType]))
			}
			if lastToken.TokenType != Number {
				return nil, errors.New(fmt.Sprintf("[+-*/] need type of arg is number but actual is %s", TokenTypeDict[token.TokenType]))
			}
			if token.Value.(string) == fmt.Sprintf("%c", Mul) || token.Value.(string) == fmt.Sprintf("%c", Div) {
				if index == len(tokenList)-1 {
					return nil, errors.New(fmt.Sprintf("args missing ? %s", TokenTypeDict[token.TokenType]))
				}
				resultToken, err := calc(lastToken, token, tokenList[index+1])
				if err != nil {
					return nil, err
				}
				lastToken.Value = 0.0
				token.Value = "+"
				tokenList[index+1] = resultToken
			}
			lastToken = nil
		}
		if token.TokenType == Compare {
			comp = true
			if lastToken == nil {
				return nil, errors.New(fmt.Sprintf("args missing ? %s", TokenTypeDict[token.TokenType]))
			}
			if index == len(tokenList)-1 {
				return nil, errors.New(fmt.Sprintf("args missing ? %s", TokenTypeDict[token.TokenType]))
			}
			resultToken, err := calc(lastToken, token, tokenList[index+1])
			if err != nil {
				return nil, err
			}
			tokenList[index+1] = resultToken
		}
		lastToken = token
	}

	if oper && comp {
		return nil, errors.New("compare and operate can't calc at same time")
	}
	if oper {
		// 计算需要将+ - 处理结束
		for i := 1; i < len(tokenList); i += 2 {
			resultToken, err := calc(tokenList[i-1], tokenList[i], tokenList[i+1])
			if err != nil {
				return nil, err
			}
			tokenList[i+1] = resultToken
		}
	}
	return tokenList[len(tokenList)-1], nil
}

func calc(a, b, c *Token) (*Token, error) {
	result := &Token{}
	if b.TokenType == Operator {
		result.TokenType = Number
		number1 := decimal.NewFromFloat(a.Value.(float64))
		number2 := decimal.NewFromFloat(c.Value.(float64))
		switch fmt.Sprintf("%s", b.Value) {
		case "+":
			result.Value, _ = number1.Add(number2).Float64()
		case "-":
			result.Value, _ = number1.Sub(number2).Float64()
		case "*":
			result.Value, _ = number1.Mul(number2).Float64()
		case "/":
			result.Value, _ = number1.Div(number2).Float64()
		default:
			return nil, NonsupportCalc
		}
		return result, nil
	}
	if b.TokenType == Compare {
		result.TokenType = Bool
		if a.TokenType != Number || c.TokenType != Number {
			switch fmt.Sprintf("%v", b.Value) {
			case "=":
				result.Value = fmt.Sprintf("%v", a.Value) == fmt.Sprintf("%v", b.Value)
			case "!=":
				result.Value = fmt.Sprintf("%v", a.Value) == fmt.Sprintf("%v", b.Value)
			default:
				return nil, NonsupportCalc
			}
			return result, nil
		}
		number1 := decimal.NewFromFloat(a.Value.(float64))
		number2 := decimal.NewFromFloat(c.Value.(float64))
		switch fmt.Sprintf("%s", b.Value) {
		case ">":
			result.Value = number1.GreaterThan(number2)
		case ">=":
			result.Value = number1.GreaterThanOrEqual(number2)
		case "<":
			result.Value = number1.LessThan(number2)
		case "<=":
			result.Value = number1.LessThanOrEqual(number2)
		case "=":
			result.Value = number1.Equal(number2)
		case "!=":
			result.Value = !number1.Equal(number2)
		default:
			return nil, NonsupportCalc
		}
		return result, nil
	}
	return nil, NonsupportCalc
}
