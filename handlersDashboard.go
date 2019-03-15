package main

import (
	"net/http"

	"github.com/dacoursey/pm2077/models"
)

/////
// Dashboard Section
/////

// Retrieve the Dashboard
func (a *App) getDashboard(w http.ResponseWriter, r *http.Request) {

	data := struct {
		Title       string
		Active      string
		HeaderTitle string
		UserRole    string
	}{
		Title:       "Dashboard",
		Active:      "Dashboard",
		HeaderTitle: "Project Dashboard",
		UserRole:    a.getUserRole(r),
	}

	getTemplate(w, r, "dashboard", data)
}

// getDashboardNotifications will return all system notifications in JSON format
func (a *App) getDashboardNotifications(w http.ResponseWriter, r *http.Request) {

	SystemNotifications, err := models.GetNotifications(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := struct {
		Objects []models.SystemNotification
	}{
		Objects: SystemNotifications,
	}

	respondWithJSON(w, 200, data)
}

// getDashboardNumberOfProjects will retrieve the number of projects from the db
func (a *App) getDashboardNumberOfProjects(w http.ResponseWriter, r *http.Request) {

	Projects, err := models.GetProjects(a.DB, 0, 999)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type ProjectCount struct {
		Count int `json:"count"`
	}

	data := ProjectCount{Count: len(Projects)}

	respondWithJSON(w, 200, data)

}

// getDashboardHappyCustomer will retrieve the number of happy customers
func (a *App) getDashboardHappyCustomer(w http.ResponseWriter, r *http.Request) {

	Customers, err := models.GetCustomers(a.DB, 0, 999)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type HappyCount struct {
		Count int `json:"count"`
	}

	data := HappyCount{Count: len(Customers)}

	respondWithJSON(w, 200, data)

}

// getDashboardCompletedTask will retrieve the number of completed task
func (a *App) getDashboardCompletedTask(w http.ResponseWriter, r *http.Request) {

	UserTasks, err := models.GetCompletedTasks(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	type TaskCount struct {
		Count int `json:"count"`
	}

	data := TaskCount{Count: len(UserTasks)}

	respondWithJSON(w, 200, data)

}

// getDashboardResources will get all resources from the db.
func (a *App) getDashboardResources(w http.ResponseWriter, r *http.Request) {

	Resources, err := models.GetResources(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := struct {
		Objects []models.Resource
	}{
		Objects: Resources,
	}

	respondWithJSON(w, 200, data)
}

// getDashboardProjectTasksForUser will get all project tasks for a particular user
func (a *App) getDashboardProjectTasksForUser(w http.ResponseWriter, r *http.Request) {

	var userID int

	session, err := store.Get(r, "gonv")
	if err != nil {
		userID = 0
	} else {
		userID = session.Values["userID"].(int)
	}

	Task := models.ProjectTask{UserID: userID}

	UserTasks, err := Task.GetTasksByUserID(a.DB, 0)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := struct {
		Tasks []models.ProjectTask
	}{
		Tasks: UserTasks,
	}

	respondWithJSON(w, 200, data)
}
