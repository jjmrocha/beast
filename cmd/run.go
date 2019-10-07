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

	"github.com/jjmrocha/beast/client"
	"github.com/jjmrocha/beast/config"
	"github.com/jjmrocha/beast/control"
	"github.com/jjmrocha/beast/data"
	"github.com/jjmrocha/beast/report"
	"github.com/jjmrocha/beast/request"
)

func Run(nRequests, nParallel int, fileName, configFile, dataFile string) {
	printTest(fileName, configFile, dataFile, nRequests, nParallel)
	http := createHTTPClient(configFile)
	control := control.New(nRequests, nParallel)
	requests := generateRequests(fileName, dataFile, nRequests)

	go func() {
		for _, request := range requests {
			control.WaitForSlot()
			go func(r *client.BRequest) {
				defer control.Done()
				control.Push(http.Execute(r))
			}(request)
		}
	}()

	go control.CloseWhenDone()
	stats := report.NewStats(nParallel)
	progress := report.NewBar(nRequests)

	for response := range control.OutputChannel() {
		stats.Update(response)
		progress.Update()
	}

	stats.Print()
}

func printTest(fileName, configFile, dataFile string, nRequests, nParallel int) {
	fmt.Printf("=== Request ===\n")
	fmt.Printf("Request template: %v\n", fileName)

	if dataFile != "" {
		fmt.Printf("Sample Data: %v\n", dataFile)
	}

	if configFile != "" {
		fmt.Printf("Configuration: %v\n", configFile)
	}

	fmt.Printf("Number of requests: %v\n", nRequests)
	fmt.Printf("Number of concurrent requests: %v\n", nParallel)
	fmt.Printf("=== Test ===\n")
}

func createHTTPClient(configFile string) *client.BClient {
	config := readConfig(configFile)
	return client.HTTP(config)
}

func readConfig(configFile string) *config.Config {
	if configFile == "" {
		return config.Default()
	}

	return config.Read(configFile)
}

func readData(dataFile string) *data.Data {
	if dataFile == "" {
		return nil
	}

	return data.Read(dataFile)
}

func generateRequests(fileName, dataFile string, nRequests int) []*client.BRequest {
	requestTemplate := request.Read(fileName)
	data := readData(dataFile)
	requests, err := requestTemplate.Generate(nRequests, data)
	if err != nil {
		log.Fatalf("Error generating requests: %v\n", err)
	}
	return requests
}
