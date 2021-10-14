package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
)

var (
	hostport = flag.String("a", ":8080", "address to listen on")
	serveDir = flag.String("d", ".", "serve directory")
)

func init() {
	flag.Parse()
	log.SetFlags(log.Lshortfile | log.LUTC)
}

func main() {
	listAddresses()
	serveHTTP()
}

func listAddresses() {
	const httpPrefix = "http://"

	host, port, err := net.SplitHostPort(*hostport)
	if err != nil {
		log.Fatalf("parse address: %v", err)
	}

	ip := net.ParseIP(host)
	if host != "" && (ip == nil || !ip.IsUnspecified()) {
		fmt.Println(httpPrefix + net.JoinHostPort(host, port))
		return
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		log.Fatalf("list addresses: %v", err)
	}

	for _, addr := range addrs {
		addr, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}

		ip := addr.IP
		if ip.IsLinkLocalUnicast() || ip.IsMulticast() {
			continue
		}

		fmt.Println(httpPrefix + net.JoinHostPort(ip.String(), port))
	}
}

func serveHTTP() {
	err := http.ListenAndServe(*hostport, http.FileServer(http.Dir(*serveDir)))
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("serve http: %v", err)
	}
}
