package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/dacoursey/pm2077/models"
)

// Creds is a collection for storing info entered during login.
type Creds struct {
	Username string
	Password string
}

// TemplateCache is used to hold templates loaded at start time for faster processing.
type TemplateCache struct {
	cacheSingle map[string]*template.Template
	cacheFiles  map[string]*template.Template
	cacheGlob   map[string]*template.Template
}

// newTemplateCache is the constructor to build the maps for the template cache.
func newTemplateCache() *TemplateCache {
	t := TemplateCache{}
	t.cacheSingle = make(map[string]*template.Template)
	t.cacheFiles = make(map[string]*template.Template)
	t.cacheGlob = make(map[string]*template.Template)
	return &t
}

// Templates is a global map for template storage.
// This feels weird just shoving it in here.
var Templates *TemplateCache

func loadPage(title string) (*models.Page, error) {
	_, caller, _, _ := runtime.Caller(0)
	filename := path.Dir(caller) + "/pages/" + title + ".txt"

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &models.Page{Title: title, Body: body}, nil
}

// func renderTemplateOld(w http.ResponseWriter, tmpl string, p *models.Page) {
// 	t, _ := template.ParseFiles("templates/" + tmpl + ".html")
// 	t.Execute(w, p)
// }

// func renderTemplate(tmpl string) (*template.Template, error) {
// 	_, caller, _, _ := runtime.Caller(0)
// 	filename := path.Dir(caller) + "/templates/" + tmpl + ".html"
// 	t, err := template.ParseFiles(filename)
// 	return t, err
// }

// func viewHandler(w http.ResponseWriter, r *http.Request) {
// 	//title := r.URL.Path[len("/view/"):]
// 	// p, _ := loadPage(title)
// 	//t, _ := renderTemplate("view")
// }

// func loginHandler(w http.ResponseWriter, r *http.Request) {
// 	err := r.ParseForm()
// 	if err != nil {
// 		// YOLO
// 	}

// 	c := Creds{r.PostFormValue("username"), r.PostFormValue("password")}
// 	t, _ := template.ParseFiles("templates/login.html")
// 	// Ship it.
// 	t.Execute(w, c)
// }

/////
// Template Rendering
/////

// Parse loads one template file.
func (t *TemplateCache) Parse(name, text string) (*template.Template, error) {
	tt, ok := t.cacheSingle[name]
	if ok {
		return tt, nil
	}

	// parse the template
	tt, err := template.New(name).Parse(text)
	if err != nil {
		panic(err)
	}

	// cache it
	t.cacheSingle[name] = tt
	return tt, nil
}

// ParseFiles loads multiple template files.
func (t *TemplateCache) ParseFiles(name string, filenames ...string) (*template.Template, error) {
	tt, ok := t.cacheSingle[name]
	if ok {
		return tt, nil
	}

	// parse the template
	tt, err := template.New(name).ParseFiles(filenames...)
	if err != nil {
		panic(err)
	}

	// cache it
	t.cacheSingle[name] = tt
	return tt, nil
}

// ParseGlob loads template data from a glob.
func (t *TemplateCache) ParseGlob(name, pattern string) (*template.Template, error) {
	tt, ok := t.cacheSingle[name]
	if ok {
		return tt, nil
	}

	// parse the template
	tt, err := template.New(name).ParseGlob(pattern)
	if err != nil {
		panic(err)
	}

	// cache it
	t.cacheSingle[name] = tt
	return tt, nil
}

func loadTemplates() {
	box := rice.MustFindBox("templates")

	Templates = newTemplateCache()

	err := box.Walk("", func(path string, fi os.FileInfo, err error) error {

		if !strings.HasPrefix(path, ".") && path != "" && !fi.IsDir() {

			Templates.Parse(path, box.MustString(path))
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}

// Render templates using embedded binary data storage
func getTemplate(w http.ResponseWriter, r *http.Request, requestedFile string, context interface{}) {

	// // Find our box, if this fails we panic.
	// templateBox := rice.MustFindBox("templates")

	// Find our TemplateCache
	tpl := Templates.cacheSingle[requestedFile+".html"]

	// // Get the layout contents as string
	// layoutString, err := templateBox.String("_layout.html")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Convert the layout string to a template
	// layout, err := template.ParseGlob(layoutString)
	// layout.Parse(templateBox.String("_header.html"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // parse and execute the template
	// //template, err := template.New("message").Parse(templateString)
	// fullTpl, err := layout.Parse(templateString)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// if fullTpl != nil {
	// 	err := fullTpl.Execute(w, context)
	// 	if err != nil {
	// 		log.Println(err)
	// 	}
	// } else {
	// 	w.WriteHeader(http.StatusNotFound)
	// }

	if tpl != nil {
		err := tpl.Execute(w, context)
		if err != nil {
			log.Println(err)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

/////////////////////////////////////////////////////////////////
// Beyond here is the backup of the direct FS access version
// This will be removed when the previous section is working
// These require FS symlinks for the template dirs
/////////////////////////////////////////////////////////////////

// // Render the templates that use _layout.html
// func getTemplateBackup(w http.ResponseWriter, r *http.Request, requestedFile string, context interface{}) {
// 	templates := populateTemplates()
// 	template := templates[requestedFile+".html"]

// 	if template != nil {
// 		err := template.Execute(w, context)
// 		if err != nil {
// 			log.Println(err)
// 		}
// 	} else {
// 		w.WriteHeader(404)
// 	}
// }

// // Retrieve and parse the templates
// func populateTemplatesBackup() map[string]*template.Template {
// 	result := make(map[string]*template.Template)
// 	const basePath = "templates"
// 	layout := template.Must(template.ParseFiles(basePath + "/_layout.html"))
// 	template.Must(layout.ParseFiles(basePath+"/_sideBar.html", basePath+"/_header.html"))
// 	dir, err := os.Open(basePath + "/content")
// 	if err != nil {
// 		panic("Failed to open template blocks directory: " + err.Error())
// 	}
// 	fis, err := dir.Readdir(-1)
// 	if err != nil {
// 		panic("Failed to read contents of content directory: " + err.Error())
// 	}
// 	for _, fi := range fis {
// 		f, err := os.Open(basePath + "/content/" + fi.Name())
// 		if err != nil {
// 			panic("Failed to open template '" + fi.Name() + "'")
// 		}
// 		content, err := ioutil.ReadAll(f)
// 		if err != nil {
// 			panic("Failed to read content from file '" + fi.Name() + "'")
// 		}
// 		f.Close()
// 		tmpl := template.Must(layout.Clone())
// 		_, err = tmpl.Parse(string(content))
// 		if err != nil {
// 			panic("Failed to parse contents of '" + fi.Name() + "' as template")
// 		}
// 		result[fi.Name()] = tmpl
// 	}
// 	return result
//}
