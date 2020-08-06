/*
 * Copyright 2019-20 Joaquim Rocha <jrocha@gmailbox.org> and Contributors
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

package report

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jjmrocha/beast/client"
)

const datetimeFormat = "2006-01-02 15:04:05.999"

// Output is the handler for CVS output
type Output struct {
	outputFile string
	file       *os.File
	writer     *csv.Writer
}

// NewOutput creates new Output handler
func NewOutput(fileName string) Output {
	if fileName == "" {
		return Output{
			outputFile: fileName,
		}
	}

	fileHandler, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("Error creating output file %s: %v\n", fileName, err)
	}

	csvWriter := csv.NewWriter(fileHandler)
	csvWriter.Write(csvHeader())

	return Output{
		outputFile: fileName,
		file:       fileHandler,
		writer:     csvWriter,
	}
}

// Write appends the response to the output file
func (o Output) Write(r *client.Response) {
	if o.outputFile == "" {
		return
	}

	record := encode(r)
	o.writer.Write(record)
}

// Close closes the CSV file
func (o Output) Close() {
	if o.outputFile == "" {
		return
	}

	defer o.file.Close()
	o.writer.Flush()

	fmt.Printf("===== Output File =====\n")
	fmt.Printf("Output file '%s' was successfully generated\n", o.outputFile)
}

func csvHeader() []string {
	return []string{
		"Timestamp",
		"Request",
		"Result",
		"StatusCode",
		"IsSuccess",
		"Duration",
	}
}

func encode(r *client.Response) []string {
	var timestamp = r.Timestamp.Format(datetimeFormat)
	var request = r.Request
	var result = "Executed"
	var statusCode = ""
	var isSuccess = "false"
	var duration = ""

	if r.Duration.Nanoseconds() > 0 {
		duration = strconv.FormatInt(r.Duration.Milliseconds(), 10)
	}

	if r.IsClientError() {
		result = r.ClientError()
	} else {
		statusCode = strconv.Itoa(r.StatusCode)

		if r.IsSuccess() {
			isSuccess = "true"
		}
	}

	return []string{
		timestamp,
		request,
		result,
		statusCode,
		isSuccess,
		duration,
	}
}
