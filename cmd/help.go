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

const help = `The Beast - Stress testing for RESTful APIs
Usage:
	beast [help]
	beast generate [-m <http method>] [url] <script>
	beast run [-n <number of requests>] [-c <number of concurrent requests>] <script>
Where:
	generate Creates a script template, using user provided parameters
	         -m     string HTTP method (default "GET")
	         url    string Endpoint to be tested
	         script string JSON file with details about the request to test

	run      Executes a script and presents a report with execution results
	         -c     int    Number of concurrent requests (default 1)
	         -n     int    Number of requests (default 1)
	         script string JSON file with details about the request to test 
`

func Help() {
	fmt.Println(help)
}
