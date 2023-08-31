package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/ice-bit/wheezy/log"
	"github.com/ice-bit/wheezy/model"
)

type utenteTmplType struct {
	Utente model.Utente
	Errori []string
}

type inversoTmplType struct {
	Reverse model.Inverso
	Errori  []string
}

type nilTmplType struct{}

type Entity interface {
	utenteTmplType | inversoTmplType | nilTmplType
}

func renderTemplate[T Entity](res http.ResponseWriter, view string, entity T) {
	switch any(entity).(type) {
	case nilTmplType:
		{
			t, _ := template.ParseFiles(
				view,
				"views/partials/header.tmpl",
				"views/partials/navbar.tmpl",
				"views/partials/footer.tmpl")

			t.Execute(res, nil)
		}
	default:
		{
			t, _ := template.ParseFiles(
				view,
				"views/partials/header.tmpl",
				"views/partials/navbar.tmpl",
				"views/partials/footer.tmpl")

			t.Execute(res, entity)
		}
	}
}

func RootHandler(res http.ResponseWriter, req *http.Request) {
	if req.URL.Path != "/" {
		errorNotFound(res, req)
		return
	}

	switch req.Method {
	case http.MethodGet:
		{
			log.InfoLogger.Printf("incoming %s request at %s", req.Method, req.URL.Path)
			renderTemplate[nilTmplType](res, "views/pages/index.tmpl", nilTmplType{})
		}
	case http.MethodPost:
		{
			log.InfoLogger.Printf("incoming %s request at %s", req.Method, req.URL.Path)

			// Valida i campi del form
			errors := formValidator(req)

			if len(errors) > 0 {
				renderTemplate[utenteTmplType](res, "views/pages/index.tmpl", utenteTmplType{
					Utente: model.Utente{},
					Errori: errors,
				})
			} else {
				// Estrai il form dalla request
				utente := model.Utente{
					Cognome:      req.FormValue("cognome"),
					Nome:         req.FormValue("nome"),
					Sesso:        req.FormValue("sesso"),
					LuogoNascita: req.FormValue("luogoNascita"),
					GiornoNascita: func() uint {
						v, _ := strconv.ParseUint(req.FormValue("giornoNascita"), 10, 64)
						return uint(v)
					}(),
					MeseNascita: func() uint {
						v, _ := strconv.ParseUint(req.FormValue("meseNascita"), 10, 64)
						return uint(v)
					}(),
					AnnoNascita: req.FormValue("annoNascita"),
					CodFiscale:  "",
					Errore:      "",
				}

				// Normalizza campi
				utente.Nome = normalizeField(utente.Nome)
				utente.Cognome = normalizeField(utente.Cognome)
				utente.LuogoNascita = normalizeBirthPlace(utente.LuogoNascita)

				// Estrai il codice fiscale
				utente, err := model.EstraiCodFiscale(utente)
				if err != nil {
					renderTemplate[utenteTmplType](res, "views/pages/index.tmpl", utenteTmplType{
						Utente: model.Utente{},
						Errori: []string{err.Error()},
					})
				} else {
					renderTemplate[utenteTmplType](
						res,
						"views/pages/index.tmpl", utenteTmplType{
							Utente: utente,
							Errori: nil,
						})
				}
			}
		}
	default:
		http.Error(res, fmt.Sprintf("Cannot %s %s", req.Method, req.URL.Path), http.StatusMethodNotAllowed)
	}
}

func ReverseHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		{
			log.InfoLogger.Printf("incoming %s request at %s", req.Method, req.URL.Path)
			renderTemplate[nilTmplType](res, "views/pages/reverse.tmpl", nilTmplType{})
		}
	case http.MethodPost:
		{
			log.InfoLogger.Printf("incoming %s request at %s", req.Method, req.URL.Path)

			// Estrai e valida il campo del codice fiscale
			codFiscale := req.FormValue("codFiscale")
			if codFiscale == "" {
				renderTemplate[inversoTmplType](res, "views/pages/reverse.tmpl", inversoTmplType{
					model.Inverso{},
					[]string{"inserire il codice fiscale"},
				})
			}

			// Normalizza il campo
			codFiscale = strings.ToUpper(strings.TrimSpace(codFiscale))
			// Estrai l'utente
			utente, err := model.EstraiInverso(codFiscale)
			if err != nil {
				renderTemplate[inversoTmplType](res, "views/pages/reverse.tmpl", inversoTmplType{
					model.Inverso{},
					[]string{err.Error()},
				})
			} else {
				renderTemplate[inversoTmplType](res, "views/pages/reverse.tmpl", inversoTmplType{
					utente,
					nil,
				})
			}
		}
	default:
		http.Error(res, fmt.Sprintf("Cannot %s %s", req.Method, req.URL.Path), http.StatusMethodNotAllowed)
	}
}

func AboutHandler(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case http.MethodGet:
		{
			log.InfoLogger.Printf("incoming %s request at %s", req.Method, req.URL.Path)
			renderTemplate[nilTmplType](res, "views/pages/about.tmpl", nilTmplType{})
		}
	default:
		http.Error(res, fmt.Sprintf("Cannot %s %s", req.Method, req.URL.Path), http.StatusMethodNotAllowed)
	}
}

func errorNotFound(res http.ResponseWriter, req *http.Request) {
	log.WarnLogger.Printf("cannot %s %s", req.Method, req.URL.Path)
	http.Error(res, fmt.Sprintf("Cannot %s %s", req.Method, req.URL.Path), http.StatusNotFound)
}

func formValidator(req *http.Request) []string {
	var (
		cognome       = req.FormValue("cognome")
		nome          = req.FormValue("nome")
		sesso         = req.FormValue("sesso")
		luogoNascita  = req.FormValue("luogoNascita")
		giornoNascita = req.FormValue("giornoNascita")
		meseNascita   = req.FormValue("meseNascita")
		annoNascita   = req.FormValue("annoNascita")
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
			emptyFields = append(emptyFields, "inserire il "+fieldName)
		}
	}

	return emptyFields
}

func normalizeField(s string) string {
	s = strings.TrimSpace(s)

	if len(s) > 0 {
		return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
	}

	return s
}

func normalizeBirthPlace(s string) string {
	s = strings.TrimSpace(s)
	words := strings.Fields(s)
	normalizedWords := make([]string, 0, len(words))

	for i, word := range words {
		if i == 0 {
			normalizedWords = append(normalizedWords, strings.ToUpper(word[:1])+strings.ToLower(word[1:]))
		} else {
			normalizedWords = append(normalizedWords, strings.ToLower(word))
		}
	}

	return strings.Join(normalizedWords, " ")
}
