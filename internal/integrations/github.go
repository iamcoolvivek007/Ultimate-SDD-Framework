package integrations

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v60/github"
	"golang.org/x/oauth2"
)

// GitHubIntegration handles GitHub repository operations
type GitHubIntegration struct {
	client *github.Client
	ctx    context.Context
	owner  string
	repo   string
	token  string
}

// NewGitHubIntegration creates a new GitHub integration
func NewGitHubIntegration(token, owner, repo string) *GitHubIntegration {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	return &GitHubIntegration{
		client: client,
		ctx:    ctx,
		owner:  owner,
		repo:   repo,
		token:  token,
	}
}

// ValidateToken checks if the GitHub token is valid
func (ghi *GitHubIntegration) ValidateToken() error {
	_, _, err := ghi.client.Users.Get(ghi.ctx, "")
	if err != nil {
		return fmt.Errorf("invalid GitHub token: %w", err)
	}
	return nil
}

// GetPullRequests gets recent pull requests
func (ghi *GitHubIntegration) GetPullRequests(state string, limit int) ([]*github.PullRequest, error) {
	opts := &github.PullRequestListOptions{
		State: state,
		ListOptions: github.ListOptions{
			PerPage: limit,
		},
	}

	prs, _, err := ghi.client.PullRequests.List(ghi.ctx, ghi.owner, ghi.repo, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get pull requests: %w", err)
	}

	return prs, nil
}

// CreatePullRequestReview performs AI-powered code review on a PR
func (ghi *GitHubIntegration) CreatePullRequestReview(prNumber int, reviewBody string) error {
	review := &github.PullRequestReviewRequest{
		Body:  &reviewBody,
		Event: github.String("COMMENT"), // Can be "APPROVE", "REQUEST_CHANGES", or "COMMENT"
	}

	_, _, err := ghi.client.PullRequests.CreateReview(ghi.ctx, ghi.owner, ghi.repo, prNumber, review)
	if err != nil {
		return fmt.Errorf("failed to create PR review: %w", err)
	}

	return nil
}

// GetPullRequestFiles gets the files changed in a PR
func (ghi *GitHubIntegration) GetPullRequestFiles(prNumber int) ([]*github.CommitFile, error) {
	opts := &github.ListOptions{
		PerPage: 100,
	}

	files, _, err := ghi.client.PullRequests.ListFiles(ghi.ctx, ghi.owner, ghi.repo, prNumber, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get PR files: %w", err)
	}

	return files, nil
}

// CreateIssue creates a new GitHub issue
func (ghi *GitHubIntegration) CreateIssue(title, body string, labels []string) (*github.Issue, error) {
	issue := &github.IssueRequest{
		Title:  &title,
		Body:   &body,
		Labels: &labels,
	}

	createdIssue, _, err := ghi.client.Issues.Create(ghi.ctx, ghi.owner, ghi.repo, issue)
	if err != nil {
		return nil, fmt.Errorf("failed to create issue: %w", err)
	}

	return createdIssue, nil
}

// GetRepositoryInfo gets basic repository information
func (ghi *GitHubIntegration) GetRepositoryInfo() (*github.Repository, error) {
	repo, _, err := ghi.client.Repositories.Get(ghi.ctx, ghi.owner, ghi.repo)
	if err != nil {
		return nil, fmt.Errorf("failed to get repository info: %w", err)
	}

	return repo, nil
}

// CreateBranchProtection sets up branch protection rules
func (ghi *GitHubIntegration) CreateBranchProtection(branch string) error {
	strict := true
	enforceAdmins := true

	protection := &github.ProtectionRequest{
		RequiredStatusChecks: &github.RequiredStatusChecks{
			Strict:   strict,
			Contexts: &[]string{"continuous-integration/travis-ci"},
		},
		RequiredPullRequestReviews: &github.PullRequestReviewsEnforcementRequest{
			RequiredApprovingReviewCount: 1,
		},
		EnforceAdmins: enforceAdmins,
		Restrictions:  nil,
	}

	_, _, err := ghi.client.Repositories.UpdateBranchProtection(ghi.ctx, ghi.owner, ghi.repo, branch, protection)
	if err != nil {
		return fmt.Errorf("failed to create branch protection: %w", err)
	}

	return nil
}

// GetContributors gets repository contributors
func (ghi *GitHubIntegration) GetContributors() ([]*github.Contributor, error) {
	opts := &github.ListContributorsOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
	}

	contributors, _, err := ghi.client.Repositories.ListContributors(ghi.ctx, ghi.owner, ghi.repo, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get contributors: %w", err)
	}

	return contributors, nil
}

// CreateWorkflowFile creates a GitHub Actions workflow file
func (ghi *GitHubIntegration) CreateWorkflowFile(workflowName, content string) error {
	workflowsPath := ".github/workflows"
	fileName := fmt.Sprintf("%s.yml", strings.ReplaceAll(strings.ToLower(workflowName), " ", "_"))

	// This would typically create the file locally and commit it
	// For now, we'll just validate the content
	if !strings.Contains(content, "name:") {
		return fmt.Errorf("invalid workflow file: missing name field")
	}

	if !strings.Contains(content, "on:") {
		return fmt.Errorf("invalid workflow file: missing trigger configuration")
	}

	fmt.Printf("✅ GitHub Actions workflow '%s' validated successfully\n", workflowName)
	fmt.Printf("📁 Would create: %s/%s\n", workflowsPath, fileName)

	return nil
}