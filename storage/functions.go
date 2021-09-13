package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/artofimagination/mysql-resources-db-go-service/models"
)

var ErrResourcesMissing = errors.New("This resources is missing or old value is the same as new")

func rollbackWithErrorStack(tx *sql.Tx, errorStack error) error {
	if err := tx.Rollback(); err != nil {
		errorString := fmt.Sprintf("%s\n%s\n", errorStack.Error(), err.Error())
		return errors.Wrap(errors.WithStack(errors.New(errorString)), "Failed to rollback changes")
	}
	return errorStack
}

const addResourceQuery = `
	INSERT INTO 
	resources(id, category, content) 
	VALUES 
	(UUID_TO_BIN(?), ?, CAST(CONVERT(? USING utf8) AS JSON))
`

func addResource(resource *models.Resource, tx *sql.Tx) error {
	// Execute transaction
	_, err := tx.Exec(addResourceQuery, resource.ID, resource.Category, resource.Content)
	if err != nil {
		return rollbackWithErrorStack(tx, errors.WithStack(err))
	}

	return nil
}

const updateResourceQuery = `
	UPDATE resources 
	SET content = CAST(CONVERT(? USING utf8) AS JSON), category = ? 
	WHERE id = UUID_TO_BIN(?)
`

func updateResource(resource *models.Resource, tx *sql.Tx) error {
	result, err := tx.Exec(updateResourceQuery, resource.Content, resource.Category, resource.ID)
	if err != nil {
		return rollbackWithErrorStack(tx, errors.WithStack(err))
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return rollbackWithErrorStack(tx, errors.WithStack(err))
	}

	if affected == 0 {
		return rollbackWithErrorStack(tx, ErrResourcesMissing)
	}

	return nil
}

const getResourceByIDQuery = `
	SELECT BIN_TO_UUID(id), category, content 
	FROM resources 
	WHERE id = UUID_TO_BIN(?)
`

func (mySQL *MySQL) getResourceByID(resourceID uuid.UUID) (*models.Resource, error) {
	resource := &models.Resource{}

	tx, err := mySQL.db.Begin()
	if err != nil {
		return nil, err
	}

	result := tx.QueryRow(getResourceByIDQuery, resourceID)

	err = result.Scan(&resource.ID, &resource.Category, &resource.Content)
	switch {
	case err == sql.ErrNoRows:
		if errRb := tx.Commit(); errRb != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, rollbackWithErrorStack(tx, errors.WithStack(err))
	default:
	}

	return resource, tx.Commit()
}

const deleteResourceQuery = `
	DELETE FROM resources 
	WHERE id=UUID_TO_BIN(?)
`

func deleteResource(resourceID string, tx *sql.Tx) error {
	result, err := tx.Exec(deleteResourceQuery, resourceID)
	if err != nil {
		return rollbackWithErrorStack(tx, errors.WithStack(err))
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return rollbackWithErrorStack(tx, errors.WithStack(err))
	}

	if affected == 0 {
		return rollbackWithErrorStack(tx, ErrResourcesMissing)
	}
	return nil
}

var GetResourcesByIDsQuery = "SELECT BIN_TO_UUID(id), category, content FROM resources WHERE id IN (UUID_TO_BIN(?)"

func getResourcesByIDs(IDs []uuid.UUID, tx *sql.Tx) ([]models.Resource, error) {
	query := GetResourcesByIDsQuery + strings.Repeat(",UUID_TO_BIN(?)", len(IDs)-1) + ") ORDER BY created_at DESC"
	interfaceList := make([]interface{}, len(IDs))
	for i := range IDs {
		interfaceList[i] = IDs[i]
	}
	rows, err := tx.Query(query, interfaceList...)
	if err != nil {
		return nil, rollbackWithErrorStack(tx, errors.WithStack(err))
	}

	defer func() {
		_ = rows.Close()
	}()

	resources := make([]models.Resource, 0)
	for rows.Next() {
		resource := models.Resource{}
		err := rows.Scan(&resource.ID, &resource.Category, &resource.Content)
		if err != nil {
			return nil, rollbackWithErrorStack(tx, errors.WithStack(err))
		}
		resources = append(resources, resource)
	}
	err = rows.Err()
	if err != nil {
		return nil, rollbackWithErrorStack(tx, errors.WithStack(err))
	}

	if len(resources) == 0 {
		return nil, sql.ErrNoRows
	}

	return resources, nil
}

