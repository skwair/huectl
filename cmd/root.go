package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/skwair/huectl/pkg/config"
	"github.com/skwair/huectl/pkg/hue"
	"github.com/spf13/cobra"
)

func Huectl() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "huectl",
		Short: "huectl controls a Philips Hue installation",
	}

	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newInitCmd())

	lightsCmd := newLightsCmd()
	rootCmd.AddCommand(lightsCmd)

	lightsCmd.AddCommand(newListLightsCmd())
	lightsCmd.AddCommand(newSetLightStateCmd())
	lightsCmd.AddCommand(newToggleLightCmd())

	return rootCmd
}

func must(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func setupClient() (*hue.Client, error) {
	cfg, err := config.Read()
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("huectl is not initialized, please run `huectl init` first")
		}

		return nil, err
	}

	client := hue.NewClient(cfg.BridgeURL, cfg.ClientID, hue.WithCertFingerprint(cfg.CertFingerprint))

	return client, nil
}

func expectLightID() cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("at least one light id is required, e.g.: `%s 1`", cmd.CommandPath())
		}

		return nil
	}
}
