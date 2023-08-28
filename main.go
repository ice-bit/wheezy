package main

import (
	"flag"
	"net/http"

	"github.com/ice-bit/wheezy/controller"
)

func main() {
	var host string
	var port string

	// Leggi host e porta da CLI
	flag.StringVar(&host, "l", "127.0.0.1", "Specify listening address.")
	flag.StringVar(&port, "p", "9000", "Specify listening port.")
	flag.Parse()

	// Definisci le rotte
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", controller.RootHandler)
	http.HandleFunc("/reverse", controller.ReverseHandler)
	http.HandleFunc("/about", controller.AboutHandler)

	// Avvia il server
	http.ListenAndServe(host+":"+port, nil)
}
