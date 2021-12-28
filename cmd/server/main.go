package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rclancey/dnssd/v2"
)

func main() {
	h, err := makeServer("sample", 8192)
	if err != nil {
		log.Fatal(err)
	}
	info := map[string]string{
		"ssl": "false",
	}
	m, err := dnssd.Register("sample", 8192, info)
	if err != nil {
		log.Fatal(err)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	handled := false
	go func() {
		err := h.ListenAndServe()
		if err != nil {
			if !handled {
				log.Fatal(err)
			}
		}
	}()
	<-stop
	log.Println("shutting down")
	m() // shutdown dnssd
	ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
	defer cancel()
	err = h.Shutdown(ctx)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("exiting")
}

func makeServer(name string, port int) (*http.Server, error) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	}
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", port),
		Handler: http.HandlerFunc(handler),
	}
	return server, nil
}
