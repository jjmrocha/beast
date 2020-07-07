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

package cmd

import (
	"fmt"
	"log"
	"runtime"

	"github.com/jjmrocha/beast/client"
	"github.com/jjmrocha/beast/config"
	"github.com/jjmrocha/beast/control"
	"github.com/jjmrocha/beast/data"
	"github.com/jjmrocha/beast/report"
	"github.com/jjmrocha/beast/template"
)

var errorGeneratingRequestResponse = &client.Response{
	StatusCode: -100,
}

// Run implements the `beast run ...` command
func Run(nRequests, nParallel int, fileName, configFile, dataFile string) {
	printSystem()
	printTest(fileName, configFile, dataFile, nRequests, nParallel)
	fmt.Printf("===== Preparing =====\n")
	ctrl := control.New(nRequests, nParallel)
	httpClient := createHTTPClient(configFile, nParallel)
	generators := createRequestGenerators(fileName, dataFile, nRequests)
	fmt.Printf("===== Executing =====\n")

	go func() {
		for _, gnt := range generators {
			ctrl.WaitForSlot()
			go func(g *template.Generator) {
				defer ctrl.Finish()
				request, err := g.Request()
				if err != nil {
					log.Printf("Error generating request for %s: %v\n", g.Log(), err)
					ctrl.Push(errorGeneratingRequestResponse)
					return
				}
				ctrl.WaitToExecute()
				defer ctrl.FinishExecution()
				ctrl.Push(httpClient.Execute(request))
			}(gnt)
		}
	}()

	go ctrl.CloseWhenDone()
	stats := report.NewStats(nParallel, report.NewBar(nRequests))

	for response := range ctrl.OutputChannel() {
		stats.Update(response)
	}

	stats.PrintStats()
}

func printSystem() {
	fmt.Printf("===== System =====\n")
	fmt.Printf("Operating System: %v\n", runtime.GOOS)
	fmt.Printf("System Architecture: %v\n", runtime.GOARCH)
	fmt.Printf("Logical CPUs: %v\n", runtime.NumCPU())
}

func printTest(fileName, configFile, dataFile string, nRequests, nParallel int) {
	fmt.Printf("===== Test =====\n")
	fmt.Printf("Request template: %v\n", fileName)

	if dataFile != "" {
		fmt.Printf("Sample Data: %v\n", dataFile)
	}

	if configFile != "" {
		fmt.Printf("Configuration: %v\n", configFile)
	}

	fmt.Printf("Number of requests: %v\n", nRequests)
	fmt.Printf("Number of concurrent requests: %v\n", nParallel)
}

func createHTTPClient(configFile string, nParallel int) *client.Client {
	cfg := readConfig(configFile)
	return client.NewClient(cfg, nParallel)
}

func readConfig(configFile string) *config.Config {
	if configFile == "" {
		return config.Default()
	}

	fmt.Println("- Reading configuration")
	return config.Read(configFile)
}

func readData(dataFile string) *data.Data {
	if dataFile == "" {
		return nil
	}

	fmt.Println("- Loading data file")
	return data.Read(dataFile)
}

func createRequestGenerators(fileName, dataFile string, nRequests int) []*template.Generator {
	data := readData(dataFile)
	fmt.Println("- Loading request template")
	tmpl := template.Read(fileName)
	fmt.Println("- Generating requests")
	generators, err := tmpl.BuildGenerators(nRequests, data)
	if err != nil {
		log.Fatalf("Error generating requests: %v\n", err)
	}
	return generators
}
