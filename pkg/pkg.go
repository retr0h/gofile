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

// Package pkg responsible for downloading go packages from a manifest named
// gofile.yml.
package pkg

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/caarlos0/spin"
	"github.com/ghodss/yaml"
	"github.com/logrusorgru/aurora"
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

var (
	jsonSchemaValidator = gojsonschema.Validate
)

// Package containing the go package details.  All fields are required unless
// otherwise specified.
type Package struct {
	URL string `yaml:"url"`
}

// Packages contains a list of `Package` structs initalized by the cli
// via the `--filename` flag.
type Packages struct {
	Packages []Package
	Debug    bool // Debug option set from CLI with debug state.
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
	result, err := jsonSchemaValidator(schemaLoader, documentLoader)
	if err != nil {
		return err
	}

	// Build schema validation failures.
	if !result.Valid() {
		var errstrings []string
		for _, desc := range result.Errors() {
			err := fmt.Errorf("%s", desc)
			errstrings = append(errstrings, err.Error())
		}

		return errors.New(strings.Join(errstrings, "\n"))
	}

	return nil
}

// Install loops through the `Packages` struct and calls `go get` against
// the resulting package.
func (p *Packages) Install() error {
	for _, pkg := range p.Packages {
		if !p.Debug {
			s := spin.New("%s ")
			s.Set(spin.Spin8)
			s.Start()
			defer s.Stop()
			// Allow the spinner to show when the install returns too quickly.
			time.Sleep(5 * time.Millisecond)
		}
		fmt.Printf("Installing: %s\n", aurora.Cyan(pkg.URL))

		goCmdArgs := []string{"get"}
		if p.Debug {
			goCmdArgs = append(goCmdArgs, "-v")
		}
		goCmdArgs = append(goCmdArgs, pkg.URL)
		if err := p.RunCmd("go", goCmdArgs...); err != nil {
			return err
		}
	}
	return nil
}

// RunCmd execute the provided command with args.
func (p *Packages) RunCmd(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if p.Debug {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		commands := strings.Join(cmd.Args, " ")
		msg := fmt.Sprintf("COMMAND: %s", aurora.Colorize(commands, aurora.BlackFg|aurora.RedBg))
		fmt.Println(msg)
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
