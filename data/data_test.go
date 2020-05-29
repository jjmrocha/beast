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
	"reflect"
	"testing"
)

func TestRead(t *testing.T) {
	// given
	expectedData := &Data{
		fields: []string{"A", "B"},
		records: []Record{
			{"A": "a1", "B": "b1"},
			{"A": "a2", "B": "b2"},
		},
	}
	// when
	result := Read("../testdata/data.csv")
	// then
	if !reflect.DeepEqual(result, expectedData) {
		t.Errorf("got %v expected %v", result, expectedData)
	}
}

func TestNext(t *testing.T) {
	// given
	expectedFirstRecord := Record{
		"A": "a1",
		"B": "b1",
	}
	expectedSecondRecord := Record{
		"A": "a2",
		"B": "b2",
	}
	dt := &Data{
		fields:  []string{"A", "B"},
		records: []Record{expectedFirstRecord, expectedSecondRecord},
	}
	// when
	firstRecord := dt.Next()
	secondRecord := dt.Next()
	thirdRecord := dt.Next()
	// then
	if !reflect.DeepEqual(firstRecord, &expectedFirstRecord) {
		t.Errorf("got %v expected %v", firstRecord, &expectedFirstRecord)
	}

	if !reflect.DeepEqual(secondRecord, &expectedSecondRecord) {
		t.Errorf("got %v expected %v", secondRecord, &expectedSecondRecord)
	}

	if !reflect.DeepEqual(thirdRecord, &expectedFirstRecord) {
		t.Errorf("got %v expected %v", thirdRecord, &expectedFirstRecord)
	}
}
