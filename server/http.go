package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const devAssets = false

//go:generate go run github.com/mjibson/esc -private -o assets.go -pkg server -prefix ../assets ../assets

// Httpd is an http server
type Httpd struct {
	qe DeviceInfo
}

func (h *Httpd) mkIndexHandler(mc *Cache) func(http.ResponseWriter, *http.Request) {
	mainTemplate, err := _escFSString(devAssets, "/index.html")
	if err != nil {
		log.Fatal("failed to load embedded template: " + err.Error())
	}
	t, err := template.New("mbmd").Parse(string(mainTemplate))
	if err != nil {
		log.Fatal("failed to create main page template: ", err.Error())
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
			log.Fatal("failed to render main page: ", err.Error())
		}
	})
}

func (h *Httpd) allDevicesHandler(
	mc *Cache, readingsProvider func(id string) (*Readings, error),
) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids := mc.SortedIDs()
		res := make([]apiData, 0)
		for _, id := range ids {
			readings, err := readingsProvider(id)
			if err != nil {
				// Skip this meter, it will simply not be displayed
				continue
			}

			data := apiData{device: id, readings: readings}
			res = append(res, data)
		}

		if len(res) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "all meters are inactive")
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("failed to encode JSON: %s", err.Error())
		}
	})
}

func (h *Httpd) singleDeviceHandler(
	mc *Cache, readingsProvider func(id string) (*Readings, error),
) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, ok := vars["id"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		readings, err := readingsProvider(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, err.Error())
			return
		}

		data := apiData{device: id, readings: readings}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("failed to encode JSON %s", err.Error())
		}
	})
}

// mkSocketHandler attaches status handler to uri
func (h *Httpd) mkStatusHandler(s *Status) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(s); err != nil {
			log.Printf("failed to encode JSON: %s", err.Error())
		}
	})
}

// mkSocketHandler attaches websocket handler to uri
func (h *Httpd) mkSocketHandler(hub *SocketHub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ServeWebsocket(hub, w, r)
	}
}

type debugLogger struct {
	pattern string
}

func (d debugLogger) Write(p []byte) (n int, err error) {
	s := string(p)
	if strings.Contains(s, d.pattern) {
		debug.PrintStack()
	}
	return os.Stderr.Write(p)
}

// jsonHandler is a middleware that decorates responses with JSON and CORS headers
func jsonHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	})
}

// NewHttpd creates HTTP daemon
func NewHttpd(qe DeviceInfo) *Httpd {
	return &Httpd{qe: qe}
}

// Run executes the http server
func (h *Httpd) Run(
	mc *Cache,
	hub *SocketHub,
	s *Status,
	url string,
) {
	log.Printf("starting API at %s", url)
	router := mux.NewRouter().StrictSlash(true)

	// static
	router.HandleFunc("/", h.mkIndexHandler(mc))

	// individual handlers per folder
	for _, folder := range []string{"js", "css"} {
		prefix := fmt.Sprintf("/%s/", folder)
		router.PathPrefix(prefix).Handler(http.StripPrefix(prefix, http.FileServer(_escDir(devAssets, prefix))))
	}

	// api
	api := router.PathPrefix("/api").Subrouter()
	api.HandleFunc("/last", h.allDevicesHandler(mc, mc.Current))
	api.HandleFunc("/last/{id:[a-zA-Z0-9.]+}", h.singleDeviceHandler(mc, mc.Current))
	api.HandleFunc("/avg", h.allDevicesHandler(mc, mc.Average))
	api.HandleFunc("/avg/{id:[a-zA-Z0-9.]+}", h.singleDeviceHandler(mc, mc.Average))
	api.HandleFunc("/status", h.mkStatusHandler(s))
	api.Use(jsonHandler)

	// websocket
	router.HandleFunc("/ws", h.mkSocketHandler(hub))

	// debug logger
	_ = log.New(debugLogger{"superfluous"}, "", 0)

	srv := http.Server{
		Addr:         url,
		Handler:      handlers.CompressHandler(jsonHandler(router)),
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		// ErrorLog: debug,
	}

	srv.SetKeepAlivesEnabled(true)
	log.Fatal(srv.ListenAndServe())
}
