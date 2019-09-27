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
package cmd

import (
	"log"
	"net/http"
	"sync"

	"github.com/jjmrocha/beast"
)

func Run(nRequests, nParallel int, fileName string) {
	req := readRequest(fileName)
	output := make(chan *beast.BResponse, nRequests)
	client := beast.HttpClient()
	semaphore := beast.NewSemaphore(nParallel)
	var wg sync.WaitGroup
	wg.Add(nRequests)

	go func() {
		for i := 0; i < nRequests; i++ {
			semaphore.Acquire()
			go func(r *http.Request) {
				defer wg.Done()
				output <- beast.Execute(client, r)
				semaphore.Release()
			}(req)
		}
	}()

	go func() {
		wg.Wait()
		close(output)
	}()

	report := new(beast.Report)

	for response := range output {
		report.Update(response)
	}

	report.Print()
}

func readRequest(fileName string) *http.Request {
	script := beast.ReadScript(fileName)
	request, err := beast.Convert(script)
	if err != nil {
		log.Fatalf("Invalid request %v: %v\n", *script, err)
	}
	return request
}
