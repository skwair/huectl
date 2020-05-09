package cmd

import (
	"fmt"
	"math"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func newLightsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "lights",
		Aliases: []string{"light", "l"},
		Short:   "Manage Hue light bulbs",
		Args:    cobra.NoArgs,
		// If called with no sub-command, list lights instead of printing help.
		Run: func(*cobra.Command, []string) { must(runListLightsCmd()) },
	}
}

func newListLightsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "List available lights",
		Args:    cobra.NoArgs,
		Run:     func(*cobra.Command, []string) { must(runListLightsCmd()) },
	}
}

func runListLightsCmd() error {
	client, err := setupClient()
	if err != nil {
		return fmt.Errorf("unable to setup Hue client: %w", err)
	}

	lights, err := client.Lights()
	if err != nil {
		return fmt.Errorf("unable to list lights: %w", err)
	}

	tw := tabwriter.NewWriter(os.Stdout, 0, 4, 4, ' ', 0)
	defer tw.Flush()

	fmt.Fprintln(tw, "ID\tNAME\tON\tREACHABLE\tBRIGHTNESS (%)\tHUE")

	for _, light := range lights {
		bri := math.Round(float64(light.State.Bri) / 254 * 100)

		fmt.Fprintf(tw, "%s\t%s\t%t\t%t\t%d\t%d\n", light.ID, light.Name, light.State.On, light.State.Reachable, int(bri), light.State.Hue)
	}

	return nil
}
