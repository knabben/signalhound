/* Copyright © 2024 Amim Knabben */
package cmd

import (
	"fmt"
	"github.com/knabben/stalker/pkg/testgrid"
	"github.com/knabben/stalker/pkg/tui"
	"log"

	"github.com/spf13/cobra"
)

// issueCmd represents the issue command
var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Generate a GitHub issue template for failing tests",
	RunE:  RunIssue,
}

func init() {
	rootCmd.AddCommand(issueCmd)
}

func RunIssue(cmd *cobra.Command, args []string) error {
	for _, dashboard := range testBoards {
		summary, err := tg.FetchSummary(dashboard)
		if err != nil {
			log.Fatal("error fetching the summary", err)
		}
		for tab, dashboard := range *summary.Dashboards {
			if hasStatus(dashboard.OverallStatus, []string{testgrid.FAILING_STATUS}) {
				table, err := tg.FetchTable(dashboard.DashboardName, tab)
				if err != nil {
					_ = fmt.Errorf("error fetching table : %s", err)
					continue
				}
				issue := tui.NewDashboardIssue(summary.URL, tab, dashboard, table)
				issue.RenderTemplate()
			}
		}
	}
	return nil
}
