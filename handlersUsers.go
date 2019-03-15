package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/dacoursey/pm2077/models"
	"github.com/gorilla/mux"
)

/////
// Application Users
/////

// getUsers will get all app users from the db.
func (a *App) getUsers(w http.ResponseWriter, r *http.Request) {

	if a.getUserRole(r) != "Admin" {
		http.NotFound(w, r)
		return
	}

	Users, err := models.GetUsers(a.DB, 0, 999)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := struct {
		Title       string
		Objects     []models.User
		Active      string
		HeaderTitle string
		UserRole    string
	}{
		Title:       "User List",
		Objects:     Users,
		Active:      "Settings",
		HeaderTitle: "Applicaiton Users",
		UserRole:    a.getUserRole(r),
	}

	getTemplate(w, r, "viewAllUsers", data)
}

// createUser will handle the setting up the template for creating a new user
func (a *App) createUser(w http.ResponseWriter, r *http.Request) {

	if a.getUserRole(r) != "Admin" {
		http.NotFound(w, r)
		return
	}

	User := models.User{ID: 0, Username: "", Password: "", Role: ""}

	Roles, err := models.GetAppRoles(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := struct {
		Title       string
		Active      string
		HeaderTitle string
		UserRole    string
		User        models.User
		Roles       []models.Role
	}{
		Title:       "User - Create",
		Active:      "Settings",
		HeaderTitle: "Applicaiton Users",
		UserRole:    a.getUserRole(r),
		User:        User,
		Roles:       Roles,
	}

	getTemplate(w, r, "saveUsers", data)
}

// editUser will handle the setting up the template for editing a new user
func (a *App) editUser(w http.ResponseWriter, r *http.Request) {

	if a.getUserRole(r) != "Admin" {
		http.NotFound(w, r)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Applicaiton User ID")
		return
	}

	User := models.User{ID: int(id)}
	if err := User.GetUser(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Application User not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	Roles, err := models.GetAppRoles(a.DB)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := struct {
		Title       string
		Active      string
		HeaderTitle string
		UserRole    string
		User        models.User
		Roles       []models.Role
	}{
		Title:       "User - Edit",
		Active:      "Settings",
		HeaderTitle: "Applicaiton Users",
		UserRole:    a.getUserRole(r),
		User:        User,
		Roles:       Roles,
	}

	getTemplate(w, r, "saveUsers", data)
}

// saveUser will update or add a application user
func (a *App) saveUser(w http.ResponseWriter, r *http.Request) {

	if a.getUserRole(r) != "Admin" {
		http.NotFound(w, r)
		return
	}

	id := r.FormValue("id")
	userName := r.FormValue("Username")
	password := r.FormValue("Password")
	roleID := r.FormValue("RoleID")

	err := validatePassword(password)
	if err != nil {
		msg := ""
		switch err.Error() {
		case "pwdShort":
			msg = "Password%20length%20rules%20not%20met,%20please%20try%20again."
		case "pwdNotComplex":
			msg = "Password complexity rules not met, please try again."
		}

		if msg == "" {
			http.Redirect(w, r, "/user/"+id, http.StatusFound)
		} else {
			http.Redirect(w, r, "/users?err="+msg, http.StatusFound)
		}

		return
	}

	userID, err := strconv.ParseInt(id, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Applicaiton User ID")
		return
	}

	m := models.User{ID: int(userID), Username: userName, Password: password, Role: roleID}

	if userID > 0 {
		if err := m.UpdateUser(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		if err := m.CreateUser(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	http.Redirect(w, r, "/users", http.StatusFound)

}

func (a *App) deleteUser(w http.ResponseWriter, r *http.Request) {

	if a.getUserRole(r) != "Admin" {
		http.NotFound(w, r)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Application User ID")
		return
	}

	m := models.User{ID: int(id)}

	if err := m.DeleteUser(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, "/users", http.StatusFound)

}
