package calc_engnieer

import (
	"errors"
	"fmt"
)

// calcFormula 计算公式
func calcFormula(formulaToken *Token, ts *tokenStream) (*Token, error) {
	formulaName := formulaToken.Value.(string)
	formulaFunc := formulaDict[formulaName]
	if formulaFunc == nil {
		return nil, errors.New(fmt.Sprintf("unsupport formula :%s", formulaName))
	}
	return formulaFunc(ts)
}

const (
	and = "AND"
	or  = "OR"
)

var formulaDict = map[string]func(ts *tokenStream) (*Token, error){
	and: fAnd,
	or:  fOr,
}

// fAnd AND函数
func fAnd(ts *tokenStream) (*Token, error) {
	// 忽略 (
	mustBeBool := true
	boolList := make([]bool, 0)
	for ts.canRead() {
		token := ts.read()
		if mustBeBool && token.TokenType != Bool {
			return nil, errors.New(fmt.Sprintf("Formula Add need args (Bool,Bool)"))
		}
		if token.TokenType == Separator {
			mustBeBool = true
			continue
		}
		if token.TokenType == Bool {
			boolList = append(boolList, token.Value.(bool))
			mustBeBool = false
		}
		if token.TokenType == ClauseClose {
			break
		}
	}
	if len(boolList) <= 0 {
		return nil, errors.New(fmt.Sprintf("Formula Add must have Bool args but it is empty"))
	}
	result := &Token{
		TokenType: Bool,
		Value:     true,
	}
	for _, item := range boolList {
		if item == false {
			result.Value = false
		}
	}
	// 所有的都必须是true才能通过
	return result, nil
}

// fOr Or函数
func fOr(ts *tokenStream) (*Token, error) {
	// 忽略 (
	mustBeBool := true
	boolList := make([]bool, 0)
	for ts.canRead() {
		token := ts.read()
		if mustBeBool && token.TokenType != Bool {
			return nil, errors.New(fmt.Sprintf("Formula Or need args (Bool,Bool)"))
		}
		if token.TokenType == Separator {
			mustBeBool = true
			continue
		}
		if token.TokenType == Bool {
			boolList = append(boolList, token.Value.(bool))
			mustBeBool = false
		}
		if token.TokenType == ClauseClose {
			break
		}
	}
	if len(boolList) <= 0 {
		return nil, errors.New(fmt.Sprintf("Formula Or must have Bool args but it is empty"))
	}
	result := &Token{
		TokenType: Bool,
		Value:     false,
	}
	// 只要有一个是true就行
	for _, item := range boolList {
		if item == true {
			result.Value = true
			break
		}
	}
	return result, nil
}
