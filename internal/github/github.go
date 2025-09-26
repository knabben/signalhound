package github

import (
	"context"
	"errors"
	"strings"

	g3 "github.com/google/go-github/v65/github"
	g4 "github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

// projectID is the Kubernetes CI Signal board - CI Signal (SIG Release / Release Team)
const projectID = "PVT_kwDOAM_34M4AAThW"

// GitHub is a GitHub connection management and general metadata holder
type GitHub struct {
	// ctx is shared context
	ctx context.Context

	// ClientV3 (deprecated) is the official GitHub API v3 client
	ClientV3 *g3.Client
	// owner is a global repository owner
	owner string
	// workflowFile is the global workflow file used to extrac and trigger runs
	workflowFile string
	// branch is the global reference used to trigget new runs
	branch string

	// ClientV4 is the official GitHub API v4 client
	ClientV4 *g4.Client
}

// Repository represents a repo abstractions and runs for the workflow
type Repository struct {
	github *GitHub
	// repo is the GitHub repository object
	repo *g3.Repository
	// runs holds the latest scraped runs for a workflow
	runs []*g3.WorkflowRun
}

type Interface interface {
	CreateDraftIssue(title, body string) error
	GetRepositories(filter string, perPage int) ([]*g3.Repository, error)
}

type RepositoryInterface interface {
	getWorkflow() (*g3.Workflow, error)
	TriggerNewRun() error
	GetWorkflowRuns(perPage int) ([]*g3.WorkflowRun, error)
}

// NewGithub returns a new GitHub object with metadata set
func NewGithub(ctx context.Context, token string) Interface {
	httpClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token}))
	return &GitHub{ctx: ctx, ClientV4: g4.NewClient(httpClient)}
}

// NewRepository returns a new internal Repository abstraction
func NewRepository(github *GitHub, repo *g3.Repository) RepositoryInterface {
	return &Repository{github: github, repo: repo}
}

func (r *Repository) TriggerNewRun() error {
	if r.github.ClientV3 == nil {
		return errors.New("github client is nil")
	}
	event := g3.CreateWorkflowDispatchEventRequest{Ref: r.github.branch}
	_, err := r.github.ClientV3.Actions.CreateWorkflowDispatchEventByFileName(
		r.github.ctx, r.github.owner, *r.repo.Name, r.github.workflowFile, event,
	)
	return err
}

// GetWorkflowRuns returns the list of workflows for a specific repository
func (r *Repository) GetWorkflowRuns(perPage int) ([]*g3.WorkflowRun, error) {
	if r.github.ClientV3 == nil {
		return nil, errors.New("github client is nil")
	}

	opts := &g3.ListWorkflowRunsOptions{ListOptions: g3.ListOptions{PerPage: perPage}}
	runs, _, err := r.github.ClientV3.Actions.ListWorkflowRunsByFileName(
		r.github.ctx, r.github.owner, *r.repo.Name, r.github.workflowFile, opts,
	)
	if err != nil {
		return nil, err
	}
	r.runs = runs.WorkflowRuns
	return r.runs, nil
}

func (r *Repository) getWorkflow() (*g3.Workflow, error) {
	if r.github.ClientV3 == nil {
		return nil, errors.New("github client is nil")
	}

	workflow, _, err := r.github.ClientV3.Actions.GetWorkflowByFileName(
		r.github.ctx, r.github.owner, r.repo.GetName(), r.github.workflowFile,
	)
	return workflow, err
}

// CreateDraftIssue creates a new issue draft issue in the board with a
// specific test issue template.
func (g *GitHub) CreateDraftIssue(title, body string) error {
	var mutationDraft struct {
		AddProjectV2DraftIssue struct {
			ProjectItem struct {
				ID g4.ID
			}
		} `graphql:"addProjectV2DraftIssue(input: $input)"`
	}
	bodyInput := g4.String(body)
	inputDraft := g4.AddProjectV2DraftIssueInput{
		ProjectID: projectID,
		Title:     g4.String(title),
		Body:      &bodyInput,
	}
	// mutation for creating a draft issue in the Signal CI board
	if err := g.ClientV4.Mutate(context.Background(), &mutationDraft, inputDraft, nil); err != nil {
		return err
	}

	itemID := mutationDraft.AddProjectV2DraftIssue.ProjectItem.ID

	// iterate in the draft issue items and update each field for
	// latest view. Single select field (like K8s Version) needs
	// to be updated automatically.

	fields := map[string]g4.String{
		"PVTSSF_lADOAM_34M4AAThWzgAKb0I": g4.String("fb5648b7"), // v1.32 - K8s Release - ProjectV2SingleSelectField
		"PVTSSF_lADOAM_34M4AAThWzgAtHfg": g4.String("3edeeefb"), // issue-tracking - view - ProjectV2SingleSelectField
		"PVTSSF_lADOAM_34M4AAThWzgAKbaA": g4.String("179de113"), // drafting - Status - ProjectV2SingleSelectField
	}

	var mutationUpdate struct {
		UpdateProjectV2ItemFieldValue struct {
			ClientMutationID string
		} `graphql:"updateProjectV2ItemFieldValue(input: $input)"`
	}
	for key, field := range fields {
		if err := g.ClientV4.Mutate(context.Background(), &mutationUpdate, g4.UpdateProjectV2ItemFieldValueInput{
			ProjectID: projectID,
			ItemID:    itemID,
			FieldID:   key,
			Value:     g4.ProjectV2FieldValue{SingleSelectOptionID: &field},
		}, nil); err != nil {
			return err
		}
	}
	return nil
}

// GetRepositories returns the list of filtered repositories by the filter arguments
func (g *GitHub) GetRepositories(filter string, perPage int) (filteredRepos []*g3.Repository, err error) {
	if g.ClientV3 == nil {
		return nil, errors.New("github client is nil")
	}

	opts := &g3.RepositoryListByAuthenticatedUserOptions{ListOptions: g3.ListOptions{PerPage: perPage}}
	for {
		var (
			repositories []*g3.Repository
			resp         *g3.Response
		)
		repositories, resp, err = g.ClientV3.Repositories.ListByAuthenticatedUser(g.ctx, opts)
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
