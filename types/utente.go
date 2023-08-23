package types

type Utente struct {
	Cognome       string
	Nome          string
	Sesso         string
	LuogoNascita  string
	GiornoNascita uint8
	MeseNascita   uint8
	AnnoNascita   uint64
	CodFiscale    string
}
