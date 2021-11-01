package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gorilla/mux"
	"github.com/shemming/http-server/proxy"
)

func main() {
	ctx := context.Background()
	done := make(chan os.Signal, 1)

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
	address := "127.0.0.1:9000"
	proxy := proxy.NewProxy()

	r := mux.NewRouter()
	r.Handle("/", http.HandlerFunc(proxy.HelloWorld)).Methods(http.MethodPost)

	log.Printf("Starting http server...")
	server := &http.Server{
		Handler: r,
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
