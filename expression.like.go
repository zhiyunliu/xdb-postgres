package postgres

import (
	"fmt"
	"strings"

	"github.com/zhiyunliu/glue/xdb"
)

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

	normalize := func(exprName xdb.ExprName, param xdb.DBParam, value any) (newVal any, err xdb.MissError) {
		return escapeLikeValue(exprName.GetOper(), value), nil
	}

	operList := []xdb.Operator{
		xdb.NewOperator("like", likecallback, normalize),
		xdb.NewOperator("%like", likecallback, normalize),
		xdb.NewOperator("like%", likecallback, normalize),
		xdb.NewOperator("%like%", likecallback, normalize),

		xdb.NewOperator("notlike", notlikecallback, normalize),
		xdb.NewOperator("%notlike", notlikecallback, normalize),
		xdb.NewOperator("notlike%", notlikecallback, normalize),
		xdb.NewOperator("%notlike%", notlikecallback, normalize),
	}
	return operList
}
