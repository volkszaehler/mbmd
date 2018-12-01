package sdm630

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strings"

	MQTT "github.com/eclipse/paho.mqtt.golang"

	. "github.com/gonium/gosdm630/internal/meters"
)

const (
	version   = "3.0.1"
	nodeTopic = "meter"
)

type HomieRunner struct {
	*MqttClient
	rootTopic string
}

// Run MQTT client publisher
func (m *HomieRunner) Run(in QuerySnipChannel, rate int) {
	rateMap := make(RateMap)

	for {
		snip := <-in
		topic := fmt.Sprintf("%s/%s/%s/%s",
			m.rootTopic,
			m.DeviceTopic(snip.DeviceId),
			nodeTopic,
			snip.IEC61850)

		if rateMap.Allowed(rate, topic) {
			message := fmt.Sprintf("%.3f", snip.Value)
			go m.Publish(topic, false, message)
		} else {
			if m.verbose {
				log.Printf("MQTT: skipped %s, rate to high", topic)
			}
		}
	}
}

// Register subcribes GoSDM as discoverable device
func (m *HomieRunner) Register(rootTopic string, meters map[uint8]*Meter) {
	// mqttOpts.SetWill(m.homieTopic("$state"), "lost", byte(m.mqttQos), true)
	m.rootTopic = stripSlash(rootTopic)

	// devices
	for _, meter := range meters {
		m.publishMeter(meter)
	}
}

// stripSlash removes trailing slash
func stripSlash(s string) string {
	if s[len(s)-1:] == "/" {
		s = s[:len(s)-1]
	}
	return s
}

func (m *HomieRunner) publishMeter(meter *Meter) {
	// clear retained messages
	subTopic := m.DeviceTopic(meter.DeviceId)
	m.unpublish(subTopic)

	// device
	m.publish(subTopic+"/$homie", version)
	m.publish(subTopic+"/$name", "GoSDM")
	m.publish(subTopic+"/$state", "ready")
	// m.publish(subTopic+"/$implementation", "GoSDM")

	// node
	m.publish(subTopic+"/$nodes", nodeTopic)

	subTopic = fmt.Sprintf("%s/%s", subTopic, nodeTopic)
	m.publish(subTopic+"/$name", meter.Producer.GetMeterType())
	m.publish(subTopic+"/$type", "meter")

	// properties
	m.publishProperties(subTopic, meter)
}

func (m *HomieRunner) publishProperties(subtopic string, meter *Meter) {
	meterOps := meter.Producer.Produce()
	sort.Slice(meterOps, func(a, b int) bool {
		return meterOps[a].IEC61850.String() < meterOps[b].IEC61850.String()
	})

	properties := make([]string, len(meterOps))
	re, _ := regexp.Compile(`^(.+) \((.+)\)$`)

	for i, operation := range meterOps {
		property := strings.ToLower(operation.IEC61850.String())
		properties[i] = property

		description := operation.IEC61850.Description()
		matches := re.FindStringSubmatch(description)
		if len(matches) == 3 {
			// strip unit from name
			description = matches[1]
		}

		propertySubtopic := fmt.Sprintf("%s/%s", subtopic, property)
		m.publish(propertySubtopic+"/$name", description)

		if len(matches) == 3 {
			unit := matches[2]
			m.publish(propertySubtopic+"/$unit", unit)
			m.publish(propertySubtopic+"/$datatype", "float")
		}
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
		log.Printf("MQTT: cleaning %s", topic)
	}

	m.client.Subscribe(topic, byte(m.mqttQos), func(c MQTT.Client, msg MQTT.Message) {
		topic := msg.Topic()
		m.Publish(topic, true, []byte{})
	})
}
