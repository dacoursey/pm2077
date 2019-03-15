package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/dacoursey/pm2077/models"
	"github.com/gorilla/mux"
)

/////
// Project Section
/////

// Retrieve all projects from the db.
func (a *App) getProjects(w http.ResponseWriter, r *http.Request) {
	count, err := strconv.Atoi(r.FormValue("count"))
	start, err := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	Projects, err := models.GetProjects(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := struct {
		Title       string
		Objects     []models.Project
		Active      string
		HeaderTitle string
		UserRole    string
	}{
		Title:       "Project List",
		Objects:     Projects,
		Active:      "Project",
		HeaderTitle: "Project",
		UserRole:    a.getUserRole(r),
	}

	getTemplate(w, r, "viewAllProjects", data)
}

// Retrieve one project from the db.
func (a *App) getProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid project ID")
		return
	}

	p := models.Project{ID: id}
	if err := p.GetProject(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Project not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	data := struct {
		Title       string
		Object      models.Project
		Active      string
		HeaderTitle string
		UserRole    string
	}{
		Title: "Project Details",
		Object: models.Project{
			ID:        p.ID,
			CustID:    p.CustID,
			Name:      p.Name,
			StartDate: p.StartDate,
			Hours:     p.Hours,
		},
		Active:      "Project",
		HeaderTitle: "Project",
		UserRole:    a.getUserRole(r),
	}

	getTemplate(w, r, "viewProject", data)
}

// Edit the currently seledted project.
func (a *App) editProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Project ID")
		return
	}

	p := models.Project{ID: int(id)}
	if err := p.GetProject(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Project not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	data := struct {
		Title       string
		Object      models.Project
		Active      string
		HeaderTitle string
		UserRole    string
	}{
		Title: "Project - Edit",
		Object: models.Project{
			ID:        p.ID,
			CustID:    p.CustID,
			Name:      p.Name,
			StartDate: p.StartDate,
			Hours:     p.Hours,
		},
		Active:      "Projects",
		HeaderTitle: "Project",
		UserRole:    a.getUserRole(r),
	}

	getTemplate(w, r, "editProject", data)
}

// Create a new project in the database.
func (a *App) createProject(w http.ResponseWriter, r *http.Request) {

	m := models.Project{
		ID:        0,
		CustID:    0,
		Name:      "",
		StartDate: "",
		Hours:     0,
	}

	data := struct {
		Title       string
		Object      models.Project
		Active      string
		HeaderTitle string
		UserRole    string
	}{
		Title:       "Project - Create",
		Object:      m,
		Active:      "Projects",
		HeaderTitle: "Project",
		UserRole:    a.getUserRole(r),
	}

	getTemplate(w, r, "editProject", data)

}

// Save the current project back to the database.
func (a *App) saveProject(w http.ResponseWriter, r *http.Request) {

	id := r.FormValue("id")
	cID := r.FormValue("custid")
	name := r.FormValue("name")
	startDate := r.FormValue("startdate")
	h := r.FormValue("hours")

	i, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Project ID")
		return
	}
	projectID := int(i)

	i, err = strconv.ParseInt(cID, 10, 0)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Customer ID")
		return
	}
	custID := int(i)

	i, err = strconv.ParseInt(h, 10, 0)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Hours")
		return
	}
	hours := int(i)

	m := models.Project{ID: projectID, CustID: custID, Name: name, StartDate: startDate, Hours: hours}

	if projectID > 0 {
		if err := m.UpdateProject(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		if err := m.CreateProject(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	http.Redirect(w, r, "/projects", http.StatusFound)

}

// Delete a project from the database.
func (a *App) deleteProject(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Project ID")
		return
	}

	m := models.Project{ID: int(id)}

	if err := m.DeleteProject(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, "/projects", http.StatusFound)
}
