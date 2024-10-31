package testgrid

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_FetchSummary(t *testing.T) {
	var dashboard = "dashboard-test"
	response := fmt.Sprintf(`
	{
	  "%v": {
		"overall_status": "PASSING",
		"dashboard_name": "sig-release-master-blocking"
	  }
	}
	`, dashboard)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	}))
	defer server.Close()
	tg := NewTestGrid(server.URL)
	summary, err := tg.FetchSummary("dashboard-test")
	assert.NoError(t, err)
	assert.Contains(t, summary.URL, server.URL)
	assert.Len(t, *summary.Dashboards, 1)
	for name, dash := range *summary.Dashboards {
		assert.Equal(t, name, dashboard)
		assert.Equal(t, dash.OverallStatus, PASSING_STATUS)
	}
}

func Test_FetchTable(*testing.T) {}
