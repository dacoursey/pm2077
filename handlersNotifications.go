package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/dacoursey/pm2077/models"
	"github.com/gorilla/mux"
)

/////
// System Notification Section
/////

// Retrieve all System Notification from the db.
func (a *App) getSystemNotifications(w http.ResponseWriter, r *http.Request) {

	SystemNotifications, err := models.GetNotifications(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := struct {
		Title       string
		Objects     []models.SystemNotification
		Active      string
		HeaderTitle string
		UserRole    string
	}{
		Title:       "System Notification List",
		Objects:     SystemNotifications,
		Active:      "Settings",
		HeaderTitle: "System Notification",
		UserRole:    a.getUserRole(r),
	}

	getTemplate(w, r, "viewAllSystemNotifications", data)
}

// Retrieve one System Notification from the db for edit
func (a *App) editSystemNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid System Notification ID")
		return
	}

	n := models.SystemNotification{ID: int32(id)}
	if err := n.GetNotificationByID(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "System Notification not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	data := struct {
		Title       string
		Object      models.SystemNotification
		Active      string
		HeaderTitle string
		UserRole    string
	}{
		Title: "System Notification - Edit",
		Object: models.SystemNotification{
			ID:      n.ID,
			Title:   n.Title,
			Message: n.Message,
		},
		Active:      "Settings",
		HeaderTitle: "System Notification",
		UserRole:    a.getUserRole(r),
	}

	getTemplate(w, r, "editSystemNotification", data)
}

// createSystemNotification will setup template for new notification
func (a *App) createSystemNotification(w http.ResponseWriter, r *http.Request) {

	m := models.SystemNotification{ID: 0, Title: "", Message: ""}

	data := struct {
		Title       string
		Object      models.SystemNotification
		Active      string
		HeaderTitle string
		UserRole    string
	}{
		Title:       "System Notification - Create",
		Object:      m,
		Active:      "Settings",
		HeaderTitle: "System Notification",
		UserRole:    a.getUserRole(r),
	}

	getTemplate(w, r, "editSystemNotification", data)

}

// saveSystemNotification will update or add a system notification record
func (a *App) saveSystemNotification(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	title := r.FormValue("title")
	message := r.FormValue("message")

	i, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid System Notification ID")
		return
	}
	notificationID := int32(i)

	m := models.SystemNotification{ID: notificationID, Title: title, Message: message}

	if notificationID > 0 {
		if err := m.UpdateNotificationByID(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		if err := m.CreateNotification(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	http.Redirect(w, r, "/notifications", http.StatusFound)

}

// deleteSystemNotification a System Notification by the ID
func (a *App) deleteSystemNotification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid System Notification ID")
		return
	}

	m := models.SystemNotification{ID: int32(id)}

	if err := m.DeleteNotificationByID(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, "/notifications", http.StatusFound)

}
