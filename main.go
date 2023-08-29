package main

import (
	"net/http"
	"os"

	"github.com/ice-bit/wheezy/controller"
	"github.com/ice-bit/wheezy/log"
)

func main() {
	var (
		host       = os.Getenv("WHEEZY_LISTEN_ADDRESS")
		port       = os.Getenv("WHEEZY_LISTEN_PORT")
		redis_host = os.Getenv("WHEEZY_REDIS_ADDRESS")
		redis_port = os.Getenv("WHEEZY_REDIS_PORT")
	)

	// Controlla che le env vars siano definite
	if host == "" || port == "" || redis_host == "" || redis_port == "" {
		log.ErrLogger.Printf("environment variables not configured")
		panic("Environment variables not configured")
	}

	// Definisci le rotte
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", controller.RootHandler)
	http.HandleFunc("/reverse", controller.ReverseHandler)
	http.HandleFunc("/about", controller.AboutHandler)

	// Avvia il server
	log.InfoLogger.Printf("server listening on http://%s", (host + ":" + port))
	http.ListenAndServe(host+":"+port, nil)
}
