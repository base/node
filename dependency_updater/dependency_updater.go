package main

import (
	"context"
	"encoding/json"
	"errors" // Added for advanced error handling
	"fmt"
	"log"
	"os"
	"os/exec"
	"slices"
	"strings"
	"time"

	"github.com/ethereum-optimism/optimism/op-service/retry"
	"github.com/google/go-github/v72/github"
	"github.com/urfave/cli/v3"
)

// Info holds the version and repository metadata for a single dependency.
type Info struct {
	Tag string `json:"tag,omitempty"`
	Commit string `json:"commit"`
	TagPrefix string `json:"tagPrefix,omitempty"`
	Owner string `json:"owner"`
	Repo string `json:"repo"`
	Branch string `json:"branch,omitempty"`
	Tracking string `json:"tracking"`
}

// VersionUpdateInfo describes a specific version change for a dependency.
type VersionUpdateInfo struct {
	Repo    string
	From    string
	To      string
	DiffUrl string
}

type Dependencies = map[string]*Info

func main() {
	cmd := &cli.Command{
		Name:  "updater",
		Usage: "Updates repository dependency versions defined in versions.json (e.g., for Dockerfiles)",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:      "token",
				Usage:     "Authentication token for the Github API. Must be set using export GITHUB_TOKEN.",
				Sources:   cli.EnvVars("GITHUB_TOKEN"),
				Required:  true,
			},
			&cli.StringFlag{
				Name:      "repo",
				Usage:     "Specifies the local path of the repository containing versions.json.",
				Required:  true,
			},
			&cli.BoolFlag{
				Name:      "commit",
				Usage:     "Stages changes, creates a commit message, and runs 'git commit -am'.",
				Required:  false,
			},
			&cli.BoolFlag{
				Name:      "github-action",
				Usage:     "Specifies whether the tool is running inside a GitHub Actions workflow (outputs to GITHUB_OUTPUT).",
				Required:  false,
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			if err := runUpdater(cmd.String("token"), cmd.String("repo"), cmd.Bool("commit"), cmd.Bool("github-action")); err != nil {
				// Use errors.Wrap or fmt.Errorf("%w") for proper error chain tracing
				return fmt.Errorf("failed to run updater: %w", err)
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		// Use log.Fatalf for cleaner error termination
		log.Fatalf("Updater failed: %v", err)
	}
}

// runUpdater orchestrates the entire update process: reading JSON, fetching updates,
// writing back, creating environment variables, and optionally committing changes.
func runUpdater(token string, repoPath string, commit bool, githubAction bool) error {
	// Remove unnecessary 'var err error' declaration. Use short declarations.
	
	f, err := os.ReadFile(repoPath + "/versions.json")
	if err != nil {
		return fmt.Errorf("error reading versions JSON at %s: %w", repoPath+"/versions.json", err)
	}

	client := github.NewClient(nil).WithAuthToken(token)
	ctx := context.Background()

	var dependencies Dependencies
	if err := json.Unmarshal(f, &dependencies); err != nil {
		return fmt.Errorf("error unmarshalling versions JSON to dependencies: %w", err)
	}

	updatedDependencies := make([]VersionUpdateInfo, 0)

	// Process each dependency with retry mechanism
	for dependencyType := range dependencies {
		var updatedDependency VersionUpdateInfo
		
		// Use a dedicated, named context for retries if needed, otherwise Background is fine.
		// Using a more informative error message in the Do0 call.
		if err := retry.Do0(context.Background(), 3, retry.Fixed(1*time.Second), func() error {
			updatedDependency, err = getAndUpdateDependency(
				ctx,
				client,
				dependencyType,
				repoPath,
				dependencies,
			)
			if err != nil {
				return fmt.Errorf("attempt failed for %s: %w", dependencyType, err)
			}
			return nil
		}); err != nil {
			return fmt.Errorf("failed after all retries for %s: %w", dependencyType, err)
		}

		// Only append if an update was actually found.
		if updatedDependency != (VersionUpdateInfo{}) {
			updatedDependencies = append(updatedDependencies, updatedDependency)
		}
	}

	// Create versions.env file
	if err := createVersionsEnv(repoPath, dependencies); err != nil {
		return fmt.Errorf("error creating versions.env: %w", err)
	}

	// Commit changes if requested AND if updates were found.
	if (commit || githubAction) && len(updatedDependencies) > 0 {
		if err := createCommitMessage(updatedDependencies, repoPath, githubAction); err != nil {
			return fmt.Errorf("error creating commit message: %w", err)
		}
	}

	return nil
}

// createCommitMessage generates and executes a git commit or writes output for GitHub Actions.
func createCommitMessage(updatedDependencies []VersionUpdateInfo, repoPath string, githubAction bool) error {
	var repos []string
	descriptionLines := []string{
		"### Dependency Updates",
	}

	// Building description lines
	for _, dependency := range updatedDependencies {
		repo, tag := dependency.Repo, dependency.To
		descriptionLines = append(descriptionLines, fmt.Sprintf("**%s** - %s: [diff](%s)", repo, tag, dependency.DiffUrl))
		repos = append(repos, dependency.Repo) // Use dependency.Repo for clarity
	}
	
	// Ensure unique repo names in the title
	slices.Sort(repos)
	repos = slices.Compact(repos)
	
	commitTitle := "chore: updated " + strings.Join(repos, ", ")
	commitDescription := strings.Join(descriptionLines, "\n")
	
	if githubAction {
		if err := writeToGithubOutput(commitTitle, commitDescription); err != nil { // repoPath is not needed here
			return fmt.Errorf("error writing to GitHub output: %w", err)
		}
	} else {
		// Use '-m' flag multiple times for multi-line messages, or pipe the description.
		// Using exec.Command with multiple arguments for safety and clarity.
		cmd := exec.Command("git", "commit", "-m", commitTitle, "-m", commitDescription)
		
		// Set working directory for safety
		cmd.Dir = repoPath
		
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to run git commit: %s (Output: %s)", err, string(output))
		}
	}
	return nil
}

