package server

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	"github.com/volkszaehler/mbmd/meters"
)

const (
	specVersion = "4.0"
	nodeTopic   = "meter"
	timeout     = 500 * time.Millisecond
)

// HomieRunner publishes query results as homie mqtt topics
type HomieRunner struct {
	options   *MQTT.ClientOptions
	qos       byte
	verbose   bool
	rootTopic string
	qe        DeviceInfo
	cc        <-chan ControlSnip
	meters    map[string]*homieMeter
}

type homieMeter struct {
	*MqttClient
	rootTopic string
	meter     string
	online    bool
	observed  map[meters.Measurement]bool
}

// NewHomieRunner create new runner for homie IoT spec
func NewHomieRunner(qe DeviceInfo, cc <-chan ControlSnip, options *MQTT.ClientOptions, qos byte, rootTopic string, verbose bool) *HomieRunner {
	hr := &HomieRunner{
		options:   options,
		qos:       qos,
		verbose:   verbose,
		rootTopic: rootTopic,
		qe:        qe,
		cc:        cc,
		meters:    make(map[string]*homieMeter),
	}

	return hr
}

// cloneOptions creates a copy of the relevant mqtt options
func (hr *HomieRunner) cloneOptions() *MQTT.ClientOptions {
	opt := MQTT.NewClientOptions()

	opt.SetUsername(hr.options.Username)
	opt.SetPassword(hr.options.Password)
	opt.SetClientID(hr.options.ClientID)
	opt.SetCleanSession(hr.options.CleanSession)
	opt.SetAutoReconnect(hr.options.AutoReconnect)

	for _, b := range hr.options.Servers {
		opt.AddBroker(b.String())
	}

	return opt
}

// Run MQTT client publisher
func (hr *HomieRunner) Run(in <-chan QuerySnip) {
	defer hr.unregister() // cleanup topics

	for {
		select {
		case snip, chanOpen := <-in:
			if !chanOpen {
				return // channel closed
			}
			meter, ok := hr.meters[snip.Device]
			if !ok {
				meter = hr.createMeter(snip)
			}
			// publish actual message
			meter.publishMessage(snip)
		case snip := <-hr.cc:
			if meter, ok := hr.meters[snip.Device]; ok {
				meter.status(snip.Status.Online)
			}
		}
	}
}

// createMeter creates new homie device
func (hr *HomieRunner) createMeter(snip QuerySnip) *homieMeter {
	// new client with unique id
	options := hr.cloneOptions()
	clientID := fmt.Sprintf("%s-%s", options.ClientID, mqttDeviceTopic(snip.Device))
	options.SetClientID(clientID)

	lwt := fmt.Sprintf("%s/%s/$state", hr.rootTopic, mqttDeviceTopic(snip.Device))
	options.SetWill(lwt, "lost", hr.qos, true)

	client := NewMqttClient(options, hr.qos, hr.verbose)

	// add meter and publish
	meter := newHomieMeter(client, hr.rootTopic, snip.Device)
	hr.meters[snip.Device] = meter

	d := hr.qe.DeviceDescriptorByID(snip.Device)
	meter.publishMeter(d)

	return meter
}

// unregister unpublishes device information
func (hr *HomieRunner) unregister() {
	for _, meter := range hr.meters {
		meter.unregister()
	}
}

// newHomieMeter creates meter on given mqtt client
func newHomieMeter(client *MqttClient, rootTopic string, meter string) *homieMeter {
	hm := &homieMeter{
		MqttClient: client,
		rootTopic:  rootTopic,
		meter:      meter,
		observed:   make(map[meters.Measurement]bool),
	}
	return hm
}

// unregister clears meter's retained messages
func (hr *homieMeter) unregister() {
	subTopic := mqttDeviceTopic(hr.meter)
	hr.publish(subTopic+"/$state", "disconnected")
	hr.Client.Disconnect(uint(timeout.Milliseconds()))
}

