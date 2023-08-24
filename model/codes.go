package model

import (
	"database/sql"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Dato un comune, estrai il codice catastale corrispondente
func EstraiCodiceCatastale(comune string) string {
	db, err := sql.Open("sqlite3", "codes.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT Code FROM codCatastali WHERE UPPER(City) = ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	var code string
	err = stmt.QueryRow(strings.ToUpper(comune)).Scan(&code)
	if err != nil {
		if err == sql.ErrNoRows {
			return ""
		} else {
			panic(err)
		}
	}

	return code
}

// Dato un codice catastale, estrai il comune corrispondente
func EstraiComune(codCatastale string) string {
	db, err := sql.Open("sqlite3", "codes.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT City FROM codCatastali WHERE Code = ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	var city string
	err = stmt.QueryRow(strings.ToUpper(codCatastale)).Scan(&city)
	if err != nil {
		if err == sql.ErrNoRows {
			return ""
		} else {
			panic(err)
		}
	}

	return city
}

// Data una nazione, estrai il codice nazionale
func EstraiCodiceNazione(nazione string) string {
	db, err := sql.Open("sqlite3", "codes.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT Code FROM codNazioni WHERE UPPER(State) = ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	var code string
	err = stmt.QueryRow(strings.ToUpper(nazione)).Scan(&code)
	if err != nil {
		if err == sql.ErrNoRows {
			return ""
		} else {
			panic(err)
		}
	}

	return code
}

// Dato il codice di una nazione, estrai la nazione corrispondente
func EstraiNazione(codNazione string) string {
	db, err := sql.Open("sqlite3", "codes.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT State FROM codNazioni WHERE Code = ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	var state string
	err = stmt.QueryRow(strings.ToUpper(codNazione)).Scan(&state)
	if err != nil {
		if err == sql.ErrNoRows {
			return ""
		} else {
			panic(err)
		}
	}

	return state
}
