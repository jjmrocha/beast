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

package data

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

type Record map[string]string

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
		record := make(map[string]string)

		for i := 0; i < fieldsNumber; i++ {
			record[d.fields[i]] = columns[i]
		}

		d.records = append(d.records, record)
	}
}

func (d *Data) Next() *Record {
	d.current++

	if d.current == len(d.records) {
		d.current = 0
	}

	return &d.records[d.current]
}

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
