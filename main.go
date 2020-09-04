package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tcassaert/bmwcd_exporter/bmwcd"
)

var (
	password = flag.String("password", "", "BMW Connected Drive password")
	port     = flag.String("port", "9744", "Exporter port")
	username = flag.String("username", "", "BMW Connected Drive username")
)

func main() {
	flag.Parse()

	if *username == "" {
		flag.Usage()
		log.Fatal("ERROR: Please provide a username")
	}

	if *password == "" {
		flag.Usage()
		log.Fatal("ERROR: Please provide a password")
	}

	go bmwcd.StartPolling(*username, *password)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%s", *port), nil)
}
