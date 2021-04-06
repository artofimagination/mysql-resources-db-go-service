package dbcontrollers

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/artofimagination/mysql-resources-db-go-service/models"
	"github.com/artofimagination/mysql-resources-db-go-service/mysqldb"
	"github.com/google/uuid"
)

var ErrResourceNotFound = errors.New("The selected resource not found")
var ErrResourceAlreadyExists = errors.New("The resource already exists")
var ErrResourceHasTooManyAttachements = errors.New("The resource has too many attachements")

var ErrDuplicateEntrySubString = "Duplicate entry"

// MaxContentItems describes the maximum number or resources to upload to a resources an attachement
var MaxContentItems = 2

func (c *MYSQLController) AddResource(resource *models.Resource) error {
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	category, err := c.DBFunctions.GetCategoryByName(models.CategoryContent)
	if err != nil {
		return err
	}

	if len(resource.Content) > MaxContentItems {
		return ErrResourceHasTooManyAttachements
	}

	for k, v := range resource.Content {
		if k != models.LocationKey {
			resourceItem, err := models.NewResource(k, category.ID, v)
			if err != nil {
				return mysqldb.RollbackWithErrorStack(tx, err)
			}
			if err := c.DBFunctions.AddResource(resourceItem, tx); err != nil {
				if strings.Contains(err.Error(), ErrDuplicateEntrySubString) {
					return ErrResourceAlreadyExists
				}
				return err
			}
		}
	}

	if err := c.DBFunctions.AddResource(resource, tx); err != nil {
		if strings.Contains(err.Error(), ErrDuplicateEntrySubString) {
			return ErrResourceAlreadyExists
		}
		return err
	}

	return c.DBConnector.Commit(tx)
}

func (c *MYSQLController) GetResourcesByCategory(category int) ([]models.Resource, error) {
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	resources, err := c.DBFunctions.GetResourcesByCategory(category, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	return resources, c.DBConnector.Commit(tx)
}

func (c *MYSQLController) GetResourcesByIDs(IDs []uuid.UUID) ([]models.Resource, error) {
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	resources, err := c.DBFunctions.GetResourcesByIDs(IDs, tx)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	return resources, c.DBConnector.Commit(tx)
}

func (c *MYSQLController) GetResourceByID(ID uuid.UUID) (*models.Resource, error) {
	resources, err := c.DBFunctions.GetResourceByID(ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrResourceNotFound
		}
		return nil, err
	}

	return resources, nil
}

func (c *MYSQLController) UpdateResource(resource *models.Resource) error {
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	if len(resource.Content) > MaxContentItems {
		return mysqldb.RollbackWithErrorStack(tx, ErrResourceHasTooManyAttachements)
	}

	resourceFromDB, err := c.DBFunctions.GetResourceByID(resource.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return mysqldb.RollbackWithErrorStack(tx, ErrResourceNotFound)
		}
		return err
	}

	category, err := c.DBFunctions.GetCategoryByName(models.CategoryContent)
	if err != nil {
		return mysqldb.RollbackWithErrorStack(tx, err)
	}

	for k, v := range resource.Content {
		if _, ok := resourceFromDB.Content[k]; !ok {
			resourceItem, err := models.NewResource(k, category.ID, v)
			if err != nil {
				return mysqldb.RollbackWithErrorStack(tx, err)
			}
			if err := c.DBFunctions.AddResource(resourceItem, tx); err != nil {
				if strings.Contains(err.Error(), ErrDuplicateEntrySubString) {
					return ErrResourceAlreadyExists
				}
				return err
			}
		}
	}

	if err := c.DBFunctions.UpdateResource(resource, tx); err != nil {
		if err == mysqldb.ErrResourcesMissing {
			return ErrResourceNotFound
		}
		return err
	}

	return c.DBConnector.Commit(tx)
}

func (c *MYSQLController) DeleteResource(resource *models.Resource) error {
	tx, err := c.DBConnector.ConnectSystem()
	if err != nil {
		return err
	}

	for k := range resource.Content {
		if k != models.LocationKey {
			if err := c.DBFunctions.DeleteResource(k, tx); err != nil {
				if err == mysqldb.ErrResourcesMissing {
					return ErrResourceNotFound
				}
				return err
			}
		}
	}

	if err := c.DBFunctions.DeleteResource(resource.ID.String(), tx); err != nil {
		if err == mysqldb.ErrResourcesMissing {
			return ErrResourceNotFound
		}
		return err
	}

	return c.DBConnector.Commit(tx)
}

func (c *MYSQLController) GetCategories() ([]models.Category, error) {
	categories, err := c.DBFunctions.GetCategories()
	if err != nil {
		return nil, err
	}

	return categories, nil
}
