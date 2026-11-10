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

package template

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Read reads an HTTP request template from a file
func Read(fileName string) (*Template, error) {
	data, err := os.ReadFile(fileName) //nolint:gosec // reads the template file the user named
	if err != nil {
		return nil, err
	}

	if isJSON(fileName) {
		return readJSON(fileName, data)
	}

	return readYAML(fileName, data), nil
}

// Write writes an HTTP request template to a file
func Write(fileName string, tmpl *Template) {
	if isJSON(fileName) {
		writeJSON(fileName, tmpl)
	} else {
		writeYAML(fileName, tmpl)
	}
}

func writeFile(data []byte, fileName string) {
	if err := os.WriteFile(fileName, data, 0o600); err != nil {
		log.Printf("Error writing template to file %s: %v\n", fileName, err)
	}
}

// JSON

func isJSON(fileName string) bool {
	lowerCaseFileName := strings.ToLower(fileName)
	return strings.HasSuffix(lowerCaseFileName, ".json")
}

func readJSON(fileName string, data []byte) (*Template, error) {
	var tmpl Template
	if err := json.Unmarshal(data, &tmpl); err != nil {
		log.Printf("Invalid JSON template file %s: %v\n", fileName, err)
	}

	body, read, err := externalBody(tmpl.Body)
	if err != nil {
		return nil, err
	}

	if read {
		tmpl.Body = body
	}

	return &tmpl, nil
}

func externalBody(body string) (string, bool, error) {
	if strings.HasPrefix(body, "@") {
		fileName := body[1:]
		data, err := os.ReadFile(fileName) //nolint:gosec // reads the body file the user named in the template
		if err != nil {
			return "", false, fmt.Errorf("external body: %w", err)
		}

		return string(data), true, nil
	}

	return "", false, nil
}

func writeJSON(fileName string, tmpl *Template) {
	data, err := json.MarshalIndent(tmpl, "", "\t")
	if err != nil {
		log.Printf("Error encoding request %v to JSON: %v\n", tmpl, err)
	}

	writeFile(data, fileName)
}

// YAML

type templateY struct {
	Method   string            `yaml:"method"`
	Endpoint string            `yaml:"endpoint"`
	Headers  map[string]string `yaml:"headers,omitempty"`
	Body     string            `yaml:"request-body,omitempty"`
}

func toYamlTemplate(tmpl *Template) *templateY {
	return &templateY{
		Method:   tmpl.Method,
		Endpoint: tmpl.Endpoint,
		Headers:  toHeaderMap(tmpl.Headers),
		Body:     tmpl.Body,
	}
}

func toHeaderMap(headers []Header) map[string]string {
	if headers == nil {
		return nil
	}

	headerMap := make(map[string]string)

	for _, header := range headers {
		headerMap[header.Key] = header.Value
	}

	return headerMap
}

func fromYamlTemplate(tmply *templateY) *Template {
	return &Template{
		Method:   tmply.Method,
		Endpoint: tmply.Endpoint,
		Headers:  fromHeaderMap(tmply.Headers),
		Body:     tmply.Body,
	}
}

func fromHeaderMap(headers map[string]string) []Header {
	if headers == nil {
		return nil
	}

	tHeaders := make([]Header, 0, len(headers))

	for key, value := range headers {
		hd := Header{
			Key:   key,
			Value: value,
		}
		tHeaders = append(tHeaders, hd)
	}

	return tHeaders
}

func writeYAML(fileName string, tmpl *Template) {
	tmply := toYamlTemplate(tmpl)

	data, err := yaml.Marshal(tmply)
	if err != nil {
		log.Printf("Error encoding request %v to YAML: %v\n", tmpl, err)
	}

	writeFile(data, fileName)
}

func readYAML(fileName string, data []byte) *Template {
	var tmply templateY
	if err := yaml.Unmarshal(data, &tmply); err != nil {
		log.Printf("Invalid YAML template file %s: %v\n", fileName, err)
	}

	return fromYamlTemplate(&tmply)
}
