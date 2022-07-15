package database

import (
	"database/sql"

	_ "github.com/denisenkom/go-mssqldb"
)

// sqlserver://cash:Mille.2021@192.168.1.10:1433?database=ADB_MILLEFRUTTISRL
// sqlserver://sa:recall@192.168.0.15:1433?database=ADB_DEMO

func Open() (*sql.DB, error) {
	db, err := sql.Open("sqlserver", "sqlserver://cash:Mille.2021@192.168.1.10:1433?database=ADB_MILLEFRUTTISRL")
	if err != nil {
		return nil, err
	}
	return db, nil
}
