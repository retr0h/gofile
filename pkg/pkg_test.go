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
	"errors"
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
	"github.com/xeipuuv/gojsonschema"
)

var p Packages

func TestValidateSchemaHasErrorReturnsError(t *testing.T) {
	originalJSONSchemaValidator := jsonSchemaValidator
	jsonSchemaValidator = func(gojsonschema.JSONLoader, gojsonschema.JSONLoader) (*gojsonschema.Result, error) {
		return nil, errors.New("Failed to load schema")
	}
	defer func() { jsonSchemaValidator = originalJSONSchemaValidator }()

	err := p.validate([]byte(""))
	assert.Equal(t, "Failed to load schema", err.Error())
}

func TestValidateWithoutRootArrayReturnsError(t *testing.T) {
	var data = `
---
foo:
`
	jsonData, _ := yaml.YAMLToJSON([]byte(data))
	err := p.validate([]byte(jsonData))
	assert.Error(t, err)

	expectedError := errors.New("(root): Invalid type. Expected: array, given: object")
	assert.Equal(t, expectedError, err)
}

func TestValidateWithoutStringReturnsError(t *testing.T) {
	var data = `
---
- url:
`
	jsonData, _ := yaml.YAMLToJSON([]byte(data))
	err := p.validate([]byte(jsonData))
	assert.Error(t, err)

	messages := []string{
		"0.url: Invalid type. Expected: string, given: null",
	}
	for _, msg := range messages {
		assert.Contains(t, err.Error(), msg)
	}
}

func TestValidateWithoutURLKeyReturnsError(t *testing.T) {
	var data = `
---
- foo: github.com/simeji/jid/cmd/jid
`
	jsonData, _ := yaml.YAMLToJSON([]byte(data))
	err := p.validate([]byte(jsonData))
	assert.Error(t, err)

	expectedError := errors.New("url: url is required")
	assert.Equal(t, expectedError, err)
}

func TestValidate(t *testing.T) {
	var data = `
---
- url: https://example.com/user/repo.git
`
	jsonData, _ := yaml.YAMLToJSON([]byte(data))
	err := p.validate([]byte(jsonData))
	assert.NoError(t, err)
}
