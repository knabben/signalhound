package tui

import (
	"github.com/knabben/stalker/pkg/testgrid"
	"strings"
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
