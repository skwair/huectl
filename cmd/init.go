package cmd

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/skwair/huectl/pkg/config"
	"github.com/skwair/huectl/pkg/hue"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "Initializes huectl, connecting to a local Hue bridge and creating a new user",
		Args:  cobra.NoArgs,
		Run:   func(*cobra.Command, []string) { must(runInitCmd()) },
	}
}

func runInitCmd() error {
	cfgPath, err := config.AbsolutePath()
	if err != nil {
		return err
	}

	if _, err = os.Stat(cfgPath); err == nil {
		return fmt.Errorf("huectl already initialized; configuration found at %q", cfgPath)
	}

	fmt.Println("Searching for a Hue bridge on your local network...")

	httpClient := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	bridges, err := hue.DiscoverBridges(httpClient)
	if err != nil {
		return fmt.Errorf("unable to discover Hue bridges: %w", err)
	}

	if len(bridges) == 0 {
		return errors.New("no Hue bridge found on your local network")
	}

	// TODO: if multiple bridges are found, ask the user which bridge the CLI should connect to.
	selectedBridge := bridges[0]

	fmt.Printf("Found Hue bridge %q at: %s\n", selectedBridge.Name, selectedBridge.IPAddr)
	fmt.Println("Registering new user, please press the button on the bridge then press `Enter`")
	fmt.Scanln()

	deviceType := "huectl"
	hn, err := os.Hostname()
	if err == nil {
		deviceType += "#" + hn
	}

	clientID, err := hue.RegisterUser(httpClient, selectedBridge.IPAddr, deviceType)
	if err != nil {
		return fmt.Errorf("unable to register new user: %w", err)
	}

	cfg := &config.Config{
		BridgeID:        selectedBridge.ID,
		BridgeURL:       fmt.Sprintf("https://%s", bridges[0].IPAddr),
		ClientID:        clientID,
		CertFingerprint: selectedBridge.CertFingerprint,
	}

	fmt.Printf("Saving configuration to %q\n", cfgPath)

	if err = saveConfig(cfgPath, cfg); err != nil {
		return fmt.Errorf("unable to save configuration: %w", err)
	}

	return nil
}

func saveConfig(path string, cfg *config.Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("unable to create configuration directory: %w", err)
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to create configuration file: %w", err)
	}

	if err = yaml.NewEncoder(f).Encode(cfg); err != nil {
		return fmt.Errorf("unable to serialize configuration: %w", err)
	}

	return nil
}
