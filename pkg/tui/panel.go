package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/knabben/stalker/pkg/testgrid"
	"strings"
	"time"
)

var summaryRegex = `(?<TABS>\d+ of \d+) (?<PERCENT>\(\d+\.\d+%\)) \w.* \((\d+ of \d+) or (\w.*) cells\)`
var keyStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))

var bold = lipgloss.NewStyle().Bold(true).
	Foreground(lipgloss.Color("#FAFAFA")).
	Background(lipgloss.Color("#a8323a")).
	Width(250).TabWidth(2)

func Render(tab string, dashboard *testgrid.Dashboard, url string) error {
	var style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#c9c9c9")).
		Background(lipgloss.Color("#744278")).
		Width(250).TabWidth(2)
	url = strings.ReplaceAll(strings.ReplaceAll(url, " ", "%20"), "/summary", "#"+tab)
	fmt.Println(style.Render(fmt.Sprintf("* [%s] %s \t %s", dashboard.DashboardName, tab, url)))
	fmt.Println(style.Render(fmt.Sprintf("%s", style.Render(dashboard.Status))))
	// render tests content
	if err := RenderTests(tab, dashboard); err != nil {
		return err
	}
	fmt.Println("\n")
	return nil
}

func RenderTests(tab string, dashboard *testgrid.Dashboard) error {
	tg := testgrid.NewTestGrid()
	result, err := tg.FetchTable(dashboard.DashboardName, tab)
	if err != nil {
		return err
	}
	for _, test := range result.Tests {
		item := fmt.Sprintf("%s -- \n", bold.Render(test.Name))
		fmt.Println(item, RenderStatuses(&test, result.Timestamps))
	}
	return nil
}

func RenderStatuses(test *testgrid.Test, timestamps []int64) (text string) {

	for i, t := range test.ShortTexts {
		if t != "" {
			tm := time.Unix(timestamps[i]/1000, 0)
			msg := test.Messages[i]
			text += fmt.Sprintf("%s %s %s\n", keyStyle.Render(t), tm, msg)
		}
	}
	return
}
