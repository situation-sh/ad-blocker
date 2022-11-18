package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/miekg/dns"
)

type myCustomServer struct {
	addr string
}

// ServeDNS is used to implement the Handler interface
func (m myCustomServer) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {

	query := r.Question[0].Name
	fmt.Println(r.Question[0].Qtype)

	// id := r.Id
	fmt.Println(query)
	if query == "www.situation.sh." {
		// Redirect to google.com 216.58.214.174
		m := new(dns.Msg)
		m.SetReply(r)
		// A = 1, AAAA = 28, HTTPS = 65

		// if r.Question[0].Qtype == 65 {
		// 	m.Answer = []dns.RR{&dns.HTTPS{: net.IP{216, 58, 214, 174}, Hdr: dns.RR_Header{Name: "google.com.", Rrtype: 1, Class: 1, Ttl: 3600}}}
		// 	} else if r.Question[0].Qtype == 28 {
		// 		m.Answer = []dns.RR{&dns.AAAA{A: net.IP{216, 58, 214, 174}, Hdr: dns.RR_Header{Name: "google.com.", Rrtype: 1, Class: 1, Ttl: 3600}}}
		// 		} else {
		// 			m.Answer = []dns.RR{&dns.A{A: net.IP{216, 58, 214, 174}, Hdr: dns.RR_Header{Name: "google.com.", Rrtype: 1, Class: 1, Ttl: 3600}}}
		// }

		m.Answer = []dns.RR{&dns.A{A: net.IP{216, 58, 214, 174}, Hdr: dns.RR_Header{Name: "google.com.", Rrtype: 1, Class: 1, Ttl: 3600}}}
		m.RecursionAvailable = true
		m.RecursionDesired = true
		m.Authoritative = true
		m.AuthenticatedData = true

		fmt.Println(m)
		err := w.WriteMsg(m)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
	} else {
		// Forward to google dns service
		response, err := dns.Exchange(r, "1.1.1.1:53")
		if err != nil {
			fmt.Printf("%v\n", err)
		} else {
			err := w.WriteMsg(response)
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}
	}

	// fmt.Println(message)

}

func (m myCustomServer) parsePort() string {
	return strings.Split(m.addr, ":")[1]
}

func (m myCustomServer) parseIp() string {
	return strings.Split(m.addr, ":")[0]
}

func (m myCustomServer) starterMessage() {
	fmt.Println("#### The Dns Server is hosted at => ", m.parseIp())
	fmt.Println("#### The Dns Server is listenning on port => ", m.parsePort())
}

func main() {
	// Create a server
	server := myCustomServer{addr: "127.0.0.1:53"}
	// Initialize the Dns server object from the server above
	dnsServer := dns.Server{Addr: server.addr, Net: "udp", Handler: server, NotifyStartedFunc: server.starterMessage}

	err := dnsServer.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	} else {
		dnsServer.NotifyStartedFunc()
	}
}

// 1. Pouvoir instancier un server Dns
// Utiliser l'objet serveur dans la lib miekg : Address, network,  à dispositions sur l'objet.
// Creer un objet bidon qui implemente l'interface Handler (grace a serveDns())
// NotifyStartedFunc => Message d'intro

// Utiliser la function ListenAndServe sur le server.
// 1. Forwarder la requete vers un autre serveur
// 2. Regarder la question Dns
// 3. Faire la réponse (regarder comment remplir le packetDns de retour)
