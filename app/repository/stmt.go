package repository

import (
	"context"
	"database/sql"
	"log"
)

func mustCreateStmt(ctx context.Context, conn *sql.DB, query string) *sql.Stmt {
	stmt, err := conn.PrepareContext(ctx, query)
	if err != nil {
		log.Panicln("unable to create sql prepared statement:", err.Error())
	}
	return stmt
}
