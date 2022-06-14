package database

import (
	"database/sql"
	"time"
)

type Cfg struct {
	URI     string
	Name    string
	Timeout time.Duration
}

func Open(cfg Cfg) (*sql.DB, error) {
	db, err := sql.Open("sqlserver", "sqlserver://administrator:buffaesbuffa@192.168.0.15:1433")
	if err != nil {
		return nil, err
	}
	return db, nil
}
