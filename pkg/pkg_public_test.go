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

package pkg_test

import (
	"errors"
	"os"
	"path"
	"path/filepath"
	"testing"

	capturer "github.com/kami-zh/go-capturer"
	"github.com/retr0h/gofile/pkg"
	"github.com/stretchr/testify/assert"
)

var p pkg.Packages

func TestUnmarshalYAMLDoesNotParseYAMLAndReturnsError(t *testing.T) {
	data := `
---
%foo:
`
	err := p.UnmarshalYAML([]byte(data))
	want := "yaml: line 3: found unexpected non-alphabetical character"

	assert.Equal(t, want, err.Error())
	assert.Error(t, err)
}

func TestUnmarshalYAMLDoesNotValidateYAMLAndReturnsError(t *testing.T) {
	data := `
---
foo: bar
`
	err := p.UnmarshalYAML([]byte(data))
	want := errors.New("(root): Invalid type. Expected: array, given: object")

	assert.Equal(t, want, err)
}

func TestUnmarshalYAML(t *testing.T) {
	data := `
---
- url: github.com/simeji/jid/cmd/jid
`
	err := p.UnmarshalYAML([]byte(data))
	want := "github.com/simeji/jid/cmd/jid"

	assert.NoError(t, err)
	assert.Equal(t, want, p.Packages[0].URL)
}

func TestUnmarshalYAMLFileReturnsErrorWithMissingFile(t *testing.T) {
	filename := "missing.yml"

	err := p.UnmarshalYAMLFile(filename)
	want := "open missing.yml: no such file or directory"

	assert.Equal(t, want, err.Error())
	assert.Error(t, err)
}

func TestUnmarshalYAMLFile(t *testing.T) {
	filename := path.Join("..", "test", "gofile.yml")
	p.UnmarshalYAMLFile(filename)

	assert.NotNil(t, p.Packages[0].URL)
}

func TestInstall(t *testing.T) {
	data := `
---
- url: github.com/golang/example/hello
`
	pkgDir := getGolangExamplePackageDir()
	p.UnmarshalYAML([]byte(data))
	err := p.Install()

	assert.NoError(t, err)
	assert.DirExists(t, pkgDir)

	defer os.RemoveAll(pkgDir)
}

func TestInstallDebugAddsVFlag(t *testing.T) {
	data := `
---
- url: github.com/golang/example/hello
`
	p := pkg.Packages{
		Debug: true,
	}
	p.UnmarshalYAML([]byte(data))
	got := capturer.CaptureStdout(func() {
		err := p.Install()
		assert.NoError(t, err)
	})
	want := "Installing: \x1b[36mgithub.com/golang/example/hello\x1b[0m\nCOMMAND: \x1b[30;41mgo get -v github.com/golang/example/hello\x1b[0m\n"

	assert.Equal(t, want, got)

	defer os.RemoveAll(getGolangExamplePackageDir())
}

func TestInstallReturnsErrorWhenRunCmdErrors(t *testing.T) {
	data := `
---
- url: invalid.
`
	p := pkg.Packages{
		Debug: true,
	}
	p.UnmarshalYAML([]byte(data))

	err := p.Install()
	assert.Error(t, err)
}

func TestRunCommand(t *testing.T) {
	got := capturer.CaptureStdout(func() {
		err := p.RunCmd("ls")
		assert.NoError(t, err)
	})

	assert.Empty(t, got)
}

func TestRunCommandPrintsStreamingStdout(t *testing.T) {
	p := pkg.Packages{
		Debug: true,
	}
	got := capturer.CaptureStdout(func() {
		err := p.RunCmd("echo", "-n", "foo")
		assert.NoError(t, err)
	})
	want := "COMMAND: \x1b[30;41mecho -n foo\x1b[0m\nfoo"

	assert.Equal(t, want, got)
}

func TestRunCommandPrintsStreamingStderr(t *testing.T) {
	p := pkg.Packages{
		Debug: true,
	}
	got := capturer.CaptureStderr(func() {
		err := p.RunCmd("cat", "foo")
		assert.Error(t, err)
	})
	want := "cat: foo: No such file or directory\n"

	assert.Equal(t, want, got)
}

func TestRunCommandReturnsError(t *testing.T) {
	err := p.RunCmd("false")

	assert.Error(t, err)
}

func getGolangExamplePackageDir() string {
	return filepath.Join(os.Getenv("GOPATH"), "src", "github.com", "golang", "example")

}
