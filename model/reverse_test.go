package model

import "testing"

func TestInversoAnnoNascita19(t *testing.T) {
	var inverso Inverso

	inverso.CodFiscale = "RSOMRA79D65I726V"
	expected := "1979"
	got := inverso.estraiAnnoNascita().AnnoNascita

	if expected != got {
		t.Errorf("error: expected %s got %s", expected, got)
	}
}

func TestInversoAnnoNascita20(t *testing.T) {
	var inverso Inverso

	inverso.CodFiscale = "RSOMRA05D65I726V"
	expected := "2005"
	got := inverso.estraiAnnoNascita().AnnoNascita

	if expected != got {
		t.Errorf("error: expected %s got %s", expected, got)
	}
}

func TestInversoMeseNascita(t *testing.T) {
	var inverso Inverso

	inverso.CodFiscale = "RSOMRA05D65I726V"
	expected := 4
	got := inverso.estraiMeseNascita().MeseNascita

	if expected != int(got) {
		t.Errorf("error: expected %d got %d", expected, got)
	}
}

func TestInversoGiornoNascitaDonna(t *testing.T) {
	var inverso Inverso

	inverso.CodFiscale = "RSOMRA05D65I726V"
	expected := 25
	got := inverso.estraiGiornoNascita().GiornoNascita

	if expected != int(got) {
		t.Errorf("error: expected %d got %d", expected, got)
	}
}

func TestInversoGiornoNascitaUomo(t *testing.T) {
	var inverso Inverso

	inverso.CodFiscale = "RSOMRA05D10I726V"
	expected := 10
	got := inverso.estraiGiornoNascita().GiornoNascita

	if expected != int(got) {
		t.Errorf("error: expected %d got %d", expected, got)
	}
}

func TestInversoGiornoNascitaInvalid(t *testing.T) {
	var inverso Inverso

	inverso.CodFiscale = "RSOMRA05D80I726V"
	expected := 0
	got := inverso.estraiGiornoNascita().GiornoNascita

	if expected != int(got) || inverso.Errore == "" {
		t.Errorf("error: expected %d got %d", expected, got)
	}
}

func TestInversoSessoUomo(t *testing.T) {
	var inverso Inverso

	inverso.CodFiscale = "RSOMRA05D10I726V"
	expected := "maschile"
	got := inverso.estraiSesso().Sesso

	if expected != got {
		t.Errorf("error: expected %s got %s", expected, got)
	}
}

func TestInversoSessoDonna(t *testing.T) {
	var inverso Inverso

	inverso.CodFiscale = "RSOMRA05D45I726V"
	expected := "femminile"
	got := inverso.estraiSesso().Sesso

	if expected != got {
		t.Errorf("error: expected %s got %s", expected, got)
	}
}
