package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/proemergotech/log/v3"

	"github.com/artofimagination/mysql-resources-db-go-service/models"
	httpModels "github.com/artofimagination/mysql-resources-db-go-service/models/http"
	"github.com/artofimagination/mysql-resources-db-go-service/models/myerrors"
)

func DLiveRHTTPErrorHandler(err error, eCtx echo.Context) {
	defer func() {
		sc := eCtx.Response().Status
		if sc >= 400 && sc < 500 {
			log.Warn(eCtx.Request().Context(), err.Error(), "error", err)
		} else {
			log.Error(eCtx.Request().Context(), err.Error(), "error", err)
		}
	}()

	if eErr, ok := err.(*echo.HTTPError); ok {
		sc := eErr.Code

		switch sc {
		case http.StatusNotFound:
			err = myerrors.WithFields(errors.Wrap(errors.WithStack(eErr), "route not found"), models.HTTPCode, http.StatusNotFound)
		case http.StatusMethodNotAllowed:
			err = myerrors.WithFields(errors.Wrap(errors.WithStack(eErr), "method not allowed"), models.HTTPCode, http.StatusMethodNotAllowed)
		default:
			err = myerrors.WithFields(errors.Wrap(errors.WithStack(eErr), "semantic error"), models.HTTPCode, http.StatusInternalServerError)
		}
	}

	statusCode := http.StatusInternalServerError
	httpCode := myerrors.Field(err, models.HTTPCode)
	if httpCode != nil {
		statusCode = httpCode.(int)
	}

	_ = eCtx.JSON(statusCode, httpModels.ResponseData{
		Error: errors.WithStack(err).Error(),
	})
}
