// - Creating a data structure with load and sve methods
// - Using the net/http package to build web applications
// - Using the html/template package to process HTML templates
// - Using the regexp package to validate user input
// Reference: https://go.dev/doc/articles/wiki/
//
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var (
	templates = template.Must(template.ParseFiles("edit.html", "view.html"))
)

func renderTemplate(w http.ResponseWriter, tmpl string, page *Page) {
	// version 1
	// r, err := template.ParseFiles(tmpl + ".html")
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// err = r.Execute(w, page)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// }

	// version 2 with cache
	err := templates.ExecuteTemplate(w, tmpl+".html", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func defaultHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Default page, your request path: %s", r.URL.Path)
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/view/"):]
	page, err := loadPage(title)
	if err != nil {
		// fmt.Fprintf(w, "404 - Page not found: %s", title)
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}

	// version 1
	// fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", page.Title, page.Body)

	// version 2
	// t, _ := template.ParseFiles("view.html")
	// t.Execute(w, page)

	// version 3
	renderTemplate(w, "view", page)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	page, err := loadPage(title)
	if err != nil {
		page = &Page{Title: title}
	}

	// fmt.Fprintf(w, "<h1>Editing %s</h1>"+
	// 	"<form action=\"/save/%s\" method=\"POST\">"+
	// 	"<textarea name=\"body\">%s</textarea><br>"+
	// 	"</form>",
	// 	page.Title, page.Title, page.Body,
	// )

	// t, _ := template.ParseFiles("edit.html")
	// t.Execute(w, page)

	renderTemplate(w, "edit", page)
}

func saveHandler(w http.ResponseWriter, r *http.Request) {
	title := r.URL.Path[len("/edit/"):]
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	if err := p.Save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	http.HandleFunc("/", defaultHandler)
	http.HandleFunc("/view/", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
