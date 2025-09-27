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
	for _, summary := range summaries {
		table, err := tg.FetchTable(summary, minFailure, minFlake)
		if err != nil {
			fmt.Println(fmt.Errorf("error fetching table : %s", err))
			continue
		}
		dashboardTab := NewDashboardTab(&summary, table)
		if len(dashboardTab.Tests) > 0 {
			dashboardTabs = append(dashboardTabs, dashboardTab)
		}
	}
	return
}

func NewDashboardTab(dashboard *v1alpha1.DashboardSummary, table *v1alpha1.TestGroup) *DashboardTab {
	aggregation := fmt.Sprintf("%s#%s", dashboard.DashboardName, dashboard.TabName)
	icon := ":large_purple_square:"
	if dashboard.OverallStatus == v1alpha1.FAILING_STATUS {
		icon = ":large_red_square:"
	}
	return &DashboardTab{
		URL:       dashboard.TabURL,
		Tab:       dashboard.TabName,
		Dashboard: dashboard,
		BoardURL:  testgrid.CleanHTMLCharacters(fmt.Sprintf("https://testgrid.k8s.io/%s&exclude-non-failed-tests=", aggregation)),
		BoardHash: aggregation,
		State:     dashboard.OverallStatus,
		Tests:     renderTable(table),
		Icon:      icon,
	}
}

func renderTable(table *v1alpha1.TestGroup) (tests []*TabTest) {
	for _, test := range table.Tests {
		testName := test.Name
		if strings.Contains(test.Name, e2eSuitePrefix) {
			testName = prow.GetRegexParameter(testRegex, testName)["TEST"]
		}
		errMessage, _, firstFailure := test.RenderStatuses(table.Timestamps)
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
	return tests
}
