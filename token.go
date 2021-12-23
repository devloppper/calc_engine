package calc_engnieer

// Token 令牌
type Token struct {
	TokenType
	Value interface{}
}

type TokenType int

const (
	None     = iota // 无状态
	Number          // 数值
	String          // 字符串
	Variable        // 变量 A2 [只能是字母开头，且只能包含字母和数字]

	Function // 公式

	Separator   // 分隔符 ,
	ClauseOpen  // 左括号 (
	ClauseClose // 右括号 )

	Operator // 运算符
	Compare  // 比较符号
	Bool     // 布尔类型
	Null     // 空类型
)

var TokenTypeDict = map[TokenType]string{
	None:        "None",
	Number:      "Number",
	String:      "String",
	Variable:    "Variable",
	Function:    "Function",
	Separator:   "Separator",
	ClauseOpen:  "ClauseOpen",
	ClauseClose: "ClauseClose",
	Operator:    "Operator",
	Compare:     "Compare",
}
