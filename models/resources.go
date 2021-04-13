package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

var LocationKey = "location"

var CategoryContent = "Content"

type Resource struct {
	ID       uuid.UUID  `json:"id" validate:"required"`
	Category int        `json:"category" validate:"required"`
	Content  ContentMap `json:"content" validate:"required"`
}

type ContentMap map[string]string

func (cm *ContentMap) Scan(src interface{}) error {
	var source []byte
	_cm := make(map[string]string)

	switch s := src.(type) {
	case []uint8:
		source = []byte(s)
	case nil:
		return nil
	default:
		return errors.New("incompatible type for StringInterfaceMap")
	}
	err := json.Unmarshal(source, &_cm)
	if err != nil {
		return err
	}
	*cm = ContentMap(_cm)
	return nil
}

func (cm ContentMap) Value() (driver.Value, error) {
	if len(cm) == 0 {
		return nil, nil
	}
	j, err := json.Marshal(cm)
	if err != nil {
		return nil, err
	}
	return driver.Value([]byte(j)), nil
}

type Category struct {
	ID          int    `json:"id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func SetField(content ContentMap, keyString string, field string) {
	content[keyString] = field
}

func GetField(content ContentMap, keyString string, defaultValue string) string {
	field, ok := content[keyString]
	if !ok {
		return defaultValue
	}
	return field
}

func ClearAsset(content ContentMap, keyString string) error {
	if _, ok := content[keyString]; !ok {
		return errors.New("unknown content reference")
	}
	delete(content, keyString)
	return nil
}

func NewResource(idString string, category int, location string) (*Resource, error) {
	id, err := uuid.Parse(idString)
	if err != nil {
		return nil, err
	}

	content := make(ContentMap)
	content[LocationKey] = location

	resource := &Resource{
		ID:       id,
		Category: category,
		Content:  content,
	}

	return resource, nil
}
