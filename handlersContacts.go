package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/dacoursey/pm2077/models"
	"github.com/gorilla/mux"
)

/////
// Contacts Section
/////

// Retrieve all contacts from the db.
func (a *App) getContacts(w http.ResponseWriter, r *http.Request) {
	// REMOVE ALL OF THIS WHEN WRAPPER AUTH IS WORKING
	// Check valid session
	// err := getSession(w, r)
	// if err != nil {
	// 	respondWithError(w, http.StatusInternalServerError, err.Error())
	// 	return
	// }

	count, err := strconv.Atoi(r.FormValue("count"))
	start, err := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	Contacts, err := models.GetContacts(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := struct {
		Title       string
		Objects     []models.Contact
		Active      string
		HeaderTitle string
		UserRole    string
	}{
		Title:       "Contacts List",
		Objects:     Contacts,
		Active:      "Contact",
		HeaderTitle: "Contact",
		UserRole:    a.getUserRole(r),
	}

	getTemplate(w, r, "viewAllContacts", data)
}

// Retrieve one contact from the db.
func (a *App) getContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid contact ID")
		return
	}

	p := models.Contact{ID: id}
	if err := p.GetContact(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Contact not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	data := struct {
		Title       string
		Object      models.Contact
		Active      string
		HeaderTitle string
		UserRole    string
	}{
		Title: "Contact Details",
		Object: models.Contact{
			ID:     p.ID,
			CustID: p.CustID,
			FName:  p.FName,
			LName:  p.LName,
			Email:  p.Email,
		},
		Active:      "Contact",
		HeaderTitle: "Contact",
		UserRole:    a.getUserRole(r),
	}

	getTemplate(w, r, "viewContact", data)
}

// Edit the currently selected contact.
func (a *App) editContact(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Contact ID")
		return
	}

	c := models.Contact{ID: int(id)}
	if err := c.GetContact(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Contact not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	data := struct {
		Title       string
		Object      models.Contact
		Active      string
		HeaderTitle string
		UserRole    string
	}{
		Title: "Project - Edit",
		Object: models.Contact{
			ID:     c.ID,
			CustID: c.CustID,
			FName:  c.FName,
			LName:  c.LName,
			Email:  c.Email,
		},
		Active:      "Contacts",
		HeaderTitle: "Contacts",
		UserRole:    a.getUserRole(r),
	}

	getTemplate(w, r, "editContact", data)
}

// Create a new contact in the database.
func (a *App) createContact(w http.ResponseWriter, r *http.Request) {

	c := models.Contact{
		ID:     0,
		CustID: 0,
		FName:  "",
		LName:  "",
		Email:  "",
	}

	data := struct {
		Title       string
		Object      models.Contact
		Active      string
		HeaderTitle string
		UserRole    string
	}{
		Title:       "Contact - Create",
		Object:      c,
		Active:      "Contacts",
		HeaderTitle: "Contacts",
		UserRole:    a.getUserRole(r),
	}

	getTemplate(w, r, "editContact", data)

}

// Save the current contact back to the database.
func (a *App) saveContact(w http.ResponseWriter, r *http.Request) {

	id := r.FormValue("id")
	cID := r.FormValue("custid")
	fname := r.FormValue("fname")
	lname := r.FormValue("lname")
	email := r.FormValue("email")

	i, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Contact ID")
		return
	}
	contactID := int(i)

	i, err = strconv.ParseInt(cID, 10, 0)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Customer ID")
		return
	}
	custID := int(i)

	c := models.Contact{ID: contactID, CustID: custID, FName: fname, LName: lname, Email: email}

	if contactID > 0 {
		if err := c.UpdateContact(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		if err := c.CreateContact(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	http.Redirect(w, r, "/contacts", http.StatusFound)

}

// Delete a contact from the database.
func (a *App) deleteContact(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Contact ID")
		return
	}

	m := models.Contact{ID: int(id)}

	if err := m.DeleteContact(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, "/contacts", http.StatusFound)
}
