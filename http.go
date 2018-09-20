package sdm630

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
	"github.com/jcuga/golongpoll"
)

const (
	SECONDS_BETWEEN_STATUSUPDATE = 1
)

// Generate the embedded assets using https://github.com/aprice/embed
//go:generate go run github.com/aprice/embed/cmd/embed -c "embed.json"

func MkIndexHandler(mc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	loader := GetEmbeddedContent()
	mainTemplate, err := loader.GetContents("/index.html")
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

func MkLastAllValuesHandler(mc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		ids := mc.GetSortedIDs()
		lasts := ReadingSlice{}
		for _, id := range ids {
			reading, err := mc.GetLast(id)
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

func MkLastSingleValuesHandler(mc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		last, err := mc.GetLast(byte(id))
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

func MkLastMinuteAvgSingleHandler(mc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
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

func MkLastMinuteAvgAllHandler(mc *MeasurementCache) func(http.ResponseWriter, *http.Request) {
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

func MkStatusHandler(s *Status) func(http.ResponseWriter, *http.Request) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		s.Update()
		if err := json.NewEncoder(w).Encode(s); err != nil {
			log.Printf("Failed to create JSON representation of measurements: %s", err.Error())
		}
	})
}

func MkSocketHandler(hub *SocketHub) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ServeWebsocket(hub, w, r)
	}
}

type Firehose struct {
	lpManager  *golongpoll.LongpollManager
	in         QuerySnipChannel
	statstream chan string
}

func NewFirehose(inChannel QuerySnipChannel, status *Status, verbose bool) *Firehose {
	options := golongpoll.Options{}
	// see https://github.com/jcuga/golongpoll/blob/master/longpoll.go#L81
	//options := golongpoll.Options{
	//	LoggingEnabled:                 false,
	//	MaxLongpollTimeoutSeconds:      60,
	//	MaxEventBufferSize:             250,
	//	EventTimeToLiveSeconds:         60,
	//	DeleteEventAfterFirstRetrieval: false,
	//}
	if verbose {
		options.LoggingEnabled = true
	}
	manager, err := golongpoll.StartLongpoll(options)
	if err != nil {
		log.Fatalf("Failed to create firehose longpoll manager: %q", err)
	}
	// Attach a goroutine that will push meter status information
	// periodically
	var statusstream = make(chan string)
	go func() {
		for {
			time.Sleep(SECONDS_BETWEEN_STATUSUPDATE * time.Second)
			status.Update()
			if bytes, err := json.Marshal(status); err == nil {
				statusstream <- string(bytes)
			}
		}
	}()
	return &Firehose{
		lpManager:  manager,
		in:         inChannel,
		statstream: statusstream,
	}
}

func (f *Firehose) Run() {
	for {
		select {
		case snip := <-f.in:
			f.lpManager.Publish("meterupdate", snip)
		case statupdate := <-f.statstream:
			f.lpManager.Publish("statusupdate", statupdate)
		}
	}
}

func (f *Firehose) GetHandler() func(w http.ResponseWriter, r *http.Request) {
	return f.lpManager.SubscriptionHandler
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
	firehose *Firehose,
	hub *SocketHub,
	s *Status,
	url string,
) {
	router := mux.NewRouter().StrictSlash(true)

	// static
	router.HandleFunc("/", MkIndexHandler(mc))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/",
		GetEmbeddedContent()))

	// api
	router.HandleFunc("/last", serveJson(MkLastAllValuesHandler(mc)))
	router.HandleFunc("/last/{id:[0-9]+}", serveJson(MkLastSingleValuesHandler(mc)))
	router.HandleFunc("/minuteavg", serveJson(MkLastMinuteAvgAllHandler(mc)))
	router.HandleFunc("/minuteavg/{id:[0-9]+}", serveJson(MkLastMinuteAvgSingleHandler(mc)))
	router.HandleFunc("/status", serveJson(MkStatusHandler(s)))

	// longpoll
	if firehose != nil {
		router.HandleFunc("/firehose", firehose.GetHandler())
	}

	// websocket
	router.HandleFunc("/ws", MkSocketHandler(hub))

	srv := http.Server{
		Addr:         url,
		Handler:      handlers.CompressHandler(router),
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}

	srv.SetKeepAlivesEnabled(true)
	log.Fatal(srv.ListenAndServe())
}
