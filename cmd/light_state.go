package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/skwair/harmony/optional"
	"github.com/skwair/huectl/pkg/hue"
	"github.com/spf13/cobra"
)

type setLightStateFlags struct {
	On         bool
	Brightness int
	Hue        int
}

const setLightStateExample = `
	# Switch on the light 1 and set its brightness to 75%
	huectl light set 1 --on --bri=75

	# Set the color of the light 3 to blue
	huectl light set 3 --hue=46920`

func newSetLightStateCmd() *cobra.Command {
	var flags setLightStateFlags

	cmd := &cobra.Command{
		Use:     "set ID [flags]",
		Short:   "Set the state of lights",
		Example: setLightStateExample,
		Args:    expectLightID(),
		Run:     func(cmd *cobra.Command, args []string) { must(runSetLightStateCmd(cmd, args, &flags)) },
	}

	cmd.Flags().BoolVar(&flags.On, "on", false, "Sets the on/off state of the light")
	cmd.Flags().IntVar(&flags.Brightness, "bri", 0, "Brightness percentage to set the light to")
	cmd.Flags().IntVar(&flags.Hue, "hue", 0, "Color to set the light to, ranges from 0 to 65535")

	return cmd
}

func runSetLightStateCmd(cmd *cobra.Command, args []string, flags *setLightStateFlags) error {
	if cmd.Flags().NFlag() == 0 {
		return errors.New("no flags provided; nothing to do")
	}

	client, err := setupClient()
	if err != nil {
		return fmt.Errorf("unable to setup Hue client: %w", err)
	}

	for _, arg := range args {
		var req hue.SetLightStateRequest
		if cmd.Flags().Changed("on") {
			req.On = optional.NewBool(flags.On)
		}

		if cmd.Flags().Changed("bri") {
			bri := 254 / 100 * flags.Brightness
			req.Bri = optional.NewInt(bri)
		}

		if cmd.Flags().Changed("hue") {
			req.Hue = optional.NewInt(flags.Hue)
		}

		if err = client.SetLightState(arg, &req); err != nil {
			fmt.Fprintf(os.Stderr, "unable to toggle light %q: %v\n", arg, err)
			continue
		}
	}

	return nil
}
