package main

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/dacoursey/pm2077/models"
	"github.com/gorilla/sessions"
)

/////
// Auth Section
/////

var store = sessions.NewCookieStore([]byte("our-little-secret"))

func authn(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	fn := func(w http.ResponseWriter, r *http.Request) {
		redirect := r.URL.RequestURI()
		session, err := store.Get(r, "gonv")

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if val, ok := session.Values["sessionid"].(string); ok {
			// if val is a string
			switch val {
			case "":
				http.Redirect(w, r, "/login?redirect="+redirect, http.StatusFound)
			}
		} else {
			// if val is not a string type
			http.Redirect(w, r, "/login?redirect="+redirect, http.StatusFound)
		}

		f(w, r)
	}

	return http.HandlerFunc(fn)
}

// Retrieve login document.
func (a *App) getLogin(w http.ResponseWriter, r *http.Request) {
	redirect := r.URL.Query().Get("redirect")

	data := struct {
		Redirect string
	}{
		Redirect: redirect,
	}

	getTemplate(w, r, "login", data)
}

// Process login form submission.
func (a *App) processLogin(w http.ResponseWriter, r *http.Request) {
	redirect := r.URL.Query().Get("redirect")
	user := r.FormValue("username")
	pass := r.FormValue("password")

	var errorMessage string

	u := models.User{Username: user, Password: pass}

	// Get the user account by user name
	if err := u.GetUserByCreds(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			errorMessage = "Invalid user name or password.  Are you sure you have an account with us?"
		default:
			errorMessage = "Ouch." + err.Error()
		}
	} else {
		if u.Password != pass {
			errorMessage = "The password you entered is incorrect.  Please try again."
		}
	}

	if len(errorMessage) > 0 {
		http.Redirect(w, r, "/login?err="+errorMessage+"&redirect="+redirect, http.StatusFound)
	} else {
		a.logUserIn(w, r, u)
	}
}

// Set up the session and role cookie for the user
func (a *App) logUserIn(w http.ResponseWriter, r *http.Request, u models.User) {

	// err := writeSession(u, w, r)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// Get and set the user's role
	switch u.Role {
	case "Admin":
		a.setRoleCookie(w, "Admin")
	case "Project Manager":
		a.setRoleCookie(w, "PM")
	case "User":
		a.setRoleCookie(w, "Users")
	default:
		a.setRoleCookie(w, "User")
	}

	// Set Session User ID
	session, err := store.Get(r, "gonv")
	session.Values["sessionid"] = "weak"
	session.Values["userID"] = u.ID
	session.Save(r, w)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	redirect := r.URL.Query().Get("redirect")
	if redirect == "" {
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	} else {
		http.Redirect(w, r, redirect, http.StatusFound)
	}
}

// Handles the user logout
func (a *App) processLogout(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}

// Pull the cookie store to investigate the presence of a valid session.
func getSession(w http.ResponseWriter, r *http.Request) error {
	session, err := store.Get(r, "gonv")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}

	if val, ok := session.Values["sessionid"].(string); ok {
		// if val is a string
		switch val {
		case "":
			http.Redirect(w, r, "/login", http.StatusFound)
		}
	} else {
		// if val is not a string type
		http.Redirect(w, r, "/login", http.StatusFound)
	}

	return err
}

// Write a new session and put it in the cookie store.
// func writeSession(p models.User, w http.ResponseWriter, r *http.Request) error {
// 	session, err := store.Get(r, "gonv")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return err
// 	}

// 	// Set some session values.
// 	session.Values["sessionid"] = "weak"

// 	// Save it before we write to the response/return from the handler.
// 	err = session.Save(r, w)
// 	if err != nil {
// 		respondWithError(w, http.StatusInternalServerError, err.Error())
// 	}

// 	return err
// }

/////
// Cookies
/////

func (a *App) setRoleCookie(w http.ResponseWriter, role string) {

	expiration := time.Now().Add(365 * 24 * time.Hour)

	cookie := http.Cookie{Name: "DesalinatorRole", Value: role, Expires: expiration}
	http.SetCookie(w, &cookie)

}

func (a *App) getUserRole(r *http.Request) string {
	cookie, err := r.Cookie("DesalinatorRole")

	if err != nil {
		return "guest"
	}

	return cookie.Value
}
