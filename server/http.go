package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"runtime"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const (
	devAssets = false
)

var (
	Version = "unknown version"
	Commit  = "unknown commit"
)

//go:generate go run github.com/mjibson/esc -private -o assets.go -pkg server -prefix ../assets ../assets

type Httpd struct {
}

func (h *Httpd) mkIndexHandler(mc *Cache) func(http.ResponseWriter, *http.Request) {
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

func (h *Httpd) mkLastAllValuesHandler(mc *Cache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		ids := mc.SortedIDs()
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

func (h *Httpd) mkLastSingleValuesHandler(mc *Cache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		last, err := mc.GetCurrent(id)
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

func (h *Httpd) mkLastMinuteAvgSingleHandler(mc *Cache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, ok := vars["id"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		avg, err := mc.GetMinuteAvg(id)
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

func (h *Httpd) mkLastMinuteAvgAllHandler(mc *Cache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		ids := mc.SortedIDs()
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

func (h *Httpd) mkStatusHandler(s *Status) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(s); err != nil {
			log.Printf("Failed to create JSON representation of measurements: %s", err.Error())
		}
	})
}

func (h *Httpd) mkSocketHandler(hub *SocketHub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ServeWebsocket(hub, w, r)
	}
}

// serveJSON decorates handler with required headers
func (h *Httpd) serveJSON(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		f(w, r)
	}
}

func (h *Httpd) Run(
	mc *Cache,
	hub *SocketHub,
	s *Status,
	url string,
) {
	log.Printf("Starting API at %s", url)

	router := mux.NewRouter().StrictSlash(true)

	// static
	router.HandleFunc("/", h.mkIndexHandler(mc))

	// individual handlers per folder
	for _, folder := range []string{"js", "css"} {
		prefix := fmt.Sprintf("/%s/", folder)
		router.PathPrefix(prefix).Handler(http.StripPrefix(prefix, http.FileServer(_escDir(devAssets, prefix))))
	}

	// api
	router.HandleFunc("/last", h.serveJSON(h.mkLastAllValuesHandler(mc)))
	router.HandleFunc("/last/{id:[a-z@0-9]+}", h.serveJSON(h.mkLastSingleValuesHandler(mc)))
	router.HandleFunc("/minuteavg", h.serveJSON(h.mkLastMinuteAvgAllHandler(mc)))
	router.HandleFunc("/minuteavg/{id:[a-z@0-9]+}", h.serveJSON(h.mkLastMinuteAvgSingleHandler(mc)))
	router.HandleFunc("/status", h.serveJSON(h.mkStatusHandler(s)))

	// websocket
	router.HandleFunc("/ws", h.mkSocketHandler(hub))

	srv := http.Server{
		Addr:         url,
		Handler:      handlers.CompressHandler(router),
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	srv.SetKeepAlivesEnabled(true)
	log.Fatal(srv.ListenAndServe())
}
