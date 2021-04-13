package service

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/proemergotech/log/v3"

	"github.com/artofimagination/mysql-resources-db-go-service/models"
	httpModels "github.com/artofimagination/mysql-resources-db-go-service/models/http"
	"github.com/artofimagination/mysql-resources-db-go-service/models/myerrors"
	"github.com/artofimagination/mysql-resources-db-go-service/storage"
)

func (s *Service) AddResource(ctx context.Context, resource *models.Resource) error {
	// todo: add resource should be return with the added resource
	log.Debug(ctx, "Add resource")
	// Execute function
	if err := s.mySQLStorage.AddResource(resource); err != nil {
		return myerrors.WithFields(errors.Wrap(err, "mysql error"), models.HTTPCode, http.StatusInternalServerError)
	}

	return nil
}

func (s *Service) GetResourceByID(ctx context.Context, req *httpModels.GetResourceByIDRequest) (*models.Resource, error) {
	log.Debug(ctx, "Getting resource by id")

	resource, err := s.mySQLStorage.GetResourceByID(req.UUID)
	if err != nil {
		if err.Error() == storage.ErrResourceNotFound.Error() {
			return nil, myerrors.WithFields(err, models.HTTPCode, http.StatusAccepted)
		}
		return nil, myerrors.WithFields(err, models.HTTPCode, http.StatusInternalServerError)
	}

	return resource, nil
}

func (s *Service) UpdateResource(ctx context.Context, resource *models.Resource) error {
	// todo: update should return with the updated resource
	log.Debug(ctx, "Updating resource")

	if err := s.mySQLStorage.UpdateResource(resource); err != nil {
		if err.Error() == storage.ErrResourceNotFound.Error() {
			return myerrors.WithFields(err, models.HTTPCode, http.StatusAccepted)
		}
		return myerrors.WithFields(err, models.HTTPCode, http.StatusInternalServerError)
	}

	return nil
}

func (s *Service) DeleteResource(ctx context.Context, resource *models.Resource) error {
	log.Debug(ctx, "Deleting resource")

	if err := s.mySQLStorage.DeleteResource(resource); err != nil {
		if err.Error() == storage.ErrResourceNotFound.Error() {
			return myerrors.WithFields(err, models.HTTPCode, http.StatusAccepted)
		}
		return myerrors.WithFields(err, models.HTTPCode, http.StatusInternalServerError)
	}

	_, err := s.mySQLStorage.GetResourceByID(resource.ID)
	if err != nil && err.Error() != storage.ErrResourceNotFound.Error() {
		return myerrors.WithFields(err, models.HTTPCode, http.StatusInternalServerError)
	}

	return nil
}

func (s *Service) GetCategories(ctx context.Context) (*httpModels.CategoriesResponse, error) {
	log.Debug(ctx, "Getting categories")

	categories, err := s.mySQLStorage.GetCategories()
	if err != nil {
		return nil, myerrors.WithFields(err, models.HTTPCode, http.StatusInternalServerError)
	}

	return &httpModels.CategoriesResponse{
		Categories: categories,
	}, nil
}

func (s *Service) GetResourcesByCategory(ctx context.Context, req *httpModels.GetResourcesByCategoryRequest) (*httpModels.ResourcesResponse, error) {
	log.Debug(ctx, "Getting multiple resources by category")

	resources, err := s.mySQLStorage.GetResourcesByCategory(req.Category)
	if err != nil {
		if err.Error() == storage.ErrResourceNotFound.Error() {
			return nil, myerrors.WithFields(err, models.HTTPCode, http.StatusAccepted)
		}
		return nil, myerrors.WithFields(err, models.HTTPCode, http.StatusInternalServerError)
	}

	return &httpModels.ResourcesResponse{
		Resources: resources,
	}, nil
}

func (s *Service) GetResourcesByIDs(_ context.Context, req *httpModels.GetResourcesByIDsRequest) (*httpModels.ResourcesResponse, error) {
	resources, err := s.mySQLStorage.GetResourcesByIDs(req.UUIDs)
	if err != nil {
		if err.Error() == storage.ErrResourceNotFound.Error() {
			return nil, myerrors.WithFields(err, models.HTTPCode, http.StatusAccepted)
		}
		return nil, myerrors.WithFields(err, models.HTTPCode, http.StatusInternalServerError)
	}

	return &httpModels.ResourcesResponse{
		Resources: resources,
	}, nil
}
