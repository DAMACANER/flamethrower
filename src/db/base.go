package db

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
)

type BaseRepo struct {
	QueryBuilder StmtBuilder
}

type Repo interface {
	ExtractVars() map[string]interface{}
}

type StmtBuilder interface {
	ToSql() (string, []interface{}, error)
	Limit(limit uint64) squirrel.SelectBuilder
	Offset(offset uint64) squirrel.SelectBuilder
}

func (b *BaseRepo) Query() (*sql.Rows, error) {
	sql, args, err := b.QueryBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	rows, err := DB.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (b *BaseRepo) Paginate(page, pageSize uint64) *BaseRepo {
	b.QueryBuilder.Limit(pageSize).Offset(page * pageSize)
	return b
}
