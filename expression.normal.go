package postgres

import (
	"reflect"
	"strings"

	"github.com/lib/pq"
	"github.com/zhiyunliu/glue/xdb"
)

var _ xdb.ExpressionBuildCallback = normalExpressBuildCallback

func normalExpressBuildCallback(item xdb.ExpressionValuer, state xdb.SqlState, param xdb.DBParam) (expression string, err xdb.MissError) {
	var (
		phName   string
		propName = item.GetPropName()
	)
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

	if !strings.EqualFold(item.GetSymbol().Name(), xdb.SymbolReplace) {
		if !checkValueIsArray(value) {
			phName = state.AppendExpr(propName, value)
		} else {
			phName = state.AppendExpr(propName, pq.Array(value))
		}
	}
	operCallback, ok := item.GetOperatorCallback()
	if !ok {
		err = xdb.NewMissOperError(item.GetOper())
		return
	}
	return operCallback(item, param, phName, value), nil
}

func checkValueIsArray(value any) bool {
	refVal := reflect.ValueOf(value)
	if !(refVal.Kind() == reflect.Array ||
		refVal.Kind() == reflect.Slice) {
		return false
	}
	return true
}
