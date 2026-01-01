package postgres

import (
	"fmt"
	"strings"

	"github.com/zhiyunliu/glue/xdb"
)

var _ xdb.ExpressionBuildCallback = likeExpressBuildCallback

func likeExpressBuildCallback(item xdb.ExpressionValuer, state xdb.SqlState, param xdb.DBParam) (expression string, err xdb.MissError) {

	propName := item.GetPropName()
	value, err := param.GetVal(propName)
	if err != nil {
		//没有值，并且是可空
		if item.GetSymbol().IsDynamic() {
			return "", nil
		}
		return
	}
	err = nil
	if xdb.CheckIsNil(value) && item.GetSymbol().IsDynamic() {
		return
	}

	value = escapeLikeValue(item.GetOper(), value)

	phName := state.AppendExpr(propName, value)

	operCallback, ok := item.GetOperatorCallback()
	if !ok {
		err = xdb.NewMissOperError(item.GetOper())
		return
	}
	return operCallback(item, param, phName, value), nil
}

func escapeLikeValue(oper string, value any) string {
	val := fmt.Sprint(value)
	val = escapeLike(val)
	if strings.HasPrefix(oper, "%") {
		val = "%" + val
	}
	if strings.HasSuffix(oper, "%") {
		val = val + "%"
	}
	return val
}
func escapeLike(term string) string {
	// 重要：先转义反斜杠，再转义其他字符
	replacer := strings.NewReplacer(
		"\\", "\\\\", // 将 \ 替换为 \\
		"%", "\\%", // 将 % 替换为 \%
		"_", "\\_", // 将 _ 替换为 \_
	)
	return replacer.Replace(term)
}

func buildLikeOperators() []xdb.Operator {

	likecallback := func(item xdb.ExpressionValuer, param xdb.DBParam, phName string, _ any) string {
		return fmt.Sprintf("%s %s like %s", item.GetSymbol().Concat(), item.GetFullfield(), phName)
	}
	notlikecallback := func(item xdb.ExpressionValuer, param xdb.DBParam, phName string, _ any) string {
		return fmt.Sprintf("%s %s not like %s", item.GetSymbol().Concat(), item.GetFullfield(), phName)
	}

	operList := []xdb.Operator{
		xdb.NewOperator("like", likecallback),
		xdb.NewOperator("%like", likecallback),
		xdb.NewOperator("like%", likecallback),
		xdb.NewOperator("%like%", likecallback),

		xdb.NewOperator("notlike", notlikecallback),
		xdb.NewOperator("%notlike", notlikecallback),
		xdb.NewOperator("notlike%", notlikecallback),
		xdb.NewOperator("%notlike%", notlikecallback),
	}
	return operList
}
