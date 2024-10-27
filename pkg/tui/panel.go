package tui

import (
	"fmt"
	"github.com/charmbracelet/lipgloss"
	"github.com/knabben/stalker/pkg/testgrid"
	"strings"
	"time"
)

var (
	summaryRegex = `(?<TABS>\d+ of \d+) (?<PERCENT>\(\d+\.\d+%\)) \w.* \((\d+ of \d+) or (\w.*) cells\)`
	keyStyle     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
	bold         = lipgloss.NewStyle().Bold(true).
		Foreground(lipgloss.Color("#ff8787")).
		Width(200).TabWidth(4)
	style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#300a57")).
		Background(lipgloss.Color("#fbf7ff")).
		Width(250).TabWidth(2)
)

func Render(tab string, dashboard *testgrid.Dashboard, url string) error {
	// render tests content
	if err := RenderTests(tab, dashboard, url); err != nil {
		return err
	}
	return nil
}

func RenderTests(tab string, dashboard *testgrid.Dashboard, url string) error {
	var (
		threshold   = 3
		flakeStatus = make(map[string]string)
		isFlake     bool
	)
	url = strings.ReplaceAll(strings.ReplaceAll(url, "/summary", "#"+tab+"&exclude-non-failed-tests="), " ", "%20")

	tg := testgrid.NewTestGrid()
	result, err := tg.FetchTable(dashboard.DashboardName, tab)
	if err != nil {
		return err
	}

	for _, test := range result.Tests {
		status, num := RenderStatuses(&test, result.Timestamps)
		if num >= threshold {
			item := fmt.Sprintf("\t*** %s\n", test.Name)
			flakeStatus[item] = status
			isFlake = true
		}
	}

	if isFlake {
		icon := "🟪"
		if dashboard.OverallStatus == testgrid.FAILING_STATUS {
			icon = "🟥"
		}
		fmt.Println(style.Render(fmt.Sprintf(" [%s] * [%s] %s \t %s", icon, dashboard.DashboardName, tab, url)))
		fmt.Println(style.Render(fmt.Sprintf("%s", style.Render(dashboard.Status))))
		fmt.Println("")
		for item, status := range flakeStatus {
			fmt.Println(bold.Render(item))
			fmt.Println(status)
		}
	}

	return nil
}

func RenderStatuses(test *testgrid.Test, timestamps []int64) (text string, num int) {
	for i, t := range test.ShortTexts {
		if t != "" {
			tm := time.Unix(timestamps[i]/1000, 0)
			msg := test.Messages[i]
			text += fmt.Sprintf("\t%s %s %s\n", keyStyle.Render(t), tm, msg)
			num += 1
		}
	}
	return
}
