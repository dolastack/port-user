package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func GenDocs(cmd *cobra.Command) {
	// Create docs directory if not exists
	if err := os.MkdirAll("docs", os.ModePerm); err != nil {
		log.Fatalf("Error creating docs directory: %v", err)
	}

	// Generate man pages
	fmt.Println("Generating man pages...")
	err := doc.GenManTree(cmd, nil, "docs/")
	if err != nil {
		log.Fatalf("Error generating man pages: %v", err)
	}
}
