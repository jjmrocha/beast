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

package command

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/jjmrocha/beast/internal/client"
	"github.com/jjmrocha/beast/internal/config"
	"github.com/jjmrocha/beast/internal/control"
	"github.com/jjmrocha/beast/internal/data"
	"github.com/jjmrocha/beast/internal/report"
	"github.com/jjmrocha/beast/internal/template"
)

type RunOptions struct {
	Requests    int
	Duration    int
	Concurrency int
	Template    string
	Config      string
	Data        string
	Output      string
}

// Run implements the `beast run ...` command
func Run(ctx context.Context, opts RunOptions) {
	fmt.Printf("===== System =====\n")
	fmt.Printf("Operating System: %v\n", runtime.GOOS)
	fmt.Printf("System Architecture: %v\n", runtime.GOARCH)
	fmt.Printf("Logical CPUs: %v\n", runtime.NumCPU())

	fmt.Printf("===== Test =====\n")
	fmt.Printf("Request template: %v\n", opts.Template)
	if opts.Data != "" {
		fmt.Printf("Sample Data: %v\n", opts.Data)
	}
	if opts.Config != "" {
		fmt.Printf("Configuration: %v\n", opts.Config)
	}
	if opts.Requests > 0 {
		fmt.Printf("Number of requests: %v\n", opts.Requests)
	} else {
		fmt.Printf("Test duration: %v seconds\n", opts.Duration)
	}
	fmt.Printf("Number of concurrent requests: %v\n", opts.Concurrency)

	fmt.Printf("===== Preparing =====\n")
	httpClient := createHTTPClient(opts.Config, opts.Concurrency)
	tmpl := readTemplate(opts.Template)
	data := readData(opts.Data)

	fmt.Printf("===== Executing =====\n")
	limits := control.Limits{Requests: opts.Requests, Duration: opts.Duration, Workers: opts.Concurrency}
	results := control.Start(ctx, limits, httpClient, tmpl, data)

	output, err := report.NewOutput(opts.Output)
	if err != nil {
		log.Fatalf("Error creating output file %s: %v\n", opts.Output, err)
	}

	stats := report.NewStats(report.NewBar(opts.Requests, opts.Duration))
	for response := range results {
		stats.Update(response)
		output.Write(response)
	}

	stats.PrintStats(os.Stdout)
	output.Close(os.Stdout)
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
