package main

import (
	"context"
	"encoding/json"
	"fmt"
	"slices"
	"time"

	"github.com/ethereum-optimism/optimism/op-service/retry"
	"github.com/google/go-github/v72/github"
	"github.com/urfave/cli/v3"

	"log"
	"os"
	"os/exec"
	"strings"
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
	DiffUrl string
}

type Dependencies = map[string]*Info

func main() {
	cmd := &cli.Command{
		Name:  "updater",
		Usage: "Updates the dependencies in the geth, nethermind and reth Dockerfiles",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "token",
				Usage:    "Auth token used to make requests to the Github API must be set using export",
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
			err := updater(cmd.String("token"), cmd.String("repo"), cmd.Bool("commit"), cmd.Bool("github-action"))
			if err != nil {
				return fmt.Errorf("failed to run updater: %s", err)
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func updater(token string, repoPath string, commit bool, githubAction bool) error {
	var err error
	var dependencies Dependencies
	var updatedDependencies []VersionUpdateInfo

	f, err := os.ReadFile(repoPath + "/versions.json")
	if err != nil {
		return fmt.Errorf("error reading versions JSON: %s", err)
	}

	client := github.NewClient(nil).WithAuthToken(token)
	ctx := context.Background()

	err = json.Unmarshal(f, &dependencies)
	if err != nil {
		return fmt.Errorf("error unmarshalling versions JSON to dependencies: %s", err)
	}

	for dependency := range dependencies {
		var updatedDependency VersionUpdateInfo
		err := retry.Do0(context.Background(), 3, retry.Fixed(1*time.Second), func() error {
			updatedDependency, err = getAndUpdateDependency(
				ctx,
				client,
				dependency,
				repoPath,
				dependencies,
			)
			return err
		})
		if err != nil {
			return fmt.Errorf("error getting and updating version/commit for "+dependency+": %s", err)
		}

		if updatedDependency != (VersionUpdateInfo{}) {
			updatedDependencies = append(updatedDependencies, updatedDependency)
		}
	}

	e := createVersionsEnv(repoPath, dependencies)
	if e != nil {
		return fmt.Errorf("error creating versions.env: %s", e)
	}

	if (commit && updatedDependencies != nil) || (githubAction && updatedDependencies != nil) {
		err := createCommitMessage(updatedDependencies, repoPath, githubAction)
		if err != nil {
			return fmt.Errorf("error creating commit message: %s", err)
		}
	}

	return nil
}

func createCommitMessage(updatedDependencies []VersionUpdateInfo, repoPath string, githubAction bool) error {
	var repos []string
	descriptionLines := []string{
		"### Dependency Updates",
	}

	commitTitle := "chore: updated "

	for _, dependency := range updatedDependencies {
		repo, tag := dependency.Repo, dependency.To
		descriptionLines = append(descriptionLines, fmt.Sprintf("**%s** - %s:  [diff](%s)", repo, tag, dependency.DiffUrl))
		repos = append(repos, repo)
	}
	commitDescription := strings.Join(descriptionLines, "\n")
	commitTitle += strings.Join(repos, ", ")

	if githubAction {
		err := writeToGithubOutput(commitTitle, commitDescription, repoPath)
		if err != nil {
			return fmt.Errorf("error creating git commit message: %s", err)
		}
	} else {
		cmd := exec.Command("git", "commit", "-am", commitTitle, "-m", commitDescription)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to run git commit -m: %s", err)
		}
	}
	return nil
}

func getAndUpdateDependency(ctx context.Context, client *github.Client, dependencyType string, repoPath string, dependencies Dependencies) (VersionUpdateInfo, error) {
	version, commit, updatedDependency, err := getVersionAndCommit(ctx, client, dependencies, dependencyType)
	if err != nil {
		return VersionUpdateInfo{}, err
	}
	if updatedDependency != (VersionUpdateInfo{}) {
		e := updateVersionTagAndCommit(commit, version, dependencyType, repoPath, dependencies)
		if e != nil {
			return VersionUpdateInfo{}, fmt.Errorf("error updating version tag and commit: %s", e)
		}
	}

	return updatedDependency, nil
}

func getVersionAndCommit(ctx context.Context, client *github.Client, dependencies Dependencies, dependencyType string) (string, string, VersionUpdateInfo, error) {
	var selectedTag *github.RepositoryTag
	var commit string
	var diffUrl string
	var updatedDependency VersionUpdateInfo
	options := &github.ListOptions{Page: 1}
	currentTag := dependencies[dependencyType].Tag
	tagPrefix := dependencies[dependencyType].TagPrefix

	if dependencies[dependencyType].Tracking == "tag" || dependencies[dependencyType].Tracking == "release" {
		// Collect all valid tags across all pages, then find the max version
		var validTags []*github.RepositoryTag
		trackingMode := dependencies[dependencyType].Tracking

		for {
			tags, resp, err := client.Repositories.ListTags(
				ctx,
				dependencies[dependencyType].Owner,
				dependencies[dependencyType].Repo,
				options)

			if err != nil {
				return "", "", VersionUpdateInfo{}, fmt.Errorf("error getting tags: %s", err)
			}

			for _, tag := range tags {
				// Skip if tagPrefix is set and doesn't match
				if tagPrefix != "" && !strings.HasPrefix(*tag.Name, tagPrefix) {
					continue
				}

				// Filter based on tracking mode:
				// - "release": only stable releases (no prerelease suffix)
				// - "tag": releases and RC versions only (exclude -synctest, -alpha, etc.)
				if trackingMode == "release" {
					if !IsReleaseVersion(*tag.Name, tagPrefix) {
						continue
					}
				} else if trackingMode == "tag" {
					if !IsReleaseOrRCVersion(*tag.Name, tagPrefix) {
						continue
					}
				}

				// Check if this is a valid upgrade (not a downgrade)
				if err := ValidateVersionUpgrade(currentTag, *tag.Name, tagPrefix); err != nil {
					continue
				}

				validTags = append(validTags, tag)
			}

			if resp.NextPage == 0 {
				break
			}
			options.Page = resp.NextPage
		}

		// Find the maximum version among valid tags
		for _, tag := range validTags {
			// Skip if this tag can't be parsed
			if _, err := ParseVersion(*tag.Name, tagPrefix); err != nil {
				log.Printf("Skipping unparseable tag %s: %v", *tag.Name, err)
				continue
			}

			if selectedTag == nil {
				selectedTag = tag
				continue
			}

			cmp, err := CompareVersions(*tag.Name, *selectedTag.Name, tagPrefix)
			if err != nil {
				log.Printf("Error comparing versions %s and %s: %v", *tag.Name, *selectedTag.Name, err)
				continue
			}
			if cmp > 0 {
				selectedTag = tag
			}
		}

		// If no valid version found, keep current version
		if selectedTag == nil {
			log.Printf("No valid upgrade found for %s, keeping %s", dependencyType, currentTag)
			return currentTag, dependencies[dependencyType].Commit, VersionUpdateInfo{}, nil
		}

		if *selectedTag.Name != currentTag {
			diffUrl = generateGithubRepoUrl(dependencies, dependencyType) + "/compare/" +
				currentTag + "..." + *selectedTag.Name
		}

		// Get commit SHA from the tag
		commit = *selectedTag.Commit.SHA
	}

	if diffUrl != "" {
		updatedDependency = VersionUpdateInfo{
			dependencies[dependencyType].Repo,
			dependencies[dependencyType].Tag,
			*selectedTag.Name,
			diffUrl,
		}
	}

	if dependencies[dependencyType].Tracking == "branch" {
		branchCommit, _, err := client.Repositories.ListCommits(
			ctx,
			dependencies[dependencyType].Owner,
			dependencies[dependencyType].Repo,
			&github.CommitsListOptions{
				SHA: dependencies[dependencyType].Branch,
			},
		)
		if err != nil {
			return "", "", VersionUpdateInfo{}, fmt.Errorf("error listing commits for "+dependencyType+": %s", err)
		}
		commit = *branchCommit[0].SHA
		if dependencies[dependencyType].Commit != commit {
			from, to := dependencies[dependencyType].Commit, commit
			diffUrl = fmt.Sprintf("%s/compare/%s...%s", generateGithubRepoUrl(dependencies, dependencyType), from, to)
			updatedDependency = VersionUpdateInfo{
				dependencies[dependencyType].Repo,
				dependencies[dependencyType].Tag,
				commit,
				diffUrl,
			}
		}
	}

	if selectedTag != nil {
		return *selectedTag.Name, commit, updatedDependency, nil
	}

	return "", commit, updatedDependency, nil
}

func updateVersionTagAndCommit(
	commit string,
	tag string,
	dependencyType string,
	repoPath string,
	dependencies Dependencies) error {
	dependencies[dependencyType].Tag = tag
	dependencies[dependencyType].Commit = commit
	err := writeToVersionsJson(repoPath, dependencies)
	if err != nil {
		return fmt.Errorf("error writing to versions "+dependencyType+": %s", err)
	}

	return nil
}

func writeToVersionsJson(repoPath string, dependencies Dependencies) error {
	// formatting json
	updatedJson, err := json.MarshalIndent(dependencies, "", "	  ")
	if err != nil {
		return fmt.Errorf("error marshaling dependencies json: %s", err)
	}

	e := os.WriteFile(repoPath+"/versions.json", updatedJson, 0644)
	if e != nil {
		return fmt.Errorf("error writing to versions.json: %s", e)
	}

	return nil
}

func createVersionsEnv(repoPath string, dependencies Dependencies) error {
	envLines := []string{}

	for dependency := range dependencies {
		repoUrl := generateGithubRepoUrl(dependencies, dependency) + ".git"

		dependencyPrefix := strings.ToUpper(dependency)

		if dependencies[dependency].Tracking == "branch" {
			dependencies[dependency].Tag = dependencies[dependency].Branch
		}

		envLines = append(envLines, fmt.Sprintf("export %s_%s=%s",
			dependencyPrefix, "TAG", dependencies[dependency].Tag))

		envLines = append(envLines, fmt.Sprintf("export %s_%s=%s",
			dependencyPrefix, "COMMIT", dependencies[dependency].Commit))

		envLines = append(envLines, fmt.Sprintf("export %s_%s=%s",
			dependencyPrefix, "REPO", repoUrl))
	}

	slices.Sort(envLines)

	file, err := os.Create(repoPath + "/versions.env")
	if err != nil {
		return fmt.Errorf("error creating versions.env file: %s", err)
	}
	defer file.Close()

	_, err = file.WriteString(strings.Join(envLines, "\n"))
	if err != nil {
		return fmt.Errorf("error writing to versions.env file: %s", err)
	}

	return nil
}

func writeToGithubOutput(title string, description string, repoPath string) error {
	file := os.Getenv("GITHUB_OUTPUT")
	f, err := os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open GITHUB_OUTPUT file: %s", err)
	}
	defer f.Close()

	titleToWrite := fmt.Sprintf("%s=%s\n", "TITLE", title)
	_, err = f.WriteString(titleToWrite)
	if err != nil {
		return fmt.Errorf("failed to write to GITHUB_OUTPUT file: %s", err)
	}

	delimiter := "EOF"
	descToWrite := fmt.Sprintf("%s<<%s\n%s\n%s\n", "DESC", delimiter, description, delimiter)
	_, err = f.WriteString(descToWrite)
	if err != nil {
		return fmt.Errorf("failed to write to GITHUB_OUTPUT file: %s", err)
	}

	return nil
}

func generateGithubRepoUrl(dependencies Dependencies, dependencyType string) string {
	return "https://github.com/" + dependencies[dependencyType].Owner + "/" + dependencies[dependencyType].Repo
}
