package model

import (
	"context"
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ice-bit/wheezy/log"
	"github.com/redis/go-redis/v9"
)

type Utente struct {
	Cognome       string
	Nome          string
	Sesso         string
	LuogoNascita  string
	GiornoNascita uint
	MeseNascita   uint
	AnnoNascita   string
	CodFiscale    string
	Errore        string
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
	var codCognome string
	consonantiCognome := estraiConsonanti(utente.Cognome)
	if len(consonantiCognome) >= 3 {
		codCognome = consonantiCognome[:3]
	} else {
		codCognome = consonantiCognome
	}

	// Se le consonanti sono minori di tre, estrai pure le vocali
	if len(consonantiCognome) < 3 {
		vocaliCognome := estraiVocali(utente.Cognome)
		codCognome += vocaliCognome
		if len(codCognome) >= 3 {
			codCognome = codCognome[:3]
		}
	}

	// Se il risultato < 3(i.e. il cognome e' di due caratteri), aggiungi 'x'
	if len(codCognome) < 3 {
		codCognome += "x"
	}

	utente.CodFiscale += codCognome

	return utente
}

func (utente *Utente) estraiNome() *Utente {
	var codNome string
	// Estrai le consonanti dal nome
	consonantiNome := estraiConsonanti(utente.Nome)
	// Se le consonanti sono >= 4, estrai la prima, la terza e la quarta
	if len(consonantiNome) >= 4 {
		codNome = string(consonantiNome[0])
		codNome += string(consonantiNome[2])
		codNome += string(consonantiNome[3])
		utente.CodFiscale += codNome

		return utente
	}

	// Altrimenti prendi le prime tre consonanti in ordine
	if len(consonantiNome) == 3 {
		codNome = consonantiNome[:3]
	} else {
		codNome = consonantiNome
	}

	// Se le consonanti sono minori di tre, estrai pure le vocali
	if len(consonantiNome) < 3 {
		vocaliNome := estraiVocali(utente.Nome)
		codNome += vocaliNome
		if len(codNome) >= 3 {
			codNome = codNome[:3]
		}
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
	mappaMesi := map[uint]byte{
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
		utente.CodFiscale += strconv.FormatUint(uint64(giornoNascita), 10)

		return utente
	}

	// Se il risultato finale <= 9, anteporre uno '0' al risultato
	if giornoNascita < 10 {
		utente.CodFiscale += "0"
	}

	utente.CodFiscale += strconv.FormatUint(uint64(giornoNascita), 10)

	return utente
}

func (utente *Utente) estraiLuogoNascita() *Utente {
	// Estrai l'host e la porta di redis dalle
	// variabili d'ambiente
	var (
		host = os.Getenv("WHEEZY_REDIS_ADDRESS")
		port = os.Getenv("WHEEZY_REDIS_PORT")
	)

	if host == "" || port == "" {
		log.ErrLogger.Printf("environment variables not configured")
		panic("Environment variables not configured")
	}

	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     (host + ":" + port),
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
					log.ErrLogger.Printf(err.Error())
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
					utente.Errore = "il luogo di nascita selezionato non esiste"
				}
			}
		} else {
			log.ErrLogger.Printf(err.Error())
			panic(err)
		}
	} else { // Se il codice esiste nella cache, aggiorna il CF
		utente.CodFiscale += codLuogoNascita
	}

	return utente
}

func (utente *Utente) estraiCodiceControllo() *Utente {
	utente.CodFiscale = strings.ToUpper(utente.CodFiscale)
	// Separa i caratteri in posizione dispari da quelli in posizione pari
	caratteriDispari := estraiCaratteriDispari(utente.CodFiscale)
	caratteriPari := estraiCaratteriPari(utente.CodFiscale)

	// Mappa dei valori dispari
	mappaDispari := map[string]uint{
		"0": 1, "1": 0, "2": 5, "3": 7, "4": 9, "5": 13,
		"6": 15, "7": 17, "8": 19, "9": 21, "A": 1, "B": 0,
		"C": 5, "D": 7, "E": 9, "F": 13, "G": 15, "H": 17,
		"I": 19, "J": 21, "K": 2, "L": 4, "M": 18, "N": 20,
		"O": 11, "P": 3, "Q": 6, "R": 8, "S": 12, "T": 14,
		"U": 16, "V": 10, "W": 22, "X": 25, "Y": 24, "Z": 23,
	}

	// Mappa dei valori pari
	mappaPari := map[string]uint{
		"0": 0, "1": 1, "2": 2, "3": 3, "4": 4, "5": 5,
		"6": 6, "7": 7, "8": 8, "9": 9, "A": 0, "B": 1,
		"C": 2, "D": 3, "E": 4, "F": 5, "G": 6, "H": 7,
		"I": 8, "J": 9, "K": 10, "L": 11, "M": 12, "N": 13,
		"O": 14, "P": 15, "Q": 16, "R": 17, "S": 18, "T": 19,
		"U": 20, "V": 21, "W": 22, "X": 23, "Y": 24, "Z": 25,
	}

	// Mappa del carattere di controllo
	mappaControllo := map[uint]string{
		0: "A", 1: "B", 2: "C", 3: "D", 4: "E", 5: "F",
		6: "G", 7: "H", 8: "I", 9: "J", 10: "K", 11: "L",
		12: "M", 13: "N", 14: "O", 15: "P", 16: "Q", 17: "R",
		18: "S", 19: "T", 20: "U", 21: "V", 22: "W", 23: "X",
		24: "Y", 25: "Z",
	}

	var (
		sommaDispari    uint
		sommaPari       uint
		valoreControllo uint
	)

	// Somma i valori dispari associati a ciascun carattere
	for idx, val := range caratteriDispari {
		if idx == 0 {
			sommaDispari = mappaDispari[val]
		} else {
			sommaDispari += mappaDispari[val]
		}
	}

	// Somma i valori pari associati a ciascun carattere
	for idx, val := range caratteriPari {
		if idx == 0 {
			sommaPari = mappaPari[val]
		} else {
			sommaPari += mappaPari[val]
		}
	}

	// Somma i due risultati parziali ed esegui la divisone modulo 26
	valoreControllo = ((sommaPari + sommaDispari) % 26)

	// Mappa il valore di controllo al relativo carattere
	utente.CodFiscale += mappaControllo[valoreControllo]

	return utente
}

func EstraiCodFiscale(utente Utente) (Utente, error) {
	result := utente.
		estraiCognome().
		estraiNome().
		estraiAnnoNascita().
		estraiMeseNascita().
		estraiGiornoNascita().
		estraiLuogoNascita().
		estraiCodiceControllo()

	if utente.Errore != "" {
		// Se ci sono errori, resetta il codice fiscale
		// in modo tale che l'engine non renderizzi il template
		// del risultato
		utente.CodFiscale = ""
		return utente, errors.New(utente.Errore)
	}

	return *result, nil
}
