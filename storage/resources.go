package storage

import (
	"database/sql"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/artofimagination/mysql-resources-db-go-service/models"
)

var ErrResourceNotFound = errors.New("The selected resource not found")
var ErrResourceAlreadyExists = errors.New("The resource already exists")
var ErrResourceHasTooManyAttachments = errors.New("The resource has too many attachements")

var ErrDuplicateEntrySubString = "Duplicate entry"

// MaxContentItems describes the maximum number or resources to upload to a resources an attachement
var MaxContentItems = 2

func (mySQL *MySQL) AddResource(resource *models.Resource) (err error) {

	tx, err := mySQL.db.Begin()
	if err != nil {
		return errors.WithStack(err)
	}

	category, err := mySQL.getCategoryByName(models.CategoryContent)
	if err != nil {
		return errors.WithStack(err)
	}

	if len(resource.Content) > MaxContentItems {
		return errors.WithStack(ErrResourceHasTooManyAttachments)
	}

	for k, v := range resource.Content {
		if k != models.LocationKey {
			resourceItem, err := models.NewResource(k, category.ID, v)
			if err != nil {
				return errors.WithStack(err)
			}
			if err := addResource(resourceItem, tx); err != nil {
				if strings.Contains(err.Error(), ErrDuplicateEntrySubString) {
					return errors.WithStack(ErrResourceAlreadyExists)
				}
				return err
			}
		}
	}

	if err := addResource(resource, tx); err != nil {
		if strings.Contains(err.Error(), ErrDuplicateEntrySubString) {
			return errors.WithStack(ErrResourceAlreadyExists)
		}
		return errors.WithStack(err)
	}

	return tx.Commit()
}

func (mySQL *MySQL) GetResourcesByCategory(category int) ([]models.Resource, error) {
	tx, err := mySQL.db.Begin()
	if err != nil {
		return nil, err
	}

	resources, err := getResourcesByCategory(category, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	return resources, tx.Commit()
}

func (mySQL *MySQL) GetResourcesByIDs(IDs []uuid.UUID) ([]models.Resource, error) {
	tx, err := mySQL.db.Begin()
	if err != nil {
		return nil, err
	}

	resources, err := getResourcesByIDs(IDs, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	return resources, tx.Commit()
}

func (mySQL *MySQL) GetResourceByID(ID uuid.UUID) (*models.Resource, error) {
	resources, err := mySQL.getResourceByID(ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	return resources, nil
}

func (mySQL *MySQL) UpdateResource(resource *models.Resource) error {
	tx, err := mySQL.db.Begin()
	if err != nil {
		return err
	}

	if len(resource.Content) > MaxContentItems {
		return rollbackWithErrorStack(tx, errors.WithStack(ErrResourceHasTooManyAttachments))
	}

	resourceFromDB, err := mySQL.getResourceByID(resource.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return rollbackWithErrorStack(tx, errors.WithStack(ErrResourceNotFound))
		}
		return err
	}

	category, err := mySQL.getCategoryByName(models.CategoryContent)
	if err != nil {
		return rollbackWithErrorStack(tx, errors.WithStack(err))
	}

	for k, v := range resource.Content {
		if _, ok := resourceFromDB.Content[k]; !ok {
			resourceItem, err := models.NewResource(k, category.ID, v)
			if err != nil {
				return rollbackWithErrorStack(tx, errors.WithStack(err))
			}
			if err := addResource(resourceItem, tx); err != nil {
				if strings.Contains(err.Error(), ErrDuplicateEntrySubString) {
					return ErrResourceAlreadyExists
				}
				return err
			}
		}
	}

	if err := updateResource(resource, tx); err != nil {
		if err == ErrResourcesMissing {
			return ErrResourceNotFound
		}
		return err
	}

	return tx.Commit()
}

func (mySQL *MySQL) DeleteResource(id uuid.UUID, content models.ContentMap) error {
	tx, err := mySQL.db.Begin()
	if err != nil {
		return err
	}

	for k := range content {
		if k != models.LocationKey {
			if err := deleteResource(k, tx); err != nil {
				if err == ErrResourcesMissing {
					return ErrResourceNotFound
				}
				return err
			}
		}
	}

	if err := deleteResource(id.String(), tx); err != nil {
		if err == ErrResourcesMissing {
			return ErrResourceNotFound
		}
		return err
	}

	return tx.Commit()
}

func (mySQL *MySQL) GetCategories() ([]models.Category, error) {
	categories, err := mySQL.getCategories()
	if err != nil {
		return nil, err
	}

	return categories, nil
}
