package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/miekg/dns"
)

func setBlackList(mux *dns.ServeMux) {
	// HINT: Retrieve the blackList from : "https://hosts.anudeep.me/mirror/adservers.txt"
	// Url and filename
	url, filepath := "https://hosts.anudeep.me/mirror/adservers.txt", "./adservers.txt"

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		log.Fatal(resp.Status)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Loop over all lines, populate the list Array with domains only
	f, _ := os.Open(filepath)
	scanner := bufio.NewScanner(f)
	list := []string{}
	for scanner.Scan() {
		line := strings.Join(strings.Split(scanner.Text(), "0.0.0.0 "), "")
		list = append(list, line)
	}
	// remove the 10 first Line of the file
	list = list[10:]

	// Create a handle function for each domains of the list and attach it to the HandleFunc of the mux.
	for i := 0; i < len(list); i++ {
		mux.HandleFunc(list[i], interceptRequest)

	}
	// Create a handle function for request to be forwarded
	mux.HandleFunc(".", forwardRequest)
}

// Intercept the query to redirect to httpforever
func interceptRequest(w dns.ResponseWriter, r *dns.Msg) {
	// Retrieved the query domain name
	domain := r.Question[0].Name
	fmt.Println(" \n#### Intercepted query  =>", domain)
	// Initiliaze an answer from type dns.Msg
	m := dns.Msg{}
	m.SetReply(r)
	// Populate the answer based on the query type (A, AAAA or HTTPS)
	switch r.Question[0].Qtype {
	case 1: // Type A
		m.Answer = []dns.RR{&dns.A{A: net.IP{172, 67, 181, 181}, Hdr: dns.RR_Header{Name: domain, Rrtype: 1, Class: 1, Ttl: 5}}}
	case 28: // Type AAAA
		m.Answer = []dns.RR{&dns.AAAA{AAAA: net.ParseIP("2a00:1450:4007:80b::2003"), Hdr: dns.RR_Header{Name: "google.com.", Rrtype: 28, Class: 1, Ttl: 5}}}
		// case 65: // Type HTTPS
		// 		m.Answer = []dns.RR{&dns.HTTPS{: net.IP{216, 58, 214, 174}, Hdr: dns.RR_Header{Name: "google.com.", Rrtype: 1, Class: 1, Ttl: 3600}}}
	default:
		m.Answer = []dns.RR{&dns.A{A: net.IP{172, 67, 181, 181}, Hdr: dns.RR_Header{Name: domain, Rrtype: 1, Class: 1, Ttl: 5}}}
	}

	fmt.Println(" \n#### Custom Answer  =>", m.Answer[0].String())
	// Send back the answer to the client
	if err := w.WriteMsg(&m); err != nil {
		fmt.Printf("ERROR : Answer could not be sent back to the client :\n => %v\n", err)
	}
}

// Forward the request to the google dns resolver
func forwardRequest(w dns.ResponseWriter, r *dns.Msg) {
	domain := r.Question[0].Name
	fmt.Println(" \n#### Forwarded query  =>", domain)

	// Forward to google dns service
	if response, err := dns.Exchange(r, "1.1.1.1:53"); err != nil {
		fmt.Printf("ERROR : Query could not be forwarded to the Google dns resolver :\n => %v\n", err)
	} else {
		if err := w.WriteMsg(response); err != nil {
			fmt.Printf("ERROR : Answer from the google dns resolver could not be sent back to the client :\n => %v\n", err)
		}
	}
}

func main() {
	// TODO: Parse Flags et valeur par default

	// Set an adress for the dns server.
	address := "127.0.0.1:53"

	// Initialize the mux (Rooting sytem for incoming http adress. Stands on top of the server). It will be in charge of the request dispatch.
	mux := dns.NewServeMux()
	setBlackList(mux)

	// Initialize the dns server.
	dnsServer := dns.Server{Addr: address, Net: "udp", Handler: mux, NotifyStartedFunc: func() {
		fullAddress, port, err := net.SplitHostPort(address)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("#### The Dns Server is hosted at => ", fullAddress)
		fmt.Println("#### The Dns Server is listenning on port => ", port)
	}}

	if err := dnsServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
