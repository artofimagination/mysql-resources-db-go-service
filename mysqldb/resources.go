package mysqldb

import (
	"database/sql"
	"encoding/json"
	"strings"

	"github.com/artofimagination/mysql-resources-db-go-service/models"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

var ErrResourcesMissing = errors.New("This resources is missing or old value is the same as new")

var AddResourceQuery = "INSERT INTO resources (id, category, content) VALUES (UUID_TO_BIN(?), ?, ?)"

func (*MYSQLFunctions) AddResource(resource *models.Resource, tx *sql.Tx) error {
	// Prepare data
	binary, err := json.Marshal(resource.Content)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	// Execute transaction
	_, err = tx.Exec(AddResourceQuery, resource.ID, resource.Category, binary)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	return nil
}

var UpdateResourceQuery = "UPDATE resources set content = ?, category = ? where id = UUID_TO_BIN(?)"

func (f *MYSQLFunctions) UpdateResource(resource *models.Resource, tx *sql.Tx) error {
	binary, err := json.Marshal(resource.Content)
	if err != nil {
		return err
	}

	result, err := tx.Exec(UpdateResourceQuery, binary, resource.Category, resource.ID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	if affected == 0 {
		return RollbackWithErrorStack(tx, ErrResourcesMissing)
	}

	return nil
}

var GetResourceByIDQuery = "SELECT BIN_TO_UUID(id), category, content FROM resources WHERE id = UUID_TO_BIN(?)"

func (f *MYSQLFunctions) GetResourceByID(resourceID *uuid.UUID) (*models.Resource, error) {
	resource := &models.Resource{}

	tx, err := f.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	result := tx.QueryRow(GetResourceByIDQuery, resourceID)

	content := []byte{}
	err = result.Scan(&resource.ID, &resource.Category, &content)
	switch {
	case err == sql.ErrNoRows:
		if errRb := tx.Commit(); errRb != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}

	if err := json.Unmarshal(content, &resource.Content); err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	return resource, tx.Commit()
}

var DeleteResourceQuery = "DELETE FROM resources WHERE id=UUID_TO_BIN(?)"

func (*MYSQLFunctions) DeleteResource(resourceID string, tx *sql.Tx) error {
	result, err := tx.Exec(DeleteResourceQuery, resourceID)
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return RollbackWithErrorStack(tx, err)
	}

	if affected == 0 {
		return RollbackWithErrorStack(tx, ErrResourcesMissing)
	}
	return nil
}

var GetResourcesByIDsQuery = "SELECT BIN_TO_UUID(id), category, content FROM resources WHERE id IN (UUID_TO_BIN(?)"

func (*MYSQLFunctions) GetResourcesByIDs(IDs []uuid.UUID, tx *sql.Tx) ([]models.Resource, error) {
	query := GetResourcesByIDsQuery + strings.Repeat(",UUID_TO_BIN(?)", len(IDs)-1) + ") ORDER BY created_at DESC"
	interfaceList := make([]interface{}, len(IDs))
	for i := range IDs {
		interfaceList[i] = IDs[i]
	}
	rows, err := tx.Query(query, interfaceList...)
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	defer rows.Close()

	resources := make([]models.Resource, 0)
	for rows.Next() {
		content := []byte{}
		resource := models.Resource{}
		err := rows.Scan(&resource.ID, &resource.Category, &content)
		if err != nil {
			return nil, RollbackWithErrorStack(tx, err)
		}
		if err := json.Unmarshal(content, &resource.Content); err != nil {
			return nil, RollbackWithErrorStack(tx, err)
		}
		resources = append(resources, resource)
	}
	err = rows.Err()
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	if len(resources) == 0 {
		return nil, sql.ErrNoRows
	}

	return resources, nil
}

var GetResourceByCategoryQuery = "SELECT BIN_TO_UUID(id), category, content FROM resources WHERE category = ? ORDER BY created_at DESC"

func (*MYSQLFunctions) GetResourcesByCategory(category int, tx *sql.Tx) ([]models.Resource, error) {
	rows, err := tx.Query(GetResourceByCategoryQuery, category)
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	defer rows.Close()

	resources := make([]models.Resource, 0)
	for rows.Next() {
		content := []byte{}
		resource := models.Resource{}
		err := rows.Scan(&resource.ID, &resource.Category, &content)
		if err != nil {
			return nil, RollbackWithErrorStack(tx, err)
		}
		if err := json.Unmarshal(content, &resource.Content); err != nil {
			return nil, RollbackWithErrorStack(tx, err)
		}
		resources = append(resources, resource)
	}
	err = rows.Err()
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	if len(resources) == 0 {
		return nil, sql.ErrNoRows
	}

	return resources, nil
}

var GetCategoryByNameQuery = "SELECT id, name, description FROM categories WHERE name = ?"

func (f *MYSQLFunctions) GetCategoryByName(name string) (*models.Category, error) {
	category := &models.Category{}

	tx, err := f.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	result := tx.QueryRow(GetCategoryByNameQuery, name)

	err = result.Scan(&category.ID, &category.Name, &category.Description)
	switch {
	case err == sql.ErrNoRows:
		if errRb := tx.Commit(); errRb != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}

	return category, tx.Commit()
}

var GetCategoryByIDQuery = "SELECT id, name, description FROM categories WHERE id = ?"

func (f *MYSQLFunctions) GetCategoryByID(id int) (*models.Category, error) {
	category := &models.Category{}

	tx, err := f.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	result := tx.QueryRow(GetCategoryByIDQuery, id)

	err = result.Scan(&category.ID, &category.Name, &category.Description)
	switch {
	case err == sql.ErrNoRows:
		if errRb := tx.Commit(); errRb != nil {
			return nil, err
		}
		return nil, sql.ErrNoRows
	case err != nil:
		return nil, RollbackWithErrorStack(tx, err)
	default:
	}

	return category, tx.Commit()
}

var GetCategorsQuery = "SELECT id, name, description FROM categories"

func (f *MYSQLFunctions) GetCategories() ([]models.Category, error) {
	tx, err := f.DBConnector.ConnectSystem()
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(GetCategorsQuery)
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	defer rows.Close()

	categories := make([]models.Category, 0)
	for rows.Next() {
		category := models.Category{}
		err := rows.Scan(&category.ID, &category.Name, &category.Description)
		if err != nil {
			return nil, RollbackWithErrorStack(tx, err)
		}
		categories = append(categories, category)
	}
	err = rows.Err()
	if err != nil {
		return nil, RollbackWithErrorStack(tx, err)
	}

	if len(categories) == 0 {
		return nil, sql.ErrNoRows
	}

	return categories, nil
}
