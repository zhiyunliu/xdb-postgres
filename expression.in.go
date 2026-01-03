package postgres

import (
	"fmt"
	"reflect"

	"github.com/lib/pq"
	"github.com/zhiyunliu/glue/xdb"
)

var _ xdb.ExpressionBuildCallback = inExpressBuildCallback

func inExpressBuildCallback(item xdb.ExpressionValuer, state xdb.SqlState, param xdb.DBParam) (expression string, err xdb.MissError) {
	value, err := param.GetVal(item.GetPropName())
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

	normalizeCall, ok := item.GetOperValueNormalizeCallback()
	if !ok {
		err = xdb.NewMissOperError(item.GetOper())
		return
	}
	value, err = normalizeCall(item, param, value)
	if err != nil {
		return
	}

	phName := state.AppendExpr(item, value)
	operCallback, ok := item.GetOperExprCallback()
	if !ok {
		err = xdb.NewMissOperError(item.GetOper())
		return
	}
	return operCallback(item, param, phName, value), nil
}

func buildInOperators() []xdb.Operator {

	inCallback := func(item xdb.ExpressionValuer, param xdb.DBParam, phName string, value any) string {
		return fmt.Sprintf("%s %s = any(%s)", item.GetSymbol().Concat(), item.GetFullfield(), phName)
	}
	notinCallback := func(item xdb.ExpressionValuer, param xdb.DBParam, phName string, value any) string {
		return fmt.Sprintf("%s %s != all(%s)", item.GetSymbol().Concat(), item.GetFullfield(), phName)
	}
	normalize := func(exprName xdb.ExprName, param xdb.DBParam, value any) (newVal any, err xdb.MissError) {
		var val any
		switch t := value.(type) {
		case []int8, []int, []int16, []int32, []int64, []uint, []uint16, []uint32, []uint64:
			length := reflect.ValueOf(t).Len()
			if length == 0 {
				return
			}
			val = pq.Array(t)
		case []string:
			if len(t) <= 0 {
				return
			}
			val = pq.StringArray(t)
		case []byte:
			return "", xdb.NewMissDataTypeError(exprName.GetPropName())
		default:
			refVal := reflect.ValueOf(value)
			if !(refVal.Kind() == reflect.Array ||
				refVal.Kind() == reflect.Slice) {
				return "", xdb.NewMissDataTypeError(exprName.GetPropName())
			}
			arrayLen := refVal.Len()
			if arrayLen <= 0 {
				return
			}
			tmpStrArray := make([]string, arrayLen)
			for i := 0; i < arrayLen; i++ {
				ele := refVal.Index(i)
				tmpStrArray[i] = fmt.Sprint(ele.Interface())
			}
			val = pq.StringArray(tmpStrArray)
		}
		return val, nil
	}

	operList := []xdb.Operator{
		xdb.NewOperator("in", inCallback, normalize),
		xdb.NewOperator("notin", notinCallback, normalize),
	}

	return operList
}
