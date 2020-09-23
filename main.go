package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/log"
	"github.com/prometheus/common/version"
	"github.com/tcassaert/bmwcd_exporter/bmwcd"
)

var (
	help     = flag.Bool("help", false, "Print help message")
	logLevel = flag.String("log.level", "INFO", "Amount of logs displayed")
	password = flag.String("password", "", "BMW Connected Drive password")
	port     = flag.String("port", "9744", "Exporter port")
	region   = flag.String("region", "rest_of_world", "Region of the Connected Drive account (cn, rest_of_world, us)")
	username = flag.String("username", "", "BMW Connected Drive username")
)

func main() {
	flag.Parse()

	if *help == true {
		flag.Usage()
		os.Exit(0)
	}

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

	collector := bmwcd.NewCollector(*username, *password, *region)
	prometheus.MustRegister(collector)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			 <head><title>BMW Connected Drive Exporter</title></head>
			 <body>
			 <h1>BMW Connected Drive Exporter</h1>
			 <p><a href='/metrics'>Metrics</a></p>
			 </body>
			 </html>`))
	})
	http.ListenAndServe(fmt.Sprintf(":%s", *port), nil)
}
