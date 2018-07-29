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

package utils_test

import (
	"testing"

	capturer "github.com/kami-zh/go-capturer"
	"github.com/retr0h/gofile/utils"
	"github.com/stretchr/testify/assert"
)

func TestPrintError(t *testing.T) {
	got := capturer.CaptureStderr(func() {
		utils.PrintError("foo")
	})
	want := "\x1b[31mERROR\x1b[0m: foo\n"

	assert.Equal(t, want, got)
}

func TestPrintErrorAndExit(t *testing.T) {
	oldOsExit := utils.OsExit
	defer func() { utils.OsExit = oldOsExit }()

	var got int
	myExit := func(code int) {
		got = code
	}

	utils.OsExit = myExit

	out := capturer.CaptureStderr(func() {
		utils.PrintErrorAndExit("foo")
	})

	assert.Equal(t, "\x1b[31mERROR\x1b[0m: foo\n", out)
	assert.Equal(t, 1, got)

	utils.PrintErrorAndExit("foo", 5)
	assert.Equal(t, 5, got)
}
