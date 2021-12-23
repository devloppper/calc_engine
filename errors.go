package calc_engnieer

import "errors"

var (
	IllegalExpressionStr = errors.New("empty expression str")
	HasIllegalClause     = errors.New("has illegal clause")
	NonsupportCalc       = errors.New("this is calc is not supported")
)
