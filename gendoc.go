// +build ignore
package main

import (
	"log"

	"github.com/padfed/padfed-cli/cmd"
	"github.com/spf13/cobra/doc"
)

func main() {
	err := doc.GenMarkdownTree(cmd.Command(), "doc")
	if err != nil {
		log.Fatal(err)
	}
}
