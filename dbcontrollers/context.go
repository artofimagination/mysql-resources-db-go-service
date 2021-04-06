package dbcontrollers

import (
	"fmt"

	"github.com/artofimagination/mysql-resources-db-go-service/mysqldb"
)

type MYSQLController struct {
	DBFunctions *mysqldb.MYSQLFunctions
	DBConnector *mysqldb.MYSQLConnector
}

func NewDBController(
	mySQLDBAddress string,
	mySQLDBPort int,
	mySQLDBUser,
	mySQLDBPassword,
	mySQLDBName,
	mySQLDBMigrationDirectory string) (*MYSQLController, error) {

	dbConnection := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true",
		mySQLDBUser,
		mySQLDBPassword,
		mySQLDBAddress,
		mySQLDBPort,
		mySQLDBName)

	dbConnector := &mysqldb.MYSQLConnector{
		DBConnection:       dbConnection,
		MigrationDirectory: mySQLDBMigrationDirectory,
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