// getAndUpdateDependency fetches the latest version/commit and updates the dependencies map if necessary.
func getAndUpdateDependency(ctx context.Context, client *github.Client, dependencyType string, repoPath string, dependencies Dependencies) (VersionUpdateInfo, error) {
	version, commit, updatedDependency, err := getVersionAndCommit(ctx, client, dependencies, dependencyType)
	if err != nil {
		return VersionUpdateInfo{}, err // Error is already wrapped in getVersionAndCommit
	}
	
	// Check if an update was found (VersionUpdateInfo is not zero-value)
	if updatedDependency != (VersionUpdateInfo{}) {
		if err := updateVersionTagAndCommit(commit, version, dependencyType, repoPath, dependencies); err != nil {
			return VersionUpdateInfo{}, fmt.Errorf("error updating version tag and commit: %w", err)
		}
	}

	return updatedDependency, nil
}

// getVersionAndCommit fetches the latest version (tag or branch head) and its corresponding commit SHA.
func getVersionAndCommit(ctx context.Context, client *github.Client, dependencies Dependencies, dependencyType string) (string, string, VersionUpdateInfo, error) {
	depInfo := dependencies[dependencyType]
	repoURL := generateGithubRepoURL(depInfo.Owner, depInfo.Repo)
	
	var latestTag string
	var latestCommitSHA string
	var diffURL string
	var updatedDependency VersionUpdateInfo

	// === TAG TRACKING LOGIC ===
	if depInfo.Tracking == "tag" {
		version, commitSHA, diff, err := getLatestTag(ctx, client, depInfo)
		if err != nil {
			return "", "", VersionUpdateInfo{}, fmt.Errorf("error fetching latest tag for %s: %w", depInfo.Repo, err)
		}
		latestTag = version
		latestCommitSHA = commitSHA
		diffURL = diff
	}

	// === BRANCH TRACKING LOGIC ===
	if depInfo.Tracking == "branch" {
		commitSHA, err := getLatestBranchCommit(ctx, client, depInfo)
		if err != nil {
			return "", "", VersionUpdateInfo{}, fmt.Errorf("error fetching latest commit for branch %s/%s: %w", depInfo.Owner, depInfo.Repo, err)
		}
		latestCommitSHA = commitSHA
		latestTag = depInfo.Branch // For tracking purposes, 'Tag' is the branch name

		if depInfo.Commit != latestCommitSHA {
			// When tracking a branch, the 'From' is the old commit, 'To' is the new commit.
			// DiffURL points to the commit range comparison.
			diffURL = fmt.Sprintf("%s/compare/%s...%s", repoURL, depInfo.Commit, latestCommitSHA)

			// The 'To' in VersionUpdateInfo is the new commit SHA when tracking a branch
			updatedDependency = VersionUpdateInfo{
				Repo:    depInfo.Repo,
				From:    depInfo.Commit,
				To:      latestCommitSHA,
				DiffUrl: diffURL,
			}
		}
	}
	
	// Create VersionUpdateInfo for tag tracking if an update was found.
	if depInfo.Tracking == "tag" && diffURL != "" && latestTag != depInfo.Tag {
		updatedDependency = VersionUpdateInfo{
			Repo:    depInfo.Repo,
			From:    depInfo.Tag,
			To:      latestTag,
			DiffUrl: diffURL,
		}
	}

	// Return the latest known tag/branch and its commit SHA.
	// If tracking tag, latestTag contains the tag name. If tracking branch, it contains the branch name.
	return latestTag, latestCommitSHA, updatedDependency, nil
}

// getLatestTag retrieves the latest release tag matching the prefix, if any.
func getLatestTag(ctx context.Context, client *github.Client, depInfo *Info) (string, string, string, error) {
	options := &github.ListOptions{Page: 1, PerPage: 100} // Increase PerPage for efficiency

	for {
		releases, resp, err := client.Repositories.ListReleases(
			ctx,
			depInfo.Owner,
			depInfo.Repo,
			options)

		if err != nil {
			return "", "", "", fmt.Errorf("error getting releases for %s: %w", depInfo.Repo, err)
		}

		for _, release := range releases {
			tagName := *release.TagName
			if depInfo.TagPrefix == "" || strings.HasPrefix(tagName, depInfo.TagPrefix) {
				// Found the latest release (either no prefix or matching prefix)
				
				// Get commit SHA for the tag
				commit, _, err := client.Repositories.GetCommit(
					ctx,
					depInfo.Owner,
					depInfo.Repo,
					"refs/tags/"+tagName,
					nil) // ListOptions not needed here
				if err != nil {
					return "", "", "", fmt.Errorf("error getting commit for tag %s: %w", tagName, err)
				}
				commitSHA := *commit.SHA
				
				var diffURL string
				if tagName != depInfo.Tag {
					// Update found, generate diff URL
					repoURL := generateGithubRepoURL(depInfo.Owner,
