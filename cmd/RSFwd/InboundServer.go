package main

import (
	"time"
	"github.com/miekg/dns"
	"sync"
	log "github.com/Sirupsen/logrus"
	"os"
)

type Server struct {
	Addr	string
	Soreuseport	bool
	ConnTimeout	time.Duration
	debug	bool

}


func NewServer(addr string)(*Server, error){
	s := &Server{Addr: addr, ConnTimeout: 5*time.Second}

	return s, nil
}

func (s *Server)Run(){
	mux := dns.NewServeMux();
	mux.Handle(".", s)

	wg := new(sync.WaitGroup)
	wg.Add(2)

	log.Info("Start overture on " + s.Addr)

	for _, p := range [2]string{"tcp", "udp"} {
		go func(p string) {
			err := dns.ListenAndServe(s.Addr, p, mux)
			if err != nil {
				log.Fatal("Listen "+p+" failed: ", err)
				os.Exit(1)
			}
		}(p)
	}

}

func (s *Server)GetDefaultFwd()(string){
	return "8.8.8.8:53"
}

func (s *Server)Fetch(fwd string, r *dns.Msg)(*dns.Msg){
	c := new(dns.Client)
	in, _, _ := c.Exchange(r, fwd)

	return in
}

// ServeDNS is the entry point for every request to the address that s
// is bound to. It acts as a multiplexer for the requests zonename as
// defined in the request so that the correct zone

func (s *Server) ServeDNS(w dns.ResponseWriter, r *dns.Msg) {
	if r == nil || len(r.Question) == 0 || r.MsgHdr.Response == true {
		return
	}

	fwd := s.GetDefaultFwd()

	in := s.Fetch(fwd, r)

	w.WriteMsg(in)
}