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

// TODO:  Exporter toutes les constantes

// TODO:  IF adservers.txt is in the repo, THEN ne pas fetcher le fichier
// TODO:  IF -update flag is true THEN appeler une fonction de fetch adservers.txt
// TODO:  Remplacer les Fmt.Print par la lib "zeroLog"
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
		mux.HandleFunc(list[i], blockRequest)

	}
}

// TODO:  Retourner une erreur DNS clean lorsqu'une requete est bloqués
// TODO:  Remplacer les Fmt.Print par la lib "zeroLog"
func blockRequest(w dns.ResponseWriter, r *dns.Msg) {
	// Retrieved the query domain name
	domain := r.Question[0].Name
	fmt.Println(" \n#### BLOCKED query  =>", domain)
	// Initiliaze an answer from type dns.Msg
	m := dns.Msg{}
	m.SetReply(r)
	// Populate the answer based on the query type (A, AAAA or HTTPS)
	switch r.Question[0].Qtype {
	case dns.TypeA:
		m.Answer = []dns.RR{&dns.A{A: net.IP{172, 67, 181, 181}, Hdr: dns.RR_Header{Name: domain, Rrtype: 1, Class: 1, Ttl: 5}}}
	case dns.TypeAAAA:
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

// Intercept the query to redirect to httpforever
// TODO:  Remplacer les Fmt.Print par la lib "zeroLog"
func redirectRequest(w dns.ResponseWriter, r *dns.Msg) {
	// Retrieved the query domain name
	domain := r.Question[0].Name
	fmt.Println(" \n#### INTERCEPTED query  =>", domain)
	// Initiliaze an answer from type dns.Msg
	m := dns.Msg{}
	m.SetReply(r)
	// Populate the answer based on the query type (A, AAAA or HTTPS)
	switch questionType := r.Question[0].Qtype; questionType {
	case dns.TypeA: // Type A
		m.Answer = []dns.RR{&dns.A{A: net.IP{172, 67, 181, 181}, Hdr: dns.RR_Header{Name: domain, Rrtype: questionType, Class: dns.ClassINET, Ttl: 5}}}
	case dns.TypeAAAA: // Type AAAA
		m.Answer = []dns.RR{&dns.AAAA{AAAA: net.ParseIP("2a00:1450:4007:80b::2003"), Hdr: dns.RR_Header{Name: "google.com.", Rrtype: questionType, Class: 1, Ttl: 5}}}
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
// TODO:  Remplacer les Fmt.Print par la lib "zeroLog"
func forwardRequest(w dns.ResponseWriter, r *dns.Msg) {
	domain := r.Question[0].Name
	fmt.Println(" \n#### FORWARDED query  =>", domain)

	// Forward to google dns service
	if response, err := dns.Exchange(r, "1.1.1.1:53"); err != nil {
		fmt.Printf("ERROR : Query could not be forwarded to the Google dns resolver :\n => %v\n", err)
	} else {
		if err := w.WriteMsg(response); err != nil {
			fmt.Printf("ERROR : Answer from the google dns resolver could not be sent back to the client :\n => %v\n", err)
		}
	}
}
