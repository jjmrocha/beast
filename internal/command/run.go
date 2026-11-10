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

package command

import (
	"fmt"
	"log"
	"runtime"

	"github.com/jjmrocha/beast/internal/client"
	"github.com/jjmrocha/beast/internal/config"
	"github.com/jjmrocha/beast/internal/control"
	"github.com/jjmrocha/beast/internal/data"
	"github.com/jjmrocha/beast/internal/report"
	"github.com/jjmrocha/beast/internal/template"
)

// Run implements the `beast run ...` command
func Run(nRequests, tDuration int, nParallel int, fileName, configFile, dataFile string, outputFile string) {
	fmt.Printf("===== System =====\n")
	fmt.Printf("Operating System: %v\n", runtime.GOOS)
	fmt.Printf("System Architecture: %v\n", runtime.GOARCH)
	fmt.Printf("Logical CPUs: %v\n", runtime.NumCPU())

	fmt.Printf("===== Test =====\n")
	fmt.Printf("Request template: %v\n", fileName)
	if dataFile != "" {
		fmt.Printf("Sample Data: %v\n", dataFile)
	}
	if configFile != "" {
		fmt.Printf("Configuration: %v\n", configFile)
	}
	if nRequests > 0 {
		fmt.Printf("Number of requests: %v\n", nRequests)
	} else {
		fmt.Printf("Test duration: %v seconds\n", tDuration)
	}
	fmt.Printf("Number of concurrent requests: %v\n", nParallel)

	fmt.Printf("===== Preparing =====\n")
	httpClient := createHTTPClient(configFile, nParallel)
	tmpl := readTemplate(fileName)
	data := readData(dataFile)

	fmt.Printf("===== Executing =====\n")
	ctrl := control.New(nRequests, tDuration, nParallel)
	ctrl.AsyncExecute(httpClient, tmpl, data)

	stats := report.NewStats(report.NewBar(nRequests, tDuration), outputFile)
	for response := range ctrl.OutputChannel() {
		stats.Update(response)
	}

	stats.PrintStats()
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
	cfg, err := config.Read(configFile)
	if err != nil {
		log.Fatalf("Error reading config file %s: %v\n", configFile, err)
	}
	return cfg
}

func readData(dataFile string) *data.Data {
	if dataFile == "" {
		return nil
	}

	fmt.Println("- Loading data file")
	dt, err := data.Read(dataFile)
	if err != nil {
		log.Fatalf("Error reading data file %s: %v\n", dataFile, err)
	}
	return dt
}

func readTemplate(fileName string) *template.CompiledTemplate {
	fmt.Println("- Loading request template")
	tmpl, err := template.Read(fileName)
	if err != nil {
		log.Fatalf("Error reading template file %s: %v\n", fileName, err)
	}
	tmplC, err := tmpl.Compile()
	if err != nil {
		log.Fatalf("Error compiling template: %v\n", err)
	}
	return tmplC
}
