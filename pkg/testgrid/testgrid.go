package testgrid

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var TestGridURL = "https://testgrid.k8s.io"

type TestGrid struct {
	TestGridURL string
}

type TestGridInterface interface {
	FetchSummary(dashboard string) (*Summary, error)
	FetchTable(dashboard, tab string) (*TestGroup, error)
}

func NewTestGrid() TestGridInterface {
	return &TestGrid{TestGridURL: TestGridURL}
}

func (t TestGrid) FetchSummary(dashboard string) (*Summary, error) {
	var (
		dashboardList = &DashboardMap{}
		data          []byte
	)
	url := fmt.Sprintf("%s/%s/summary", t.TestGridURL, strings.ReplaceAll(dashboard, " ", "%20"))
	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if data, err = io.ReadAll(response.Body); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, dashboardList); err != nil {
		return nil, err
	}
	return &Summary{URL: url, Dashboards: dashboardList}, nil
}

func (t TestGrid) FetchTable(dashboard, tab string) (*TestGroup, error) {
	var (
		testGroup = &TestGroup{}
		data      []byte
	)
	url := fmt.Sprintf("%s/%s/table?tab=%s&exclude-non-failed-tests=&dashboard=%s", t.TestGridURL, dashboard, tab, dashboard)
	response, err := http.Get(strings.ReplaceAll(url, " ", "%20"))
	if err != nil {
		return nil, err
	}
	if data, err = io.ReadAll(response.Body); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(data, testGroup); err != nil {
		return nil, err
	}
	return testGroup, nil
}
