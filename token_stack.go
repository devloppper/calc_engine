package calc_engnieer

type tokenStack struct {
	tokenList []*Token
	pos       int
}

func newTokenStack() *tokenStack {
	return &tokenStack{
		tokenList: make([]*Token, 0),
		pos:       0,
	}
}

func (t *tokenStack) push(token *Token) {
	t.tokenList = append(t.tokenList, token)
	t.pos++
}

func (t *tokenStack) pop() *Token {
	if t.pos <= 0 {
		return nil
	}
	t.pos--
	tempToken := t.tokenList[t.pos]
	t.tokenList = t.tokenList[0 : t.pos+1]
	return tempToken
}
