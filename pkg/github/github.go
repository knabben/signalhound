package github

import (
	"context"
	gh "github.com/google/go-github/v65/github"
	"strings"
)

const OWNER = "kubernetes"

// Github is a Github connections and general metadata holder
type Github struct {
	// ctx is shared context
	ctx context.Context
	// client is the official github client
	client *gh.Client
	// owner is a global repository owner
	owner string
	// workflowFile is the global workflow file used to extrac and trigger runs
	workflowFile string
	// branch is the global reference used to trigget new runs
	branch string
}

// Repository represents a repo abstractions and runs for the workflow
type Repository struct {
	github *Github
	// repo is the github repository object
	repo *gh.Repository
	// runs holds the latest scraped runs for a workflow
	runs []*gh.WorkflowRun
}

type GitHubInterface interface {
	GetRepositories(filter string, perPage int) ([]*gh.Repository, error)
}

type RepositoryInterface interface {
	getWorkflow() (*gh.Workflow, error)
	TriggerNewRun() error
	GetWorkflowRuns(perPage int) ([]*gh.WorkflowRun, error)
	GetWorkflowLogs(runID string) string
}

// NewGithub returns a new Github object with metadata set
func NewGithub(ctx context.Context, client *gh.Client, workflowFile, owner, branch string) GitHubInterface {
	return &Github{ctx: ctx, client: client, workflowFile: workflowFile, owner: owner, branch: branch}
}

// NewRepository returns a new internal Repository abstraction
func NewRepository(github *Github, repo *gh.Repository) RepositoryInterface {
	return &Repository{github: github, repo: repo}
}

func (r *Repository) TriggerNewRun() error {
	// todo(knabben) - Inputs must be provide by specialized workflows, check from template
	event := gh.CreateWorkflowDispatchEventRequest{Ref: r.github.branch}
	_, err := r.github.client.Actions.CreateWorkflowDispatchEventByFileName(r.github.ctx, r.github.owner, *r.repo.Name, r.github.workflowFile, event)
	return err
}

// GetWorkflowRuns returns the list of workflows for a specific repository
func (r *Repository) GetWorkflowRuns(perPage int) ([]*gh.WorkflowRun, error) {
	opts := &gh.ListWorkflowRunsOptions{ListOptions: gh.ListOptions{PerPage: perPage}}
	runs, _, err := r.github.client.Actions.ListWorkflowRunsByFileName(r.github.ctx, r.github.owner, *r.repo.Name, r.github.workflowFile, opts)
	if err != nil {
		return nil, err
	}
	r.runs = runs.WorkflowRuns
	return r.runs, nil
}

func (r *Repository) getWorkflow() (*gh.Workflow, error) {
	workflow, _, err := r.github.client.Actions.GetWorkflowByFileName(r.github.ctx, r.github.owner, r.repo.GetName(), r.github.workflowFile)
	return workflow, err
}

// GetWorkflowLogs extract Logs from the latest job and print
func (r *Repository) GetWorkflowLogs(runID string) string {
	var parsedURL string
	// todo(knabben) - missing implementation
	/*
		jobs, _, _ := client.Actions.ListWorkflowJobs(ctx, OWNER, repo.GetName(), runID, &gh.ListWorkflowJobsOptions{})
		for _, j := range jobs.Jobs {
			jobID := j.ID
			parsedURL, _, _ := client.Actions.GetWorkflowJobLogs(ctx, OWNER, repo.GetName(), *jobID, 1)
		}
	*/
	return parsedURL
}

// GetRepositories returns the list of filtered repositories by the filter arguments
func (g *Github) GetRepositories(filter string, perPage int) (filteredRepos []*gh.Repository, err error) {
	opts := &gh.RepositoryListByAuthenticatedUserOptions{ListOptions: gh.ListOptions{PerPage: perPage}}
	for {
		var (
			repositories []*gh.Repository
			resp         *gh.Response
		)
		repositories, resp, err = g.client.Repositories.ListByAuthenticatedUser(g.ctx, opts)
		if err != nil {
			return nil, err
		}

		for _, repo := range repositories {
			if strings.Contains(*repo.Name, filter) {
				filteredRepos = append(filteredRepos, repo)
			}
		}
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	return filteredRepos, nil
}
