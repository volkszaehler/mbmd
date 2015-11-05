package sdm630

import (
	"fmt"
)

type TextDumper struct {
	datastream ReadingChannel
}

func NewTextDumper(ds ReadingChannel) *TextDumper {
	return &TextDumper{datastream: ds}
}

func (td *TextDumper) ConsumeData() {
	for {
		readings := <-td.datastream
		fmt.Printf("%s\r\n", &readings)
	}
}
