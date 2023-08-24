package model

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type Error struct {
	Codice  uint16
	Message string
}

type Utente struct {
	Cognome       string
	Nome          string
	Sesso         string
	LuogoNascita  string
	GiornoNascita int
	MeseNascita   int
	AnnoNascita   string
	CodFiscale    string
	Errore        Error
}

func estraiConsonanti(s string) string {
	vocali := "aeiou"
	s = strings.ToLower(s)
	var result []rune

	for _, c := range s {
		if !strings.ContainsRune(vocali, c) {
			result = append(result, c)
		}
	}

	return string(result)
}

func estraiVocali(s string) string {
	vocali := "aeiou"
	s = strings.ToLower(s)
	var result []rune

	for _, c := range s {
		if strings.ContainsRune(vocali, c) {
			result = append(result, c)
		}
	}

	return string(result)
}

func estraiCaratteriPari(s string) []string {
	var result []string

	for i, c := range s {
		if (i+1)%2 == 0 {
			result = append(result, string(c))
		}
	}

	return result
}

func estraiCaratteriDispari(s string) []string {
	var result []string

	for i, c := range s {
		if (i+1)%2 != 0 {
			result = append(result, string(c))
		}
	}

	return result
}

func (utente *Utente) estraiCognome() *Utente {
	// Estrai le prime tre consonanti dal cognome
	consonantiCognome := estraiConsonanti(utente.Cognome)
	codCognome := consonantiCognome[:3]

	// Se le consonanti sono minori di tre, estrai pure le vocali
	if len(consonantiCognome) < 3 {
		vocaliCognome := estraiVocali(utente.Cognome)
		codCognome += vocaliCognome
		codCognome = codCognome[:3]
	}

	// Se il risultato < 3(i.e. il cognome e' di due caratteri), aggiungi 'x'
	if len(codCognome) < 3 {
		codCognome += "x"
	}

	utente.CodFiscale += codCognome

	return utente
}

func (utente *Utente) estraiNome() *Utente {
	// Estrai le consonanti dal nome
	consonantiNome := estraiConsonanti(utente.Nome)
	// Se le consonanti sono >= 4, estrai la prima, la terza e la quarta
	if len(consonantiNome) >= 4 {
		codNome := string(consonantiNome[0])
		codNome += string(consonantiNome[2])
		codNome += string(consonantiNome[3])
		utente.CodFiscale += codNome

		return utente
	}

	// Altrimenti prende le prime tre consonanti in ordine
	codNome := consonantiNome[:3]
	// Se le consonanti sono minori di tre, estrai pure le vocali
	if len(consonantiNome) < 3 {
		vocaliNome := estraiVocali(utente.Nome)
		codNome += vocaliNome
		codNome = codNome[:3]
	}

	// Se il risultato < 3(i.e. il nome e' di due caratteri), aggiungi 'x'
	if len(codNome) < 3 {
		codNome += "x"
	}

	utente.CodFiscale += codNome

	return utente
}

func (utente *Utente) estraiAnnoNascita() *Utente {
	utente.CodFiscale += utente.AnnoNascita[len(utente.AnnoNascita)-2:]

	return utente
}

func (utente *Utente) estraiMeseNascita() *Utente {
	// Mappa ciascun mese al corrispondente valore
	mappaMesi := map[int]byte{
		1:  'A', // Gennaio
		2:  'B', // Febbraio
		3:  'C', // Marzo
		4:  'D', // Aprile
		5:  'E', // Maggio
		6:  'H', // Giugno
		7:  'L', // Luglio
		8:  'M', // Agosto
		9:  'P', // Settembre
		10: 'R', // Ottobre
		11: 'S', // Novembre
		12: 'T', // Dicembre
	}

	// Ritorna il valore corrispondente al mese scelto
	utente.CodFiscale += string(mappaMesi[utente.MeseNascita])

	return utente
}

func (utente *Utente) estraiGiornoNascita() *Utente {
	giornoNascita := utente.GiornoNascita

	// Se il soggetto e' una donna, sommare 40 al giorno di nascita
	if utente.Sesso == "femminile" {
		giornoNascita += 40
		utente.CodFiscale += strconv.Itoa(giornoNascita)

		return utente
	}

	// Se il risultato finale <= 9, anteporre uno '0' al risultato
	if giornoNascita < 10 {
		utente.CodFiscale += "0"
	}

	utente.CodFiscale += strconv.Itoa(giornoNascita)

	return utente
}

func (utente *Utente) estraiLuogoNascita() *Utente {
	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// Cerca il codice del luogo di nascita nella cache
	codLuogoNascita, err := redisClient.Get(ctx, strings.ToUpper(utente.LuogoNascita)).Result()
	// Se non viene trovato, estrailo dal database
	if err != nil {
		if err == redis.Nil {
			codCatastale := EstraiCodiceCatastale(utente.LuogoNascita)

			// Se il codice catastale e' stato trovato, salvalo nella cache
			if codCatastale != "" {
				err := redisClient.Set(
					ctx,
					strings.ToUpper(utente.LuogoNascita),
					codCatastale,
					24*time.Hour).Err() // Invalida la chiave dopo 24 ore

				if err != nil {
					panic(err)
				}

				// Aggiorna il codice fiscale
				utente.CodFiscale += codCatastale
			} else {
				// Se invece il codice catastale non e' stato trovato, prova a
				// cercare il codice della nazione
				codNazione := EstraiCodiceNazione(utente.LuogoNascita)

				// Se il codice della nazione esiste, aggiorna il codice fiscale
				if codNazione != "" {
					utente.CodFiscale += codNazione
				} else {
					// Altrimenti, se non e' stato trovato nemmeno il codice
					// della nazione, ritorna un errore
					utente.Errore = Error{
						Codice:  400,
						Message: "Il luogo di nascita selezionato non esiste",
					}
				}
			}
		} else {
			panic(err)
		}
	} else { // Se il codice esiste nella cache, aggiorna il CF
		utente.CodFiscale += codLuogoNascita
	}

	return utente
}

func EstraiCodFiscale(utente Utente) (string, Error) {
	result := utente.
		estraiCognome().
		estraiNome().
		estraiAnnoNascita().
		estraiMeseNascita().
		estraiGiornoNascita().
		estraiLuogoNascita()

	fmt.Println(result.CodFiscale)

	if result.Errore.Message != "" {
		return "", result.Errore
	}

	return result.CodFiscale, Error{}
}
