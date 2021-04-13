package storage

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/proemergotech/log/v3"
	migrate "github.com/rubenv/sql-migrate"
)

type MySQL struct {
	db *sqlx.DB
}

func NewMySQL(db *sqlx.DB) *MySQL {
	return &MySQL{
		db: db,
	}
}

func (mySQL *MySQL) BootstrapSystem(migrationDirectory string) error {
	fmt.Printf("Executing MYSQL migration\n")
	migrations := &migrate.FileMigrationSource{
		Dir: migrationDirectory,
	}
	fmt.Printf("Getting migration files\n")

	var n int
	var err error
	for retryCount := 20; retryCount > 0; retryCount-- {
		n, err = migrate.Exec(mySQL.db.DB, "mysql", migrations, migrate.Up)
		if err == nil {
			break
		}
		time.Sleep(time.Second)
		log.Debug(context.Background(), "Failed to execute migration %s. Retrying...\n", "err", err.Error())
	}

	if err != nil {
		return errors.Wrap(errors.WithStack(err), "Migration failed after multiple retries.")
	}
	fmt.Printf("Applied %d migration(s)!\n", n)
	return nil
}
