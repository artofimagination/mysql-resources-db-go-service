package restcontrollers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/artofimagination/mysql-resources-db-go-service/dbcontrollers"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type RESTController struct {
	DBController *dbcontrollers.MYSQLController
}

const (
	ResourceURIAdd                   = "/add-resource"
	ResourceURIGetByID               = "/get-resource-by-id"
	ResourceURIGetMultipleByIDs      = "/get-resources-by-ids"
	ResourceURIGetMultipleByCategory = "/get-resources-by-category"
	ResourceURIUpdate                = "/update-resource"
	ResourceURIDelete                = "/delete-resource"
	ResourceURIGetCategories         = "/get-categories"
)

var DataOK = "OK"

type ResponseWriter struct {
	http.ResponseWriter
}

type Request struct {
	*http.Request
}

type ResponseData struct {
	Error string      `json:"error" validation:"required"`
	Data  interface{} `json:"data" validation:"required"`
}

func (w ResponseWriter) writeError(message string, statusCode int) {
	response := &ResponseData{
		Error: message,
	}
	w.writeResponse(response, statusCode)
}

func (w ResponseWriter) writeData(data interface{}, statusCode int) {
	response := &ResponseData{
		Data: data,
	}
	w.writeResponse(response, statusCode)
}

func (w ResponseWriter) writeResponse(response *ResponseData, statusCode int) {
	b, err := json.Marshal(response)
	if err != nil {
		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	fmt.Fprint(w, string(b))
}

func checkRequestType(requestTypeString string, w ResponseWriter, r *Request) error {
	if r.Method != requestTypeString {
		w.WriteHeader(http.StatusBadRequest)
		errorString := fmt.Sprintf("Invalid request type %s", r.Method)
		return errors.New(errorString)
	}
	return nil
}

func decodePostData(w ResponseWriter, r *Request) ([]byte, error) {
	if err := checkRequestType(http.MethodPost, w, r); err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = errors.Wrap(errors.WithStack(err), "Failed to decode request json")
		return nil, err
	}

	return data, nil
}

func parseIDList(r *Request) ([]uuid.UUID, error) {
	ids, ok := r.URL.Query()["ids"]
	if !ok || len(ids[0]) < 1 {
		return nil, errors.New("Missing 'ids'")
	}

	idList := make([]uuid.UUID, 0)
	for _, idString := range ids {
		id, err := uuid.Parse(idString)
		if err != nil {
			return nil, errors.New("Invalid 'ids'")
		}
		idList = append(idList, id)
	}

	return idList, nil
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hi! I am a user database server!")
}

func makeHandler(fn func(ResponseWriter, *Request)) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		r := &Request{request}
		w := ResponseWriter{writer}
		fn(w, r)
	}
}

func NewRESTController() (*RESTController, error) {
	dbController, err := dbcontrollers.NewDBController()
	if err != nil {
		return nil, err
	}

	restController := &RESTController{
		DBController: dbController,
	}

	http.HandleFunc("/", sayHello)
	http.HandleFunc(ResourceURIAdd, makeHandler(restController.addResource))
	http.HandleFunc(ResourceURIGetByID, makeHandler(restController.getResourceByID))
	http.HandleFunc(ResourceURIGetMultipleByIDs, makeHandler(restController.getResourcesByIDs))
	http.HandleFunc(ResourceURIGetMultipleByCategory, makeHandler(restController.getResourcesByCategory))
	http.HandleFunc(ResourceURIUpdate, makeHandler(restController.updateResource))
	http.HandleFunc(ResourceURIDelete, makeHandler(restController.deleteResource))
	http.HandleFunc(ResourceURIGetCategories, makeHandler(restController.getCategories))

	return restController, nil
}
