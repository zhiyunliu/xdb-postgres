package postgres

import (
	"github.com/zhiyunliu/glue/xdb"
)

var (
	DefaultDbTypeHandler = []xdb.StmtDbTypeHandler{
		// &varcharHandler{},
		// &varcharMaxHandler{},
		// &nvarcharMaxHandler{},
		// &tvpHandler{},
		// &contribxdb.StmtDbTypeOutputHandler{},
	}
)