const getResourceByCategoryQuery = `
	SELECT BIN_TO_UUID(id), category, content 
	FROM resources 
	WHERE category = ? 
	ORDER BY created_at DESC
`

func getResourcesByCategory(category int, tx *sql.Tx) ([]models.Resource, error) {
	rows, err := tx.Query(getResourceByCategoryQuery, category)
	if err != nil {
		return nil, rollbackWithErrorStack(tx, errors.WithStack(err))
	}

	defer func() {
		_ = rows.Close()
	}()

	resources := make([]models.Resource, 0)
	for rows.Next() {
		resource := models.Resource{}
		err := rows.Scan(&resource.ID, &resource.Category, &resource.Content)
		if err != nil {
			return nil, rollbackWithErrorStack(tx, errors.WithStack(err))
		}
		resources = append(resources, resource)
	}
	err = rows.Err()
	if err != nil {
		return nil, rollbackWithErrorStack(tx, errors.WithStack(err))
	}

	if len(resources) == 0 {
		return nil, sql.ErrNoRows
	}

	return resources, nil
}

var GetCategoryByNameQuery = "SELECT id, name, description FROM categories WHERE name = ?"

func (mySQL *MySQL) getCategoryByName(name string) (*models.Category, error) {
	category := &models.Category{}

	tx, err := mySQL.db.Begin()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	result := tx.QueryRow(GetCategoryByNameQuery, name)

	err = result.Scan(&category.ID, &category.Name, &category.Description)
	switch {
	case err == sql.ErrNoRows:
		if errRb := tx.Commit(); errRb != nil {
			return nil, errors.WithStack(err)
		}
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, rollbackWithErrorStack(tx, errors.WithStack(err))
	default:
	}

	return category, tx.Commit()
}

const getCategoryByIDQuery = `
	SELECT id, name, description 
	FROM categories WHERE id = ?
`

func (mySQL *MySQL) GetCategoryByID(id int) (*models.Category, error) {
	category := &models.Category{}

	tx, err := mySQL.db.Begin()
	if err != nil {
		return nil, err
	}

	result := tx.QueryRow(getCategoryByIDQuery, id)

	err = result.Scan(&category.ID, &category.Name, &category.Description)
	switch {
	case err == sql.ErrNoRows:
		if errRb := tx.Commit(); errRb != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, rollbackWithErrorStack(tx, errors.WithStack(err))
	default:
	}

	return category, tx.Commit()
}

const getCategorsQuery = `
	SELECT id, name, description 
	FROM categories
`

func (mySQL *MySQL) getCategories() ([]models.Category, error) {
	tx, err := mySQL.db.Begin()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(getCategorsQuery)
	if err != nil {
		return nil, rollbackWithErrorStack(tx, errors.WithStack(err))
	}

	defer func() {
		_ = rows.Close()
	}()

	categories := make([]models.Category, 0)
	for rows.Next() {
		category := models.Category{}
		err := rows.Scan(&category.ID, &category.Name, &category.Description)
		if err != nil {
			return nil, rollbackWithErrorStack(tx, errors.WithStack(err))
		}
		categories = append(categories, category)
	}
	err = rows.Err()
	if err != nil {
		return nil, rollbackWithErrorStack(tx, errors.WithStack(err))
	}

	if len(categories) == 0 {
		return nil, sql.ErrNoRows
	}

	return categories, nil
}
