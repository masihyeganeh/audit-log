package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/pkg/errors"
	"time"
)

type repository struct {
	database *sql.DB
	isOpen   bool
}

var models = map[string]interface{}{}

func ConnectAndCreateRepository(ctx context.Context, connectionString string) (*repository, bool, error) {
	var db *sql.DB
	var err error

	// try connecting to datastore 3 times with sleep (docker-compose won't wait for db to be ready before starting us)
	for i := 0; i < 3; i++ {
		db, err = connect(connectionString)
		if err != nil {
			time.Sleep(3 * time.Second)
		} else {
			break
		}
	}

	if err != nil {
		return nil, false, err
	}

	repo := &repository{
		database: db,
		isOpen:   true,
	}

	err, isFirstRun := repo.setup(ctx)
	if err != nil {
		return nil, isFirstRun, errors.Wrap(err, "could not execute setup script of datastore")
	}

	return repo, isFirstRun, nil
}

func connect(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("clickhouse", connectionString)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			return nil, fmt.Errorf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			return nil, err
		}
	}

	return db, nil
}

func (r *repository) Close() error {
	if !r.isOpen {
		return nil
	}

	r.isOpen = false
	return r.database.Close()
}
