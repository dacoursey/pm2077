package main

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/dacoursey/pm2077/models"
	"github.com/gorilla/mux"
)

/////
// Customer Section
/////

// Retrieve all customers from the db.
func (a *App) getCustomers(w http.ResponseWriter, r *http.Request) {

	count, err := strconv.Atoi(r.FormValue("count"))
	start, err := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	Customers, err := models.GetCustomers(a.DB, start, count)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	data := struct {
		Title       string
		Objects     []models.Customer
		Active      string
		HeaderTitle string
		UserRole    string
	}{
		Title:       "Customers List",
		Objects:     Customers,
		Active:      "Customer",
		HeaderTitle: "Customer",
		UserRole:    a.getUserRole(r),
	}

	getTemplate(w, r, "viewAllCustomers", data)

}

// Retrieve one customer from the db.
func (a *App) getCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid customer ID")
		return
	}

	c := models.Customer{ID: id}
	if err := c.GetCustomer(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Customer not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	data := struct {
		Title       string
		Object      models.Customer
		Active      string
		HeaderTitle string
		UserRole    string
	}{
		Title: "Customer Profile",
		Object: models.Customer{
			ID:   c.ID,
			Name: c.Name,
			POC:  c.POC,
		},
		Active:      "Customer",
		HeaderTitle: "Customer",
		UserRole:    a.getUserRole(r),
	}

	getTemplate(w, r, "viewCustomer", data)

}

// Edit the currently selected customer.
func (a *App) editCustomer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Customer ID")
		return
	}

	c := models.Customer{ID: int(id)}
	if err := c.GetCustomer(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Customer not found")
		default:
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	data := struct {
		Title       string
		Object      models.Customer
		Active      string
		HeaderTitle string
		UserRole    string
	}{
		Title: "Project - Edit",
		Object: models.Customer{
			ID:   c.ID,
			Name: c.Name,
			POC:  c.POC,
		},
		Active:      "Customers",
		HeaderTitle: "Customers",
		UserRole:    a.getUserRole(r),
	}

	getTemplate(w, r, "editCustomer", data)
}

// Create a new customer in the database.
func (a *App) createCustomer(w http.ResponseWriter, r *http.Request) {

	m := models.Customer{ID: 0, Name: "", POC: 0}

	data := struct {
		Title       string
		Object      models.Customer
		Active      string
		HeaderTitle string
		UserRole    string
	}{
		Title:       "Customer - Create",
		Object:      m,
		Active:      "Customers",
		HeaderTitle: "Customers",
		UserRole:    a.getUserRole(r),
	}

	getTemplate(w, r, "editCustomer", data)
}

// Save the current customer back to the database.
func (a *App) saveCustomer(w http.ResponseWriter, r *http.Request) {
	id := r.FormValue("id")
	name := r.FormValue("name")
	poc := r.FormValue("poc")

	i, err := strconv.ParseInt(id, 10, 0)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Customer ID")
		return
	}
	customerID := int(i)

	p, err := strconv.ParseInt(poc, 10, 0)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Customer POC")
		return
	}
	pocID := int(p)

	c := models.Customer{ID: customerID, Name: name, POC: pocID}

	if customerID > 0 {
		if err := c.UpdateCustomer(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		if err := c.CreateCustomer(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	http.Redirect(w, r, "/customers", http.StatusFound)
}

// Delete a customer from the database.
func (a *App) deleteCustomer(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Customer ID")
		return
	}

	m := models.Customer{ID: int(id)}

	if err := m.DeleteCustomer(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	http.Redirect(w, r, "/customers", http.StatusFound)

}
