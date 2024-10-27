/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/knabben/stalker/pkg/testgrid"
	"github.com/knabben/stalker/pkg/tui"
	"log"

	"github.com/spf13/cobra"
)

// abstractCmd represents the abstract command
var abstractCmd = &cobra.Command{
	Use:   "abstract",
	Short: "Summarize the board status and present the flake or failing ones",
	RunE: func(cmd *cobra.Command, args []string) error {
		tg := testgrid.NewTestGrid()
		for _, dashboard := range []string{"sig-release-master-blocking", "sig-release-master-informing"} {
			summary, err := tg.FetchSummary(dashboard)
			if err != nil {
				log.Fatal("ERROR ", err)
			}

			var counter = 0
			for tab, dashboard := range *summary.Dashboards {
				if hasStatus(dashboard.OverallStatus, []string{testgrid.FAILING_STATUS, testgrid.FLAKY_STATUS}) {
					if err := tui.Render(tab, dashboard, summary.URL, counter+1); err != nil {
						return err
					}
					counter += 1
				}
			}
		}

		return nil
	},
}

func hasStatus(boardStatus string, statuses []string) bool {
	for _, status := range statuses {
		if boardStatus == status {
			return true
		}
	}
	return false
}

func init() {
	rootCmd.AddCommand(abstractCmd)
}
