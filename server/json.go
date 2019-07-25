package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
)

// data combines readings with associated device id for JSON encoding
// using kvslice it ensured order export of the readings map
type data struct {
	device   string
	readings Readings
}

type kvslice []kv

func (os kvslice) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("{")

	for i, kv := range os {
		if i != 0 {
			buf.WriteString(",")
		}
		// marshal key
		key, err := json.Marshal(kv.key)
		if err != nil {
			return nil, err
		}
		buf.Write(key)
		buf.WriteString(":")
		// marshal value
		val, err := json.Marshal(kv.val)
		if err != nil {
			return nil, err
		}
		buf.Write(val)
	}

	buf.WriteString("}")
	return buf.Bytes(), nil
}

type kv struct {
	key string
	val interface{}
}

func (o kv) MarshalJSON() ([]byte, error) {
	switch o.val.(type) {
	case float64:
		return []byte(fmt.Sprintf("%g", o.val)), nil
	default:
		return json.Marshal(o.val)
	}
}

// MarshalJSON creates device api json for export
func (d data) MarshalJSON() ([]byte, error) {
	res := kvslice{
		{"device", d.device},
		{"timestamp", d.readings.Timestamp},
		{"unix", d.readings.Timestamp.Unix()},
	}

	if d.readings.Values == nil {
		return json.Marshal(res)
	}

	for m, v := range d.readings.Values {
		if math.IsNaN(v) {
			// safeguard for NaN values - should only happen in simluation mode
			log.Printf("skipping unexpected NaN value for %s", m)
			continue
		}
		k := strings.ToLower(m.String())
		res = append(res, kv{k, v})
	}

	return json.Marshal(res)
}
