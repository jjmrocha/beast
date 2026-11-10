/*
 * Copyright 2020 Joaquim Rocha <jrocha@gmailbox.org> and Contributors
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
	"errors"
	"os"
	"strconv"

	"github.com/jjmrocha/beast/internal/client"
)

type csvReport struct {
	FileName string
	file     *os.File
	writer   *csv.Writer
	err      error
}

func newCsvReport(fileName string) (*csvReport, error) {
	fileHandler, err := os.Create(fileName) //nolint:gosec // writes the output file the user named
	if err != nil {
		return nil, err
	}

	csvWriter := csv.NewWriter(fileHandler)
	_ = csvWriter.Write(csvHeader())

	return &csvReport{
		FileName: fileName,
		file:     fileHandler,
		writer:   csvWriter,
	}, nil
}

func (csv *csvReport) write(response *client.Response) {
	record := encode(response)
	_ = csv.writer.Write(record)
}

func (csv *csvReport) flush() {
	csv.writer.Flush()
	csv.err = errors.Join(csv.writer.Error(), csv.file.Close())
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

func encode(response *client.Response) []string {
	var timestamp = response.Timestamp.Format(datetimeFormat)
	var request = response.Request
	var result = "Executed"
	var statusCode = ""
	var isSuccess = "false"
	var duration = ""

	if response.Duration.Nanoseconds() > 0 {
		duration = strconv.FormatInt(response.Duration.Milliseconds(), 10)
	}

	if response.IsClientError() {
		result = response.ClientError()
	} else {
		statusCode = strconv.Itoa(response.StatusCode)

		if response.IsSuccess() {
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
