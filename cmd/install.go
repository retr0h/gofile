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

package cmd

import (
	"fmt"

	"github.com/retr0h/gofile/pkg"
	"github.com/retr0h/gofile/utils"
	"github.com/spf13/cobra"
)

var (
	fileName string
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install gofile packages",
	RunE: func(cmd *cobra.Command, args []string) error {
		p := pkg.Packages{
			Debug: debug,
		}

		if err := p.UnmarshalYAMLFile(fileName); err != nil {
			msg := fmt.Sprintf("An error occurred unmarshalling '%s'.\n%s\n", fileName, err)
			utils.PrintErrorAndExit(msg)
		}

		if err := p.Install(); err != nil {
			msg := fmt.Sprintf("An error occurred installing packages.\n%s\n", err)
			utils.PrintErrorAndExit(msg)
		}

		return nil
	},
}

func init() {
	installCmd.PersistentFlags().StringVarP(&fileName, "filename", "f", "gofile.yml", "Path to gofile")
	rootCmd.AddCommand(installCmd)
}
