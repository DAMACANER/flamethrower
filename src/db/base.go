package db

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
)

type BaseRepo struct {
	QueryBuilder  squirrel.SelectBuilder
	InsertBuilder squirrel.InsertBuilder
}

type Repo interface {
	ExtractVars() map[string]interface{}
}

type QueryBuilder interface {
	ToSql() (string, []interface{}, error)
	Limit(limit uint64) squirrel.SelectBuilder
	Offset(offset uint64) squirrel.SelectBuilder
}

type InsertBuilder interface {
	ToSql() (string, []interface{}, error)
	Columns(columns ...string) squirrel.InsertBuilder
	Values(values ...interface{}) squirrel.InsertBuilder
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
	b.QueryBuilder = squirrel.Select("*")
	return rows, nil
}

func (b *BaseRepo) Paginate(page, pageSize uint64) *BaseRepo {
	b.QueryBuilder = b.QueryBuilder.Limit(pageSize).Offset(page * pageSize)
	return b
}

// OrderBy
//
// Squirrel automatically orders in ascending order.
// Specify descending with colymn name or other fields like COLLATE NOCASE for lowercase.
func (b *BaseRepo) OrderBy(orderBys ...string) *BaseRepo {
	b.QueryBuilder = b.QueryBuilder.OrderBy(orderBys...)
	return b
}
