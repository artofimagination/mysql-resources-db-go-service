package service

import "github.com/artofimagination/mysql-resources-db-go-service/dbcontrollers"

type Service struct {
	dbController *dbcontrollers.MYSQLController
}

func NewService(dbController *dbcontrollers.MYSQLController) *Service {
	return &Service{
		dbController: dbController,
	}
}
