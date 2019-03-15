package main

import (
	"database/sql"
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/dacoursey/pm2077/models"
	"github.com/gorilla/mux"
)

/////
// Project Task
/////

// getProjectTasks will get all project tasks from the db.
func (a *App) getProjectTasks(w http.ResponseWriter, r *http.Request) {

	if a.getUserRole(r) != "PM" {
		http.NotFound(w, r)
		return
	}

	ProjectTask, err := models.GetTasks(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := struct {
		Title       string
		Active      string
		HeaderTitle string
		UserRole    string
		Objects     []models.ProjectTask
	}{
		Title:       "Project Task List",
		Active:      "Task",
		HeaderTitle: "Project Task",
		UserRole:    a.getUserRole(r),
		Objects:     ProjectTask,
	}

	getTemplate(w, r, "viewAllProjectTask", data)
}

// getUserTasks will get all project tasks from the db.
func (a *App) getUserTasks(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid User ID")
		return
	}

	Task := models.ProjectTask{UserID: int(id)}
	ProjectTask, err := Task.GetTasksForUser(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := struct {
		Title       string
		Active      string
		HeaderTitle string
		UserRole    string
		Objects     []models.ProjectTask
	}{
		Title:       "Your Open and Closed Task",
		Active:      "Dashboard",
		HeaderTitle: "Project Task",
		UserRole:    a.getUserRole(r),
		Objects:     ProjectTask,
	}

	getTemplate(w, r, "viewAllUserTasks", data)
}

// createProjectTask will handle the setting up the template for creating a new task
func (a *App) createProjectTask(w http.ResponseWriter, r *http.Request) {

	if a.getUserRole(r) != "PM" {
		http.NotFound(w, r)
		return
	}

	Task := models.ProjectTask{ID: 0, UserID: 0, ProjectID: 0}

	Projects, err := models.GetProjects(a.DB, 0, 999)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Users, err := models.GetUsers(a.DB, 0, 999)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := struct {
		Title       string
		Active      string
		HeaderTitle string
		UserRole    string
		Task        models.ProjectTask
		Users       []models.User
		Projects    []models.Project
		UserID      int
	}{
		Title:       "Project Task - Create",
		Active:      "Task",
		HeaderTitle: "Project Task",
		UserRole:    a.getUserRole(r),
		Task:        Task,
		Users:       Users,
		Projects:    Projects,
		UserID:      3,
	}

	getTemplate(w, r, "saveProjectTask", data)
}

// editProjectTask will handle the setting up the template for updating a new task
func (a *App) editProjectTask(w http.ResponseWriter, r *http.Request) {

	if a.getUserRole(r) != "PM" {
		http.NotFound(w, r)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Project Task ID")
		return
	}

	Task := models.ProjectTask{ID: int(id)}
	if err := Task.GetTaskByID(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Project Task not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	Projects, err := models.GetProjects(a.DB, 0, 999)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	Users, err := models.GetUsers(a.DB, 0, 999)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := struct {
		Title       string
		Active      string
		HeaderTitle string
		UserRole    string
		Task        models.ProjectTask
		Users       []models.User
		Projects    []models.Project
		UserID      int
	}{
		Title:       "Project Task - Create",
		Active:      "Task",
		HeaderTitle: "Project Task",
		UserRole:    a.getUserRole(r),
		Task:        Task,
		Users:       Users,
		Projects:    Projects,
		UserID:      3,
	}

	getTemplate(w, r, "saveProjectTask", data)
}

// saveTask will update or add a project task
func (a *App) saveProjectTask(w http.ResponseWriter, r *http.Request) {

	if a.getUserRole(r) != "PM" {
		http.NotFound(w, r)
		return
	}

	tenMB := int64(10000000)
	r.ParseMultipartForm(tenMB)

	id := r.FormValue("id")
	task := r.FormValue("task")
	userID := r.FormValue("userId")
	projectID := r.FormValue("projectId")
	isCompleted := r.FormValue("isCompleted")
	file, handler, err := r.FormFile("uploadfile")

	b64Filename := ""
	if err == nil {
		filename := "./uploads/" + handler.Filename
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		defer file.Close()

		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		defer f.Close()
		io.Copy(f, file)

		barrayFilename := []byte(filename)
		b64Filename = base64.StdEncoding.EncodeToString(barrayFilename)
	}

	taskID, _ := strconv.ParseInt(id, 10, 32)
	uid, _ := strconv.ParseInt(userID, 10, 32)
	pid, _ := strconv.ParseInt(projectID, 10, 32)

	var comp int
	if isCompleted == "on" {
		comp = 1
	} else {
		comp = 0
	}

	m := models.ProjectTask{ID: int(taskID), Title: task, UserID: int(uid), ProjectID: int(pid), IsCompleted: comp, AttachmentPath: b64Filename}

	if taskID > 0 {
		if err := m.UpdateTask(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		if err := m.CreateTask(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	http.Redirect(w, r, "/tasks", http.StatusFound)

}

// deleteProjectTask handles the deletion of a project task
func (a *App) deleteProjectTask(w http.ResponseWriter, r *http.Request) {

	if a.getUserRole(r) != "PM" {
		http.NotFound(w, r)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Project Task ID")
		return
	}

	m := models.ProjectTask{ID: int(id)}

	if err := m.DeleteTaskByID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, "/tasks", http.StatusFound)

}

func (a *App) getAttachment(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["path"]

	if !ok || len(keys) < 1 {
		http.Redirect(w, r, "/tasks", http.StatusFound)
		return
	}

	b64Path := string(keys[0])
	barrayPath, err := base64.StdEncoding.DecodeString(b64Path)
	path := string(barrayPath)
	if err != nil {
		http.Redirect(w, r, "/tasks", http.StatusFound)
	}

	http.ServeFile(w, r, path)
}
