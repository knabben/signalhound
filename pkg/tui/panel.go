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
	Foreground(lipgloss.Color("#ff8787")).
	Width(200).TabWidth(4)

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
	return nil
}

func RenderTests(tab string, dashboard *testgrid.Dashboard) error {
	tg := testgrid.NewTestGrid()
	result, err := tg.FetchTable(dashboard.DashboardName, tab)
	if err != nil {
		return err
	}
	for _, test := range result.Tests {
		status, num := RenderStatuses(&test, result.Timestamps)
		item := fmt.Sprintf("\t*** %s\n", test.Name)
		if num >= 3 {
			item = fmt.Sprintf("\t🟪 [FLAKE] %s\n", test.Name)
		}
		fmt.Println(bold.Render(item))
		fmt.Println(status)
	}
	return nil
}

func RenderStatuses(test *testgrid.Test, timestamps []int64) (text string, num int) {
	for i, t := range test.ShortTexts {
		if t != "" {
			tm := time.Unix(timestamps[i]/1000, 0)
			msg := test.Messages[i]
			text += fmt.Sprintf("\t\t%s %s %s\n", keyStyle.Render(t), tm, msg)
			num += 1
		}
	}
	return
}
