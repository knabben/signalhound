/* Copyright © 2024 Amim Knabben */
package cmd

import (
	"fmt"
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

		issue := tui.RenderFromSummary(summary, brokenStatus)
		fmt.Println(issue)
	}
	return nil
}
