// - Creating a data structure with load and sve methods
// - Using the net/http package to build web applications
// - Using the html/template package to process HTML templates
// - Using the regexp package to validate user input
// Reference: https://go.dev/doc/articles/wiki/
//
package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	tmplDir = "tmpl"
	// When parsing multiple files with the same name in different directories,
	// the last one mentioned will be the one that results.
	// For instance, ParseFiles("a/foo", "b/foo") stores "b/foo" as the template
	// named "foo", while "a/foo" is unavailable.
	templates = template.Must(template.ParseFiles(
		tmplDir+"/header.html",
		tmplDir+"/default.html",
		tmplDir+"/edit.html",
		tmplDir+"/view.html",
	))
	validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")
)

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}
	// the Title is the second subexpression
	return m[2], nil
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title, err := getTitle(w, r)
		if err != nil {
			return
		}
		fn(w, r, title)
	}
}

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
	// just passing the filename is enough, not need to pass the full path
	err := templates.ExecuteTemplate(w, tmpl+".html", page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// defaultHandler show a list of page to view
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	// version 1
	// fmt.Fprintf(w, "Default page, your request path: %s", r.URL.Path)

	// version 2
	entries, err := os.ReadDir(dataDir)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	pageNames := make([]string, len(entries))
	for i, v := range entries {
		pageNames[i] = strings.Split(v.Name(), ".")[0]
	}

	err = templates.ExecuteTemplate(w, "default.html", pageNames)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
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

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
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

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
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
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
