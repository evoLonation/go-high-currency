package common

import (
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func NewMysqlDB(dataSource string) *sqlx.DB {
	db, err := sqlx.Open("mysql", dataSource)
	if err != nil {
		log.Fatal(errors.Wrap(err, "open db error"))
	}
	return db
}
