package testgrid

import (
	"fmt"
	"time"
)

type Summary struct {
	URL        string
	Dashboards *DashboardMap
}

type DashboardMap map[string]*Dashboard

func (d DashboardMap) Add(key string, value *Dashboard) {
	d[key] = value
}

func (d DashboardMap) Get(key string) *Dashboard {
	return d[key]
}

const (
	PASSING_STATUS = "PASSING"
	FAILING_STATUS = "FAILING"
	FLAKY_STATUS   = "FLAKY"
)

type Dashboard struct {
	Alert               string        `json:"alert"`
	LastRunTimestamp    int           `json:"last_run_timestamp"`
	LastUpdateTimestamp int           `json:"last_update_timestamp"`
	LatestGreen         string        `json:"latest_green"`
	OverallStatus       string        `json:"overall_status"`
	OverallStatusIcon   string        `json:"overall_status_icon"`
	Status              string        `json:"status"`
	Tests               []interface{} `json:"tests"`
	DashboardName       string        `json:"dashboard_name"`
	BugUrl              string        `json:"bug_url"`
}

type TestGroup struct {
	TestGroupName string `json:"test-group-name"`
	Query         string `json:"query"`
	Status        string `json:"status"`
	PhaseTimer    struct {
		Phases []string  `json:"phases"`
		Delta  []float64 `json:"delta"`
		Total  float64   `json:"total"`
	} `json:"phase-timer"`
	Cached  bool   `json:"cached"`
	Summary string `json:"summary"`
	Bugs    struct {
	} `json:"bugs"`
	Changelists       []string   `json:"changelists"`
	ColumnIds         []string   `json:"column_ids"`
	CustomColumns     [][]string `json:"custom-columns"`
	ColumnHeaderNames []string   `json:"column-header-names"`
	Groups            []string   `json:"groups"`
	Metrics           []string   `json:"metrics"`
	Tests             []Test
	RowIds            []string    `json:"row_ids"`
	Timestamps        []int64     `json:"timestamps"`
	Clusters          interface{} `json:"clusters"`
	TestIdMap         interface{} `json:"test_id_map"`
	IdMap             struct {
	} `json:"idMap"`
	TestMetadata struct {
	} `json:"test-metadata"`
	StaleTestThreshold    int    `json:"stale-test-threshold"`
	NumStaleTests         int    `json:"num-stale-tests"`
	AddTabularNamesOption bool   `json:"add-tabular-names-option"`
	ShowTabularNames      bool   `json:"show-tabular-names"`
	Description           string `json:"description"`
	BugComponent          int    `json:"bug-component"`
	CodeSearchPath        string `json:"code-search-path"`
	OpenTestTemplate      struct {
		Url     string `json:"url"`
		Name    string `json:"name"`
		Options struct {
		} `json:"options"`
	} `json:"open-test-template"`
	FileBugTemplate struct {
		Url     string `json:"url"`
		Name    string `json:"name"`
		Options struct {
			Body  string `json:"body"`
			Title string `json:"title"`
		} `json:"options"`
	} `json:"file-bug-template"`
	AttachBugTemplate struct {
		Url     string `json:"url"`
		Name    string `json:"name"`
		Options struct {
		} `json:"options"`
	} `json:"attach-bug-template"`
	ResultsUrlTemplate struct {
		Url     string `json:"url"`
		Name    string `json:"name"`
		Options struct {
		} `json:"options"`
	} `json:"results-url-template"`
	CodeSearchUrlTemplate struct {
		Url     string `json:"url"`
		Name    string `json:"name"`
		Options struct {
		} `json:"options"`
	} `json:"code-search-url-template"`
	AboutDashboardUrl string `json:"about-dashboard-url"`
	OpenBugTemplate   struct {
		Url     string `json:"url"`
		Name    string `json:"name"`
		Options struct {
		} `json:"options"`
	} `json:"open-bug-template"`
	ContextMenuTemplate struct {
		Url     string `json:"url"`
		Name    string `json:"name"`
		Options struct {
		} `json:"options"`
	} `json:"context-menu-template"`
	ColumnDiffLinkTemplates interface{} `json:"column-diff-link-templates"`
	ResultsText             string      `json:"results-text"`
	LatestGreen             string      `json:"latest-green"`
	TriageEnabled           bool        `json:"triage-enabled"`
	Notifications           interface{} `json:"notifications"`
	OverallStatus           int         `json:"overall-status"`
}

type Test struct {
	Name         string        `json:"name"`
	OriginalName string        `json:"original-name"`
	Alert        interface{}   `json:"alert"`
	LinkedBugs   []interface{} `json:"linked_bugs"`
	Messages     []string      `json:"messages"`
	ShortTexts   []string      `json:"short_texts"`
	Statuses     []Statuses    `json:"statuses"`
	Target       string        `json:"target"`
	UserProperty interface{}   `json:"user_property"`
}

type Statuses struct {
	Count int `json:"count"`
	Value int `json:"value"`
}

func (te *Test) RenderStatuses(timestamps []int64) (string, int, int) {
	var firstFailure, text, failures = -1, "", 0
	for i, s := range te.ShortTexts {
		if s != "" {
			if firstFailure < 0 {
				firstFailure = i
			}
			tm := time.Unix(timestamps[i]/1000, 0)
			text += fmt.Sprintf("\t%s %s %s\n", s, tm, te.Messages[i])
			failures += 1
		}
	}
	return text, failures, firstFailure
}
