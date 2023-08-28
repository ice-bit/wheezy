package model

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lithammer/fuzzysearch/fuzzy"
)

type Inverso struct {
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

func (inverso *Inverso) estraiCognome() *Inverso {
	iniziali := inverso.CodFiscale[:3]
	fileCognomi := "cognomi.txt"
	fileContent, err := os.ReadFile(fileCognomi)
	if err != nil {
		panic(err)
	}

	cognomi := strings.Split(string(fileContent), "\n")
	risultatoRicerca := fuzzy.RankFindFold(iniziali, cognomi)

	for _, match := range risultatoRicerca {
		// Estrai i risultati con una distanza di Hamming bassa
		if match.Distance <= 7 {
			inverso.Cognome += (match.Target + ", ")
		}
	}

	// Rimuovi l'ultima virgola dal risultato
	if len(inverso.Cognome) > 0 {
		inverso.Cognome = inverso.Cognome[:len(inverso.Cognome)-2]
	}

	return inverso
}

func (inverso *Inverso) estraiNome() *Inverso {
	iniziali := inverso.CodFiscale[3:6]
	sesso := inverso.estraiSesso().Sesso
	var fileNomi string

	if sesso == "maschile" {
		fileNomi = "nomi_maschili.txt"
	} else {
		fileNomi = "nomi_femminili.txt"
	}

	fileContent, err := os.ReadFile(fileNomi)
	if err != nil {
		panic(err)
	}

	nomi := strings.Split(string(fileContent), "\n")
	risultatoRicerca := fuzzy.RankFindFold(iniziali, nomi)

	for _, match := range risultatoRicerca {
		// Estrai i risultati con una distanza di Hamming bassa
		if match.Distance <= 10 {
			inverso.Nome += (match.Target + ", ")
		}
	}

	// Rimuovi l'ultima virgola dal risultato
	if len(inverso.Nome) > 0 {
		inverso.Nome = inverso.Nome[:len(inverso.Nome)-2]
	}

	return inverso
}

func (inverso *Inverso) estraiAnnoNascita() *Inverso {
	annoNascita := inverso.CodFiscale[6:8]
	annoCorrente := time.Now().Year()

	// Se le cifre dell'anno di nascita sono maggiori
	// dell'anno corrente, allora anteponi '19' all'anno di nascita.
	// Altrimenti anteponi '20'
	if v, _ := strconv.Atoi(("20" + annoNascita)); v > annoCorrente {
		inverso.AnnoNascita = "19" + annoNascita
	} else {
		inverso.AnnoNascita = "20" + annoNascita
	}

	return inverso
}

func (inverso *Inverso) estraiMeseNascita() *Inverso {
	mappaMesi := map[string]uint{
		"A": 1,  // Gennaio
		"B": 2,  // Febbraio
		"C": 3,  // Marzo
		"D": 4,  // Aprile
		"E": 5,  // Maggio
		"H": 6,  // Giugno
		"L": 7,  // Luglio
		"M": 8,  // Agosto
		"P": 9,  // Settembre
		"R": 10, // Ottobre
		"S": 11, // Novembre
		"T": 12, // Dicembre
	}

	inverso.MeseNascita = mappaMesi[string(inverso.CodFiscale[8])]

	return inverso
}

func (inverso *Inverso) estraiGiornoNascita() *Inverso {
	giornoNascita := inverso.CodFiscale[9:11]

	// Se il giorno di nascita e' compreso tra '41' e '71'
	// si tratta di un soggetto di sesso femminile.
	// Pertanto si sottrae '40' dal risultato
	if v, _ := strconv.Atoi(giornoNascita); v >= 41 && v <= 71 {
		inverso.GiornoNascita = uint(v - 40)
	} else if v >= 1 && v <= 31 {
		inverso.GiornoNascita = uint(v)
	} else {
		inverso.Errore = "il giorno di nascita del codice fiscale risulta invalido"
	}

	return inverso
}

func (inverso *Inverso) estraiSesso() *Inverso {
	giornoNascita, err := strconv.Atoi(inverso.CodFiscale[9:11])
	if err != nil {
		panic(err)
	}

	if giornoNascita >= 41 && giornoNascita <= 71 {
		inverso.Sesso = "femminile"
	} else {
		inverso.Sesso = "maschile"
	}

	return inverso
}

func EstraiInverso(codFiscale string) (Inverso, error) {
	var inverso Inverso

	// Un codice fiscale valido ha lunghezza pari a 16
	if len(codFiscale) != 16 {
		return Inverso{}, errors.New("codice fiscale invalido")
	}

	inverso.CodFiscale = codFiscale
	result := inverso.
		estraiCognome().
		estraiNome().
		estraiAnnoNascita().
		estraiMeseNascita().
		estraiGiornoNascita()

	if result.Errore != "" {
		return Inverso{}, errors.New(result.Errore)
	}

	return *result, nil
}
