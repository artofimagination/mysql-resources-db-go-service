package dbcontrollers

import (
	"errors"
	"fmt"
	"os"

	"github.com/artofimagination/mysql-resources-db-go-service/mysqldb"
)

type MYSQLController struct {
	DBFunctions *mysqldb.MYSQLFunctions
	DBConnector *mysqldb.MYSQLConnector
}

func NewDBController() (*MYSQLController, error) {
	address := os.Getenv("MYSQL_DB_ADDRESS")
	if address == "" {
		return nil, errors.New("MYSQL DB address not defined")
	}
	port := os.Getenv("MYSQL_DB_PORT")
	if address == "" {
		return nil, errors.New("MYSQL DB port not defined")
	}
	username := os.Getenv("MYSQL_DB_USER")
	if address == "" {
		return nil, errors.New("MYSQL DB username not defined")
	}
	pass := os.Getenv("MYSQL_DB_PASSWORD")
	if address == "" {
		return nil, errors.New("MYSQL DB password not defined")
	}
	dbName := os.Getenv("MYSQL_DB_NAME")
	if address == "" {
		return nil, errors.New("MYSQL DB name not defined")
	}

	migrationDirectory := os.Getenv("MYSQL_DB_MIGRATION_DIR")
	if migrationDirectory == "" {
		return nil, errors.New("MYSQL DB migration folder not defined")
	}

	dbConnection := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		username,
		pass,
		address,
		port,
		dbName)

	dbConnector := &mysqldb.MYSQLConnector{
		DBConnection:       dbConnection,
		MigrationDirectory: migrationDirectory,
	}

	controller := &MYSQLController{
		DBFunctions: &mysqldb.MYSQLFunctions{
			DBConnector: dbConnector,
		},
		DBConnector: dbConnector,
	}

	if err := controller.DBConnector.BootstrapSystem(); err != nil {
		return nil, err
	}

	return controller, nil
}
