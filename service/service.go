package service

import (
	"github.com/artofimagination/mysql-resources-db-go-service/storage"
)

type Service struct {
	mySQLStorage *storage.MySQL
}

func NewService(mySQLStorage *storage.MySQL) *Service {
	return &Service{
		mySQLStorage: mySQLStorage,
	}
}
