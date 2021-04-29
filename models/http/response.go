package http

import "github.com/artofimagination/mysql-resources-db-go-service/models"

type CategoriesResponse struct {
	Categories []models.Category `json:"categories"`
}

type ResourcesResponse struct {
	Resources []models.Resource `json:"resources"`
}

type ResponseData struct {
	Error string      `json:"error" validation:"required"`
	Data  interface{} `json:"data" validation:"required"`
}
