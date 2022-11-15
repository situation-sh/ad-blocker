package main

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

type myCustomServer struct {
	addr string
}

// ServeDNS is used to implement the Handler interface
func (m myCustomServer) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	if r.Question != nil {
		// Construire la reponse en commencant par un forward sur 8.8.8.8
		//  resourceRecord := []dns.RR{"8.8.8.8"}
		// 	fmt.Printf("%v\n", r)
		// 	r.Answer = append(r.Answer, []dns.RR{"8.8.8.8"} )
		w.WriteMsg(r)
	}
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
	server := myCustomServer{addr: "127.0.0.1:8053"}
	// Initialize the Dns server object from the server above
	dnsServer := dns.Server{Addr: server.addr, Net: "udp", Handler: server, NotifyStartedFunc: server.starterMessage}
	// fmt.Printf("%+v \n", dnsServer)
	dnsServer.ListenAndServe()
	dnsServer.NotifyStartedFunc()
}

// 1. Pouvoir instancier un server Dns
// Utiliser l'objet serveur dans la lib miekg : Address, network,  à dispositions sur l'objet.
// Creer un objet bidon qui implemente l'interface Handler (grace a serveDns())
// NotifyStartedFunc => Message d'intro

// Utiliser la function ListenAndServe sur le server.
// 1. Forwarder la requete vers un autre serveur
// 2. Regarder la question Dns
// 3. Faire la réponse (regarder comment remplir le packetDns de retour)
