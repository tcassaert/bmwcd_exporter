package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"github.com/tcassaert/bmwcd_exporter/bmwcd"
)

var (
	password = flag.String("password", "", "BMW Connected Drive password")
	port     = flag.String("port", "9744", "Exporter port")
	username = flag.String("username", "", "BMW Connected Drive username")
	logLevel = flag.String("log.level", "INFO", "Amount of logs displayed")
)

func main() {
	flag.Parse()
	log.Base().SetLevel(*logLevel)
	log.Infoln("Starting BMW Connected Drive exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	if *username == "" {
		log.Errorln("Please provide a username")
		os.Exit(1)
	}

	if *password == "" {
		log.Errorln("Please provide a password")
		os.Exit(1)
	}

	go bmwcd.StartPolling(*username, *password)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(fmt.Sprintf(":%s", *port), nil)
}
