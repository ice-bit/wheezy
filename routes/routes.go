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
			// Valida i campi del form
			errors := formValidator(r)

			if errors != nil {
				t, _ := template.ParseFiles(
					"views/pages/index.tmpl",
					"views/partials/header.tmpl",
					"views/partials/navbar.tmpl",
					"views/partials/footer.tmpl")

				t.Execute(w, struct {
					ErrorMessages []string
				}{
					ErrorMessages: errors,
				})
			}
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

func formValidator(r *http.Request) []string {
	var (
		cognome       = r.FormValue("cognome")
		nome          = r.FormValue("nome")
		sesso         = r.FormValue("sesso")
		luogoNascita  = r.FormValue("luogoNascita")
		giornoNascita = r.FormValue("giornoNascita")
		meseNascita   = r.FormValue("meseNascita")
		annoNascita   = r.FormValue("annoNascita")
	)

	emptyFields := []string{}

	fieldsToCheck := map[string]string{
		"cognome":       cognome,
		"nome":          nome,
		"sesso":         sesso,
		"luogoNascita":  luogoNascita,
		"giornoNascita": giornoNascita,
		"meseNascita":   meseNascita,
		"annoNascita":   annoNascita,
	}

	for fieldName, fieldValue := range fieldsToCheck {
		if fieldValue == "" {
			emptyFields = append(emptyFields, fieldName)
		}
	}

	return emptyFields
}
