/*
 * Copyright 2019 Joaquim Rocha <jrocha@gmailbox.org> and Contributors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

// Package data provides functions to read and manipulate CSV files used on the generation of requests
package data

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
)

// Record represents a line on the CSV file
type Record map[string]string

func (r *Record) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("{")
	for key, value := range *r {
		if buffer.Len() > 1 {
			buffer.WriteString(", ")
		}

		keyValue := fmt.Sprintf("%v: %v", key, value)
		buffer.WriteString(keyValue)
	}
	buffer.WriteString("}")
	return buffer.String()
}

// Data contains the content of the CSV file
type Data struct {
	fields  []string
	records []Record
	current int
}

func (d *Data) add(columns []string) {
	if d.fields == nil {
		d.fields = columns
	} else {
		fieldsNumber := len(d.fields)
		record := NewRecord()

		for i := 0; i < fieldsNumber; i++ {
			record[d.fields[i]] = columns[i]
		}

		d.records = append(d.records, record)
	}
}

// NewRecord creates a new record
func NewRecord() Record {
	return make(map[string]string)
}

// Next loops through the records
func (d *Data) Next() *Record {
	if d.current == len(d.records) {
		d.current = 0
	}

	next := d.records[d.current]
	d.current++
	return &next
}

// Read reads the content of the CSV file
func Read(fileName string) *Data {
	csvfile, err := os.Open(fileName)
	if err != nil {
		log.Fatalf("Error reading data file %s: %v\n", fileName, err)
	}

	defer csvfile.Close()
	reader := csv.NewReader(csvfile)
	data := Data{
		records: make([]Record, 0),
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatalf("Error reading data file %s: %v\n", fileName, err)
		}

		data.add(record)
	}

	return &data
}
