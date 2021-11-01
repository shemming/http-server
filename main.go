package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/shemming/http-server/proxy"
)

func main() {
	ctx := context.Background()
	done := make(chan os.Signal, 1)

	// seed random to get different strings at each launch of service
	rand.Seed(time.Now().UnixNano())

	signal.Notify(done,
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	httpServer := startHTTPServer(ctx)

	<-done

	defer close(done)
	err := httpServer.Shutdown(ctx)
	if err != nil {
		log.Fatal("error shutting down http server")
	}

	fmt.Println("exiting")
}

func startHTTPServer(ctx context.Context) *http.Server {
	address := "127.0.0.1:4000"

	r := mux.NewRouter()
	proxy := proxy.NewProxy(r)

	log.Printf("Starting http server...")
	server := &http.Server{
		Handler: proxy.Router,
		Addr:    address,
	}

	// run the server
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	log.Printf("Listening on: %s", address)

	return server
}
