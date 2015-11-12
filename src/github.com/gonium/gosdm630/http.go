package sdm630

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/olekukonko/tablewriter"
	"log"
	"net/http"
	"time"
)

// formatFloat helper
func fF(val float32) string {
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "%.2f", val)
	return buffer.String()
}

func MkIndexHandler(hc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		v := hc.GetLast()
		fmt.Fprintf(w, "Last measurement taken %s:\r\n", v.Timestamp.Format(time.RFC850))
		table := tablewriter.NewWriter(w)
		table.SetHeader([]string{"Phase", "Voltage [V]", "Current [A]", "Power [W]", "Power Factor", "Import [kWh]", "Export [kWh]"})
		table.Append([]string{"L1", fF(v.Voltage.L1), fF(v.Current.L1), fF(v.Power.L1), fF(v.Cosphi.L1), fF(v.Import.L1), fF(v.Export.L1)})
		table.Append([]string{"L2", fF(v.Voltage.L2), fF(v.Current.L2), fF(v.Power.L2), fF(v.Cosphi.L2), fF(v.Import.L2), fF(v.Export.L2)})
		table.Append([]string{"L3", fF(v.Voltage.L3), fF(v.Current.L3), fF(v.Power.L3), fF(v.Cosphi.L3), fF(v.Import.L3), fF(v.Export.L3)})
		table.Append([]string{"ALL", "n/a", fF(v.Current.L1 + v.Current.L2 + v.Current.L3), fF(v.Power.L1 + v.Power.L2 + v.Power.L3), "n/a", fF(v.Import.L1 + v.Import.L2 + v.Import.L3), fF(v.Export.L1 + v.Export.L2 + v.Export.L3)})
		table.Render()
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
