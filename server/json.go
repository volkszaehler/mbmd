package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"sort"
)

// apiData combines readings with associated device id for JSON encoding
// using kvslice it ensured order export of the readings map
type apiData struct {
	device   string
	readings *Readings
}

// MarshalJSON creates device api json for export
func (d apiData) MarshalJSON() ([]byte, error) {
	res := kvslice{
		{"Device", d.device},
		{"Timestamp", d.readings.Timestamp},
		{"Unix", d.readings.Timestamp.Unix()},
	}

	if d.readings.Values == nil {
		return json.Marshal(res)
	}

	values := kvslice{}
	for m, v := range d.readings.Values {
		if math.IsNaN(v) {
			// safeguard for NaN values - should only happen in simluation mode
			log.Printf("skipping unexpected NaN value for %s", m)
			continue
		}
		values = append(values, kv{m.String(), v})
	}

	sort.Slice(values, func(a, b int) bool {
		return values[a].key < values[b].key
	})

	return json.Marshal(append(res, values...))
}

type kvslice []kv

func (s kvslice) MarshalJSON() ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString("{")

	for i, kv := range s {
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
		var val []byte
		switch kv.val.(type) {
		case float64:
			val = []byte(fmt.Sprintf("%.5g", kv.val))
		default:
			if val, err = json.Marshal(kv.val); err != nil {
				return nil, err
			}
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
