package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var templates = template.Must(template.ParseGlob("templates/*"))

const MagentoScan = "Magento Scanner"
const MagentoScanLogin = "Magento Scanner - Login Page"
const MagentoScanAdminLogin = "Magento Scanner - Admin Login Page"
const MagentoScanVersion = "Magento Scanner - Version"

type Page struct {
	Title    string
	Host     string
	Error    bool
	Username string
	Version  string
}

// IndexHandler provides static pages depending on the request. If
// magento_version is requested, return the configured magento_version_text and
// flag the IP. Otherwise, return the index page and check whether to record the
// http.Request.
func IndexHandler(app *App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		recordAttack(app, r, MagentoScan)
		host := fmt.Sprintf("http://%s", app.SensorIP)
		p := Page{
			Title: app.Config.Magento.SiteName,
			Host:  host,
		}
		err := templates.ExecuteTemplate(w, "index.html", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func fileServe(w http.ResponseWriter, r *http.Request) {
	path := fmt.Sprintf("public%s", r.URL.Path)
	http.ServeFile(w, r, path)
}

func adminLoginHandler(app *App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		host := fmt.Sprintf("http://%s", app.SensorIP)
		p := Page{
			Title: app.Config.Magento.SiteName,
			Host:  host,
		}
		var err error
		if r.Method == "POST" {
			username := r.PostFormValue("name")
			password := r.PostFormValue("pass")
			recordCredentials(app, r, username, password)

			p.Username = username
			p.Error = true

			err = templates.ExecuteTemplate(w, "admin-login-invalid.html", p)
		} else {
			recordAttack(app, r, MagentoScanAdminLogin)
			err = templates.ExecuteTemplate(w, "admin-login.html", p)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func versionHandler(app *App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		recordAttack(app, r, MagentoScanVersion)
		host := fmt.Sprintf("http://%s", app.SensorIP)
		p := Page{
			Title:   app.Config.Magento.SiteName,
			Version: app.Config.Magento.MagentoVersionText,
			Host:    host,
		}
		err := templates.ExecuteTemplate(w, "magento_version.txt", p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// TODO: Randomize csrf token
func loginHandler(app *App) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		host := fmt.Sprintf("http://%s", app.SensorIP)
		p := Page{
			Title: app.Config.Magento.SiteName,
			Host:  host,
		}
		var err error
		if r.Method == "POST" {
			username := r.PostFormValue("name")
			password := r.PostFormValue("pass")
			recordCredentials(app, r, username, password)

			p.Username = username
			p.Error = true

			err = templates.ExecuteTemplate(w, "login-invalid.html", p)
		} else {
			recordAttack(app, r, MagentoScanLogin)
			err = templates.ExecuteTemplate(w, "login.html", p)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// routes sets up the necessary http routing for the webapp.
func routes(app *App) *http.ServeMux {
	mux := http.NewServeMux()
	// TODO: Smart regexes for routing.
	mux.HandleFunc("/", IndexHandler(app))
	mux.HandleFunc("/magento_version", versionHandler(app))
	mux.HandleFunc("/pub/", fileServe)
	mux.HandleFunc("/admin_access/", adminLoginHandler(app))      // Any un-covered route
	mux.HandleFunc("/customer/account/login/", loginHandler(app)) // Anything containing login
	return mux
}

// recordRequest will parse the http.Request and put it into a normalized format
// and then marshal to JSON. This can then be sent on an hpfeeds channel or
// logged to a file directly.
func recordAttack(app *App, r *http.Request, signature string) {
	// Populate data to send
	payload, err := app.Agave.NewHTTPAttack(signature, r)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// Marshal it to json so we can send it over hpfeeds.
	buf, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	log.Printf("%s", buf)

	// Send to hpfeeds broker
	if app.Config.Hpfeeds.Enabled {
		app.Publish <- buf
	}
}

// recordRequest will parse the http.Request and put it into a normalized format
// and then marshal to JSON. This can then be sent on an hpfeeds channel or
// logged to a file directly.
func recordCredentials(app *App, r *http.Request, username string, password string) {
	// Populate data to send
	payload, err := app.Agave.NewCredentialAttack(r, username, password)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	// Marshal it to json so we can send it over hpfeeds.
	buf, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	log.Printf("%s", buf)

	// Send to hpfeeds broker
	if app.Config.Hpfeeds.Enabled {
		app.Publish <- buf
	}
}
