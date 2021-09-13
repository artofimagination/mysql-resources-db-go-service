package rest

import (
	"net/http"
	"net/http/pprof"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/proemergotech/log/v3"
	"github.com/proemergotech/log/v3/echolog"

	"github.com/artofimagination/mysql-resources-db-go-service/models"
	httpModels "github.com/artofimagination/mysql-resources-db-go-service/models/http"
	"github.com/artofimagination/mysql-resources-db-go-service/service"
)

type controller struct {
	echoEngine *echo.Echo
	svc        *service.Service
	debugPProf bool
}

func NewController(
	echoEngine *echo.Echo,
	svc *service.Service,
	debugPProf bool,
) Controller {
	return &controller{
		echoEngine: echoEngine,
		svc:        svc,
		debugPProf: debugPProf,
	}
}

func (c *controller) Start() {

	c.echoEngine.Use(echolog.DebugMiddleware(log.GlobalLogger(), true, true))
	if c.debugPProf {
		runtime.SetBlockProfileRate(1)
		runtime.SetMutexProfileFraction(5)

		c.echoEngine.GET("/debug/pprof/*", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
		c.echoEngine.GET("/debug/pprof/cmdline", echo.WrapHandler(http.HandlerFunc(pprof.Cmdline)))
		c.echoEngine.GET("/debug/pprof/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
		c.echoEngine.GET("/debug/pprof/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
		c.echoEngine.GET("/debug/pprof/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))
	}

	c.echoEngine.GET("/", func(eCtx echo.Context) error {
		return eCtx.String(http.StatusOK, "Hi! I am a resources database server!")
	})

	c.echoEngine.GET("/healthcheck", func(eCtx echo.Context) error {
		return eCtx.NoContent(http.StatusOK)
	})

	c.echoEngine.POST("/add-resource", func(eCtx echo.Context) error {
		resource := &models.Resource{}
		if err := eCtx.Bind(resource); err != nil {
			return errors.Wrap(err, "cannot bind request")
		}

		if err := eCtx.Validate(resource); err != nil {
			return errors.Wrap(err, "request failed on validation")
		}

		_, err := c.svc.AddResource(eCtx.Request().Context(), resource)
		if err != nil {
			return err
		}

		return eCtx.JSON(http.StatusCreated, httpModels.ResponseData{Data: "OK"})
	})

	c.echoEngine.GET("/get-resource-by-id", func(eCtx echo.Context) error {
		req := &httpModels.GetResourceByIDWithQueryRequest{}
		if err := eCtx.Bind(req); err != nil {
			return err
		}

		if err := eCtx.Validate(req); err != nil {
			return err
		}

		resp, err := c.svc.GetResourceByID(eCtx.Request().Context(), req.UUID)
		if err != nil {
			return err
		}

		return eCtx.JSON(http.StatusOK, httpModels.ResponseData{Data: resp})
	})

	c.echoEngine.POST("/update-resource", func(eCtx echo.Context) error {
		resource := &models.Resource{}
		if err := eCtx.Bind(resource); err != nil {
			return err
		}

		if err := eCtx.Validate(resource); err != nil {
			return err
		}

		err := c.svc.UpdateResource(eCtx.Request().Context(), resource)
		if err != nil {
			return err
		}

		return eCtx.JSON(http.StatusCreated, httpModels.ResponseData{Data: "OK"})
	})

	c.echoEngine.POST("/delete-resource", func(eCtx echo.Context) error {
		req := &httpModels.DeleteResourceRequest{}
		if err := eCtx.Bind(req); err != nil {
			return err
		}

		if err := eCtx.Validate(req); err != nil {
			return err
		}

		err := c.svc.DeleteResource(eCtx.Request().Context(), req)
		if err != nil {
			return err
		}

		return eCtx.JSON(http.StatusOK, httpModels.ResponseData{Data: "OK"})
	})

	c.echoEngine.GET("/get-categories", func(eCtx echo.Context) error {
		resp, err := c.svc.GetCategories(eCtx.Request().Context())
		if err != nil {
			return err
		}

		return eCtx.JSON(http.StatusOK, httpModels.ResponseData{Data: resp})
	})

	c.echoEngine.GET("/get-resources-by-ids", func(eCtx echo.Context) error {
		req := &httpModels.GetResourcesByIDsRequest{}
		if err := eCtx.Bind(req); err != nil {
			return err
		}

		if err := eCtx.Validate(req); err != nil {
			return err
		}

		resp, err := c.svc.GetResourcesByIDs(eCtx.Request().Context(), req)
		if err != nil {
			return err
		}

		return eCtx.JSON(http.StatusOK, httpModels.ResponseData{Data: resp})
	})

	c.echoEngine.GET("/get-resources-by-category", func(eCtx echo.Context) error {
		req := &httpModels.GetResourcesByCategoryRequest{}
		if err := eCtx.Bind(req); err != nil {
			return err
		}

		if err := eCtx.Validate(req); err != nil {
			return err
		}

		resp, err := c.svc.GetResourcesByCategory(eCtx.Request().Context(), req)
		if err != nil {
			return err
		}

		return eCtx.JSON(http.StatusOK, httpModels.ResponseData{Data: resp})
	})

	// new endpoint format follows REST and CRUD basics
	apiRoutes := c.echoEngine.Group("/api/v1")
	apiRoutes.Use(echolog.DebugMiddleware(log.GlobalLogger(), true, true))

	resourcesRoutes := apiRoutes.Group("/resources")
	resourcesRoutes.GET("/", func(eCtx echo.Context) error {
		req := &httpModels.GetResourcesByIDsRequest{}
		if err := eCtx.Bind(req); err != nil {
			return err
		}

		if err := eCtx.Validate(req); err != nil {
			return err
		}

		resp, err := c.svc.GetResourcesByIDs(eCtx.Request().Context(), req)
		if err != nil {
			return err
		}

		return eCtx.JSON(http.StatusOK, httpModels.ResponseData{Data: resp})
	})

	resourcesRoutes.GET("/categories/:category", func(eCtx echo.Context) error {
		req := &httpModels.GetResourcesByCategoryRequest{}
		if err := eCtx.Bind(req); err != nil {
			return err
		}

		if err := eCtx.Validate(req); err != nil {
			return err
		}

		resp, err := c.svc.GetResourcesByCategory(eCtx.Request().Context(), req)
		if err != nil {
			return err
		}

		return eCtx.JSON(http.StatusOK, httpModels.ResponseData{Data: resp})
	})

	resourcesCRUDRoutes := resourcesRoutes.Group("/:resource_id")
	resourcesCRUDRoutes.GET("/", func(eCtx echo.Context) error {
		req := &httpModels.GetResourceByIDRequest{}
		if err := eCtx.Bind(req); err != nil {
			return err
		}

		if err := eCtx.Validate(req); err != nil {
			return err
		}

		resp, err := c.svc.GetResourceByID(eCtx.Request().Context(), req.UUID)
		if err != nil {
			return err
		}

		return eCtx.JSON(http.StatusOK, resp)
	})

	resourcesCRUDRoutes.POST("/", func(eCtx echo.Context) error {
		req := &httpModels.AddResourceRequest{}
		if err := eCtx.Bind(req); err != nil {
			return err
		}

		if err := eCtx.Validate(req); err != nil {
			return err
		}

		resp, err := c.svc.AddResource(eCtx.Request().Context(), req.Resource)
		if err != nil {
			return err
		}

		return eCtx.JSON(http.StatusOK, resp) // todo: if the resource not exists should be created and return with StatusCreated
	})

	resourcesCRUDRoutes.PUT("/", func(eCtx echo.Context) error {
		resource := &models.Resource{}
		if err := eCtx.Bind(resource); err != nil {
			return err
		}

		if err := eCtx.Validate(resource); err != nil {
			return err
		}

		err := c.svc.UpdateResource(eCtx.Request().Context(), resource)
		if err != nil {
			return err
		}

		return eCtx.NoContent(http.StatusCreated) // todo: updated Resource should be returned
	})

	resourcesCRUDRoutes.DELETE("/", func(eCtx echo.Context) error {
		req := &httpModels.DeleteResourceRequest{}
		if err := eCtx.Bind(req); err != nil {
			return err
		}

		if err := eCtx.Validate(req); err != nil {
			return err
		}

		err := c.svc.DeleteResource(eCtx.Request().Context(), req)
		if err != nil {
			return err
		}

		return eCtx.NoContent(http.StatusOK)
	})

	categoryRoutes := apiRoutes.Group("/categories")
	categoryRoutes.GET("/", func(eCtx echo.Context) error {
		resp, err := c.svc.GetCategories(eCtx.Request().Context())
		if err != nil {
			return err
		}

		return eCtx.JSON(http.StatusOK, httpModels.ResponseData{Data: resp})
	})
}
