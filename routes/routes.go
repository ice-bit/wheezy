package routes

import (
	"fmt"
	"net/http"
	"text/template"
)

func RootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		errorNotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		{
			t, _ := template.ParseFiles(
				"views/pages/index.tmpl",
				"views/partials/header.tmpl",
				"views/partials/navbar.tmpl",
				"views/partials/footer.tmpl")

			t.Execute(w, nil)
		}
	case http.MethodPost:
		{

		}
	default:
		http.Error(w, fmt.Sprintf("Cannot %s %s", r.Method, r.URL.Path), http.StatusMethodNotAllowed)
	}
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		{
			t, _ := template.ParseFiles(
				"views/pages/about.tmpl",
				"views/partials/header.tmpl",
				"views/partials/navbar.tmpl",
				"views/partials/footer.tmpl")

			t.Execute(w, nil)
		}
	default:
		http.Error(w, fmt.Sprintf("Cannot %s %s", r.Method, r.URL.Path), http.StatusMethodNotAllowed)
	}
}

func errorNotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("Cannot %s %s", r.Method, r.URL.Path), http.StatusNotFound)
}
