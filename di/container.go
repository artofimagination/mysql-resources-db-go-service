package di

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/proemergotech/log/v3"
	"github.com/proemergotech/log/v3/echolog"

	"github.com/artofimagination/mysql-resources-db-go-service/config"
	"github.com/artofimagination/mysql-resources-db-go-service/rest"
	"github.com/artofimagination/mysql-resources-db-go-service/service"
	"github.com/artofimagination/mysql-resources-db-go-service/storage"
	"github.com/artofimagination/mysql-resources-db-go-service/validation"
)

type Container struct {
	RestServer *rest.Server
	database   *sqlx.DB
}

func NewContainer(cfg *config.Config) (*Container, error) {
	c := &Container{}

	v, err := NewValidator()
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize validator")
	}

	c.database, err = newSQLDatabase(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize MySQL database")
	}

	mysqlStorage := storage.NewMySQL(c.database)

	err = mysqlStorage.BootstrapSystem(cfg.MySQLDBMigrationDirectory)
	if err != nil {
		return nil, errors.Wrap(err, "cannot bootstrap MySQL database")
	}

	echoEngine := newEcho(cfg.Port, v, rest.DLiveRHTTPErrorHandler)

	svc := service.NewService(mysqlStorage)

	c.RestServer = rest.NewServer(
		echoEngine,
		rest.NewController(
			echoEngine,
			svc,
			cfg.DebugPProf,
		),
	)

	return c, nil
}

func NewValidator() (*validation.Validator, error) {
	v := validator.New()

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		structTags := []string{"param", "json", "query"}
		for _, t := range structTags {
			name := strings.SplitN(field.Tag.Get(t), ",", 2)[0]
			if name != "" && name != "-" {
				return name
			}
		}
		return ""
	})

	// use it for fields with type slice and map - for these `required` isn't working as expected
	err := v.RegisterValidation("notblank", validators.NotBlank)
	if err != nil {
		return nil, err
	}

	return validation.NewValidator(v), nil
}

func newSQLDatabase(cfg *config.Config) (*sqlx.DB, error) {
	return sqlx.Open(
		"mysql",
		fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?parseTime=true&interpolateParams=true",
			cfg.MySQLDBUser,
			cfg.MySQLDBPassword,
			cfg.MySQLDBAddress,
			cfg.MySQLDBPort,
			cfg.MySQLDBName,
		),
	)
}

func newEcho(port int, validator *validation.Validator, httpErrorHandler echo.HTTPErrorHandler) *echo.Echo {
	e := echo.New()

	e.Use(echolog.RecoveryMiddleware(log.GlobalLogger()))
	e.HTTPErrorHandler = httpErrorHandler
	e.Validator = validator
	e.HideBanner = true
	e.HidePort = true

	e.Server = &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: e,
	}

	return e
}

func (c *Container) Close() {
	if err := c.database.Close(); err != nil {
		err = errors.Wrap(err, "Database graceful close failed")
		log.Warn(context.Background(), err.Error(), "error", err)
	}
}
