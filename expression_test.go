package calc_engnieer

import (
	"fmt"
	"testing"
)

func TestNewExpressFromStr(t *testing.T) {
	str := "AND({M4} = 20, OR(1 = 2, {L4} = 31))"
	paramMap := map[string]interface{}{
		"M4": 20,
		"L4": 30,
		"N4": 20,
	}
	expression, err := NewExpressFromStr(str)
	if err != nil {
		println(err.Error())
	}
	for _, t := range expression.tokens {
		fmt.Printf("\t%v", t.Value)
	}
	i, err := expression.Execute(paramMap)
	if err != nil {
		println(err.Error())
	}
	fmt.Printf("\n%v\n", i)
}
