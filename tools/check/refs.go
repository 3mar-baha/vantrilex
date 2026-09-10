package main

import (
	"fmt"
	"os"

	"vantrilex/internal/catalog"
)

// checkInstallRefs ensures every MCP entry carries a runnable install command.
func checkInstallRefs() {
	bad := 0
	for _, m := range catalog.MCPs() {
		if m.Install == "" {
			bad++
		}
	}
	if bad > 0 {
		fmt.Printf("install refs: %d entries missing command\n", bad)
		os.Exit(1)
	}
	fmt.Println("install refs: ok")
}
