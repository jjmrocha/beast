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
		cmd.WriteDefaultConfig()
	case "run":
		runCmd(os.Args[2:])
	case "generate":
		templateCmd(os.Args[2:])
	default:
		cmd.Help()
	}
}

func runCmd(args []string) {
	runOption := flag.NewFlagSet("run", flag.ExitOnError)
	nRequests := runOption.Int("n", 1, "Number of requests")
	nParallel := runOption.Int("c", 1, "Number of concurrent requests")
	runOption.Parse(args)
	nonFlagArgs := runOption.Args()

	if len(nonFlagArgs) != 1 {
		cmd.Help()
		return
	}

	fileName := nonFlagArgs[0]
	cmd.Run(*nRequests, *nParallel, fileName)
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
