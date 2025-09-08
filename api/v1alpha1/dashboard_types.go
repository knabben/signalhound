/*
Copyright 2025.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DashboardSpec defines the desired state of Dashboard.
type DashboardSpec struct {
	// Name is the dashboard name to be scrapped
	DashboardTab string `json:"dashboardTab,omitempty"`
}

// DashboardStatus defines the observed state of Dashboard.
type DashboardStatus struct {
	LastUpdate       metav1.Time        `json:"lastFetched,omitempty"`
	DashboardSummary []DashboardSummary `json:"summary,omitempty"`
}

type DashboardSummary struct {
	Alert               string `json:"alert,omitempty"`
	LastRunTimestamp    int64  `json:"last_run_timestamp,omitempty"`
	LastUpdateTimestamp int64  `json:"last_update_timestamp,omitempty"`
	LatestGreen         string `json:"latest_green,omitempty"`
	OverallStatus       string `json:"overall_status,omitempty"`
	OverallStatusIcon   string `json:"overall_status_icon,omitempty"`
	Status              string `json:"status,omitempty"`
	DashboardName       string `json:"dashboard_name,omitempty"`
}

type TestGroup struct {
	TestGroupName string `json:"test-group-name"`
	Query         string `json:"query"`
	Status        string `json:"status"`
	//PhaseTimer    struct {
	//	Phases []string  `json:"phases"`
	//	Delta  []float64 `json:"delta"`
	//	Total  float64   `json:"total"`
	//} `json:"phase-timer"`
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
	RowIds            []string `json:"row_ids"`
	Timestamps        []int64  `json:"timestamps"`
	//Clusters          interface{} `json:"clusters"`
	//TestIdMap interface{} `json:"test_id_map"`
	IdMap struct {
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
	//ColumnDiffLinkTemplates interface{} `json:"column-diff-link-templates"`
	ResultsText   string `json:"results-text"`
	LatestGreen   string `json:"latest-green"`
	TriageEnabled bool   `json:"triage-enabled"`
	//Notifications           interface{} `json:"notifications"`
	OverallStatus int `json:"overall-status"`
}

type Test struct {
	Name         string `json:"name"`
	OriginalName string `json:"original-name"`
	//Alert        interface{}   `json:"alert"`
	//LinkedBugs   []interface{} `json:"linked_bugs"`
	Messages   []string   `json:"messages"`
	ShortTexts []string   `json:"short_texts"`
	Statuses   []Statuses `json:"statuses"`
	Target     string     `json:"target"`
	//UserProperty interface{}   `json:"user_property"`
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

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Dashboard is the Schema for the dashboards API.
type Dashboard struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   DashboardSpec   `json:"spec,omitempty"`
	Status DashboardStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// DashboardList contains a list of Dashboard.
type DashboardList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Dashboard `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Dashboard{}, &DashboardList{})
}
