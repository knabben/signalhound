package prow

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ExtractJunitIFrame(t *testing.T) {
	content, err := os.ReadFile("./testdata/logs.html")
	assert.NoError(t, err)

	body := bytes.NewReader(content)
	lensURL, err := extractBuildLensURL(body)
	assert.NoError(t, err)

	req := `{"artifacts": ["build-log.txt"],"index": 1,"src": "gs/kubernetes-ci-logs/logs/ci-kubernetes-node-e2e-containerd/1855351637342687232"}`
	assert.Equal(t, fmt.Sprintf("https://prow.k8s.io/spyglass/lens/buildlog/iframe?req=%s", url.QueryEscape(req)), lensURL)
}

func Test_ExtractBuildLog(t *testing.T) {
	content, err := os.ReadFile("./testdata/buildlog.html")
	assert.NoError(t, err)

	body := bytes.NewReader(content)
	data, err := extractBuildLogs(body)
	assert.NoError(t, err)
	assert.Len(t, data.Error, 3464)
	assert.Contains(t, data.Error, "err: exit status 255")
}
