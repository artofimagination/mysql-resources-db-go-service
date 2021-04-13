package http

import (
	"github.com/google/uuid"

	"github.com/artofimagination/mysql-resources-db-go-service/models"
)

type AddResourceRequest struct {
	Resource *models.Resource `json:"resource"`
}

// @deprecated
type GetResourceByIDWithQueryRequest struct {
	UUID uuid.UUID `query:"id" validate:"required"`
}

type GetResourceByIDRequest struct {
	UUID uuid.UUID `json:"resource_id" validate:"required,uuid"`
}

type GetResourcesByCategoryRequest struct {
	Category int `query:"category" validate:"required"`
}

type GetResourcesByIDsRequest struct {
	UUIDs []uuid.UUID `query:"ids" validate:"required"`
}

type DeleteResourceRequest struct {
	ID       uuid.UUID         `json:"id" validate:"required"`
	Category int               `json:"category"`
	Content  models.ContentMap `json:"content"`
}
