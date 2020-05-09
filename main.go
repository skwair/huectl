package main

import (
	"fmt"
	"os"

	"github.com/skwair/huectl/cmd"
)

func main() {
	if err := cmd.Huectl().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
