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

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/jjmrocha/beast/internal/command"
)

const (
	usageError  = 2
	interrupted = 130
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		command.Help()
		return 0
	}

	switch args[0] {
	case "help":
		command.Help()
		return 0
	case "config":
		return configCmd(args[1:])
	case "run":
		return runCmd(args[1:])
	case "template":
		return templateCmd(args[1:])
	default:
		command.Help()
		return usageError
	}
}

func configCmd(args []string) int {
	if len(args) != 1 {
		command.Help()
		return usageError
	}

	fileName := args[0]
	if err := command.Config(fileName); err != nil {
		log.Printf("Error writing config file %s: %v\n", fileName, err) //nolint:gosec // the file name comes from the user's own command line
		return 1
	}
	return 0
}

func runCmd(args []string) int {
	runOption := flag.NewFlagSet("run", flag.ExitOnError)
	nRequests := runOption.Int("n", 0, "Number of requests")
	tDuration := runOption.Int("t", 0, "Duration of the test in seconds")
	nParallel := runOption.Int("c", 1, "Number of concurrent requests")
	configFile := runOption.String("config", "", "Config file to setup HTTP client")
	dataFile := runOption.String("data", "", "CSV file with data for request generation")
	outputFile := runOption.String("output", "", "CSV file with detailed execution results")
	_ = runOption.Parse(args)
	nonFlagArgs := runOption.Args()

	if len(nonFlagArgs) != 1 {
		command.Help()
		return usageError
	}

	if *nRequests < 0 || *nParallel <= 0 || *tDuration < 0 {
		command.Help()
		return usageError
	}

	if (*nRequests == 0 && *tDuration == 0) || (*nRequests > 0 && *tDuration > 0) {
		command.Help()
		return usageError
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	announce := context.AfterFunc(ctx, func() {
		log.Println("Stopping: waiting for requests in flight. Press Ctrl-C again to exit now.")
		stop()
	})
	defer announce()

	fileName := nonFlagArgs[0]
	command.Run(ctx, command.RunOptions{
		Requests:    *nRequests,
		Duration:    *tDuration,
		Concurrency: *nParallel,
		Template:    fileName,
		Config:      *configFile,
		Data:        *dataFile,
		Output:      *outputFile,
	})

	if ctx.Err() != nil {
		return interrupted
	}

	return 0
}

func templateCmd(args []string) int {
	templateOption := flag.NewFlagSet("template", flag.ExitOnError)
	method := templateOption.String("m", "GET", "HTTP method")
	_ = templateOption.Parse(args)
	nonFlagArgs := templateOption.Args()

	var url, fileName string

	switch len(nonFlagArgs) {
	case 1:
		fileName = nonFlagArgs[0]
	case 2:
		url = nonFlagArgs[0]
		fileName = nonFlagArgs[1]
	default:
		command.Help()
		return usageError
	}

	if err := command.Template(command.TemplateOptions{Method: *method, URL: url, File: fileName}); err != nil {
		log.Printf("Error writing template file %s: %v\n", fileName, err)
		return 1
	}
	return 0
}
