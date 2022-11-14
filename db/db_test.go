package db

import (
	"context"
	"testing"
	"time"
)

func TestDatabaseConnection(t *testing.T) {

	conn := New()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err := conn.PingContext(ctx)
	if err != nil {
		t.Errorf("unable to ping mysql database: %s", err.Error())
	}

}
