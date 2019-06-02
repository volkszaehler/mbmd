package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	SECONDS_BETWEEN_STATUSUPDATE = 1
	devAssets                    = false
)

var (
	Version = "unknown version"
	Commit  = "unknown commit"
)

//go:generate go run github.com/mjibson/esc -private -o assets.go -pkg server -prefix ../assets ../assets

func mkIndexHandler(mc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	mainTemplate, err := _escFSString(devAssets, "/index.html")
	if err != nil {
		log.Fatal("Failed to load embedded template: " + err.Error())
	}
	t, err := template.New("mbmd").Parse(string(mainTemplate))
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
			SoftwareVersion: Version,
			GolangVersion:   runtime.Version(),
		}
		err := t.Execute(w, data)
		if err != nil {
			log.Fatal("Failed to render main page: ", err.Error())
		}
	})
}

func mkLastAllValuesHandler(mc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		ids := mc.GetSortedIDs()
		lasts := ReadingSlice{}
		for _, id := range ids {
			reading, err := mc.GetCurrent(id)
			if err != nil {
				// Skip this meter, it will simply not be displayed
				continue
				//w.WriteHeader(http.StatusBadRequest)
				//fmt.Fprintf(w, err.Error())
				//return
			}
			lasts = append(lasts, *reading)
		}
		if len(lasts) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "All meters are inactive.")
			return
		}
		if err := json.NewEncoder(w).Encode(lasts); err != nil {
			log.Printf("Failed to create JSON representation of measurements: %s", err.Error())
		}
	})
}

func mkLastSingleValuesHandler(mc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		last, err := mc.GetCurrent(byte(id))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(last); err != nil {
			log.Printf("Failed to create JSON representation of measurement %s", last.String())
		}
	})
}

func mkLastMinuteAvgSingleHandler(mc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		avg, err := mc.GetMinuteAvg(byte(id))
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, err.Error())
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(avg); err != nil {
			log.Printf("Failed to create JSON representation of measurement %s", avg.String())
		}
	})
}

func mkLastMinuteAvgAllHandler(mc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		ids := mc.GetSortedIDs()
		avgs := ReadingSlice{}
		for _, id := range ids {
			reading, err := mc.GetMinuteAvg(id)
			if err != nil {
				// Skip this meter, it will simply not be displayed
				continue
			}
			avgs = append(avgs, *reading)
		}
		if len(avgs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "All meters are inactive.")
			return
		}
		if err := json.NewEncoder(w).Encode(avgs); err != nil {
			log.Printf("Failed to create JSON representation of measurements: %s", err.Error())
		}
	})
}

func mkStatusHandler(s *Status) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		s.Update()
		if err := json.NewEncoder(w).Encode(s); err != nil {
			log.Printf("Failed to create JSON representation of measurements: %s", err.Error())
		}
	})
}

func mkSocketHandler(hub *SocketHub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ServeWebsocket(hub, w, r)
	}
}

// serveJson decorates handler with required headers
func serveJson(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		f(w, r)
	}
}

func Run_httpd(
	mc *MeasurementCache,
	hub *SocketHub,
	s *Status,
	url string,
) {
	log.Printf("Starting API at %s", url)

	router := mux.NewRouter().StrictSlash(true)

	// static
	router.HandleFunc("/", mkIndexHandler(mc))

	// individual handlers per folder
	for _, folder := range []string{"js", "css"} {
		prefix := fmt.Sprintf("/%s/", folder)
		router.PathPrefix(prefix).Handler(http.StripPrefix(prefix, http.FileServer(_escDir(devAssets, prefix))))
	}

	// api
	router.HandleFunc("/last", serveJson(mkLastAllValuesHandler(mc)))
	router.HandleFunc("/last/{id:[0-9]+}", serveJson(mkLastSingleValuesHandler(mc)))
	router.HandleFunc("/minuteavg", serveJson(mkLastMinuteAvgAllHandler(mc)))
	router.HandleFunc("/minuteavg/{id:[0-9]+}", serveJson(mkLastMinuteAvgSingleHandler(mc)))
	router.HandleFunc("/status", serveJson(mkStatusHandler(s)))

	// websocket
	router.HandleFunc("/ws", mkSocketHandler(hub))

	srv := http.Server{
		Addr:         url,
		Handler:      handlers.CompressHandler(router),
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	srv.SetKeepAlivesEnabled(true)
	log.Fatal(srv.ListenAndServe())
}
