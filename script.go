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
package beast

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

func ReadScript(fileName string) *BRequest {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalf("Error reading file %s: %v\n", fileName, err)
	}

	var request BRequest
	json.Unmarshal(data, &request)
	return &request
}

func WriteScript(fileName string, request *BRequest) {
	data, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error encoding request %v to JSON: %v\n", *request, err)
	}

	err = ioutil.WriteFile(fileName, data, 0666)
	if err != nil {
		log.Printf("Error writing to file %s: %v\n", fileName, err)
	}
}
