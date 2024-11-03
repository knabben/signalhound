/* Copyright © 2024 Amim Knabben */
package cmd

import (
	"fmt"
	"github.com/knabben/stalker/pkg/testgrid"
	"github.com/knabben/stalker/pkg/tui"
	"github.com/spf13/cobra"
	"os"
)

// abstractCmd represents the abstract command
var abstractCmd = &cobra.Command{
	Use:   "abstract",
	Short: "Summarize the board status and present the flake or failing ones",
	RunE:  RunAbstract,
}

var (
	tg                   = testgrid.NewTestGrid("")
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
	var allTabs []*tui.DashboardTab
	fmt.Println("Scrapping the testgrid dashboard, wait...")
	for _, dashboard := range testBoards {
		// render each board summary
		summary, err := tg.FetchSummary(dashboard)
		if err != nil {
			return err
		}
		// renders the final board summary with tests
		fromSummary := tui.RenderFromSummary(tg.(*testgrid.TestGrid), summary, brokenStatus, minFailure, minFlake)
		allTabs = append(allTabs, fromSummary...)
	}
	return tui.RenderVisual(allTabs, token)
}
