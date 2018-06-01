// Copyright (c) 2018 John Dewey

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER
// DEALINGS IN THE SOFTWARE.

// Package pkg TODO
package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/xeipuuv/gojsonschema"
)

const pkgSchema = `
{
  "type": "array",
  "$schema": "http://json-schema.org/draft-04/schema#",
  "description": "",
  "minItems": 1,
  "uniqueItems": true,
  "items": {
    "type": "object",
    "required": [
      "url"
    ],
    "properties": {
      "url": {
        "type": "string"
      }
    }
  }
}
`

// Package containing the go package details.  All fields are required unless
// otherwise specified.
type Package struct {
	URL string `yaml:"url"`
}

// Packages contains a list of `Package` structs initalized by the cli
// via the `--filename` flag.
type Packages struct {
	Packages []Package
}

// UnmarshalYAML decodes the first YAML document found within the data byte
// slice, passes the string through a generic YAML-to-JSON converter, performs
// validation, provides the resulting JSON to json.Unmarshal, and assigns the
// decoded values to the Packages struct.
func (p *Packages) UnmarshalYAML(data []byte) error {
	jsonData, err := yaml.YAMLToJSON(data)
	if err != nil {
		return err
	}

	// Validate the jsonData against the schema.
	if err = p.validate(jsonData); err != nil {
		return err
	}

	// Unmarshal the jsonData to the `Packages` struct.
	err = json.Unmarshal(jsonData, &p.Packages)
	return err
}

// UnmarshalYAMLFile reads the file named by `filename` and passes the source
// data byte slice to `UnmarshalYAML` for decoding.
func (p *Packages) UnmarshalYAMLFile(filename string) error {
	// Open and Read the provided filename.
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	// Unmarshal the file contents.
	err = p.UnmarshalYAML([]byte(source))
	return err
}

// Validate the the data byte slice against the `pkgSchema`.
func (p *Packages) validate(data []byte) error {
	schemaLoader := gojsonschema.NewStringLoader(pkgSchema)
	documentLoader := gojsonschema.NewBytesLoader(data)
	// Validate the document against the schema.
	result, err := gojsonschema.Validate(schemaLoader, documentLoader)

	if err != nil {
		// Happens when schema has errors
		fmt.Println(err)
		return err
	}

	// Log schema validation failures.
	if !result.Valid() {
		for _, desc := range result.Errors() {
			log.Error(fmt.Sprintf("The document is not valid - %s.", desc))
		}
		return errors.New("Invalid YAML provided")
	}
	return nil
}

// Install loops through the `Packages` struct and calls `go get` against
// the resulting package.
func (p *Packages) Install() error {
	for _, pkg := range p.Packages {
		var stderr bytes.Buffer
		cmd := exec.Command("go", "get", pkg.URL)
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			log.Error(stderr.String())

			return err
		}
	}
	return nil
}
