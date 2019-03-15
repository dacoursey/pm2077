package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/dacoursey/pm2077/models"
	"github.com/gorilla/mux"
)

/////
// Internal Resource Section
/////

// getResource will get all resources from the db.
func (a *App) getResources(w http.ResponseWriter, r *http.Request) {

	Resources, err := models.GetResources(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := struct {
		Title       string
		Objects     []models.Resource
		Active      string
		HeaderTitle string
		UserRole    string
	}{
		Title:       "Internal Resource List",
		Objects:     Resources,
		Active:      "Settings",
		HeaderTitle: "Internal Resource",
		UserRole:    a.getUserRole(r),
	}

	getTemplate(w, r, "viewAllResources", data)
}

// editResource will retrieve a Resource from the db for edit
func (a *App) editResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Resource ID")
		return
	}

	n := models.Resource{ID: int32(id)}
	if err := n.GetResourceByID(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Resource not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	data := struct {
		Title       string
		Object      models.Resource
		Active      string
		HeaderTitle string
		UserRole    string
	}{
		Title: "Internal Resource - Edit",
		Object: models.Resource{
			ID:    n.ID,
			Title: n.Title,
			URL:   n.URL,
		},
		Active:      "Settings",
		HeaderTitle: "Internal Resource",
		UserRole:    a.getUserRole(r),
	}

	getTemplate(w, r, "saveResources", data)
}

// createResource will setup template for new Resource
func (a *App) createResource(w http.ResponseWriter, r *http.Request) {

	m := models.Resource{ID: 0, Title: "", URL: ""}

	data := struct {
		Title       string
		Object      models.Resource
		Active      string
		HeaderTitle string
		UserRole    string
	}{
		Title:       "Internal Resource - Create",
		Object:      m,
		Active:      "Settings",
		HeaderTitle: "Internal Resource",
		UserRole:    a.getUserRole(r),
	}

	getTemplate(w, r, "saveResources", data)

}

// saveResource will update or add a Resource record
func (a *App) saveResource(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	url := r.FormValue("url")

	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Resource ID")
		return
	}
	resourceID := int32(i)

	m := models.Resource{ID: resourceID, Title: title, URL: url}

	if resourceID > 0 {
		if err := m.UpdateResourceByID(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		if err := m.CreateResource(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	http.Redirect(w, r, "/resources", http.StatusFound)

}

// deleteResource a Resource by the ID
func (a *App) deleteResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Resource ID")
		return
	}

	m := models.Resource{ID: int32(id)}

	if err := m.DeleteResourceByID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, "/resources", http.StatusFound)

}
