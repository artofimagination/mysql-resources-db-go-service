package mysqldb

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	// Need to register mysql drivers with database/sql
	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	migrate "github.com/rubenv/sql-migrate"
)

// MYSQL database connector implementation
type MYSQLConnector struct {
	DBConnection       string
	MigrationDirectory string
}

// MYSQLFunctions represents the implementation of MYSQL data manipulation functions.
type MYSQLFunctions struct {
	DBConnector *MYSQLConnector
}

func (*MYSQLConnector) Commit(tx *sql.Tx) error {
	return tx.Commit()
}

func (*MYSQLConnector) Rollback(tx *sql.Tx) error {
	return tx.Rollback()
}

func (c *MYSQLConnector) BootstrapSystem() error {
	fmt.Printf("Executing MYSQL migration\n")
	migrations := &migrate.FileMigrationSource{
		Dir: c.MigrationDirectory,
	}
	fmt.Printf("Getting migration files\n")

	db, err := sql.Open("mysql", c.DBConnection)
	if err != nil {
		return err
	}
	fmt.Printf("DB connection open\n")

	n := 0
	for retryCount := 20; retryCount > 0; retryCount-- {
		n, err = migrate.Exec(db, "mysql", migrations, migrate.Up)
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
		log.Printf("Failed to execute migration %s. Retrying...\n", err.Error())
	}

	if err != nil {
		return errors.Wrap(errors.WithStack(err), "Migration failed after multiple retries.")
	}
	fmt.Printf("Applied %d migrations!\n", n)
	return nil
}

func RollbackWithErrorStack(tx *sql.Tx, errorStack error) error {
	if err := tx.Rollback(); err != nil {
		errorString := fmt.Sprintf("%s\n%s\n", errorStack.Error(), err.Error())
		return errors.Wrap(errors.WithStack(errors.New(errorString)), "Failed to rollback changes")
	}
	return errorStack
}

func (c *MYSQLConnector) ConnectSystem() (*sql.Tx, error) {
	db, err := sql.Open("mysql", c.DBConnection)
	if err != nil {
		return nil, err
	}

	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	return tx, nil
}
