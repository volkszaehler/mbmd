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
	app.Name = "sdm630_monitor"
	app.Usage = "SDM630 monitor"
	app.Version = sdm630.RELEASEVERSION
	app.HideVersion = true
	app.Flags = []cli.Flag{
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
		cli.IntFlag{
			Name:  "timeout, t",
			Value: 45,
			Usage: "timeout value in seconds",
		},
		cli.IntFlag{
			Name:  "device, d",
			Usage: "specify the MODBUS id of the device to monitor",
		},
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "print verbose messages",
		},
	}
	app.Action = func(c *cli.Context) {
		if !c.IsSet("device") {
			log.Fatal("No device id given -- aborting. See --help for more information")
		}
		endpointUrl :=
			fmt.Sprintf("http://%s/firehose?timeout=%d&category=%s",
				c.String("url"), c.Int("timeout"), c.String("category"))
		if c.Bool("verbose") {
			log.Printf("Client startup - will connect to %s", endpointUrl)
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
				log.Fatal("Failed to read from endpoint: ", err.Error())
			}
			rawevents, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("Failed to process message: ", err.Error())
			} else {
				// handle the events.
				var events eventResponse
				err := json.Unmarshal(rawevents, &events)
				if err != nil {
					log.Fatal("Failed to decode JSON events: ", err.Error())
				}
				for _, event := range *events.Events {
					snip := event.Data
					if snip.DeviceId == uint8(c.Int("device")) {
						if snip.IEC61850 == "WLocPhsA" {
							log.Printf("Device %d: L1 %.2f W", snip.DeviceId, snip.Value)
						}
						if snip.IEC61850 == "WLocPhsB" {
							log.Printf("Device %d: L2 %.2f W", snip.DeviceId, snip.Value)
						}
						if snip.IEC61850 == "WLocPhsC" {
							log.Printf("Device %d: L3 %.2f W", snip.DeviceId, snip.Value)
						}
					}
				}

			}
			if resp.Body != nil {
				resp.Body.Close()
			}
		}
	}
	app.Run(os.Args)
}
