package http

import (
	"github.com/google/uuid"

	"github.com/artofimagination/mysql-resources-db-go-service/models"
)

type AddResourceRequest struct {
	Resource *models.Resource `json:"resource"`
}

type GetResourceByIDRequest struct {
	UUID uuid.UUID `query:"id" json:"resource_id" validate:"required,uuid"`
}

type GetResourcesByCategoryRequest struct {
	Category int `query:"category" validate:"required"`
}

type GetResourcesByIDsRequest struct {
	UUIDs []uuid.UUID `query:"ids" validate:"required,dive,uuid"`
}
