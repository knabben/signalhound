package testgrid

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/knabben/stalker/api/v1alpha1"
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
func (t *TestGrid) FetchSummary(dashboard string) (summary []v1alpha1.DashboardSummary, err error) {
	url := fmt.Sprintf("%s/%s/summary", t.URL, cleanHTMLCharacters(dashboard))

	// request summary data from TestGrid
	var response *http.Response
	if response, err = http.Get(url); err != nil {
		return nil, fmt.Errorf("error fetching testgrid dashboard summary endpoint: %v", err)
	}

	var data []byte
	if data, err = io.ReadAll(response.Body); err != nil {
		return nil, err
	}

	// unmarshal summary data into a struct
	var dashboardList DashboardMapper
	if err = json.Unmarshal(data, &dashboardList); err != nil {
		return nil, err
	}

	// iterate and save the final value
	for dashName, dashboardSummary := range dashboardList {
		dashboardSummary.DashboardName = dashName
		summary = append(summary, *dashboardSummary)
	}
	return summary, nil
}

func cleanHTMLCharacters(str string) string {
	return strings.ReplaceAll(str, " ", "%20")
}
