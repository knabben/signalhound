/* Copyright 2025 Amim Knabben */

package cmd

import (
	"fmt"
	"os"

	"github.com/knabben/signalhound/api/v1alpha1"
	"github.com/knabben/signalhound/internal/testgrid"
	tui2 "github.com/knabben/signalhound/internal/tui"
	"github.com/spf13/cobra"
)

// abstractCmd represents the abstract command
var abstractCmd = &cobra.Command{
	Use:   "abstract",
	Short: "Summarize the board status and present the flake or failing ones",
	RunE:  RunAbstract,
}

var (
	tg                   = testgrid.NewTestGrid(testgrid.URL)
	minFailure, minFlake int
	token                string
)

func init() {
	rootCmd.AddCommand(abstractCmd)

	abstractCmd.PersistentFlags().IntVarP(&minFailure, "min-failure", "f", 2, "minimum threshold for test failures")
	abstractCmd.PersistentFlags().IntVarP(&minFlake, "min-flake", "m", 3, "minimum threshold for test flakeness")
	token = os.Getenv("GITHUB_TOKEN")
}

// RunAbstract starts the main command to scrape TestGrid.
func RunAbstract(cmd *cobra.Command, args []string) error {
	var allTabs []*tui2.DashboardTab
	fmt.Println("Scrapping the testgrid dashboard, wait...")

	// render each board summary
	for _, dashboard := range []string{"sig-release-master-blocking", "sig-release-master-informing"} {
		summaries, err := tg.FetchSummary(dashboard, v1alpha1.ERROR_STATUSES)
		if err != nil {
			return err
		}
		// renders the final board summary with tests
		fromSummary := tui2.RenderFromSummary(tg, summaries, minFailure, minFlake)
		allTabs = append(allTabs, fromSummary...)
	}

	return tui2.RenderVisual(allTabs, token)
}
