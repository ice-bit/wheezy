package model

import (
	"testing"
)

func TestEstraiConsonanti(t *testing.T) {
	str := "esempio"
	expected := "smp"
	got := estraiConsonanti(str)

	if expected != got {
		t.Errorf("error: expected %s got %s", expected, got)
	}
}

func TestEstraiVocali(t *testing.T) {
	str := "esempio"
	expected := "eeio"
	got := estraiVocali(str)

	if expected != got {
		t.Errorf("error: expected %s got %s", expected, got)
	}
}

func TestEstraiCaratteriPari(t *testing.T) {
	str := "esempio"
	expected := "smi"
	got := estraiCaratteriPari(str)

	for i, v := range got {
		if v != string(expected[i]) {
			t.Errorf("error: expected %s got %s", expected, got)
		}
	}
}

func TestEstraiCaratteriDispari(t *testing.T) {
	str := "esempio"
	expected := "eepo"
	got := estraiCaratteriDispari(str)

	for i, v := range got {
		if v != string(expected[i]) {
			t.Errorf("error: expected %s got %s", expected, got)
		}
	}
}

func TestEstraiCognome(t *testing.T) {
	var utente Utente

	utente.Cognome = "viola"
	expected := "vli"
	got := utente.estraiCognome().CodFiscale

	if expected != got {
		t.Errorf("error: expected %s got %s", expected, got)
	}
}

func TestEstraiCognomeBreve(t *testing.T) {
	var utente Utente

	utente.Cognome = "fo"
	expected := "fox"
	got := utente.estraiCognome().CodFiscale

	if expected != got {
		t.Errorf("error: expected %s got %s", expected, got)
	}
}

func TestEstraiNomeConsonantiGEQ4(t *testing.T) {
	var utente Utente

	utente.Nome = "francesco"
	expected := "fnc"
	got := utente.estraiNome().CodFiscale

	if expected != got {
		t.Errorf("error: expected %s got %s", expected, got)
	}
}

func TestEstraiNomeConsonantiEQ3(t *testing.T) {
	var utente Utente

	utente.Nome = "marco"
	expected := "mrc"
	got := utente.estraiNome().CodFiscale

	if expected != got {
		t.Errorf("error: expected %s got %s", expected, got)
	}
}

func TestEstraiNomeConsonantiLT3(t *testing.T) {
	var utente Utente

	utente.Nome = "mario"
	expected := "mra"
	got := utente.estraiNome().CodFiscale

	if expected != got {
		t.Errorf("error: expected %s got %s", expected, got)
	}
}

func TestEstraiNomeBreve(t *testing.T) {
	var utente Utente

	utente.Nome = "ed"
	expected := "dex"
	got := utente.estraiNome().CodFiscale

	if expected != got {
		t.Errorf("error: expected %s got %s", expected, got)
	}
}

func TestEstraiAnnoNascita(t *testing.T) {
	var utente Utente

	utente.AnnoNascita = "1968"
	expected := "68"
	got := utente.estraiAnnoNascita().CodFiscale

	if expected != got {
		t.Errorf("error: expected %s got %s", expected, got)
	}
}

func TestEstraiMeseNascita(t *testing.T) {
	var utente Utente

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

	utente.MeseNascita = 1
	for i := 0; i < 12; i++ {
		got := utente.estraiMeseNascita().CodFiscale
		expected := mappaMesi[utente.MeseNascita]

		if got != string(expected) {
			t.Errorf("error: expected %s got %s", string(expected), got)
		}

		utente.CodFiscale = ""
		utente.MeseNascita++
	}
}

func TestEstraiGiornoNascitaUomo(t *testing.T) {
	var utente Utente

	utente.GiornoNascita = 5
	utente.Sesso = "maschile"
	expected := "05"
	got := utente.estraiGiornoNascita().CodFiscale

	if got != expected {
		t.Errorf("error: expected %s got %s", string(expected), got)
	}
}

func TestEstraiGiornoNascitaDonna(t *testing.T) {
	var utente Utente

	utente.GiornoNascita = 25
	utente.Sesso = "femminile"
	expected := "65"
	got := utente.estraiGiornoNascita().CodFiscale

	if got != expected {
		t.Errorf("error: expected %s got %s", string(expected), got)
	}
}

func TestEstraiCodiceControllo(t *testing.T) {
	var utente Utente

	utente.CodFiscale = "RSOMRA76B45G478"
	expected := "RSOMRA76B45G478I"
	got := utente.estraiCodiceControllo().CodFiscale

	if got != expected {
		t.Errorf("error: expected %s got %s", string(expected), got)
	}
}
