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
	"encoding/csv"
	"errors"
	"io"
	"os"
)

// Data contains the content of the CSV file
type Data struct {
	fields  []string
	records []Record
	current int
}

// Read reads the content of the CSV file
func Read(fileName string) (*Data, error) {
	csvfile, err := os.Open(fileName) //nolint:gosec // reads the data file the user named
	if err != nil {
		return nil, err
	}
	defer func() { _ = csvfile.Close() }()

	reader := csv.NewReader(csvfile)
	var dt Data
	for {
		columns, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, err
		}

		dt.add(columns)
	}

	if len(dt.records) == 0 {
		return nil, ErrNoRows
	}

	return &dt, nil
}

func (d *Data) add(columns []string) {
	if d.fields == nil {
		d.fields = columns
		return
	}

	record := NewRecord()
	for i, field := range d.fields {
		record[field] = columns[i]
	}
	d.records = append(d.records, record)
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
