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

package template

import (
	"fmt"

	"github.com/jjmrocha/beast/client"
	"github.com/jjmrocha/beast/data"
)

// Generator generates BRequests
type Generator struct {
	Template *CompiledTemplate
	RecordID int
	Data     *data.Record
}

// Request uses that template and a record and returns a BRequests
func (g *Generator) Request() (*client.Request, error) {
	tmplf, err := g.Template.executeTemplate(g.RecordID, g.Data)
	if err != nil {
		return nil, err
	}

	req, err := tmplf.request()
	if err != nil {
		return nil, err
	}

	return req, nil
}

// Log generates a log message for the request
func (g *Generator) Log() string {
	return fmt.Sprintf("requestId: %v and data: %v", g.RecordID, g.Data)
}
