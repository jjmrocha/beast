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
	"flag"
	"os"

	"github.com/jjmrocha/beast/cmd"
)

func main() {
	if len(os.Args) == 1 {
		cmd.Help()
		return
	}

	switch os.Args[1] {
	case "help":
		cmd.Help()
	case "config":
		configCmd(os.Args[2:])
	case "run":
		runCmd(os.Args[2:])
	case "template":
		templateCmd(os.Args[2:])
	default:
		cmd.Help()
	}
}

func configCmd(args []string) {
	if len(args) != 1 {
		cmd.Help()
		return
	}

	fileName := args[0]
	cmd.Config(fileName)
}

func runCmd(args []string) {
	runOption := flag.NewFlagSet("run", flag.ExitOnError)
	nRequests := runOption.Int("n", 0, "Number of requests")
	tDuration := runOption.Int("t", 0, "Duration of the test in seconds")
	nParallel := runOption.Int("c", 1, "Number of concurrent requests")
	configFile := runOption.String("config", "", "Config file to setup HTTP client")
	dataFile := runOption.String("data", "", "CSV file with data for request generation")
	outputFile := runOption.String("output", "", "CVS file with detailed execution results")
	runOption.Parse(args)
	nonFlagArgs := runOption.Args()

	if len(nonFlagArgs) != 1 {
		cmd.Help()
		return
	}

	if *nRequests < 0 || *nParallel <= 0 || *tDuration < 0 {
		cmd.Help()
		return
	}

	if (*nRequests == 0 && *tDuration == 0) || (*nRequests > 0 && *tDuration > 0) {
		cmd.Help()
		return
	}

	fileName := nonFlagArgs[0]
	cmd.Run(*nRequests, *tDuration, *nParallel, fileName, *configFile, *dataFile, *outputFile)
}

func templateCmd(args []string) {
	templateOption := flag.NewFlagSet("template", flag.ExitOnError)
	method := templateOption.String("m", "GET", "HTTP method")
	templateOption.Parse(args)
	nonFlagArgs := templateOption.Args()

	var url, fileName string

	switch len(nonFlagArgs) {
	case 1:
		fileName = nonFlagArgs[0]
	case 2:
		url = nonFlagArgs[0]
		fileName = nonFlagArgs[1]
	default:
		cmd.Help()
		return
	}

	cmd.Template(*method, url, fileName)
}
