package server

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	mbmd_prometheus "github.com/volkszaehler/mbmd/prometheus"
)

// Assets is the embedded assets file system
var Assets fs.FS

// AssetsDir is the assets directory relative to the module root
const AssetsDir = "assets"

// Httpd is an http server
type Httpd struct {
	router *mux.Router
	mc     *Cache
	qe     DeviceInfo
}

func (h *Httpd) mkIndexHandler() func(http.ResponseWriter, *http.Request) {
	mainTemplate, err := fs.ReadFile(Assets, "index.html")
	if err != nil {
		log.Fatal("httpd: failed to load embedded template: " + err.Error())
	}
	t, err := template.New("mbmd").Parse(string(mainTemplate))
	if err != nil {
		log.Fatal("httpd: failed to create main page template: ", err.Error())
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
			log.Fatal("httpd: failed to render main page: ", err.Error())
		}
	})
}

func (h *Httpd) allDevicesHandler(
	readingsProvider func(id string) (*Readings, error),
) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ids := h.mc.SortedIDs()
		res := make(map[string]apiData)

		for _, id := range ids {
			readings, err := readingsProvider(id)
			if err != nil {
				// Skip this meter, it will simply not be displayed
				continue
			}

			data := apiData{readings: readings}
			res[id] = data
		}

		if len(res) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "all meters are inactive")
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Printf("httpd: failed to encode JSON: %s", err.Error())
		}
	})
}

func (h *Httpd) singleDeviceHandler(
	readingsProvider func(id string) (*Readings, error),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		id, ok := vars["id"]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		readings, err := readingsProvider(id)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, err.Error())
			return
		}

		data := apiData{readings: readings}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("httpd: failed to encode JSON %s", err.Error())
		}
	}
}

// mkSocketHandler attaches status handler to uri
func (h *Httpd) mkStatusHandler(s *Status) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(s); err != nil {
			log.Printf("httpd: failed to encode JSON: %s", err.Error())
		}
	}
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

// prometheusHandler is a middleware that prepares Prometheus metrics for collection by HTTP request
func (h *Httpd) prometheusHandler() http.Handler {
	return promhttp.HandlerFor(
		mbmd_prometheus.MBMDRegistry,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	)
}

// NewHttpd creates HTTP daemon
func NewHttpd(hub *SocketHub, s *Status, qe DeviceInfo, mc *Cache) *Httpd {
	srv := &Httpd{
		router: mux.NewRouter().StrictSlash(true),
		qe:     qe,
		mc:     mc,
	}

	// static
	static := srv.router.PathPrefix("/").Subrouter()
	static.Use(handlers.CompressHandler)

	// individual handlers per folder
	static.HandleFunc("/", srv.mkIndexHandler())
	for _, dir := range []string{"css", "js"} {
		static.PathPrefix("/" + dir).Handler(http.FileServer(http.FS(Assets)))
	}

	// Prometheus
	prom := srv.router.Path("/metrics")
	prom.Handler(srv.prometheusHandler())

	// api
	api := srv.router.PathPrefix("/api").Subrouter()
	api.Use(jsonHandler)
	api.Use(handlers.CompressHandler)

	api.HandleFunc("/last", srv.allDevicesHandler(srv.mc.Current))
	api.HandleFunc("/last/{id:[a-zA-Z0-9.]+}", srv.singleDeviceHandler(srv.mc.Current))
	api.HandleFunc("/avg", srv.allDevicesHandler(srv.mc.Average))
	api.HandleFunc("/avg/{id:[a-zA-Z0-9.]+}", srv.singleDeviceHandler(srv.mc.Average))
	api.HandleFunc("/status", srv.mkStatusHandler(s))

	// websocket
	srv.router.HandleFunc("/ws", srv.mkSocketHandler(hub))

	return srv
}

// Router returns the root router
func (h *Httpd) Router() *mux.Router {
	return h.router
}

// Run executes the http server
func (h *Httpd) Run(url string) {
	log.Printf("httpd: starting api at %s", url)

	// debug logger
	_ = log.New(debugLogger{"superfluous"}, "", 0)

	srv := http.Server{
		Addr:         url,
		Handler:      h.router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
		// ErrorLog: debug,
	}

	srv.SetKeepAlivesEnabled(true)
	log.Fatal(srv.ListenAndServe())
}
