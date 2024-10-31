/* Copyright © 2024 Amim Knabben */
package cmd

import (
	"fmt"
	"github.com/knabben/stalker/pkg/tui"
	"github.com/spf13/cobra"
)

// abstractCmd represents the abstract command
var abstractCmd = &cobra.Command{
	Use:   "abstract",
	Short: "Summarize the board status and present the flake or failing ones",
	RunE:  RunAbstract,
}

func init() {
	rootCmd.AddCommand(abstractCmd)
}

func RunAbstract(cmd *cobra.Command, args []string) error {
	var allTabs []*tui.DashboardTab
	fmt.Println("Scrapping the testgrid dashboard...")
	for _, dashboard := range testBoards {
		//render each board summary
		summary, err := tg.FetchSummary(dashboard)
		if err != nil {
			return err
		}
		//renders the final board summary with tests
		allTabs = append(allTabs, tui.RenderFromSummary(summary, brokenStatus)...)
	}
	return tui.RenderVisual(allTabs)
}
