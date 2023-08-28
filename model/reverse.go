package model

import (
	"errors"
	"os"
	"strings"

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
		if match.Distance <= 5 {
			inverso.Cognome += (match.Target + ", ")
		}
	}

	// Rimuovi l'ultima virgola dal risultato
	inverso.Cognome = inverso.Cognome[:len(inverso.Cognome)-2]

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
		estraiCognome()

	return *result, nil
}
