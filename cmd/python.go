/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"log"
	"path"

	"github.com/cinuor/vsp/util"
	"github.com/spf13/cobra"
	"html/template"
)

var pytmpl = `{
  "configurations": {
    "{{ .Name }}: Launch": {
      "adapter": "debugpy",
      "configuration": {
        "name": "{{ .Name }}: Launch",
        "type": "python",
        "request": "launch",
        "cwd": "${workspaceRoot}",
        "stopOnEntry": true,
        "console": "externalTerminal",
        "program": "{{ .Program }}"
      }
    }
  }
}`

type PythonConfig struct {
	Name    string
	Program string
}

// pythonCmd represents the python command
var pythonCmd = &cobra.Command{
	Use:     "python",
	Short:   "generate .vimspector.json for python",
	Long:    `generate .vimspector.json for python.you should specific THE RELATIVE PATH OF MAIN FILE.`,
	Example: "vsp python [-n NAME] -p PATH",
	Aliases: []string{"py", "Py", "Python"},
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatalf("error occurred: %s", err.Error())
		}
		program, err := cmd.Flags().GetString("program")
		if err != nil {
			log.Fatalf("error occurred: %s", err.Error())
		}
		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			log.Fatalf("error occurred: %s", err.Error())
		}

		p := PythonConfig{
			Name:    name,
			Program: path.Join("${workspaceRoot}", program),
		}

		var data bytes.Buffer
		tmpl := template.New(name)
		tmpl.Parse(pytmpl)
		err = tmpl.Execute(&data, p)
		if err != nil {
			log.Fatalf("Applying Value to Template Failed: %s", err.Error())
		}

		if err := util.GenerateFile(data.Bytes(), dryRun); err != nil {
			log.Fatalf("Generate .vimspector.json Failed: %s", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(pythonCmd)
	pythonCmd.Flags().StringP("name", "n", "python", "the name of the configuration, default: python")
	pythonCmd.Flags().StringP("program", "p", "", "specific the MAIN FILE PATH relatived to the .vimspector.json")
	pythonCmd.MarkFlagRequired("program")
}
