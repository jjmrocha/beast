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
	"fmt"

	"github.com/jjmrocha/beast/client"
	"github.com/jjmrocha/beast/control"
)

const datetimeFormat = "2006-01-02 15:04:05.999"

// Output is the handler for CVS output
type Output struct {
	csvReport     *csvReport
	outputChannel chan *client.Response
	buffer        *control.Fifo
	closeFile     *control.WaitCompletion
}

// NewOutput creates new Output handler
func NewOutput(fileName string) Output {
	if fileName == "" {
		return Output{}
	}

	output := Output{
		csvReport:     newCsvReport(fileName),
		outputChannel: make(chan *client.Response, 1024),
		buffer:        control.NewFifo(),
		closeFile:     control.NewWaitCompletion(),
	}

	go asyncWrite(output)
	return output
}

// Write appends the response to the output file
func (o Output) Write(response *client.Response) {
	if !o.isInUse() {
		return
	}

	o.outputChannel <- response
}

// Close closes the CSV file
func (o Output) Close() {
	if !o.isInUse() {
		return
	}

	o.closeFile.Request()
	fmt.Printf("===== Output File =====\n")
	fmt.Printf("Output file '%s' was successfully generated\n", o.csvReport.FileName)
}

func asyncWrite(o Output) {
	defer func() {
		defer o.closeFile.Completed()
		o.csvReport.flush()
		close(o.outputChannel)
	}()

	var msgIn, msgOut bool
	run := true

	for run {
		msgIn = false
		msgOut = false

		select {
		case response, ok := <-o.outputChannel:
			if ok {
				o.buffer.Push(response)
				msgIn = true
			}
		default:
			if element := o.buffer.Pop(); element != nil {
				o.csvReport.write(element.Value)
				msgOut = true
			}
		}

		run = msgIn || msgOut || !o.closeFile.WasRequested()
	}
}

func (o Output) isInUse() bool {
	return o.csvReport != nil
}
