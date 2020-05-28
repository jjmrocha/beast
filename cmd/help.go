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

import "fmt"

const version = "v2.4.0"
const help = `the Beast %v - Stress testing for RESTful APIs

Usage:
   beast [help]
   beast config <configFile>
   beast template [-m <http method>] [url] <templateFile>
   beast run [-n <number of requests>] [-c <number of concurrent requests>] 
             [-config <configFile>] [-data <dataFile>] <templateFile>

Where:
   config   Creates a file with the default parameters to setup HTTP connections
            configFile   string Name of the file to be created
			 			
   template Creates a request template file, using user-provided parameters
            -m           string HTTP method (default "GET")
            url          string Endpoint to be tested
            templateFile string JSON/YAML file with details about the request to test

   run      Executes a script and presents a report with execution results
            -c           int    Number of concurrent requests (default 1)
            -n           int    Number of requests (default 1)
            -config      string Config file to setup HTTP client
            -data        string CSV file with data for request generation
            templateFile string JSON/YAML file with details about the request to test 

`

// Help implements the `beast [help]` command
func Help() {
	fmt.Printf(help, version)
}
