package tui

import (
	"fmt"
	"github.com/knabben/stalker/pkg/testgrid"
	"regexp"
	"strings"
)

var (
	e2eSuitePrefix = `Kubernetes e2e suite.`
	testRegex      = e2eSuitePrefix + `\[It\] \[(\w.*)\] (?<TEST>\w.*)`
)

type DashboardTab struct {
	URL       string
	Tab       string
	Dashboard *testgrid.Dashboard
	BoardURL  string
	BoardHash string
	Icon      string
	State     string
	Status    string

	Tests []*TabTest
}

type TabTest struct {
	Name            string
	FirstTimestamp  int64
	LatestTimestamp int64
	TriageURL       string
	ProwURL         string
	ErrMessage      string
}

var minFailure, minFlake = 2, 3

func RenderFromSummary(tg *testgrid.TestGrid, summary *testgrid.Summary, failures []string) (dashboardTabs []*DashboardTab) {
	for tab, dashboard := range *summary.Dashboards {
		if hasStatus(dashboard.OverallStatus, failures) {
			table, err := tg.FetchTable(dashboard.DashboardName, tab)
			if err != nil {
				_ = fmt.Errorf("error fetching table : %s", err)
				continue
			}
			dashboardTab := NewDashboardTab(summary.URL, tab, dashboard, table)
			if len(dashboardTab.Tests) > 0 {
				dashboardTabs = append(dashboardTabs, dashboardTab)
			}
		}
	}
	return
}

func NewDashboardTab(URL string, tab string, dashboard *testgrid.Dashboard, table *testgrid.TestGroup) *DashboardTab {
	dash := DashboardTab{URL: URL, Tab: tab, Dashboard: dashboard}
	aggregation := fmt.Sprintf("%s#%s", dashboard.DashboardName, tab)
	dash.BoardURL = testgrid.CleanSpaces(fmt.Sprintf("https://testgrid.k8s.io/%s&exclude-non-failed-tests=", aggregation))
	dash.BoardHash = aggregation
	dash.State = dashboard.OverallStatus
	dash.Icon = ":large_purple_square:"
	if dashboard.OverallStatus == testgrid.FAILING_STATUS {
		dash.Icon = ":large_red_square:"
	}
	dash.Tests = renderTable(table, dash.State)
	return &dash
}

func renderTable(table *testgrid.TestGroup, state string) (tests []*TabTest) {
	for _, test := range table.Tests {
		testName := test.Name
		if strings.Contains(test.Name, e2eSuitePrefix) {
			testName = getParameter(testRegex, testName)["TEST"]
		}
		errMessage, failures, firstFailure := test.RenderStatuses(table.Timestamps)
		if (failures >= minFailure && state == testgrid.FAILING_STATUS) || (failures >= minFlake && state == testgrid.FLAKY_STATUS) {
			tabTest := TabTest{
				Name:            testName,
				LatestTimestamp: table.Timestamps[0],
				FirstTimestamp:  table.Timestamps[len(table.Timestamps)-1],
				ProwURL:         testgrid.CleanSpaces(fmt.Sprintf("https://prow.k8s.io/view/gs/%s/%s", table.Query, table.Changelists[firstFailure])),
				TriageURL:       testgrid.CleanSpaces(fmt.Sprintf("https://storage.googleapis.com/k8s-triage/index.html?test=%s", testgrid.CleanSpaces(testName))),
				ErrMessage:      errMessage,
			}
			tests = append(tests, &tabTest)
		}
	}
	return
}

func hasStatus(boardStatus string, statuses []string) bool {
	for _, status := range statuses {
		if boardStatus == status {
			return true
		}
	}
	return false
}

func getParameter(regEx, value string) (paramsMap map[string]string) {
	var r = regexp.MustCompile(regEx)
	match := r.FindStringSubmatch(value)
	paramsMap = make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return paramsMap
}
