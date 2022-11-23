package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"
	"github.com/miekg/dns"
)

const (
	defaultHost       = "127.0.0.1"
	port              = "53"
	redirectionOrigin = "neverssl.com"
	redirectionTarget = "httpforever.com"
)

// TODO:  Remplacer les Fmt.Print par la lib "zeroLog"
// TODO:  Exporter toutes les constantes
func main() {
	// Set & Parse Flags with default value
	hostFlagPtr := flag.String("h", defaultHost, "flag to define the host of the server")
	portFlagPtr := flag.String("p", port, "flag to define the port of the service")
	// updateFlagPtr := flag.Bool("u", false, "flag to force the fetch of adservers.txt") // If true, trigger la fonction.
	flag.Parse()

	// Set an adress for the dns server.
	address := *hostFlagPtr + ":" + *portFlagPtr

	// Initialize the mux (Rooting sytem for incoming http adress. Stands on top of the server). It will be in charge of the requests dispatch.
	// Attach the blackList
	mux := dns.NewServeMux()
	// Attacher les handlers
	setBlackList(mux)
	mux.HandleFunc(redirectionOrigin, redirectRequest)
	mux.HandleFunc(".", forwardRequest)

	// Initialize the dns server.
	dnsServer := dns.Server{Addr: address, Net: "udp", Handler: mux, NotifyStartedFunc: func() {
		fullAddress, port, err := net.SplitHostPort(address)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("#### The Dns Server is hosted at => ", fullAddress)
		fmt.Println("#### The Dns Server is listenning on port => ", port)
	}}
	fmt.Println(dnsServer)

	// Start the proxy in a go-routine. It will achieve the redirection by changing the header of the request.
	go func() {
		proxy := goproxy.NewProxyHttpServer()
		proxy.OnRequest().DoFunc(func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			if strings.HasSuffix(r.Host, redirectionOrigin) {
				r.Host = redirectionTarget
				fmt.Println(r)
			}
			return r, nil
		})
		log.Fatal(http.ListenAndServe("127.0.0.1:53120", proxy))
	}()

	// Start the dns server.
	if err := dnsServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
