package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Prints the version of huectl",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version:\t%s\n", version)
			fmt.Printf("Go Version:\t%s\n", runtime.Version())
			fmt.Printf("Built on:\t%s\n", date)
			fmt.Printf("Commit:\t\t%s\n", commit)
		},
	}
}
