package model

import (
	"fmt"
	"strings"
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
	GiornoNascita uint8
	MeseNascita   uint8
	AnnoNascita   uint64
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

func EstraiCodFiscale(utente Utente) (string, error) {
	result := utente.estraiCognome().estraiNome().CodFiscale

	fmt.Println(result)

	return result, nil
}
