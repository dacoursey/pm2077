package main

import "net/http"

/////
// Settings Page
/////

// Retrieve the SettingsPage
func (a *App) getSettings(w http.ResponseWriter, r *http.Request) {

	data := struct {
		Title       string
		Active      string
		HeaderTitle string
		UserRole    string
	}{
		Title:       "System Settings",
		Active:      "Settings",
		HeaderTitle: "Settings",
		UserRole:    a.getUserRole(r),
	}

	getTemplate(w, r, "settings", data)
}
