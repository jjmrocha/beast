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
	"github.com/jjmrocha/beast/report"
	"github.com/jjmrocha/beast/template"
)

func Run(nRequests, nParallel int, fileName, configFile string) {
	printTest(fileName, configFile, nRequests, nParallel)
	http := createHttpClient(configFile)
	control := control.New(nRequests, nParallel)
	request := readRequest(fileName)

	go func() {
		for i := 0; i < nRequests; i++ {
			control.WaitForRoom()
			go func() {
				defer control.Done()
				control.Push(http.Execute(request))
			}()
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

func printTest(fileName, configFile string, nRequests, nParallel int) {
	fmt.Printf("=== Request ===\n")
	fmt.Printf("Request template: %v\n", fileName)

	if configFile != "" {
		fmt.Printf("Configuration: %v\n", configFile)
	}

	fmt.Printf("Number of requests: %v\n", nRequests)
	fmt.Printf("Number of concurrent requests: %v\n", nParallel)
	fmt.Printf("=== Test ===\n")
}

func createHttpClient(configFile string) *client.BClient {
	config := readConfig(configFile)
	return client.Http(config)
}

func readConfig(configFile string) *config.Config {
	if configFile == "" {
		return config.Default()
	}

	return config.Read(configFile)
}

func readRequest(fileName string) *client.BRequest {
	requestTemplate := template.Read(fileName)
	request, err := requestTemplate.Generate()
	if err != nil {
		log.Fatalf("Invalid request %v: %v\n", requestTemplate, err)
	}
	return request
}
