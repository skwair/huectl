package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newToggleLightCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "toggle",
		Short: "Toggle lights",
		Args:  expectLightID(),
		Run:   func(_ *cobra.Command, args []string) { must(runLightToggleCmd(args)) },
	}
}

func runLightToggleCmd(args []string) error {
	client, err := setupClient()
	if err != nil {
		return fmt.Errorf("unable to setup Hue client: %w", err)
	}

	for _, arg := range args {
		if err = client.ToggleLight(arg); err != nil {
			fmt.Fprintf(os.Stderr, "unable to toggle light %q: %v\n", arg, err)
			continue
		}
	}

	return nil
}
