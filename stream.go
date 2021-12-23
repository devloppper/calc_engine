package calc_engnieer

type stream struct {
	charList []rune
	pos      int
	len      int
}

func newStream(str string) *stream {
	charList := make([]rune, 0)
	for _, c := range str {
		charList = append(charList, c)
	}
	return &stream{
		charList: charList,
		pos:      0,
		len:      len(charList),
	}
}

func (s *stream) read() rune {
	c := s.charList[s.pos]
	s.pos++
	return c
}

func (s *stream) canRead() bool {
	return s.pos < s.len
}

func (s *stream) back(amount int) {
	s.pos -= amount
}
