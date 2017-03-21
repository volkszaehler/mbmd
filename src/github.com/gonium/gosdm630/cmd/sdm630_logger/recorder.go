package main

import (
	"github.com/boltdb/bolt"
	"github.com/gonium/gosdm630"
	"log"
	"sync"
	"time"
)

type Recorder struct {
	sync.Mutex
	databaseFile string
	sleepSeconds int
	snips        []sdm630.QuerySnip
}

func NewRecorder(
	dbfile string,
	sleepseconds int,
) *Recorder {
	return &Recorder{
		databaseFile: dbfile,
		sleepSeconds: sleepseconds,
		snips:        []sdm630.QuerySnip{},
	}
}

func (r *Recorder) AddSnip(snip sdm630.QuerySnip) {
	r.Lock()
	defer r.Unlock()
	r.snips = append(r.snips, snip)
	log.Printf("Number of cached snips: %d", len(r.snips))
}

func (r *Recorder) Run() {
	// open database
	db, err := bolt.Open(r.databaseFile,
		0600,
		&bolt.Options{Timeout: 1 * time.Second},
	)
	if err != nil {
		log.Fatal("Cannot open database, exiting. Error was: ", err.Error())
	}
	defer db.Close()
	for {
		time.Sleep(time.Duration(r.sleepSeconds) * time.Second)
		log.Printf("Storing measurements.")
		r.Lock()
		// serialize
		r.snips = []sdm630.QuerySnip{}
		r.Unlock()
	}
}
