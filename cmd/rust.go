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

	"html/template"

	"github.com/cinuor/vsp/util"
	"github.com/spf13/cobra"
)

var rusttmpl = `{
  "configurations": {
    "{{ .Name }}": {
      "adapter": "CodeLLDB",
      "configuration": {
        "request": "launch",
        "program": "{{ .Program }}"
      }
    }
  }
}`

type RustConfig struct {
	Name    string
	Program string
}

// rustCmd represents the rust command
var rustCmd = &cobra.Command{
	Use:     "rust",
	Short:   "generate .vimspector.json for rust",
	Long:    `generate .vimspector.json for rust. you should specific the relatived filepath of target debug binaryfile`,
	Example: "vsp rust [-n NAME] -p PATH",
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

		r := RustConfig{
			Name:    name,
			Program: path.Join("${workspaceRoot}", program),
		}

		var data bytes.Buffer
		tmpl := template.New(name)
		tmpl.Parse(rusttmpl)
		if err := tmpl.Execute(&data, r); err != nil {
			log.Fatalf("Applying Value to Template Failed: %s", err.Error())
		}

		if err := util.GenerateFile(data.Bytes(), dryRun); err != nil {
			log.Fatalf("Generate .vimspector.json Failed: %s", err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(rustCmd)

	rustCmd.Flags().StringP("name", "n", "rust", "the name of the configuration, default: rust")
	rustCmd.Flags().StringP("program", "p", "", "specific the debug target binary relatived to the .vimspector.json")
	rustCmd.MarkFlagRequired("program")

}
