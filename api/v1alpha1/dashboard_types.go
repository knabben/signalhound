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
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	PASSING_STATUS = "PASSING"
	FAILING_STATUS = "FAILING"
	FLAKY_STATUS   = "FLAKY"
)

var ERROR_STATUSES = []string{FAILING_STATUS, FLAKY_STATUS}

// DashboardSpec defines the desired state of Dashboard.
type DashboardSpec struct {
	// DashboardTab is the name of the tab be scrapped from this board
	DashboardTab string `json:"dashboardTab,omitempty"`

	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default=2
	// MinFailures is the minimum number of failures to consider a test group as failing
	MinFailures int `json:"minFailures,omitempty"`

	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default=3
	// MinFlake is the minimum number of flakes to consider a test group as flaky
	MinFlakes int `json:"minFlakes,omitempty"`
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
	TabName             string `json:"tab_name,omitempty"`
	TabURL              string `json:"tab_url,omitempty"`
}

type TestGroup struct {
	TestGroupName      string     `json:"test-group-name"`
	Query              string     `json:"query"`
	Status             string     `json:"status"`
	Changelists        []string   `json:"changelists"`
	ColumnIds          []string   `json:"column_ids"`
	CustomColumns      [][]string `json:"custom-columns"`
	ColumnHeaderNames  []string   `json:"column-header-names"`
	Groups             []string   `json:"groups"`
	Tests              []Test
	RowIds             []string `json:"row_ids"`
	Timestamps         []int64  `json:"timestamps"`
	StaleTestThreshold int      `json:"stale-test-threshold"`
	NumStaleTests      int      `json:"num-stale-tests"`
	Description        string   `json:"description"`
	OverallStatus      int      `json:"overall-status"`
}

type Test struct {
	Name         string     `json:"name"`
	OriginalName string     `json:"original-name"`
	Messages     []string   `json:"messages"`
	ShortTexts   []string   `json:"short_texts"`
	Statuses     []Statuses `json:"statuses"`
	Target       string     `json:"target"`
}

type Statuses struct {
	Count int `json:"count"`
	Value int `json:"value"`
}

// RenderStatuses renders the statuses of a test into a string.
func (te *Test) RenderStatuses(timestamps []int64) (string, int, int) {
	var firstFailureIndex = -1
	var failureCount = 0
	var output strings.Builder

	for i, shortText := range te.ShortTexts {
		if shortText == "" {
			continue
		}

		if firstFailureIndex < 0 {
			firstFailureIndex = i
		}

		formattedStatus := formatTestStatus(shortText, timestamps[i], te.Messages[i])
		output.WriteString(formattedStatus)
		failureCount++
	}

	return output.String(), failureCount, firstFailureIndex
}

// formatTestStatus creates a formatted string for a single test status.
func formatTestStatus(shortText string, timestamp int64, message string) string {
	timeFormatted := time.Unix(timestamp/1000, 0)
	return fmt.Sprintf("\t%s %s %s\n", shortText, timeFormatted, message)
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
