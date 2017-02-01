package main

import (
	//"github.com/gonium/gosdm630"
	"fmt"
	"gopkg.in/urfave/cli.v1"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "sdm630_monitor"
	app.Usage = "SDM630 monitor"
	app.Version = "0.2.0"
	app.HideVersion = true
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "url, u",
			Value: "localhost:8080",
			Usage: "the URL of the server we should connect to",
		},
		cli.StringFlag{
			Name:  "category, c",
			Value: "all",
			Usage: "the firehose category to subscribe to",
		},
		cli.IntFlag{
			Name:  "timeout, t",
			Value: 45,
			Usage: "timeout value in seconds",
		},
		cli.BoolFlag{
			Name:  "verbose, v",
			Usage: "print verbose messages",
		},
	}
	app.Action = func(c *cli.Context) {
		endpointUrl :=
			fmt.Sprintf("http://%s/firehose?timeout=%d&category=%s",
				c.String("url"), c.Int("timeout"), c.String("category"))
		if c.Bool("verbose") {
			log.Printf("Client startup - will connect to %s", endpointUrl)
		}
		client := &http.Client{
			Timeout: time.Duration(c.Int("timeout")) * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 20,
				Dial: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 10 * time.Second,
			},
		}
		for {
			resp, err := client.Get(endpointUrl)
			if err != nil {
				log.Fatal("Failed to read from endpoint: ", err.Error())
			}
			events, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("Failed to process message: ", err.Error())
			} else {
				log.Printf(string(events))
			}
			if resp.Body != nil {
				resp.Body.Close()
			}
		}
	}
	app.Run(os.Args)
}
