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

package data

import (
	"bytes"
	"fmt"
)

// Record represents a line on the CSV file
type Record map[string]string

func (r *Record) String() string {
	var buffer bytes.Buffer
	buffer.WriteString("{")
	for key, value := range *r {
		if buffer.Len() > 1 {
			buffer.WriteString(", ")
		}

		keyValue := fmt.Sprintf("%v: %v", key, value)
		buffer.WriteString(keyValue)
	}
	buffer.WriteString("}")
	return buffer.String()
}

// NewRecord creates a new record
func NewRecord() Record {
	return make(map[string]string)
}
