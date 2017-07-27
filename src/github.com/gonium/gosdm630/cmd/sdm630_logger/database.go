package main

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gonium/gosdm630"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

const (
	DEFAULT_BUCKET = "records"
)

type SnipDB struct {
	sync.Mutex
	databaseFile string
	snips        []sdm630.QuerySnip
}

func NewSnipDB(
	dbfile string,
) *SnipDB {
	return &SnipDB{
		databaseFile: dbfile,
		snips:        []sdm630.QuerySnip{},
	}
}

func (r *SnipDB) AddSnip(snip sdm630.QuerySnip) {
	r.Lock()
	defer r.Unlock()
	r.snips = append(r.snips, snip)
}

func (r *SnipDB) RunRecorder(sleepSeconds int) {
	// open database
	db, err := bolt.Open(r.databaseFile,
		0600,
		&bolt.Options{Timeout: 1 * time.Second},
	)
	if err != nil {
		log.Fatal("Cannot open database, exiting. Error was: ", err.Error())
	}
	defer db.Close()

	// now: sleep, then store recorded snips.
	for {
		time.Sleep(time.Duration(sleepSeconds) * time.Second)
		r.Lock()
		log.Printf("Cached %d measurements, storing.", len(r.snips))
		// taken from https://github.com/boltdb/bolt#autoincrementing-integer-for-the-bucket
		err := db.Update(func(tx *bolt.Tx) error {
			// Retrieve the users bucket.
			// This should be created when the DB is first opened.
			b, err := tx.CreateBucketIfNotExists([]byte(DEFAULT_BUCKET))
			if err != nil {
				return fmt.Errorf("Failed to create storage bucket, error was: %s",
					err.Error())
			}
			for _, snip := range r.snips {
				// Generate ID for the reading.
				// This returns an error only if the Tx is closed or not writeable.
				// That can't happen in an Update() call so I ignore the error check.
				id, _ := b.NextSequence()
				buf, err := json.Marshal(snip)
				if err != nil {
					return fmt.Errorf("Failed to marshal data %s, error was: %s",
						snip.String(), err.Error())
				}
				err = b.Put(itob(id), buf)
				if err != nil {
					return fmt.Errorf("Failed to add snip %s to bucket, error was: %s",
						snip.String(), err.Error())
				}
			}
			return nil
		})
		if err != nil {
			log.Fatalf("Failed to store records: %s", err.Error())
		} else {
			// clear the cache
			r.snips = []sdm630.QuerySnip{}
		}
		r.Unlock()
	}
}

// itob returns an 8-byte big endian representation of v.
func itob(v uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}

func (db *SnipDB) Inspect(w io.Writer) error {
	// open database
	blt, err := bolt.Open(db.databaseFile,
		0600,
		&bolt.Options{
			Timeout:  1 * time.Second,
			ReadOnly: true,
		},
	)
	if err != nil {
		return fmt.Errorf("Cannot open database, exiting. Error was: %s."+
			" Is the database file in use by another process?", err.Error())
	}
	defer blt.Close()

	firstSnipTime := time.Now()
	var lastSnipTime time.Time
	numSnips := 0

	err = blt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DEFAULT_BUCKET))
		if b == nil {
			return fmt.Errorf("No bucket found in database - empty?")
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			numSnips += 1
			var snip sdm630.QuerySnip
			err := json.Unmarshal(v, &snip)
			if err != nil {
				return fmt.Errorf("Failed to unmarshal snip %s, error was: %s", v,
					err.Error())
			}
			if snip.ReadTimestamp.Before(firstSnipTime) {
				firstSnipTime = snip.ReadTimestamp
			}
			if snip.ReadTimestamp.After(lastSnipTime) {
				lastSnipTime = snip.ReadTimestamp
			}
		}
		return nil
	})

	fmt.Fprintf(w, "Found %d records:\n", numSnips)
	fmt.Fprintf(w, "* First recorded on %s\n", firstSnipTime)
	fmt.Fprintf(w, "* Last recorded on %s\n", lastSnipTime)
	return err
}

func (db *SnipDB) ExportCSV(csvfile string) error {
	// open database
	blt, err := bolt.Open(db.databaseFile,
		0600,
		&bolt.Options{
			Timeout:  1 * time.Second,
			ReadOnly: true,
		},
	)
	if err != nil {
		return fmt.Errorf("Cannot open database, exiting. Error was: %s."+
			" Is the database file in use by another process?", err.Error())
	}
	defer blt.Close()

	f, err := os.Create(csvfile)
	if err != nil {
		return fmt.Errorf("Cannot write to csv file, exiting. Error was: "+
			"%s.", err.Error())
	}
	defer f.Close()
	w := bufio.NewWriter(f)

	fmt.Fprintf(w, "ID\tTime\tL1\tL2\tL3\n")
	numSnips := 0
	err = blt.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(DEFAULT_BUCKET))
		if b == nil {
			return fmt.Errorf("No bucket found in database - empty?")
		}
		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			numSnips += 1
			var snip sdm630.QuerySnip
			err := json.Unmarshal(v, &snip)
			if err != nil {
				return fmt.Errorf("Failed to unmarshal snip %s, error was: %s", v,
					err.Error())
			}

			// TODO: Export all snips
			switch snip.IEC61850 {
			case "WLocPhsA":
				fmt.Fprintf(w, "%d\t%s\t%.2f\t\t\n", snip.DeviceId,
					snip.ReadTimestamp, snip.Value)
			case "WLocPhsB":
				fmt.Fprintf(w, "%d\t%s\t\t%.2f\t\n", snip.DeviceId,
					snip.ReadTimestamp, snip.Value)
			case "WLocPhsC":
				fmt.Fprintf(w, "%d\t%s\t\t\t%.2f\n", snip.DeviceId,
					snip.ReadTimestamp, snip.Value)
			default:
				continue
			}

		}
		return nil
	})
	log.Printf("Exported %d records.", numSnips)
	w.Flush()
	return err
}
