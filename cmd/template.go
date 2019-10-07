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
	"strings"

	"github.com/jjmrocha/beast/request"
)

func Template(method, url, fileName string) {
	if url == "" {
		writeTemplate(fileName)
	} else {
		writeRequest(method, url, fileName)
	}
}

func writeTemplate(fileName string) {
	req := request.TRequest{
		Method:   "Use Http method: GET/POST/PUT/DELETE",
		Endpoint: "Http URL to be invoked",
		Headers: []request.THeader{
			{Key: "User-Agent", Value: "Beast/1"},
		},
		Body: "Optional, enter body to send with POST or PUT",
	}
	request.Write(fileName, &req)
	fmt.Printf("File %s was created, please edit before use\n", fileName)
}

func writeRequest(method, url, fileName string) {
	req := request.TRequest{
		Method:   method,
		Endpoint: url,
		Headers: []request.THeader{
			{Key: "User-Agent", Value: "Beast/1"},
		},
	}

	if usesBody(method) {
		req.Body = "Optional, enter body to send with POST or PUT"
	}

	request.Write(fileName, &req)
	fmt.Printf("File %s was created for '%s %s'\n", fileName, method, url)
}

func usesBody(method string) bool {
	upperCaseMethod := strings.ToUpper(method)
	return upperCaseMethod == "PUT" || upperCaseMethod == "POST"
}
