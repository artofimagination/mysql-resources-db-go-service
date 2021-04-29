package validation

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/pkg/errors"

	"github.com/artofimagination/mysql-resources-db-go-service/models"
	"github.com/artofimagination/mysql-resources-db-go-service/models/myerrors"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator(
	validator *validator.Validate,
) *Validator {
	return &Validator{
		validator: validator,
	}
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return myerrors.WithFields(errors.Wrap(errors.WithStack(err), "validation error"), models.HTTPCode, http.StatusBadRequest)
	}

	return nil
}
