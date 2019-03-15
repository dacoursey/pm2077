package main

//go:generate echo "Embedding file data."
//go:generate go-bindata -o assets/bindata.go -pkg assets -nomemcopy public/...

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	rice "github.com/GeertJohan/go.rice"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	_ "github.com/lib/pq"
)

// App for stuff.
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize for stuff.
func (a *App) Initialize(user, password, dbname string) {
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", user, password, dbname)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	store.Options = &sessions.Options{
		Domain: "localhost",
		Path:   "/",
		MaxAge: 3600 * 8, // 8 hours
		// HttpOnly: true, // disable for this demo
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()

	loadTemplates()
}

// Build up our routing table.
func (a *App) initializeRoutes() {
	// Root
	a.Router.HandleFunc("/", authn(a.getRoot)).Methods("GET")
	// AuthZ and AuthN
	a.Router.HandleFunc("/login", a.getLogin).Methods("GET")
	a.Router.HandleFunc("/login", a.processLogin).Methods("POST")
	a.Router.HandleFunc("/logout", a.processLogout).Methods("GET")
	// Images and stuff
	a.Router.PathPrefix("/resources/").Handler(http.StripPrefix("/resources/", http.FileServer(http.Dir("./resources/"))))
	// Contacts
	a.Router.HandleFunc("/contacts", authn(a.getContacts)).Methods("GET")
	a.Router.HandleFunc("/contact/{id:[0-9]+}", authn(a.editContact)).Methods("GET")
	a.Router.HandleFunc("/contact/create", authn(a.createContact)).Methods("GET")
	a.Router.HandleFunc("/contact/save", authn(a.saveContact)).Methods("POST")
	a.Router.HandleFunc("/contact/delete/{id:[0-9]+}", authn(a.deleteContact)).Methods("GET")
	// Customers
	a.Router.HandleFunc("/customers", authn(a.getCustomers)).Methods("GET")
	a.Router.HandleFunc("/customer/{id:[0-9]+}", authn(a.editCustomer)).Methods("GET")
	a.Router.HandleFunc("/customer/create", authn(a.createCustomer)).Methods("GET")
	a.Router.HandleFunc("/customer/save", authn(a.saveCustomer)).Methods("POST")
	a.Router.HandleFunc("/customer/delete/{id:[0-9]+}", a.deleteCustomer).Methods("GET")
	// Projects
	a.Router.HandleFunc("/projects", authn(a.getProjects)).Methods("GET")
	a.Router.HandleFunc("/project/{id:[0-9]+}", authn(a.editProject)).Methods("GET")
	a.Router.HandleFunc("/project/create", authn(a.createProject)).Methods("GET")
	a.Router.HandleFunc("/project/save", authn(a.saveProject)).Methods("POST")
	a.Router.HandleFunc("/project/delete/{id:[0-9]+}", authn(a.deleteProject)).Methods("GET")
	// Dashboard
	a.Router.HandleFunc("/dashboard", authn(a.getDashboard)).Methods("GET")
	a.Router.HandleFunc("/dashboard/notifications", authn(a.getDashboardNotifications)).Methods("GET")
	a.Router.HandleFunc("/dashboard/numberofprojects", authn(a.getDashboardNumberOfProjects)).Methods("GET")
	a.Router.HandleFunc("/dashboard/numberofhappy", authn(a.getDashboardHappyCustomer)).Methods("GET")
	a.Router.HandleFunc("/dashboard/completedtask", authn(a.getDashboardCompletedTask)).Methods("GET")
	a.Router.HandleFunc("/dashboard/resources", authn(a.getDashboardResources)).Methods("GET")
	a.Router.HandleFunc("/dashboard/tasks", authn(a.getDashboardProjectTasksForUser)).Methods("GET")
	// System Notification
	a.Router.HandleFunc("/notifications", authn(a.getSystemNotifications)).Methods("GET")
	a.Router.HandleFunc("/notification/{id:[0-9]+}", authn(a.editSystemNotification)).Methods("GET")
	a.Router.HandleFunc("/notification/create", authn(a.createSystemNotification)).Methods("GET")
	a.Router.HandleFunc("/notification/save", authn(a.saveSystemNotification)).Methods("POST")
	a.Router.HandleFunc("/notification/delete/{id:[0-9]+}", authn(a.deleteSystemNotification)).Methods("GET")
	// Internal Resources
	a.Router.HandleFunc("/resources", authn(a.getResources)).Methods("GET")
	a.Router.HandleFunc("/resource/{id:[0-9]+}", authn(a.editResource)).Methods("GET")
	a.Router.HandleFunc("/resource/create", authn(a.createResource)).Methods("GET")
	a.Router.HandleFunc("/resource/save", authn(a.saveResource)).Methods("POST")
	a.Router.HandleFunc("/resource/delete/{id:[0-9]+}", authn(a.deleteResource)).Methods("GET")
	// Project Task
	a.Router.HandleFunc("/tasks", authn(a.getProjectTasks)).Methods("GET")
	a.Router.HandleFunc("/task/{id:[0-9]+}", authn(a.editProjectTask)).Methods("GET")
	a.Router.HandleFunc("/task/create", authn(a.createProjectTask)).Methods("GET")
	a.Router.HandleFunc("/task/save", authn(a.saveProjectTask)).Methods("POST")
	a.Router.HandleFunc("/task/delete/{id:[0-9]+}", authn(a.deleteProjectTask)).Methods("GET")
	a.Router.HandleFunc("/task/attachment", authn(a.getAttachment)).Methods("GET")
	a.Router.HandleFunc("/mytask/{id:[0-9]+}", authn(a.getUserTasks)).Methods("GET")
	// Settings
	a.Router.HandleFunc("/settings", authn(a.getSettings)).Methods("GET")
	// System Backup
	a.Router.HandleFunc("/backup", authn(a.getBackup)).Methods("GET")
	a.Router.HandleFunc("/backup/start", authn(a.startBackup)).Methods("POST")
	// Application Users
	a.Router.HandleFunc("/users", authn(a.getUsers)).Methods("GET")
	a.Router.HandleFunc("/user/create", authn(a.createUser)).Methods("GET")
	a.Router.HandleFunc("/user/save", authn(a.saveUser)).Methods("POST")
	a.Router.HandleFunc("/user/{id:[0-9]+}", authn(a.editUser)).Methods("GET")
	a.Router.HandleFunc("/user/delete/{id:[0-9]+}", authn(a.deleteUser)).Methods("GET")
	// Static Files
	a.Router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(rice.MustFindBox("public").HTTPBox())))
}

// Run for stuff.
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

/////
// Root Section
/////

// Retrieve root document.
func (a *App) getRoot(w http.ResponseWriter, r *http.Request) {

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

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
