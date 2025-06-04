package main

import (
	"fmt"
	"os"

	"github.com/dolastack/port-user/cmd"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "gen-docs" {
		cmd.GenDocs(cmd.GetRootCmd())
		fmt.Println("✅ Man pages generated in ./docs/")
		return
	}

	cmd.Execute()
}
