package postgres

import (
	"reflect"

	"github.com/lib/pq"
	"github.com/zhiyunliu/glue/xdb"
)

func checkValueIsArray(value any) bool {
	refVal := reflect.ValueOf(value)
	if !(refVal.Kind() == reflect.Array ||
		refVal.Kind() == reflect.Slice) {
		return false
	}
	return true
}

func buildNormalOperators(normalMatcher xdb.ExpressionMatcher) {
	operatorMap := normalMatcher.GetOperatorMap()
	normalize := func(exprName xdb.ExprName, param xdb.DBParam, value any) (newVal any, err xdb.MissError) {
		if !checkValueIsArray(value) {
			return value, nil
		}
		return pq.Array(value), nil
	}
	newOperList := []xdb.Operator{}

	operatorList := []string{"@", "&", "|"}
	for _, oper := range operatorList {
		operator, ok := operatorMap.Load(oper)
		if !ok {
			continue
		}
		newOperList = append(newOperList, xdb.NewOperator(oper, operator.Callback, normalize))
	}
	normalMatcher.GetOperatorMap().Store(newOperList...)
	return
}
