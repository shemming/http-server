package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/shemming/singularity6/proxy"
)

func main() {
	address := "127.0.0.1:9000"
	r := mux.NewRouter()
	r.HandleFunc("/", proxy.HelloWorld)

	log.Printf("Starting http server...")
	log.Printf("Listening on: %s", address)
	server := http.Server{
		Handler: r,
		Addr:    address,
	}

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

}
