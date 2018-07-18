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

package utils

import (
	"fmt"
	"os"

	"github.com/logrusorgru/aurora"
)

var (
	// OsExit should be private?!
	OsExit = os.Exit
)

// PrintError formats and prints the provided string for error messages.
func PrintError(msg string) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", aurora.Red("ERROR"), msg)
}

// PrintErrorAndExit prints the error message and os.Exits with the optionally
// provided error code.
func PrintErrorAndExit(msg string, exitCodeOptional ...int) {
	exitCode := 1
	if len(exitCodeOptional) > 0 {
		exitCode = exitCodeOptional[0]
	}

	PrintError(msg)
	OsExit(exitCode)
}
