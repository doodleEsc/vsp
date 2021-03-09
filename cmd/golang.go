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

var gotmpl = `{
"configurations": {
    "{{ .Name }}": {
      "adapter": "vscode-go",
      "configuration": {
        "request": "launch",
        "program": "{{ .Program }}",
        "mode": "debug",
        "dlvToolPath": "{{ .Delve }}",
        "trace": true,
        "env": { "GO111MODULE": "on" }
      }
    }
  }
}`

type GolangConfig struct {
	Name    string
	Program string
	Delve   string
}

// golangCmd represents the golang command
var golangCmd = &cobra.Command{
	Use:     "golang",
	Short:   "generate .vimspector.json for golang",
	Long:    `generate .vimspector.json for golang. you should specific the absolute path of Delve and the relatived filepath or relatived dirpath of main file`,
	Example: "vsp golang [-n NAME] -dlv DLVPATH -p PATH",
	Aliases: []string{"go", "Go", "Golang"},
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatalf("error occurred: %s", err.Error())
		}
		program, err := cmd.Flags().GetString("program")
		if err != nil {
			log.Fatalf("error occurred: %s", err.Error())
		}
		delve, err := cmd.Flags().GetString("delve")
		if err != nil {
			log.Fatalf("error occurred: %s", err.Error())
		}
		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			log.Fatalf("error occurred: %s", err.Error())
		}

		g := GolangConfig{
			Name:    name,
			Program: path.Join("${workspaceRoot}", program),
			Delve:   delve,
		}

		var data bytes.Buffer
		tmpl := template.New(name)
		tmpl.Parse(gotmpl)
		err = tmpl.Execute(&data, g)
		if err != nil {
			log.Fatalf("Applying Value to Template Failed: %s", err.Error())
		}

		if err := util.GenerateFile(data.Bytes(), dryRun); err != nil {
			log.Fatalf("Generate .vimspector.json Failed: %s", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(golangCmd)

	golangCmd.Flags().StringP("name", "n", "golang", "the name of the configuration, default: golang")
	golangCmd.Flags().StringP("program", "p", "", "specific the MAIN FILE PATH or DIR PATH relatived to the .vimspector.json")
	golangCmd.Flags().StringP("delve", "D", "", "specific the Delve bin file path")
	golangCmd.MarkFlagRequired("program")
	golangCmd.MarkFlagRequired("delve")
}
