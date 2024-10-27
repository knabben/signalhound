package tui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/knabben/stalker/pkg/testgrid"
	"strings"
)

var (
	summaryRegex = `(?<TABS>\d+ of \d+) (?<PERCENT>\(\d+\.\d+%\)) \w.* \((\d+ of \d+) or (\w.*) cells\)`
	testRegex    = `Kubernetes e2e suite.\[It\] \[(\w.*)\] (?<TEST>\w.*)`

	keyStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FF0000"))
	bold     = lipgloss.NewStyle().Bold(true).
			Foreground(lipgloss.Color("#ff8787")).
			Width(200).TabWidth(4)
	style = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#300a57")).
		Background(lipgloss.Color("#fbf7ff")).
		Width(250).TabWidth(2).Padding(2)
)

type DashboardIssue struct {
	URL       string
	Tab       string
	Dashboard *testgrid.Dashboard
	Table     *testgrid.TestGroup
}

func NewDashboardIssue(URL string, tab string, dashboard *testgrid.Dashboard, table *testgrid.TestGroup) *DashboardIssue {
	return &DashboardIssue{URL: URL, Tab: tab, Dashboard: dashboard, Table: table}
}

func (d *DashboardIssue) renderURL() string {
	return strings.ReplaceAll(strings.ReplaceAll(d.URL, "/summary", "#"+d.Tab+"&exclude-non-failed-tests="), " ", "%20")
}
