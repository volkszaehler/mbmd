package sdm630

import (
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func MkIndexHandler(hc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		v := hc.GetLast()
		fmt.Fprintf(w, "Last measurement taken %s:\r\n", v.Timestamp.Format(time.RFC850))
		fmt.Fprintf(w, "L1: %.2fV %.2fA %.2fW %.2fcos\r\n",
			v.L1Voltage, v.L1Current, v.L1Power, v.L1CosPhi)
		fmt.Fprintf(w, "L2: %.2fV %.2fA %.2fW %.2fcos\r\n",
			v.L2Voltage, v.L2Current, v.L2Power, v.L2CosPhi)
		fmt.Fprintf(w, "L3: %.2fV %.2fA %.2fW %.2fcos\r\n",
			v.L3Voltage, v.L3Current, v.L3Power, v.L3CosPhi)
	})
}

func MkLastValueHandler(hc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		last := hc.GetLast()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := last.JSON(w); err != nil {
			log.Printf("Failed to create JSON representation of measurement %s", last.String())
		}
	})
}

func MkLastMinuteAvgHandler(hc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		last := hc.GetMinuteAvg()
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := last.JSON(w); err != nil {
			log.Printf("Failed to create JSON representation of measurement %s", last.String())
		}
	})
}

func Run_httpd(mc *MeasurementCache, url string) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", MkIndexHandler(mc))
	router.HandleFunc("/last", MkLastValueHandler(mc))
	router.HandleFunc("/minuteavg", MkLastMinuteAvgHandler(mc))
	log.Fatal(http.ListenAndServe(url, router))
}
