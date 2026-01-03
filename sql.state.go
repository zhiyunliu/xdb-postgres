package postgres

import (
	"sync"

	"github.com/lib/pq"
	"github.com/zhiyunliu/glue/xdb"
)

var _ xdb.SqlState = &postgresSqlState{}

type postgresSqlState struct {
	xdb.SqlState
}

// 新建一个SqlState
func NewSqlState(placeHolder xdb.Placeholder) xdb.SqlState {
	return &postgresSqlState{
		SqlState: xdb.NewSqlState(placeHolder),
	}
}

func (s *postgresSqlState) BuildCache(sql string) xdb.ExpressionCache {
	return &postgresSqlTemplateCache{
		sql:   sql,
		names: s.GetNames(),
	}
}

var _ xdb.SqlStatePool = &postgresSqlStatePool{}

type postgresSqlStatePool struct {
	statePool *sync.Pool
}

func (sp *postgresSqlStatePool) Get() xdb.SqlState {
	return sp.statePool.Get().(xdb.SqlState)
}

func (sp *postgresSqlStatePool) Put(state xdb.SqlState) {
	sp.statePool.Put(state)
}

func NewStatePool(placeHolder xdb.Placeholder) xdb.SqlStatePool {
	return &postgresSqlStatePool{
		statePool: &sync.Pool{
			New: func() any {
				return NewSqlState(placeHolder)
			},
		},
	}
}

type postgresSqlTemplateCache struct {
	sql   string
	names []xdb.ExprName
}

func (stc *postgresSqlTemplateCache) Build(state xdb.SqlState, input xdb.DBParam) (sql string, err error) {
	for _, expr := range stc.names {
		value, err := input.GetVal(expr.GetPropName())
		if err != nil {
			return "", err
		}

		if !checkValueIsArray(value) {
			state.AppendExpr(expr, value)
		} else {
			state.AppendExpr(expr, pq.Array(value))
		}
	}
	return stc.sql, nil
}
