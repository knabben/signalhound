package tui

import (
	"github.com/knabben/stalker/pkg/testgrid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_NewDashboardTab(t *testing.T) {
	dash := &testgrid.Dashboard{OverallStatus: testgrid.FAILING_STATUS}
	table := &testgrid.TestGroup{
		Timestamps:  []int64{1730221517000},
		Query:       "tab1",
		Changelists: []string{"100000000"},
		Tests: []testgrid.Test{
			{
				Name:       "t1",
				ShortTexts: []string{"F"},
				Messages:   []string{"Build failed outside of test results"},
			},
		},
	}
	minFailure, minFlake := 1, 1
	dashboardTab := NewDashboardTab("url", "tab1", dash, table, minFailure, minFlake)
	assert.Equal(t, "https://testgrid.k8s.io/#tab1&exclude-non-failed-tests=", dashboardTab.BoardURL)
	assert.Equal(t, ":large_red_square:", dashboardTab.Icon)
	assert.Len(t, dashboardTab.Tests, 1)
}
