package server

import (
	"fmt"
	"log"
	"sort"
	"strings"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	. "github.com/volkszaehler/mbmd/meters"
)

const (
	specVersion = "3.0.1"
	nodeTopic   = "meter"
	timeout     = 500 * time.Millisecond
)

type HomieRunner struct {
	*MqttClient
	rootTopic string
	meters    map[uint8]*Meter
}

// Run MQTT client publisher
func (m *HomieRunner) Run(in QuerySnipChannel) {
	defer m.unregister() // cleanup topics

	for snip := range in {
		topic := fmt.Sprintf("%s/%s/%s/%s",
			m.rootTopic,
			m.DeviceTopic(snip.DeviceId),
			nodeTopic,
			strings.ToLower(snip.IEC61850.String()))

		message := fmt.Sprintf("%.3f", snip.Value)
		go m.Publish(topic, false, message)
	}
}

// Register subcribes GoSDM as discoverable device
func (m *HomieRunner) Register(rootTopic string, meters map[uint8]*Meter, qe *ModbusEngine) {
	// mqttOpts.SetWill(m.homieTopic("$state"), "lost", byte(m.mqttQos), true)
	m.rootTopic = stripSlash(rootTopic)
	m.meters = meters

	// devices
	for _, meter := range meters {
		m.publishMeter(meter, qe)
	}
}

// Register subcribes GoSDM as discoverable device
func (m *HomieRunner) unregister() {
	// devices
	for _, meter := range m.meters {
		// clear retained messages
		subTopic := m.DeviceTopic(meter.DeviceId)
		m.unpublish(subTopic)
	}
}

// stripSlash removes trailing slash
func stripSlash(s string) string {
	if s[len(s)-1:] == "/" {
		s = s[:len(s)-1]
	}
	return s
}

func (m *HomieRunner) publishMeter(meter *Meter, qe *ModbusEngine) {
	// descriptor := m.deviceDescriptor(meter, qe)

	// // clear retained messages
	// subTopic := m.DeviceTopic(meter.DeviceId)
	// m.unpublish(subTopic)

	// // device
	// m.publish(subTopic+"/$homie", specVersion)
	// m.publish(subTopic+"/$name", "GoSDM")
	// m.publish(subTopic+"/$state", "ready")
	// // m.publish(subTopic+"/$implementation", "GoSDM")

	// // node
	// m.publish(subTopic+"/$nodes", nodeTopic)

	// subTopic = fmt.Sprintf("%s/%s", subTopic, nodeTopic)
	// m.publish(subTopic+"/$name", descriptor.Manufacturer)
	// m.publish(subTopic+"/$type", descriptor.Model)

	// // properties
	// m.publishProperties(subTopic, meter, qe)
}

// func (m *HomieRunner) deviceDescriptor(meter *Meter, qe *ModbusEngine) sunspec.DeviceDescriptor {
// 	descriptor := sunspec.DeviceDescriptor{
// 		Manufacturer: meter.Producer.Type(),
// 		Model:        nodeTopic,
// 	}

// 	if sunspec, ok := meter.Producer.(SunSpecProducer); ok {
// 		op := sunspec.GetSunSpecCommonBlock()
// 		snip := QuerySnip{
// 			DeviceId:  meter.DeviceId,
// 			Operation: op,
// 		}
// 		if b, err := qe.Query(snip); err == nil {
// 			if descriptor, err = sunspec.DecodeSunSpecCommonBlock(b); err != nil {
// 				log.Println(err)
// 			}
// 		} else {
// 			log.Println(err)
// 		}
// 	}

// 	return descriptor
// }

func (m *HomieRunner) publishProperties(subtopic string, meter *Meter, qe *ModbusEngine) {
	meterOps := meter.Producer.Produce()

	// read from device to split block operations
	// TODO refactor transformation code into queryengine
	snips := make([]QuerySnip, 0)
	for _, op := range meterOps {
		snip := QuerySnip{
			DeviceId:  meter.DeviceId,
			Operation: op,
		}
		if b, err := qe.Query(snip); err == nil {
			snips = append(snips, qe.Transform(snip, b)...)
		}
	}

	// sort by measurement type
	sort.Slice(snips, func(a, b int) bool {
		return snips[a].IEC61850.String() < snips[b].IEC61850.String()
	})

	properties := make([]string, len(snips))

	for i, operation := range snips {
		property := strings.ToLower(operation.IEC61850.String())
		properties[i] = property

		description, unit := operation.IEC61850.DescriptionAndUnit()

		propertySubtopic := fmt.Sprintf("%s/%s", subtopic, property)
		m.publish(propertySubtopic+"/$name", description)

		m.publish(propertySubtopic+"/$unit", unit)
		m.publish(propertySubtopic+"/$datatype", "float")
	}
	m.publish(subtopic+"/$properties", strings.Join(properties[:], ","))
}

func (m *HomieRunner) publish(subtopic string, message string) {
	topic := m.rootTopic + "/" + subtopic
	m.Publish(topic, true, message)
}

// unpublish retained message hierarchy
func (m *HomieRunner) unpublish(subtopic string) {
	topic := fmt.Sprintf("%s/%s/#", m.rootTopic, subtopic)
	if m.verbose {
		log.Printf("MQTT: unpublish %s", topic)
	}

	var mux sync.Mutex
	tokens := make([]MQTT.Token, 0)

	mux.Lock()
	tokens = append(tokens, m.client.Subscribe(topic, byte(m.mqttQos), func(c MQTT.Client, msg MQTT.Message) {
		// we'll also receive the unpublish messages here so ignore these
		if len(msg.Payload()) == 0 {
			return // exit on unpublish message
		}

		topic := msg.Topic()
		token := m.client.Publish(topic, byte(m.mqttQos), true, []byte{})

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
	m.client.Unsubscribe(topic)

	// wait for tokens
	for _, token := range tokens {
		m.WaitForToken(token)
	}
}
