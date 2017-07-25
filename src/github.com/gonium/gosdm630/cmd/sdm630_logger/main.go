package main

import (
	"encoding/json"
	"fmt"
	"github.com/gonium/gosdm630"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

const (
	ERROR_WAITTIME_MS = 100
)

// Copied from
// https://github.com/jcuga/golongpoll/blob/master/events.go:
type lpEvent struct {
	// Timestamp is milliseconds since epoch to match javascrits Date.getTime()
	Timestamp int64  `json:"timestamp"`
	Category  string `json:"category"`
	// NOTE: Data can be anything that is able to passed to json.Marshal()
	Data sdm630.QuerySnip `json:"data"`
}

// eventResponse is the json response that carries longpoll events.
type eventResponse struct {
	Events *[]lpEvent `json:"events"`
}

func main() {
	app := cli.NewApp()
	app.Name = "sdm630_logger"
	app.Usage = "SDM630 Logger"
	app.Version = sdm630.RELEASEVERSION
	app.HideVersion = true
	// Global flags
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "print verbose messages",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "record",
			Aliases: []string{"r"},
			Usage:   "Record all measurements",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "url, u",
					Value: "localhost:8080",
					Usage: "the URL of the server we should connect to",
				},
				cli.StringFlag{
					Name:  "category, c",
					Value: "meterupdate",
					Usage: "the firehose category to subscribe to",
				},
				cli.StringFlag{
					Name:  "dbfile, f",
					Value: "log.db",
					Usage: "the database file to record to",
				},
				cli.IntFlag{
					Name:  "timeout, t",
					Value: 45,
					Usage: "timeout value in seconds",
				},
				cli.IntFlag{
					Name:  "sleeptime, s",
					Value: 60 * 2,
					Usage: "seconds to sleep between writes to disk",
				},
			},
			Action: func(c *cli.Context) {
				db := NewSnipDB(c.String("dbfile"))
				go db.RunRecorder(c.Int("sleeptime"))
				endpointUrl :=
					fmt.Sprintf("http://%s/firehose?timeout=%d&category=%s",
						c.String("url"), c.Int("timeout"), c.String("category"))
				if c.GlobalBool("verbose") {
					log.Printf("recorder startup - will connect to %s", endpointUrl)
				}
				client := &http.Client{
					Timeout: time.Duration(c.Int("timeout")) * time.Second,
					Transport: &http.Transport{
						// 0 means: no limit.
						MaxIdleConns:        0,
						MaxIdleConnsPerHost: 0,
						IdleConnTimeout:     0,
						Dial: (&net.Dialer{
							Timeout:   30 * time.Second,
							KeepAlive: time.Minute,
						}).Dial,
						TLSHandshakeTimeout: 10 * time.Second,
						DisableKeepAlives:   false,
					},
				}
				for {
					resp, err := client.Get(endpointUrl)
					if err != nil {
						log.Println("Failed to read from endpoint: ", err.Error())
						// TODO: Exponential backoff
						time.Sleep(ERROR_WAITTIME_MS * time.Millisecond)
						continue
					}
					rawevents, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						log.Println("Failed to process message: ", err.Error())
						continue
					} else {
						// handle the events.
						var events eventResponse
						err := json.Unmarshal(rawevents, &events)
						if err != nil {
							log.Println("Failed to decode JSON events: ", err.Error())
							continue
						}
						for _, event := range *events.Events {
							snip := event.Data
							if c.GlobalBool("verbose") {
								log.Printf("%s: device %d, %s: %.2f", snip.ReadTimestamp,
									snip.DeviceId, snip.IEC61850, snip.Value)
							}
							db.AddSnip(snip)
						}

					}
					if resp.Body != nil {
						resp.Body.Close()
					}
				}
			},
		},
		{
			Name:    "export",
			Aliases: []string{"e"},
			Usage:   "export all measurements from a database",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dbfile, f",
					Value: "log.db",
					Usage: "the database file that contains all stored readings",
				},
				cli.StringFlag{
					Name:  "tsv, t",
					Value: "log.tsv",
					Usage: "the TSV file to export to",
				},
			},
			Action: func(c *cli.Context) {
				if c.GlobalBool("verbose") {
					log.Printf("exporter startup")
				}
				if c.GlobalBool("verbose") {
					log.Printf("Exporting database %s into TSV file %s",
						c.String("dbfile"), c.String("tsv"))
				}
				db := NewSnipDB(c.String("dbfile"))
				err := db.ExportCSV(c.String("tsv"))
				if err != nil {
					log.Fatalf("%s", err.Error())
				}

			},
		},
		{
			Name:    "inspect",
			Aliases: []string{"i"},
			Usage:   "inspect a recorded database",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "dbfile, f",
					Value: "log.db",
					Usage: "the database file to record to",
				},
			},
			Action: func(c *cli.Context) {
				if c.GlobalBool("verbose") {
					log.Printf("Inspecting database %s", c.String("dbfile"))
				}
				db := NewSnipDB(c.String("dbfile"))
				err := db.Inspect(os.Stdout)
				if err != nil {
					log.Fatalf("%s", err.Error())
				}
			},
		},
	}
	app.Run(os.Args)
}
