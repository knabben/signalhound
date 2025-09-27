package testgrid

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/knabben/signalhound/api/v1alpha1"
)

var URL = "https://testgrid.k8s.io"

type TestGrid struct {
	URL string
}

func NewTestGrid(url string) *TestGrid {
	return &TestGrid{URL: url}
}

type DashboardMapper map[string]*v1alpha1.DashboardSummary

// FetchSummary retrieves the summary data for a given dashboard from the TestGrid
func (t *TestGrid) FetchSummary(dashboard string, filterStatus []string) (summary []v1alpha1.DashboardSummary, err error) {
	url := fmt.Sprintf("%s/%s/summary", t.URL, CleanHTMLCharacters(dashboard))

	// request summary data from TestGrid
	var response *http.Response
	if response, err = http.Get(url); err != nil {
		return nil, fmt.Errorf("error fetching testgrid dashboard summary endpoint: %v", err)
	}

	var data []byte
	if data, err = io.ReadAll(response.Body); err != nil {
		return nil, fmt.Errorf("error parsing body response: %v", err)
	}

	// unmarshal summary data into a struct
	var dashboardList DashboardMapper
	if err = json.Unmarshal(data, &dashboardList); err != nil {
		return nil, fmt.Errorf("error unmarshaling body response: %v", err)
	}

	// iterate and save the final value filtering by status
	for tabName, dashboardSummary := range dashboardList {
		if hasStatus(dashboardSummary.OverallStatus, filterStatus) {
			dashboardSummary.TabName = tabName
			dashboardSummary.TabURL = url
			summary = append(summary, *dashboardSummary)
		}
	}
	return summary, nil
}

// FetchTable returns the test group related to the tab of a dashboard
func (t *TestGrid) FetchTable(dashboard v1alpha1.DashboardSummary, minFailure, minFlake int) (*v1alpha1.TestGroup, error) {
	url := fmt.Sprintf("%s/%s/table?tab=%s&exclude-non-failed-tests=&dashboard=%s",
		t.URL, dashboard.DashboardName, dashboard.TabName, dashboard.DashboardName)
	response, err := http.Get(CleanHTMLCharacters(url))
	if err != nil {
		return nil, err
	}

	var data []byte
	if data, err = io.ReadAll(response.Body); err != nil {
		return nil, err
	}

	var testGroup = &v1alpha1.TestGroup{}
	if err = json.Unmarshal(data, testGroup); err != nil {
		return nil, err
	}

	var filteredTests []v1alpha1.Test
	for _, test := range testGroup.Tests {
		_, failures, _ := test.RenderStatuses(testGroup.Timestamps)
		if (failures >= minFailure && dashboard.OverallStatus == v1alpha1.FAILING_STATUS) || (failures >= minFlake && dashboard.OverallStatus == v1alpha1.FLAKY_STATUS) {
			filteredTests = append(filteredTests, test)
		}
	}
	testGroup.Tests = filteredTests
	return testGroup, nil
}

func hasStatus(boardStatus string, statuses []string) bool {
	for _, status := range statuses {
		if boardStatus == status {
			return true
		}
	}
	return false
}

func CleanHTMLCharacters(str string) string {
	return strings.ReplaceAll(str, " ", "%20")
}
