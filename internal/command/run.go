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
	"io"
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
func Run(ctx context.Context, w io.Writer, opts RunOptions) error {
	printSystem(w)
	printTest(w, opts)

	_, _ = fmt.Fprintf(w, "===== Preparing =====\n")
	httpClient, err := createHTTPClient(w, opts.Config, opts.Concurrency)
	if err != nil {
		return err
	}
	tmpl, err := readTemplate(w, opts.Template)
	if err != nil {
		return err
	}
	data, err := readData(w, opts.Data)
	if err != nil {
		return err
	}

	_, _ = fmt.Fprintf(w, "===== Executing =====\n")
	output, err := createOutput(opts.Output)
	if err != nil {
		return err
	}

	limits := control.Limits{Requests: opts.Requests, Duration: opts.Duration, Workers: opts.Concurrency}
	results := control.Start(ctx, limits, httpClient, tmpl, data)
	stats := report.NewStats(report.NewBar(opts.Requests, opts.Duration))
	for response := range results {
		stats.Update(response)
		output.Write(response)
	}

	stats.PrintStats(w)
	output.Close(w)
	return nil
}

func printSystem(w io.Writer) {
	_, _ = fmt.Fprintf(w, "===== System =====\n")
	_, _ = fmt.Fprintf(w, "Operating System: %v\n", runtime.GOOS)
	_, _ = fmt.Fprintf(w, "System Architecture: %v\n", runtime.GOARCH)
	_, _ = fmt.Fprintf(w, "Logical CPUs: %v\n", runtime.NumCPU())
}

func printTest(w io.Writer, opts RunOptions) {
	_, _ = fmt.Fprintf(w, "===== Test =====\n")
	_, _ = fmt.Fprintf(w, "Request template: %v\n", opts.Template)
	if opts.Data != "" {
		_, _ = fmt.Fprintf(w, "Sample Data: %v\n", opts.Data)
	}
	if opts.Config != "" {
		_, _ = fmt.Fprintf(w, "Configuration: %v\n", opts.Config)
	}
	if opts.Requests > 0 {
		_, _ = fmt.Fprintf(w, "Number of requests: %v\n", opts.Requests)
	} else {
		_, _ = fmt.Fprintf(w, "Test duration: %v seconds\n", opts.Duration)
	}
	_, _ = fmt.Fprintf(w, "Number of concurrent requests: %v\n", opts.Concurrency)
}

func createHTTPClient(w io.Writer, configFile string, nParallel int) (*client.Client, error) {
	cfg, err := readConfig(w, configFile)
	if err != nil {
		return nil, err
	}
	return client.NewClient(cfg, nParallel), nil
}

func readConfig(w io.Writer, configFile string) (*config.Config, error) {
	if configFile == "" {
		return config.Default(), nil
	}

	_, _ = fmt.Fprintln(w, "- Reading configuration")
	cfg, err := config.Read(configFile)
	if err != nil {
		return nil, fmt.Errorf("reading config file %s: %w", configFile, err)
	}
	return cfg, nil
}

func readTemplate(w io.Writer, fileName string) (*template.CompiledTemplate, error) {
	_, _ = fmt.Fprintln(w, "- Loading request template")
	tmpl, err := template.Read(fileName)
	if err != nil {
		return nil, fmt.Errorf("reading template file %s: %w", fileName, err)
	}
	tmplC, err := tmpl.Compile()
	if err != nil {
		return nil, fmt.Errorf("compiling template: %w", err)
	}
	return tmplC, nil
}

func readData(w io.Writer, dataFile string) (*data.Data, error) {
	if dataFile == "" {
		return nil, nil
	}

	_, _ = fmt.Fprintln(w, "- Loading data file")
	dt, err := data.Read(dataFile)
	if err != nil {
		return nil, fmt.Errorf("reading data file %s: %w", dataFile, err)
	}
	return dt, nil
}

func createOutput(fileName string) (report.Output, error) {
	output, err := report.NewOutput(fileName)
	if err != nil {
		return report.Output{}, fmt.Errorf("creating output file %s: %w", fileName, err)
	}
	return output, nil
}
