package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum-optimism/optimism/op-service/retry"
	"github.com/google/go-github/v72/github"
	"github.com/urfave/cli/v3"
	"slices"
	"time"

	"log"
	"os"
	"strings"
)

type Info struct {
	Tag        string `json:"tag"`
	Commit     string `json:"commit"`
	TagPrefix  string `json:"tagPrefix,omitempty"`
	Owner      string `json:"owner`
	Repo       string `json:"repo`
}

type VersionTag []struct {
	Tag string `json:"tag_name"`
}

type Commit struct {
	Commit string `json:"sha"`
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
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			err := updater(string(cmd.String("token")), string(cmd.String("repo")))
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

func updater(token string, repoPath string) error {
	var err error

	f, err := os.ReadFile(repoPath + "/versions.json")
	if err != nil {
		return fmt.Errorf("error reading versions JSON: %s", err)
	}

	client := github.NewClient(nil).WithAuthToken(token)
	ctx := context.Background()

	var dependencies Dependencies

	err = json.Unmarshal(f, &dependencies)
	if err != nil {
		return fmt.Errorf("error unmarshaling versions JSON to dependencies: %s", err)
	}

	for dependency := range dependencies {
		err := retry.Do0(context.Background(), 3, retry.Fixed(1*time.Second), func() error {
			return getAndUpdateDependency(
				ctx,
				client,
				dependency,
				repoPath,
				dependencies,
			)
		})

		if err != nil {
			return fmt.Errorf("error getting and updating version/commit for "+dependency+": %s", err)
		}
	}

	e := createVersionsEnv(repoPath, dependencies)
	if e != nil {
		return fmt.Errorf("error creating versions.env: %s", e)
	}

	return nil
}

func getAndUpdateDependency(ctx context.Context, client *github.Client, dependencyType string, repoPath string, dependencies Dependencies) error {
	version, commit, err := getVersionAndCommit(ctx, client, dependencies, dependencyType)
	if err != nil {
		return err
	}

	e := updateVersionTagAndCommit(commit, version, dependencyType, repoPath, dependencies)
	if e != nil {
		return fmt.Errorf("error updating version tag and commit: %s", e)
	}

	return nil
}

func getVersionAndCommit(ctx context.Context, client *github.Client, dependencies Dependencies, dependencyType string) (string, string, error) {
	var version *github.RepositoryRelease
	var err error
	foundPrefixVersion := false
	options := &github.ListOptions{Page: 1}

	for {
		releases, resp, err := client.Repositories.ListReleases(
			ctx,
			dependencies[dependencyType].Owner,
			dependencies[dependencyType].Repo,
			options)

		if err != nil {
			return "", "", fmt.Errorf("error getting releases: %s", err)
		}

		if dependencies[dependencyType].TagPrefix == "" {
			version = releases[0]
			break
		} else if dependencies[dependencyType].TagPrefix != ""{
			for release := range releases {
				if strings.HasPrefix(*releases[release].TagName, dependencies[dependencyType].TagPrefix) {
					version = releases[release]
					foundPrefixVersion = true
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

	commit, _, err := client.Repositories.GetCommit(
		ctx,
		dependencies[dependencyType].Owner,
		dependencies[dependencyType].Repo,
		"refs/tags/"+*version.TagName,
		&github.ListOptions{})
	if err != nil {
		return "", "", fmt.Errorf("error getting commit for "+dependencyType+": %s", err)
	}

	return *version.TagName, *commit.SHA, nil
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
		repoUrl := "https://github.com/" + 
					dependencies[dependency].Owner + "/" +
					dependencies[dependency].Repo + ".git"

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
