package calc_engnieer

type tokenStream struct {
	tokenList []*Token
	pos       int
	len       int
}

// newTokenStream 新建token数组
func newTokenStream(tokenList []*Token) *tokenStream {
	return &tokenStream{
		tokenList: tokenList,
		pos:       0,
		len:       len(tokenList),
	}
}

func (t *tokenStream) read() *Token {
	token := t.tokenList[t.pos]
	t.pos++
	return token
}

func (t *tokenStream) back(amount int) {
	t.pos -= amount
}

func (t *tokenStream) canRead() bool {
	return t.pos < t.len
}

func (t *tokenStream) formulaReturnReplace(pos int) {

}

// formulaArea 获取公式区域 a2(2,3)
func (t *tokenStream) formulaArea() *tokenStream {
	start := false
	index := 0
	for i := t.pos; i < t.len; i++ {
		if t.tokenList[i].TokenType == ClauseOpen {
			start = true
			index++
		}
		if t.tokenList[i].TokenType == ClauseClose {
			index--
		}
		if start && index == 0 {
			return newTokenStream(t.tokenList[t.pos : i+1])
		}
	}
	return t
}

// realBlankArea 读出一个完整的括号区域 (a + 2 + AND(3 * 2 ))
func (t *tokenStream) realBlankArea() (*tokenStream, int, int) {
	start := false
	index := 0
	for i := t.pos; i < t.len; i++ {
		if t.tokenList[i].TokenType == ClauseOpen {
			start = true
			index++
		}
		if t.tokenList[i].TokenType == ClauseClose {
			index--
		}
		if start && index == 0 {
			return newTokenStream(t.tokenList[t.pos+1 : i]), t.pos, i + 1
		}
	}
	return t, t.pos, t.len
}

func (t *tokenStream) replace(newTokenList []*Token, start, end int) {
	suffix := make([]*Token, 0)
	if end < t.len {
		suffix = t.tokenList[end:]
	}
	t.tokenList = append(t.tokenList[0:start], append(newTokenList, suffix...)...)
	t.len = len(t.tokenList)
}

func (t *tokenStream) reset() {
	t.pos = 1
}
