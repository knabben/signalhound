package tui

import (
	"fmt"
	"strings"

	"github.com/knabben/signalhound/api/v1alpha1"
	"github.com/knabben/signalhound/internal/prow"
	"github.com/knabben/signalhound/internal/testgrid"
)

var (
	e2eSuitePrefix = `Kubernetes e2e suite.`
	testRegex      = e2eSuitePrefix + `\[It\] \[(\w.*)\] (?<TEST>\w.*)`
)

type DashboardTab struct {
	URL       string
	Tab       string
	Dashboard *v1alpha1.DashboardSummary
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

func RenderFromSummary(tg *testgrid.TestGrid, summaries []v1alpha1.DashboardSummary, minFailure, minFlake int) (dashboardTabs []*DashboardTab) {
	for _, dashboard := range summaries {
		table, err := tg.FetchTable(dashboard.DashboardName, dashboard.TabName)
		if err != nil {
			fmt.Println(fmt.Errorf("error fetching table : %s", err))
			continue
		}
		dashboardTab := NewDashboardTab(dashboard.TabURL, dashboard.TabName, &dashboard, table, minFailure, minFlake)
		if len(dashboardTab.Tests) > 0 {
			dashboardTabs = append(dashboardTabs, dashboardTab)
		}
	}
	return
}

func NewDashboardTab(URL string, tab string, dashboard *v1alpha1.DashboardSummary, table *v1alpha1.TestGroup, minFailure, minFlake int) *DashboardTab {
	dash := DashboardTab{URL: URL, Tab: tab, Dashboard: dashboard}
	aggregation := fmt.Sprintf("%s#%s", dashboard.DashboardName, tab)
	dash.BoardURL = testgrid.CleanHTMLCharacters(fmt.Sprintf("https://testgrid.k8s.io/%s&exclude-non-failed-tests=", aggregation))
	dash.BoardHash = aggregation
	dash.State = dashboard.OverallStatus
	dash.Icon = ":large_purple_square:"
	if dash.State == v1alpha1.FAILING_STATUS {
		dash.Icon = ":large_red_square:"
	}
	dash.Tests = renderTable(table, dash.State, minFailure, minFlake)
	return &dash
}

func renderTable(table *v1alpha1.TestGroup, state string, minFailure, minFlake int) (tests []*TabTest) {
	for _, test := range table.Tests {
		testName := test.Name
		if strings.Contains(test.Name, e2eSuitePrefix) {
			testName = prow.GetRegexParameter(testRegex, testName)["TEST"]
		}
		errMessage, failures, firstFailure := test.RenderStatuses(table.Timestamps)
		if (failures >= minFailure && state == v1alpha1.FAILING_STATUS) || (failures >= minFlake && state == v1alpha1.FLAKY_STATUS) {
			tabTest := TabTest{
				Name:            testName,
				LatestTimestamp: table.Timestamps[0],
				FirstTimestamp:  table.Timestamps[len(table.Timestamps)-1],
				ProwURL:         testgrid.CleanHTMLCharacters(fmt.Sprintf("https://prow.k8s.io/view/gs/%s/%s", table.Query, table.Changelists[firstFailure])),
				TriageURL:       testgrid.CleanHTMLCharacters(fmt.Sprintf("https://storage.googleapis.com/k8s-triage/index.html?test=%s", testgrid.CleanHTMLCharacters(testName))),
				ErrMessage:      errMessage,
			}
			tests = append(tests, &tabTest)
		}
	}
	return
}
