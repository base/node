package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/ethereum-optimism/optimism/op-service/retry"
	"github.com/google/go-github/v72/github"
	"github.com/urfave/cli/v3"

	"log"
)

const (
	// Configuration constants
	maxRetries      = 3
	retryDelay      = 1 * time.Second
	filePermissions = 0644
	jsonIndent      = "  "
	
	// GitHub API constants
	initialPage = 1
)

type Info struct {
	Tag       string `json:"tag,omitempty"`
	Commit    string `json:"commit"`
	TagPrefix string `json:"tagPrefix,omitempty"`
	Owner     string `json:"owner"`
	Repo      string `json:"repo"`
	Branch    string `json:"branch,omitempty"`
	Tracking  string `json:"tracking"`
}

type VersionUpdateInfo struct {
	Repo    string
	From    string
	To      string
	DiffURL string
}

type Dependencies map[string]*Info

// UpdateResult holds the result of a single dependency update
type UpdateResult struct {
	Info  VersionUpdateInfo
	Error error
}

func main() {
	cmd := &cli.Command{
		Name:  "updater",
		Usage: "Updates the dependencies in the geth, nethermind and reth Dockerfiles",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "token",
				Usage:    "Auth token used to make requests to the Github API",
				Sources:  cli.EnvVars("GITHUB_TOKEN"),
				Required: true,
			},
			&cli.StringFlag{
				Name:     "repo",
				Usage:    "Specifies repo location to run the version updater on",
				Required: true,
			},
			&cli.BoolFlag{
				Name:     "commit",
				Usage:    "Stages updater changes and creates commit message",
				Required: false,
			},
			&cli.BoolFlag{
				Name:     "github-action",
				Usage:    "Specifies whether tool is being used through github action workflow",
				Required: false,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if err := updater(ctx, cmd.String("token"), cmd.String("repo"), cmd.Bool("commit"), cmd.Bool("github-action")); err != nil {
				return fmt.Errorf("failed to run updater: %w", err)
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func updater(ctx context.Context, token, repoPath string, commit, githubAction bool) error {
	dependencies, err := loadDependencies(repoPath)
	if err != nil {
		return fmt.Errorf("error loading dependencies: %w", err)
	}

	client := github.NewClient(nil).WithAuthToken(token)

	// Parallel processing of dependencies
	updatedDependencies, err := processAllDependencies(ctx, client, dependencies, repoPath)
	if err != nil {
		return fmt.Errorf("error processing dependencies: %w", err)
	}

	if err := createVersionsEnv(repoPath, dependencies); err != nil {
		return fmt.Errorf("error creating versions.env: %w", err)
	}

	if (commit || githubAction) && len(updatedDependencies) > 0 {
		if err := createCommitMessage(updatedDependencies, repoPath, githubAction); err != nil {
			return fmt.Errorf("error creating commit message: %w", err)
		}
	}

	return nil
}

// loadDependencies reads and parses the versions.json file
func loadDependencies(repoPath string) (Dependencies, error) {
	data, err := os.ReadFile(repoPath + "/versions.json")
	if err != nil {
		return nil, fmt.Errorf("error reading versions.json: %w", err)
	}

	var dependencies Dependencies
	if err := json.Unmarshal(data, &dependencies); err != nil {
		return nil, fmt.Errorf("error unmarshalling versions.json: %w", err)
	}

	return dependencies, nil
}

// processAllDependencies processes all dependencies in parallel
func processAllDependencies(ctx context.Context, client *github.Client, dependencies Dependencies, repoPath string) ([]VersionUpdateInfo, error) {
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		results []VersionUpdateInfo
	)

	// Channel to collect results
	resultCh := make(chan UpdateResult, len(dependencies))

	// Process each dependency concurrently
	for name := range dependencies {
		wg.Add(1)
		go func(depName string) {
			defer wg.Done()

			var result UpdateResult
			err := retry.Do0(ctx, maxRetries, retry.Fixed(retryDelay), func() error {
				info, err := getAndUpdateDependency(ctx, client, depName, repoPath, dependencies)
				result.Info = info
				return err
			})
			result.Error = err
			resultCh <- result
		}(name)
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Collect results
	for result := range resultCh {
		if result.Error != nil {
			return nil, fmt.Errorf("error processing dependency: %w", result.Error)
		}
		if result.Info != (VersionUpdateInfo{}) {
			mu.Lock()
			results = append(results, result.Info)
			mu.Unlock()
		}
	}

	return results, nil
}

func createCommitMessage(updatedDependencies []VersionUpdateInfo, repoPath string, githubAction bool) error {
	if len(updatedDependencies) == 0 {
		return nil
	}

	repos := make([]string, 0, len(updatedDependencies))
	descriptionLines := []string{"### Dependency Updates"}

	for _, dep := range updatedDependencies {
		descriptionLines = append(descriptionLines,
			fmt.Sprintf("**%s** - %s: [diff](%s)", dep.Repo, dep.To, dep.DiffURL))
		repos = append(repos, dep.Repo)
	}

	commitTitle := "chore: updated " + strings.Join(repos, ", ")
	commitDescription := strings.Join(descriptionLines, "\n")

	if githubAction {
		return writeToGithubOutput(commitTitle, commitDescription)
	}

	cmd := exec.Command("git", "commit", "-am", commitTitle, "-m", commitDescription)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run git commit: %w", err)
	}

	return nil
}

func getAndUpdateDependency(ctx context.Context, client *github.Client, dependencyType, repoPath string, dependencies Dependencies) (VersionUpdateInfo, error) {
	version, commit, updatedDependency, err := getVersionAndCommit(ctx, client, dependencies, dependencyType)
	if err != nil {
		return VersionUpdateInfo{}, err
	}

	if updatedDependency != (VersionUpdateInfo{}) {
		if err := updateVersionTagAndCommit(commit, version, dependencyType, repoPath, dependencies); err != nil {
			return VersionUpdateInfo{}, fmt.Errorf("error updating version tag and commit: %w", err)
		}
	}

	return updatedDependency, nil
}

func getVersionAndCommit(ctx context.Context, client *github.Client, dependencies Dependencies, dependencyType string) (string, string, VersionUpdateInfo, error) {
	dep := dependencies[dependencyType]

	switch dep.Tracking {
	case "tag":
		return getTagVersion(ctx, client, dep, dependencyType)
	case "branch":
		return getBranchVersion(ctx, client, dep, dependencyType)
	default:
		return "", "", VersionUpdateInfo{}, fmt.Errorf("unknown tracking type: %s", dep.Tracking)
	}
}

func getTagVersion(ctx context.Context, client *github.Client, dep *Info, dependencyType string) (string, string, VersionUpdateInfo, error) {
	options := &github.ListOptions{Page: initialPage}
	var foundVersion *github.RepositoryRelease

	for {
		releases, resp, err := client.Repositories.ListReleases(ctx, dep.Owner, dep.Repo, options)
		if err != nil {
			return "", "", VersionUpdateInfo{}, fmt.Errorf("error getting releases: %w", err)
		}

		// Find appropriate release
		if dep.TagPrefix == "" {
			foundVersion = releases[0]
			break
		}

		// Look for release with matching prefix
		for _, release := range releases {
			if strings.HasPrefix(*release.TagName, dep.TagPrefix) {
				foundVersion = release
				break
			}
		}

		if foundVersion != nil || resp.NextPage == 0 {
			break
		}
		options.Page = resp.NextPage
	}

	if foundVersion == nil {
		return "", "", VersionUpdateInfo{}, fmt.Errorf("no suitable release found")
	}

	// Get commit for tag
	versionCommit, _, err := client.Repositories.GetCommit(
		ctx, dep.Owner, dep.Repo,
		"refs/tags/"+*foundVersion.TagName,
		&github.ListOptions{})
	if err != nil {
		return "", "", VersionUpdateInfo{}, fmt.Errorf("error getting commit for %s: %w", dependencyType, err)
	}

	var updatedInfo VersionUpdateInfo
	if *foundVersion.TagName != dep.Tag {
		diffURL := generateGithubRepoURL(dep) + "/compare/" + dep.Tag + "..." + *foundVersion.TagName
		updatedInfo = VersionUpdateInfo{
			Repo:    dep.Repo,
			From:    dep.Tag,
			To:      *foundVersion.TagName,
			DiffURL: diffURL,
		}
	}

	return *foundVersion.TagName, *versionCommit.SHA, updatedInfo, nil
}

func getBranchVersion(ctx context.Context, client *github.Client, dep *Info, dependencyType string) (string, string, VersionUpdateInfo, error) {
	commits, _, err := client.Repositories.ListCommits(
		ctx, dep.Owner, dep.Repo,
		&github.CommitsListOptions{SHA: dep.Branch})
	if err != nil {
		return "", "", VersionUpdateInfo{}, fmt.Errorf("error listing commits for %s: %w", dependencyType, err)
	}

	if len(commits) == 0 {
		return "", "", VersionUpdateInfo{}, fmt.Errorf("no commits found for branch %s", dep.Branch)
	}

	latestCommit := *commits[0].SHA
	var updatedInfo VersionUpdateInfo

	if dep.Commit != latestCommit {
		diff := dep.Commit + " => " + latestCommit
		updatedInfo = VersionUpdateInfo{
			Repo:    dep.Repo,
			From:    dep.Tag,
			To:      latestCommit,
			DiffURL: diff,
		}
	}

	return "", latestCommit, updatedInfo, nil
}

func updateVersionTagAndCommit(commit, tag, dependencyType, repoPath string, dependencies Dependencies) error {
	dependencies[dependencyType].Tag = tag
	dependencies[dependencyType].Commit = commit
	return writeToVersionsJSON(repoPath, dependencies)
}

func writeToVersionsJSON(repoPath string, dependencies Dependencies) error {
	updatedJSON, err := json.MarshalIndent(dependencies, "", jsonIndent)
	if err != nil {
		return fmt.Errorf("error marshaling dependencies json: %w", err)
	}

	if err := os.WriteFile(repoPath+"/versions.json", updatedJSON, filePermissions); err != nil {
		return fmt.Errorf("error writing to versions.json: %w", err)
	}

	return nil
}

func createVersionsEnv(repoPath string, dependencies Dependencies) error {
	envLines := make([]string, 0, len(dependencies)*3)

	for name, dep := range dependencies {
		repoURL := generateGithubRepoURL(dep) + ".git"
		prefix := strings.ToUpper(name)

		tag := dep.Tag
		if dep.Tracking == "branch" {
			tag = dep.Branch
		}

		envLines = append(envLines,
			fmt.Sprintf("export %s_TAG=%s", prefix, tag),
			fmt.Sprintf("export %s_COMMIT=%s", prefix, dep.Commit),
			fmt.Sprintf("export %s_REPO=%s", prefix, repoURL),
		)
	}

	slices.Sort(envLines)

	content := strings.Join(envLines, "\n")
	if err := os.WriteFile(repoPath+"/versions.env", []byte(content), filePermissions); err != nil {
		return fmt.Errorf("error writing versions.env: %w", err)
	}

	return nil
}

func writeToGithubOutput(title, description string) error {
	file := os.Getenv("GITHUB_OUTPUT")
	if file == "" {
		return fmt.Errorf("GITHUB_OUTPUT environment variable not set")
	}

	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, filePermissions)
	if err != nil {
		return fmt.Errorf("failed to open GITHUB_OUTPUT file: %w", err)
	}
	defer f.Close()

	const delimiter = "EOF"
	output := fmt.Sprintf("TITLE=%s\nDESC<<%s\n%s\n%s\n", title, delimiter, description, delimiter)

	if _, err := f.WriteString(output); err != nil {
		return fmt.Errorf("failed to write to GITHUB_OUTPUT: %w", err)
	}

	return nil
}

func generateGithubRepoURL(dep *Info) string {
	return fmt.Sprintf("https://github.com/%s/%s", dep.Owner, dep.Repo)
}
