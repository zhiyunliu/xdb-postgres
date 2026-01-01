package postgres

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/zhiyunliu/glue/config"
	contribxdb "github.com/zhiyunliu/glue/contrib/xdb"
	"github.com/zhiyunliu/glue/contrib/xdb/expression"
	"github.com/zhiyunliu/glue/contrib/xdb/tpl"
	"github.com/zhiyunliu/glue/xdb"
)

const (
	Proto          = "postgres"
	ArgumentPrefix = "$"
)

type postgresResolver struct {
	name string
}

func (s *postgresResolver) Name() string {
	return s.name
}

func (s *postgresResolver) Resolve(connName string, setting config.Config, opts ...xdb.Option) (interface{}, error) {
	cfg := contribxdb.NewConfig(connName)
	err := setting.ScanTo(cfg.Cfg)
	if err != nil {
		return nil, fmt.Errorf("读取DB配置(%s):%w", connName, err)
	}
	return contribxdb.NewDB(Proto, cfg, opts...)
}

func init() {
	symbols := expression.DefaultSymbols

	tplMatcher := xdb.NewTemplateMatcher(
		expression.NewNormalExpressionMatcher(symbols, xdb.WithBuildCallback(normalExpressBuildCallback)),
		expression.NewCompareExpressionMatcher(symbols),
		expression.NewLikeExpressionMatcher(symbols, xdb.WithBuildCallback(likeExpressBuildCallback), xdb.WithOperator(buildLikeOperators()...)),
		expression.NewInExpressionMatcher(symbols, xdb.WithBuildCallback(inExpressBuildCallback), xdb.WithOperator(buildInOperators()...)),
	)

	tplstmpProcessor := xdb.NewStmtDbTypeProcessor(DefaultDbTypeHandler...)

	xdb.Register(&postgresResolver{name: Proto})

	seqTemplate := tpl.NewSeq(Proto, ArgumentPrefix, tplMatcher, tplstmpProcessor)
	seqTemplate.StatePool = NewStatePool(seqTemplate.Placeholder())

	_ = xdb.RegistTemplate(seqTemplate)
}
