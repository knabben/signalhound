package github

import (
	"context"
	"fmt"
	gh "github.com/google/go-github/v65/github"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"
)

const (
	baseURLPath  = "/api-v3"
	workflowFile = "build.yaml"
)

func TestGithubListRepositories(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	fakeRepos(t, mux)

	// fetch repositories
	ctx := context.Background()
	github := NewGithub(ctx, client, workflowFile, "TNZ", "vmware-master")
	repos, err := github.GetRepositories("tkgm", 100)
	// validate repository returns
	assert.Nil(t, err)
	assert.Equal(t, len(repos), 2)
}

func TestGithubGetWorkflowRuns(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	fakeRepos(t, mux)
	fakeWorkflowRuns(t, mux)

	// fetch repositories
	ctx := context.Background()
	github := NewGithub(ctx, client, workflowFile, "TNZ", "vmware-master")
	repos, err := github.GetRepositories("tkgm", 100)
	assert.Nil(t, err)
	assert.Equal(t, len(repos), 2)

	// new workflow from a repository
	repo := NewRepository(github.(*Github), repos[0])
	workflows, err := repo.GetWorkflowRuns(2)
	// validate workflows returns
	assert.Nil(t, err)
	assert.Equal(t, len(workflows), 2)
}

func TestTriggerNewRun(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()
	fakeRepos(t, mux)
	fakeTriggerRun(t, mux)

	// fetch repositories
	ctx := context.Background()
	github := NewGithub(ctx, client, workflowFile, "TNZ", "vmware-master")
	repos, err := github.GetRepositories("tkgm", 100)
	assert.Nil(t, err)

	// new workflow from a repository
	repo := NewRepository(github.(*Github), repos[0])
	assert.Nil(t, repo.TriggerNewRun())
}

func testMethod(t *testing.T, r *http.Request, want string) {
	t.Helper()
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func fakeRepos(t *testing.T, mux *http.ServeMux) {
	mux.HandleFunc("/user/repos", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `[{"id":1, "name": "tkgm_cayman_1"},{"id":2, "name": "tkgm_cayman_2"}]`)
	})
}

func fakeTriggerRun(t *testing.T, mux *http.ServeMux) {
	mux.HandleFunc("/repos/TNZ/tkgm_cayman_1/actions/workflows/build.yaml/dispatches", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
	})
}

func fakeWorkflowRuns(t *testing.T, mux *http.ServeMux) {
	mux.HandleFunc("/repos/TNZ/tkgm_cayman_1/actions/workflows/build.yaml/runs", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `{"total_count":2,"workflow_runs":[{"id":399444496,"run_number":296,"created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"},{"id":399444497,"run_number":296,"created_at":"2019-01-02T15:04:05Z","updated_at":"2020-01-02T15:04:05Z"}]}`)
	})
}

// setup sets up a test HTTP server along with a github.Client that is
// configured to talk to that test server. Tests should register handlers on
// mux which provide mock responses for the API method being tested.
func setup() (client *gh.Client, mux *http.ServeMux, serverURL string, teardown func()) {
	mux = http.NewServeMux()
	apiHandler := http.NewServeMux()
	apiHandler.Handle(baseURLPath+"/", http.StripPrefix(baseURLPath, mux))
	apiHandler.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintln(os.Stderr, "FAIL: Client.BaseURL path prefix is not preserved in the request URL:")
		http.Error(w, "Client.BaseURL path prefix is not preserved in the request URL.", http.StatusInternalServerError)
	})

	server := httptest.NewServer(apiHandler)
	client = gh.NewClient(nil)
	fakeURL, _ := url.Parse(server.URL + baseURLPath + "/")
	client.BaseURL = fakeURL
	client.UploadURL = fakeURL
	return client, mux, server.URL, server.Close
}
