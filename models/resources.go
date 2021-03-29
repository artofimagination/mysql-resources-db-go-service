package models

import (
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
		return errors.New("Unknown content reference")
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
