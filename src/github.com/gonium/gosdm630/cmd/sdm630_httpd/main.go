package main

import (
	"flag"
	"fmt"
	"github.com/gonium/gosdm630"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

var rtuDevice = flag.String("rtuDevice", "/dev/ttyUSB0", "Path to serial RTU device")
var verbose = flag.Bool("verbose", false, "Enables extensive logging")

func init() {
	flag.Parse()
}

func MkIndexHandler(hc *sdm630.MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		v := hc.GetLast()
		fmt.Fprintf(w, "Last measurement taken %s:\r\n", v.Time.Format(time.RFC850))
		fmt.Fprintf(w, "L1: %.2fV %.2fA %.2fW %.2fcos\r\n",
			v.L1Voltage, v.L1Current, v.L1Power, v.L1CosPhi)
		fmt.Fprintf(w, "L2: %.2fV %.2fA %.2fW %.2fcos\r\n",
			v.L2Voltage, v.L2Current, v.L2Power, v.L2CosPhi)
		fmt.Fprintf(w, "L3: %.2fV %.2fA %.2fW %.2fcos\r\n",
			v.L3Voltage, v.L3Current, v.L3Power, v.L3CosPhi)
	})
}

func MkLastValueHandler(hc *sdm630.MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		last := hc.GetLast()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := last.JSON(w); err != nil {
			log.Printf("Failed to create JSON representation of measurement %s", last.String())
		}
	})
}

func MkLastMinuteAvgHandler(hc *sdm630.MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		last := hc.GetMinuteAvg()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := last.JSON(w); err != nil {
			log.Printf("Failed to create JSON representation of measurement %s", last.String())
		}
	})
}

func main() {
	var rc = make(sdm630.ReadingChannel)
	// Read the SDM630 every 5 seconds
	interval := 5
	qe := sdm630.NewQueryEngine(*rtuDevice, interval, *verbose, rc)
	go qe.Produce()
	hc := sdm630.NewMeasurementCache(rc, interval)
	go hc.ConsumeData()

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", MkIndexHandler(hc))
	router.HandleFunc("/last", MkLastValueHandler(hc))
	router.HandleFunc("/minuteavg", MkLastMinuteAvgHandler(hc))
	log.Fatal(http.ListenAndServe(":8080", router))
}
