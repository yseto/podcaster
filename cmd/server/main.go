package main

import (
	"flag"
	"log"
	"net"
	"net/http"

	_ "github.com/mattn/go-sqlite3"

	"github.com/yseto/podcaster/ent"
	"github.com/yseto/podcaster/server"
	"github.com/yseto/podcaster/server/api"
)

func main() {
	port := flag.String("port", "8080", "port where to serve traffic")
	flag.Parse()

	r := http.NewServeMux()

	client, err := ent.Open("sqlite3", "test.db?_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	svr := server.NewServer(client)
	h := api.HandlerFromMux(api.NewStrictHandler(svr, nil), r)

	mw, err := server.CreateMiddleware(client)
	if err != nil {
		log.Fatalln("error creating middleware:", err)
	}
	h = mw(h)

	mw2 := server.CreateMiddlewareEmptyContext()
	h = mw2(h)

	s := &http.Server{
		Handler: h,
		Addr:    net.JoinHostPort("0.0.0.0", *port),
	}

	log.Fatal(s.ListenAndServe())
}
