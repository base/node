package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"log"
	"context"
	"github.com/urfave/cli/v3"
)

type Info struct {
	Repo string `json:"repo"`
	Tag string `json:"tag"`
	Commit string `json:"commit"`
}

type Dependencies struct {
	Op_Node Info `json:"op_node"`
	Op_Geth Info `json:"op_geth"`
	Op_Reth Info `json:"op_reth"`
	Nethermind Info `json:"nethermind"`
	Base_Reth_Node Info `json:"base_reth_node"`
}

type VersionTag []struct {
	Tag string `json:"tag_name"`
}

type Commit struct {
	Commit string `json:"sha"`
}

func main() {
	cmd := &cli.Command {
		Name: "updater",
		Usage: "Updates the dependencies in the geth, nethermind and reth Dockerfiles",
		Action: func(context.Context, *cli.Command) error {
			err := updater()
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

func updater() error {
	var err error

	f, err := os.ReadFile("../versions.json")
	if err != nil {
		return fmt.Errorf("error reading versions JSON: %s", err)
	}

	var dependencies Dependencies

	err = json.Unmarshal(f, &dependencies)
	if err != nil {
		return fmt.Errorf("error unmarshaling versions JSON to dependencies: %s", err)
	}

	token := os.Getenv("GITHUB_TOKEN")

	// updating op reth version
	err = get_and_update_version_commit(
		"https://api.github.com/repos/paradigmxyz/reth/commits/"+dependencies.Op_Reth.Tag,
		"https://api.github.com/repos/paradigmxyz/reth/releases",
		"op_reth",
		token,
		&dependencies,
	)
	if err != nil {
		return fmt.Errorf("error getting and updating version/commit for reth: %s", err)
	}

	// updating op geth version
	err = get_and_update_version_commit(
		"https://api.github.com/repos/ethereum-optimism/op-geth/commits/"+dependencies.Op_Geth.Tag,
		"https://api.github.com/repos/ethereum-optimism/op-geth/releases",
		"op_geth",
		token,
		&dependencies,
	)
	if err != nil {
		return fmt.Errorf("error getting and updating version/commit for op geth: %s", err)
	}

	// updating op node version
	err = get_and_update_version_commit(
		"https://api.github.com/repos/ethereum-optimism/optimism/commits/"+dependencies.Op_Node.Tag,
		"https://api.github.com/repos/ethereum-optimism/optimism/releases",
		"op_node",
		token,
		&dependencies,
	)
	if err != nil {
		return fmt.Errorf("error getting and updating version/commit for op node: %s", err)
	}
	// updating nethermind version
	err = get_and_update_version_commit(
		"https://api.github.com/repos/NethermindEth/nethermind/commits/"+dependencies.Nethermind.Tag,
		"https://api.github.com/repos/NethermindEth/nethermind/releases",
		"nethermind",
		token,
		&dependencies,
	)
	if err != nil {
		return fmt.Errorf("error getting and updating version/commit for nethermind: %s", err)
	}

	// updating base reth node version
	err = get_and_update_version_commit(
		"https://api.github.com/repos/base/node-reth/commits/"+dependencies.Base_Reth_Node.Tag,
		"https://api.github.com/repos/base/node-reth/releases",
		"base_reth_node",
		token,
		&dependencies,
	)
	if err != nil {
		return fmt.Errorf("error getting and updating version/commit for base reth node: %s", err)
	}

	e := create_versions_env(dependencies)
	if e != nil {
		return fmt.Errorf("error creating versions.env: %s", e)
	}

	return nil
}

func get_and_update_version_commit(
	commitUrl string,
	versionUrl string,
	dependency_type,
	token string,
	dependencies *Dependencies) error {
	tag, err := get_version_tag(versionUrl, token)
	if err != nil {
		return fmt.Errorf("error getting "+dependency_type, " version tag: %s", err)
	}

	commit, err := get_commit(commitUrl, token)
	if err != nil {
		return fmt.Errorf("error getting "+dependency_type, " commit: %s", err)
	}

	e := update_version_tag_and_commit(commit, tag, dependency_type, dependencies)
	if e != nil {
		return fmt.Errorf("error updating version tag and commit: %s", e)
	}

	return nil
}

func get_version_tag(url string, token string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating new GET request (version tag):  %s", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	jsonBody, err := get_json_body(req)
	if err != nil {
		return "", fmt.Errorf("error getting json body:  %s", err)
	}

	var versionTag VersionTag

	e := json.Unmarshal(jsonBody, &versionTag)
	if e != nil {
		return "", fmt.Errorf("error unmarshaling:  %s", e)
	}

	// loop to return most recent version of op-node
	if url == "https://api.github.com/repos/ethereum-optimism/optimism/releases" {
		incrementor := 0
		for !(strings.HasPrefix(versionTag[incrementor].Tag, "op-node")) {
			incrementor += 1
		}
		return versionTag[incrementor].Tag, nil
	}

	return versionTag[0].Tag, nil
}

func get_commit(url string, token string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating new GET request (commit): %s", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	jsonBody, err := get_json_body(req)
	if err != nil {
		return "", fmt.Errorf("error getting json body:  %s", err)
	}

	var commit Commit

	e := json.Unmarshal(jsonBody, &commit)
	if e != nil {
		fmt.Println("Error unmarshaling")
		return "", fmt.Errorf("error unmarshaling:  %s", e)
	}

	return string(commit.Commit), nil
}

func get_json_body(req *http.Request) ([]byte, error) {
	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making GET request to client:  %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received http status code when getting json body:  %s", resp.Status)
	}

	jsonBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading json body:  %s", err)
	}

	return jsonBody, nil
}

func update_version_tag_and_commit(
	commit string, 
	tag string, 
	dependency_type string, 
	dependencies *Dependencies) error {
	if dependency_type == "op_reth" {
		dependencies.Op_Reth.Tag = tag
		dependencies.Op_Reth.Commit = commit
		err := write_to_versions_env(*dependencies)
		if err != nil {
			return fmt.Errorf("error writing to versions (op reth): %s", err)
		}
	} else if dependency_type == "op_node" {
		dependencies.Op_Node.Tag = tag
		dependencies.Op_Node.Commit = commit
		err := write_to_versions_env(*dependencies)
		if err != nil {
			return fmt.Errorf("error writing to versions (op node): %s", err)
		}
	} else if dependency_type == "op_geth" {
		dependencies.Op_Geth.Tag = tag
		dependencies.Op_Geth.Commit = commit
		err := write_to_versions_env(*dependencies)
		if err != nil {
			return fmt.Errorf("error writing to versions (op geth): %s", err)
		}
	} else if dependency_type == "nethermind" {
		dependencies.Nethermind.Tag = tag
		dependencies.Nethermind.Commit = commit
		err := write_to_versions_env(*dependencies)
		if err != nil {
			return fmt.Errorf("error writing to versions (nethermind): %s", err)
		}
	} else if dependency_type == "base_reth_node" {
		dependencies.Base_Reth_Node.Tag = tag
		dependencies.Base_Reth_Node.Commit = commit
		err := write_to_versions_env(*dependencies)
		if err != nil {
			return fmt.Errorf("error writing to versions (base reth node): %s", err)
		}
	}

	return nil
}

func write_to_versions_env(dependencies Dependencies) error {
	// formatting json
	updatedJson, err := json.MarshalIndent(dependencies, "", "	  ")
	if err != nil {
		return fmt.Errorf("error Marshaling dependencies json: %s", err)
	}

	e := os.WriteFile("../versions.json", updatedJson, 0644)
	if e != nil {
		return fmt.Errorf("error writing to versions.json: %s", e)
	}

	return nil
}

func create_versions_env(dependencies Dependencies) error {
	env := "export OP_NODE_TAG=" + dependencies.Op_Node.Tag + "\n" +
		   "export OP_NODE_COMMIT=" + dependencies.Op_Node.Commit + "\n" +
		   "export OP_NODE_REPO=https://github.com/ethereum-optimism/optimism.git" + "\n\n" +
		   "export OP_GETH_TAG=" + dependencies.Op_Geth.Tag + "\n" +
		   "export OP_GETH_COMMIT=" + dependencies.Op_Geth.Commit + "\n" +
		   "export OP_GETH_REPO=https://github.com/ethereum-optimism/op-geth.git" + "\n\n" +
		   "export OP_RETH_TAG=" + dependencies.Op_Reth.Tag + "\n" +
		   "export OP_RETH_COMMIT=" + dependencies.Op_Reth.Commit + "\n" +
		   "export OP_RETH_REPO=https://github.com/paradigmxyz/reth.git" + "\n\n" +
		   "export NETHERMIND_TAG=" + dependencies.Nethermind.Tag + "\n" +
		   "export NETHERMIND_COMMIT=" + dependencies.Nethermind.Commit + "\n" +
		   "export NETHERMIND_REPO=https://github.com/NethermindEth/nethermind.git" + "\n\n" +
		   "export BASE_RETH_NODE_TAG=" + dependencies.Base_Reth_Node.Tag + "\n" +
		   "export BASE_RETH_NODE_COMMIT=" + dependencies.Base_Reth_Node.Commit + "\n" +
		   "export BASE_RETH_NODE_REPO=https://github.com/base/node-reth.git " + "\n\n"

	file, err := os.Create("../versions.env")
	if err != nil {
		return fmt.Errorf("error creating versions.env file: %s", err)
	}
	defer file.Close()

	_, err = file.WriteString(env)
	if err != nil {
		return fmt.Errorf("error writing to versions.env file: %s", err)
	}

	return nil
}