func (hr *homieMeter) publishMeter(descriptor meters.DeviceDescriptor) {
	subTopic := mqttDeviceTopic(hr.meter)

	// device
	hr.publish(subTopic+"/$homie", specVersion)
	hr.publish(subTopic+"/$name", hr.meter)
	hr.publish(subTopic+"/$state", "init")
	hr.publish(subTopic+"/$implementation", "MBMD")

	// node
	hr.publish(subTopic+"/$nodes", nodeTopic)
	hr.unpublish(subTopic, nodeTopic, "$homie", "$name", "$state", "$nodes")

	subTopic = fmt.Sprintf("%s/%s", subTopic, nodeTopic)
	hr.publish(subTopic+"/$name", descriptor.Manufacturer)
	hr.publish(subTopic+"/$type", descriptor.Model)
}

// status updates a meter's $state attibute when its online status changes
func (hr *homieMeter) status(online bool) {
	if hr.online != online {
		subTopic := mqttDeviceTopic(hr.meter)
		msg := "alert"
		if online {
			msg = "ready"
		}
		hr.publish(subTopic+"/$state", msg)
		hr.online = online
	}
}

func (hr *homieMeter) publishMessage(snip QuerySnip) {
	// make sure property is published before publishing data
	if _, ok := hr.observed[snip.Measurement]; !ok {
		hr.observed[snip.Measurement] = true
		hr.publishProperties()
	}

	// publish data
	topic := fmt.Sprintf("%s/%s/%s/%s",
		hr.rootTopic,
		mqttDeviceTopic(snip.Device),
		nodeTopic,
		strings.ToLower(snip.Measurement.String()),
	)

	message := fmt.Sprintf("%.3f", snip.Value)
	go hr.Publish(topic, false, message)
}

func (hr *homieMeter) publishProperties() {
	subtopic := fmt.Sprintf("%s/%s", mqttDeviceTopic(hr.meter), nodeTopic)

	measurements := make([]meters.Measurement, 0)
	for k := range hr.observed {
		measurements = append(measurements, k)
	}

	// sort by measurement type
	sort.Slice(measurements, func(a, b int) bool {
		return measurements[a].String() < measurements[b].String()
	})

	properties := make([]string, len(measurements))

	for i, m := range measurements {
		property := strings.ToLower(m.String())
		properties[i] = property

		description, unit := m.DescriptionAndUnit()

		propertySubtopic := fmt.Sprintf("%s/%s", subtopic, property)
		hr.publish(propertySubtopic+"/$name", description)
		hr.publish(propertySubtopic+"/$unit", unit)
		hr.publish(propertySubtopic+"/$datatype", "float")
	}

	hr.publish(subtopic+"/$properties", strings.Join(properties, ","))

	// unpublish remains attributes if any
	exceptions := []string{"$name", "$unit", "$datatype", "$properties"}
	exceptions = append(exceptions, properties...)
	hr.unpublish(subtopic, exceptions...)
}

func (hr *homieMeter) publish(subtopic string, message string) {
	topic := hr.rootTopic + "/" + subtopic
	hr.Publish(topic, true, message)
}

// unpublish retained message hierarchy
func (hr *homieMeter) unpublish(subtopic string, exceptions ...string) {
	topic := fmt.Sprintf("%s/%s/#", hr.rootTopic, subtopic)
	if hr.verbose {
		log.Printf("mqtt: unpublish %s", topic)
	}

	var mux sync.Mutex
	tokens := make([]MQTT.Token, 0)

	mux.Lock()
	tokens = append(tokens, hr.Client.Subscribe(topic, hr.qos, func(c MQTT.Client, msg MQTT.Message) {
		if len(msg.Payload()) == 0 {
			return // ignore received unpublish messages
		}

		topic := msg.Topic()

		// don't unpublish if in exception list
		for _, exception := range exceptions {
			exceptionTopic := fmt.Sprintf("%s/%s/%s", hr.rootTopic, subtopic, exception)
			if topic == exceptionTopic || strings.HasPrefix(topic, exceptionTopic+"/") {
				// log.Printf("unpublish %s -> retain (%s)", topic, exception)
				return
			}
		}

		// log.Printf("unpublish %s", topic)
		token := hr.Client.Publish(topic, hr.qos, true, []byte{})

		mux.Lock()
		defer mux.Unlock()
		tokens = append(tokens, token)
	}))
	mux.Unlock()

	// wait for timeout according to specification
	<-time.After(timeout)
	mux.Lock()
	defer mux.Unlock()

	// stop listening
	hr.Client.Unsubscribe(topic)

	// wait for tokens
	for _, token := range tokens {
		hr.WaitForToken(token)
	}
}
