package restcontrollers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/artofimagination/mysql-resources-db-go-service/dbcontrollers"
	"github.com/artofimagination/mysql-resources-db-go-service/models"
	"github.com/google/uuid"
)

func (c *RESTController) getResourceByID(w ResponseWriter, r *Request) {
	log.Println("Getting resource by id")
	if err := checkRequestType(http.MethodGet, w, r); err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	ids, ok := r.URL.Query()["id"]
	if !ok || len(ids[0]) < 1 {
		w.writeError("Url Param 'id' is missing", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(ids[0])
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	userData, err := c.DBController.GetResourceByID(&id)
	if err != nil {
		if err.Error() == dbcontrollers.ErrResourceNotFound.Error() {
			w.writeError(err.Error(), http.StatusAccepted)
			return
		}
		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}

	w.writeData(userData, http.StatusOK)
}

func (c *RESTController) addResource(w ResponseWriter, r *Request) {
	log.Println("Adding resource")
	dataBytes, err := decodePostData(w, r)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	resource := &models.Resource{}
	if err := json.Unmarshal(dataBytes, resource); err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	// Execute function
	if err := c.DBController.AddResource(resource); err != nil {
		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}

	w.writeData("OK", http.StatusCreated)
}

func (c *RESTController) updateResource(w ResponseWriter, r *Request) {
	log.Println("Updating resource")
	dataBytes, err := decodePostData(w, r)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	resource := &models.Resource{}
	if err := json.Unmarshal(dataBytes, resource); err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	// Execute function
	if err := c.DBController.UpdateResource(resource); err != nil {
		if err.Error() == dbcontrollers.ErrResourceNotFound.Error() {
			w.writeError(err.Error(), http.StatusAccepted)
			return
		}
		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}

	w.writeData("OK", http.StatusCreated)
}

func (c *RESTController) deleteResource(w ResponseWriter, r *Request) {
	log.Println("Deleting resource")
	dataBytes, err := decodePostData(w, r)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	resource := &models.Resource{}
	if err := json.Unmarshal(dataBytes, resource); err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	// Execute function
	if err := c.DBController.DeleteResource(resource); err != nil {
		if err.Error() == dbcontrollers.ErrResourceNotFound.Error() {
			w.writeError(err.Error(), http.StatusAccepted)
			return
		}
		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = c.DBController.GetResourceByID(&resource.ID)
	if err != nil && err.Error() != dbcontrollers.ErrResourceNotFound.Error() {
		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}

	w.writeData("OK", http.StatusCreated)
}

func (c *RESTController) getResourcesByIDs(w ResponseWriter, r *Request) {
	log.Println("Getting multiple resources by ids")
	if err := checkRequestType(http.MethodGet, w, r); err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	idList, err := parseIDList(r)
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	resources, err := c.DBController.GetResourcesByIDs(idList)
	if err != nil {
		if err.Error() == dbcontrollers.ErrResourceNotFound.Error() {
			w.writeError(err.Error(), http.StatusAccepted)
			return
		}
		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}
	w.writeData(resources, http.StatusOK)
}

func (c *RESTController) getResourcesByCategory(w ResponseWriter, r *Request) {
	log.Println("Getting multiple resources by category")
	if err := checkRequestType(http.MethodGet, w, r); err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	categories, ok := r.URL.Query()["category"]
	if !ok || len(categories[0]) < 1 {
		w.writeError("Url Param 'category' is missing", http.StatusBadRequest)
		return
	}

	category, err := strconv.Atoi(categories[0])
	if err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	resources, err := c.DBController.GetResourcesByCategory(category)
	if err != nil {
		if err.Error() == dbcontrollers.ErrResourceNotFound.Error() {
			w.writeError(err.Error(), http.StatusAccepted)
			return
		}
		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}

	w.writeData(resources, http.StatusOK)
}

func (c *RESTController) getCategories(w ResponseWriter, r *Request) {
	log.Println("Getting categories")
	if err := checkRequestType(http.MethodGet, w, r); err != nil {
		w.writeError(err.Error(), http.StatusBadRequest)
		return
	}

	categories, err := c.DBController.GetCategories()
	if err != nil {
		w.writeError(err.Error(), http.StatusInternalServerError)
		return
	}

	w.writeData(categories, http.StatusOK)
}
