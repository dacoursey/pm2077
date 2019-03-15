package main

import (
	"net/http"
    "os/exec"
	"strings"
	"html/template"
)

/////
// Settings Backup Page
/////

// Retrieve the Settings Backup Page
func (a *App) getBackup(w http.ResponseWriter, r *http.Request) {

	if a.getUserRole(r) != "Admin" {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Title       string
		Active      string
		HeaderTitle string
		UserRole    string
		Output		template.HTML
	}{
		Title:       "System Backup",
		Active:      "Settings",
		HeaderTitle: "System Maintenance",
		UserRole:    a.getUserRole(r),
		Output:		 template.HTML("<br \\>"),
	}

	getTemplate(w, r, "systemBackup", data)
}

// Starts the backup process
func (a *App) startBackup(w http.ResponseWriter, r *http.Request) {
	
	if a.getUserRole(r) != "Admin" {
		http.NotFound(w, r)
		return
	}

	backupType := r.FormValue("backupType")

	var output string
	
	out, err := exec.Command("bash", "-c", backupType).Output()
	if err != nil {
		output = "Error:  " + err.Error()
	}else{
		output = string(out[:])
	}
	
	data := struct {
		Title       string
		Active      string
		HeaderTitle string
		UserRole    string
		Output		template.HTML
	}{
		Title:       "System Backup",
		Active:      "Settings",
		HeaderTitle: "System Maintenance",
		UserRole:    a.getUserRole(r),
		Output:		 template.HTML(strings.Replace(output,"\n","<br \\>",-1)),
	}

	getTemplate(w, r, "systemBackup", data)
}