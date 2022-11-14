package db

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

func New(opts ...Option) *sql.DB {
	cfg := newConfig(opts...)

	conn, err := sql.Open("mysql", cfg.Source())
	if err != nil {
		log.Panicln("unable to open mysql connection:", err)
	}

	// ref: https://www.alexedwards.net/blog/configuring-sqldb
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxLifetime(time.Minute * 5)

	return conn
}
