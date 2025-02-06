package main

import (
	"context"
	"flag"
	"fmt"
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
	migrate := flag.Bool("migrate", false, "migration db on init")
	dbPath := flag.String("db-file", "test.db", "sqlite3 file path")
	flag.Parse()

	r := http.NewServeMux()

	client, err := ent.Open("sqlite3", fmt.Sprintf("%s?_fk=1", *dbPath))
	if err != nil {
		log.Fatalf("failed opening connection to postgres: %v", err)
	}
	defer client.Close()

	if *migrate {
		if err := client.Schema.Create(context.Background()); err != nil {
			log.Fatalf("failed creating schema resources: %v", err)
		}
		log.Println("migrate done")
	}

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
