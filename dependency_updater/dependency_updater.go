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
	Tag       string `json:"tag"`
	Commit    string `json:"commit"`
	TagPrefix string `json:"tagPrefix,omitempty"`
	Owner     string `json:"owner`
	Repo      string `json:"repo`
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
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			err := updater(string(cmd.String("token")), string(cmd.String("repo")), cmd.Bool("commit"))
			if err != nil {
				return fmt.Errorf("error running updater: %s", err)
			}
			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func updater(token string, repoPath string, commit bool) error {
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
		return fmt.Errorf("error unmarshaling versions JSON to dependencies: %s", err)
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

	if commit && updatedDependencies != nil {
		err := createCommitMessage(updatedDependencies)
		if err != nil {
			return fmt.Errorf("error creating commit message: %s", err)
		}
	}

	e := createVersionsEnv(repoPath, dependencies)
	if e != nil {
		return fmt.Errorf("error creating versions.env: %s", e)
	}

	return nil
}

func createCommitMessage(updatedDependencies []VersionUpdateInfo) error {
	var repos []string
	commitTitle := "chore: updated "
	commitDescription := "Updated dependencies for: \n"

	for _, dependency := range updatedDependencies {
		repo, tag := dependency.Repo, dependency.To
		commitDescription += repo + " => " + tag + " (" + dependency.DiffUrl + ")" + "\n"
		repos = append(repos, repo)
	}

	commitTitle += strings.Join(repos, ", ")
	cmd := exec.Command("git", "commit", "-am", commitTitle, "-m", commitDescription)
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running git commit -m: %s", err)
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
	var version *github.RepositoryRelease
	var err error
	var diffUrl string
	var updatedDependency VersionUpdateInfo
	foundPrefixVersion := false
	options := &github.ListOptions{Page: 1}

	for {
		releases, resp, err := client.Repositories.ListReleases(
			ctx,
			dependencies[dependencyType].Owner,
			dependencies[dependencyType].Repo,
			options)

		if err != nil {
			return "", "", VersionUpdateInfo{}, fmt.Errorf("error getting releases: %s", err)
		}

		if dependencies[dependencyType].TagPrefix == "" {
			version = releases[0]
			if *version.TagName != dependencies[dependencyType].Tag {
				diffUrl = generateGithubRepoUrl(dependencies, dependencyType) + "/compare/" +
					dependencies[dependencyType].Tag + "..." + *version.TagName
			}
			break
		} else if dependencies[dependencyType].TagPrefix != "" {
			for release := range releases {
				if strings.HasPrefix(*releases[release].TagName, dependencies[dependencyType].TagPrefix) {
					version = releases[release]
					foundPrefixVersion = true
					if *version.TagName != dependencies[dependencyType].Tag {
						diffUrl = generateGithubRepoUrl(dependencies, dependencyType) + "/compare/" +
							dependencies[dependencyType].Tag + "..." + *version.TagName
					}
					break
				}
			}
			if foundPrefixVersion {
				break
			}
			options.Page = resp.NextPage
		} else if resp.NextPage == 0 {
			break
		}
	}

	if diffUrl != "" {
		updatedDependency = VersionUpdateInfo{
			dependencies[dependencyType].Repo,
			dependencies[dependencyType].Tag,
			*version.TagName,
			diffUrl,
		}
	}

	commit, _, err := client.Repositories.GetCommit(
		ctx,
		dependencies[dependencyType].Owner,
		dependencies[dependencyType].Repo,
		"refs/tags/"+*version.TagName,
		&github.ListOptions{})
	if err != nil {
		return "", "", VersionUpdateInfo{}, fmt.Errorf("error getting commit for "+dependencyType+": %s", err)
	}

	return *version.TagName, *commit.SHA, updatedDependency, nil
}

func updateVersionTagAndCommit(
	commit string,
	tag string,
	dependencyType string,
	repoPath string,
	dependencies Dependencies) error {
	dependencies[dependencyType].Tag = tag
	dependencies[dependencyType].Commit = commit
	err := writeToVersionsEnv(repoPath, dependencies)
	if err != nil {
		return fmt.Errorf("error writing to versions "+dependencyType+": %s", err)
	}
	
	return nil
}

func writeToVersionsEnv(repoPath string, dependencies Dependencies) error {
	// formatting json
	updatedJson, err := json.MarshalIndent(dependencies, "", "	  ")
	if err != nil {
		return fmt.Errorf("error Marshaling dependencies json: %s", err)
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

func generateGithubRepoUrl(dependencies Dependencies, dependencyType string) string {
	return "https://github.com/" + dependencies[dependencyType].Owner + "/" + dependencies[dependencyType].Repo
}
