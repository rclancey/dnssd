package main

import (
	//"context"
	//"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rclancey/dnssd/v2"
)

/*
func browse(name string, timeout time.Duration) (chan dnssd.BrowseEntry, func(), error) {
	ch := make(chan dnssd.BrowseEntry, 5)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	service := fmt.Sprintf("_%s._tcp.local.", name)
	addFn := func(e dnssd.BrowseEntry) {
		ch <- e
	}
	rmvFn := func(e dnssd.BrowseEntry) {
		fmt.Printf("rremoved %#v\n", e)
	}
	go func() {
		err := dnssd.LookupType(ctx, service, addFn, rmvFn)
		if err != nil {
			if !errors.Is(err, context.Canceled) {
				log.Println(err)
			}
			close(ch)
		}
	}()
	return ch, cancel, nil
}
*/

func main() {
	name := os.Args[1]
	results := dnssd.Browse(dnssd.LookupRequest{Name: name, Timeout: time.Second})
	client := &http.Client{}
	for _, res := range results {
		if len(res.IPs[0]) == 4 {
			log.Printf("%s:%d", res.IPs[0], res.Port)
			res, err := client.Get(fmt.Sprintf("http://%s:%d/", res.IPs[0], res.Port))
			if err != nil {
				log.Println(err)
			} else {
				data, err := ioutil.ReadAll(res.Body)
				res.Body.Close()
				if err != nil {
					log.Println(err)
				} else {
					log.Println(string(data))
				}
			}
		}
		//log.Printf("got result %#v", *res)
	}
	/*
	ch, cancel, err := browse(name, time.Second * 5)
	if err != nil {
		log.Fatal(err)
	}
	*/
	/*
	ch := make(chan *mdns.ServiceEntry, 4)
	defer close(ch)
	mdns.Lookup(fmt.Sprintf("_%s._tcp", name), ch)
	timer := time.NewTimer(5 * time.Second)
	for {
		select {
		case <-timer.C:
			return
		case entry, ok := <-ch:
			if !ok {
				return
			}
			fmt.Printf("Got new entry: %#v\n", entry)
		}
	}
	*/
}
