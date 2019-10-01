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
	"fmt"
	"log"
	"sync"

	"github.com/jjmrocha/beast/client"
	"github.com/jjmrocha/beast/report"
	"github.com/jjmrocha/beast/template"
)

func Run(nRequests, nParallel int, fileName string) {
	printTest(fileName, nRequests, nParallel)
	request := readRequest(fileName)
	output := make(chan *client.BResponse, nRequests)
	http := client.Http()
	semaphore := client.NewSemaphore(nParallel)
	var wg sync.WaitGroup
	wg.Add(nRequests)

	go func() {
		for i := 0; i < nRequests; i++ {
			semaphore.Acquire()
			go func() {
				defer wg.Done()
				output <- http.Execute(request)
				semaphore.Release()
			}()
		}
	}()

	go func() {
		wg.Wait()
		close(output)
	}()

	stats := report.NewStats(nParallel)
	progress := report.NewBar(nRequests)

	for response := range output {
		stats.Update(response)
		progress.Update()
	}

	stats.Print()
}

func printTest(fileName string, nRequests, nParallel int) {
	fmt.Printf("=== Test ===\n")
	fmt.Printf("Script to execute: %v\n", fileName)
	fmt.Printf("Number of requests: %v\n", nRequests)
	fmt.Printf("Number of concurrent requests: %v\n", nParallel)
}

func readRequest(fileName string) *client.BRequest {
	requestTemplate := template.Read(fileName)
	request, err := requestTemplate.Generate()
	if err != nil {
		log.Fatalf("Invalid request %v: %v\n", requestTemplate, err)
	}
	return request
}
