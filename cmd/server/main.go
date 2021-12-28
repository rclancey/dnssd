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
	m, err := dnssd.Register("sample", 8192)
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

/*
func zeroconf(name string, port int) (func(), error) {
	cfg := dnssd.Config{
		Name: fmt.Sprintf("%s server", name),
		Type: fmt.Sprintf("_%s._tcp", name),
		Port: port,
	}
	sv, err := dnssd.NewService(cfg)
	if err != nil {
		return nil, err
	}
	rp, err := dnssd.NewResponder()
	if err != nil {
		return nil, err
	}
	_, err = rp.Add(sv)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(context.Background())
	rp.Respond(ctx)
	return cancel, nil
}
*/

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
	/*
	host, err := os.Hostname()
	if err != nil {
		return nil, nil, err
	}
	info := []string{fmt.Sprintf("%s server", name)}
	service, err := mdns.NewMDNSService(host, fmt.Sprintf("_%s._tcp", name), "", "", port, nil, info)
	if err != nil {
		return nil, nil, err
	}
	log.Printf("mdns info: %#v", *service)
	m, err := mdns.NewServer(&mdns.Config{Zone: service})
	return server, m, nil
	*/
}
