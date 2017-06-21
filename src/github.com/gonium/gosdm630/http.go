package sdm630

import (
	"bytes"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jcuga/golongpoll"
	"github.com/olekukonko/tablewriter"
	"html/template"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

// Generate the embedded assets using https://github.com/aprice/embed
//go:generate embed -c "embed.json"

// formatFloat helper
func fF(val float64) string {
	var buffer bytes.Buffer
	fmt.Fprintf(&buffer, "%.3f", val)
	return buffer.String()
}

func MkIndexHandler(hc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	loader := GetEmbeddedContent()
	mainTemplate, err := loader.GetContents("/index.tmpl")
	if err != nil {
		log.Fatal("Failed to load embedded template: " + err.Error())
	}
	t, err := template.New("gosdm630").Parse(string(mainTemplate))
	if err != nil {
		log.Fatal("Failed to create main page template: ", err.Error())
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		data := struct {
			SoftwareVersion string
			GolangVersion   string
		}{
			SoftwareVersion: RELEASEVERSION,
			GolangVersion:   runtime.Version(),
		}
		err := t.Execute(w, data)
		if err != nil {
			log.Fatal("Failed to render main page: ", err.Error())
		}
	})
}

func MkSimpleIndexHandler(hc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	loader := GetEmbeddedContent()
	mainTemplate, err := loader.GetContents("/simple.tmpl")
	if err != nil {
		log.Fatal("Failed to load embedded template: " + err.Error())
	}
	t, err := template.New("gosdm630").Parse(string(mainTemplate))
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
			table.SetHeader([]string{"Phase", "Voltage [V]", "Current [A]",
				"Power [W]", "Power Factor", "Import [kWh]", "Export [kWh]",
				"THD Voltage (Neutral) [%]"})
			table.Append([]string{"L1", fF(v.Voltage.L1), fF(v.Current.L1),
				fF(v.Power.L1), fF(v.Cosphi.L1), fF(v.Import.L1),
				fF(v.Export.L1), fF(v.THD.VoltageNeutral.L1)})
			table.Append([]string{"L2", fF(v.Voltage.L2), fF(v.Current.L2),
				fF(v.Power.L2), fF(v.Cosphi.L2), fF(v.Import.L2),
				fF(v.Export.L2), fF(v.THD.VoltageNeutral.L2)})
			table.Append([]string{"L3", fF(v.Voltage.L3), fF(v.Current.L3),
				fF(v.Power.L3), fF(v.Cosphi.L3), fF(v.Import.L3),
				fF(v.Export.L3), fF(v.THD.VoltageNeutral.L3)})
			table.Append([]string{"ALL", "n/a", fF(v.Current.L1 + v.Current.L2 + v.Current.L3),
				fF(v.Power.L1 + v.Power.L2 + v.Power.L3), "n/a", fF(v.TotalImport),
				fF(v.TotalExport), fF(v.THD.AvgVoltageNeutral)})
			table.Render()
		}
		data := struct {
			Content         string
			ReloadInterval  int
			SoftwareVersion string
		}{
			Content:         buffer.String(),
			ReloadInterval:  2,
			SoftwareVersion: RELEASEVERSION,
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
			log.Printf("Failed to create JSON representation of measurements: %s", err.Error())
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
			log.Printf("Failed to create JSON representation of measurements: %s", err.Error())
		}

	})
}

func MkStatusHandler(s *Status) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := s.UpdateAndJSON(w); err != nil {
			log.Printf("Failed to create JSON representation of measurements: %s", err.Error())
		}
	})
}

type Firehose struct {
	lpManager *golongpoll.LongpollManager
	in        QuerySnipChannel
}

func NewFirehose(inChannel QuerySnipChannel, verbose bool) *Firehose {
	// see https://github.com/jcuga/golongpoll/blob/master/longpoll.go#L81
	options := golongpoll.Options{
		LoggingEnabled:                 false,
		MaxLongpollTimeoutSeconds:      60,
		MaxEventBufferSize:             250,
		EventTimeToLiveSeconds:         60,
		DeleteEventAfterFirstRetrieval: true,
	}
	if verbose {
		options.LoggingEnabled = true
	}
	manager, err := golongpoll.StartLongpoll(options)
	if err != nil {
		log.Fatalf("Failed to create firehose longpoll manager: %q", err)
	}
	return &Firehose{
		lpManager: manager,
		in:        inChannel,
	}
}

func (f *Firehose) Run() {
	for {
		snip := <-f.in
		// TODO: Convert snip into external representation
		f.lpManager.Publish("all", snip)
	}
}

func (f *Firehose) GetHandler() func(w http.ResponseWriter, r *http.Request) {
	return f.lpManager.SubscriptionHandler
}

func Run_httpd(
	mc *MeasurementCache,
	firehose *Firehose,
	s *Status,
	url string,
) {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", MkIndexHandler(mc))
	router.HandleFunc("/simple", MkSimpleIndexHandler(mc))
	router.HandleFunc("/last", MkLastAllValuesHandler(mc))
	router.HandleFunc("/last/{id:[0-9]+}", MkLastSingleValuesHandler(mc))
	router.HandleFunc("/minuteavg", MkLastMinuteAvgAllHandler(mc))
	router.HandleFunc("/minuteavg/{id:[0-9]+}", MkLastMinuteAvgSingleHandler(mc))
	router.HandleFunc("/status", MkStatusHandler(s))
	router.HandleFunc("/firehose", firehose.GetHandler())
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		GetEmbeddedContent()))
	srv := http.Server{
		Addr:         url,
		Handler:      handlers.CompressHandler(router),
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	srv.SetKeepAlivesEnabled(true)
	log.Fatal(srv.ListenAndServe())
}
