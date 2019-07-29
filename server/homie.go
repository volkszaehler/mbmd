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
	specVersion = "3.0.1"
	nodeTopic   = "meter"
	timeout     = 500 * time.Millisecond
)

// HomieRunner publishes query results as homie mqtt topics
type HomieRunner struct {
	*MqttClient
	qe        DeviceInfo
	rootTopic string
	meters    map[string]map[meters.Measurement]bool
}

// NewHomieRunner create new runner for homie IoT spec
func NewHomieRunner(mqttClient *MqttClient, qe DeviceInfo, rootTopic string) *HomieRunner {
	hr := &HomieRunner{
		MqttClient: mqttClient,
		qe:         qe,
		rootTopic:  rootTopic,
		meters:     make(map[string]map[meters.Measurement]bool),
	}

	return hr
}

// Run MQTT client publisher
func (hr *HomieRunner) Run(in <-chan QuerySnip) {
	defer hr.unregister() // cleanup topics

	for snip := range in {
		// publish meter
		if _, ok := hr.meters[snip.Device]; !ok {
			hr.meters[snip.Device] = make(map[meters.Measurement]bool)
			d := hr.qe.DeviceDescriptorByID(snip.Device)
			hr.publishMeter(snip.Device, d)
		}

		// publish properties
		if _, ok := hr.meters[snip.Device][snip.Measurement]; !ok {
			hr.meters[snip.Device][snip.Measurement] = true
			hr.publishProperties(snip.Device)
		}

		// publish data
		topic := fmt.Sprintf("%s/%s/%s/%s",
			hr.rootTopic,
			hr.deviceTopic(snip.Device),
			nodeTopic,
			strings.ToLower(snip.Measurement.String()),
		)

		message := fmt.Sprintf("%.3f", snip.Value)
		go hr.Publish(topic, false, message)
	}
}

// unregister unpublishes device information
func (hr *HomieRunner) unregister() {
	// devices
	for meter := range hr.meters {
		// clear retained messages
		subTopic := hr.deviceTopic(meter)
		hr.unpublish(subTopic)
	}
}

func (hr *HomieRunner) publishMeter(meter string, descriptor meters.DeviceDescriptor) {
	// clear retained messages
	subTopic := hr.deviceTopic(meter)
	hr.unpublish(subTopic)

	// device
	hr.publish(subTopic+"/$homie", specVersion)
	hr.publish(subTopic+"/$name", "MBMD")
	hr.publish(subTopic+"/$state", "ready")

	// node
	hr.publish(subTopic+"/$nodes", nodeTopic)

	subTopic = fmt.Sprintf("%s/%s", subTopic, nodeTopic)
	hr.publish(subTopic+"/$name", descriptor.Manufacturer)
	hr.publish(subTopic+"/$type", descriptor.Model)
}

func (hr *HomieRunner) publishProperties(meter string) {
	subtopic := hr.deviceTopic(meter)

	measurements := make([]meters.Measurement, 0)
	for k := range hr.meters[meter] {
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

	hr.publish(subtopic+"/$properties", strings.Join(properties[:], ","))
}

func (hr *HomieRunner) publish(subtopic string, message string) {
	topic := hr.rootTopic + "/" + subtopic
	hr.Publish(topic, true, message)
}

// unpublish retained message hierarchy
func (hr *HomieRunner) unpublish(subtopic string) {
	topic := fmt.Sprintf("%s/%s/#", hr.rootTopic, subtopic)
	if hr.verbose {
		log.Printf("MQTT: unpublish %s", topic)
	}

	var mux sync.Mutex
	tokens := make([]MQTT.Token, 0)

	mux.Lock()
	tokens = append(tokens, hr.client.Subscribe(topic, byte(hr.mqttQos), func(c MQTT.Client, msg MQTT.Message) {
		if len(msg.Payload()) == 0 {
			return // ignore received unpublish messages
		}

		topic := msg.Topic()
		token := hr.client.Publish(topic, byte(hr.mqttQos), true, []byte{})

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
	hr.client.Unsubscribe(topic)

	// wait for tokens
	for _, token := range tokens {
		hr.WaitForToken(token)
	}
}
