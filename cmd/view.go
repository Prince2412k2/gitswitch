package cmd

import (
	"fmt"

	"gitswitch/internal/config"
	"gitswitch/internal/ui"
)

func runView() error {
	profiles, err := config.LoadAll()
	if err != nil {
		return fmt.Errorf("loading profiles: %w", err)
	}

	fmt.Println()
	ui.PrintProfiles(profiles)

	if len(profiles) > 0 {
		dir, _ := config.Dir()
		fmt.Printf("  stored in %s\n\n", dir)
	}

	return nil
}
