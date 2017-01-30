package sdm630

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/olekukonko/tablewriter"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

// formatFloat helper
func fF(val float64) string {
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "%.2f", val)
	return buffer.String()
}

func MkIndexHandler(hc *MeasurementCache) func(http.ResponseWriter, *http.Request) {

	const mainTemplate = `
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
		 <meta http-equiv="refresh" content="{{.ReloadInterval}}" />
    <title>GoSDM630 overview page</title>
  </head>
  <body>
		<pre>
Reloading every {{.ReloadInterval}} seconds.
{{.Content}}
		</pre>
  </body>
</html>
`
	t, err := template.New("gosdm630").Parse(mainTemplate)
	if err != nil {
		log.Fatal("Failed to create main page template: ", err.Error())
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		var buffer bytes.Buffer
		ids := hc.GetSortedIDs()
		for _, id := range ids {
			v, err := hc.GetLast(id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, err.Error())
				return
			}
			fmt.Fprintf(&buffer, "\r\n\r\nModbus ID %d, last measurement taken %s:\r\n",
				v.ModbusDeviceId, v.Timestamp.Format(time.RFC850))
			table := tablewriter.NewWriter(&buffer)
			table.SetHeader([]string{"Phase", "Voltage [V]", "Current [A]", "Power [W]", "Power Factor", "Import [kWh]", "Export [kWh]"})
			table.Append([]string{"L1", fF(v.Voltage.L1), fF(v.Current.L1), fF(v.Power.L1), fF(v.Cosphi.L1), fF(v.Import.L1), fF(v.Export.L1)})
			table.Append([]string{"L2", fF(v.Voltage.L2), fF(v.Current.L2), fF(v.Power.L2), fF(v.Cosphi.L2), fF(v.Import.L2), fF(v.Export.L2)})
			table.Append([]string{"L3", fF(v.Voltage.L3), fF(v.Current.L3), fF(v.Power.L3), fF(v.Cosphi.L3), fF(v.Import.L3), fF(v.Export.L3)})
			table.Append([]string{"ALL", "n/a", fF(v.Current.L1 + v.Current.L2 + v.Current.L3), fF(v.Power.L1 + v.Power.L2 + v.Power.L3), "n/a", fF(v.Import.L1 + v.Import.L2 + v.Import.L3), fF(v.Export.L1 + v.Export.L2 + v.Export.L3)})
			table.Render()
		}
		data := struct {
			Content        string
			ReloadInterval int
		}{
			Content:        buffer.String(),
			ReloadInterval: 2,
		}
		err := t.Execute(w, data)
		if err != nil {
			log.Fatal("Failed to render main page: ", err.Error())
		}
	})
}

func MkLastAllValuesHandler(hc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		ids := hc.GetSortedIDs()
		lasts := ReadingSlice{}
		for _, id := range ids {
			reading, err := hc.GetLast(id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, err.Error())
				return
			}
			lasts = append(lasts, *reading)
		}
		if err := lasts.JSON(w); err != nil {
			log.Printf("Failed to create JSON representation of measurements: ", err.Error())
		}
	})
}

func MkLastSingleValuesHandler(hc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		last, err := hc.GetLast(byte(id))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := last.JSON(w); err != nil {
			log.Printf("Failed to create JSON representation of measurement %s", last.String())
		}
	})
}

func MkLastMinuteAvgSingleHandler(hc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		avg, err := hc.GetMinuteAvg(byte(id))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := avg.JSON(w); err != nil {
			log.Printf("Failed to create JSON representation of measurement %s", avg.String())
		}
	})
}

func MkLastMinuteAvgAllHandler(hc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		ids := hc.GetSortedIDs()
		avgs := ReadingSlice{}
		for _, id := range ids {
			reading, err := hc.GetMinuteAvg(id)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, err.Error())
				return
			}
			avgs = append(avgs, reading)
		}
		if err := avgs.JSON(w); err != nil {
			log.Printf("Failed to create JSON representation of measurements: ", err.Error())
		}

	})
}

func MkStatusHandler(s *Status) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := s.UpdateAndJSON(w); err != nil {
			log.Printf("Failed to create JSON representation of measurements: ", err.Error())
		}
	})
}

func Run_httpd(mc *MeasurementCache, s *Status, url string) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", MkIndexHandler(mc))
	router.HandleFunc("/last", MkLastAllValuesHandler(mc))
	router.HandleFunc("/last/{id:[0-9]+}", MkLastSingleValuesHandler(mc))
	router.HandleFunc("/minuteavg", MkLastMinuteAvgAllHandler(mc))
	router.HandleFunc("/minuteavg/{id:[0-9]+}", MkLastMinuteAvgSingleHandler(mc))
	router.HandleFunc("/status", MkStatusHandler(s))
	log.Fatal(http.ListenAndServe(url, router))
}
