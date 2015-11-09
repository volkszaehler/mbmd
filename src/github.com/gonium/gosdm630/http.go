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
			v.Voltage.L1, v.Current.L1, v.Power.L1, v.Cosphi.L1)
		fmt.Fprintf(w, "L2: %.2fV %.2fA %.2fW %.2fcos\r\n",
			v.Voltage.L2, v.Current.L2, v.Power.L2, v.Cosphi.L2)
		fmt.Fprintf(w, "L3: %.2fV %.2fA %.2fW %.2fcos\r\n",
			v.Voltage.L3, v.Current.L3, v.Power.L3, v.Cosphi.L3)
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
