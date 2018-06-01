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

package pkg

import (
	"bytes"
	"os"
	"testing"

	log "github.com/Sirupsen/logrus"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
)

var p Packages

func CaptureLogOutput(fn func()) string {
	var buffer bytes.Buffer

	log.SetOutput(&buffer)
	fn()
	log.SetOutput(os.Stderr)

	return buffer.String()
}

func TestValidateWithoutRootArrayReturnsErrorAndLogsValidationError(t *testing.T) {
	var data = `
---
foo:
`

	logOutput := CaptureLogOutput(func() {
		jsonData, _ := yaml.YAMLToJSON([]byte(data))
		err := p.validate([]byte(jsonData))
		assert.Error(t, err)
	})

	msg := "The document is not valid - (root): Invalid type. Expected: array, given: object."
	assert.Contains(t, logOutput, msg)
}

func TestValidateWithoutStringReturnsErrorAndLogsValidationError(t *testing.T) {
	var data = `
---
- url:
`

	logOutput := CaptureLogOutput(func() {
		jsonData, _ := yaml.YAMLToJSON([]byte(data))
		err := p.validate([]byte(jsonData))
		assert.Error(t, err)
	})

	msgs := []string{
		"The document is not valid - 0.url: Invalid type. Expected: string, given: null.",
	}

	for _, msg := range msgs {
		assert.Contains(t, logOutput, msg)
	}
}

// Test minitems
// Test uniqueItems

func TestValidate(t *testing.T) {
	var data = `
---
- url: https://example.com/user/repo.git
`

	jsonData, _ := yaml.YAMLToJSON([]byte(data))
	err := p.validate([]byte(jsonData))
	assert.NoError(t, err)
}
