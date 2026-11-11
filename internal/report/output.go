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
	"fmt"
	"io"

	"github.com/jjmrocha/beast/internal/client"
	"github.com/jjmrocha/go-algo/queue"
)

const datetimeFormat = "2006-01-02 15:04:05.999"

// Output writes results to the CSV output file
type Output struct {
	csvReport *csvReport
	records   *queue.SyncQueue[*client.Response]
	wake      chan struct{}
	done      chan struct{}
}

// NewOutput creates new Output handler
func NewOutput(fileName string) (Output, error) {
	if fileName == "" {
		return Output{}, nil
	}

	report, err := newCsvReport(fileName)
	if err != nil {
		return Output{}, err
	}

	output := Output{
		csvReport: report,
		records:   queue.NewSyncQueue[*client.Response](),
		wake:      make(chan struct{}, 1),
		done:      make(chan struct{}),
	}

	go output.writeAll()
	return output, nil
}

// Write appends the response to the output file
func (o Output) Write(response *client.Response) {
	if !o.isInUse() {
		return
	}

	o.records.Enqueue(response)
	select {
	case o.wake <- struct{}{}:
	default:
	}
}

// Close closes the CSV file
func (o Output) Close(w io.Writer) {
	if !o.isInUse() {
		return
	}

	close(o.wake)
	<-o.done
	_, _ = fmt.Fprintf(w, "===== Output File =====\n")
	if err := o.csvReport.err; err != nil {
		_, _ = fmt.Fprintf(w, "Error writing output file '%s': %v\n", o.csvReport.FileName, err)
		return
	}
	_, _ = fmt.Fprintf(w, "Output file '%s' was successfully generated\n", o.csvReport.FileName)
}

func (o Output) writeAll() {
	defer close(o.done)
	for open := true; open; {
		_, open = <-o.wake
		for response := range o.records.Drain() {
			o.csvReport.write(response)
		}
	}
	o.csvReport.flush()
}

func (o Output) isInUse() bool {
	return o.csvReport != nil
}
