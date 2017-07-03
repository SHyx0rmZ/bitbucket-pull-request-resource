package main

import (
	"context"
	"encoding/json"
	resource "github.com/SHyx0rmZ/bitbucket-pull-request-resource"
	"github.com/SHyx0rmZ/go-bitbucket/bitbucket"
	"github.com/SHyx0rmZ/go-bitbucket/server"
	"os"
	"strconv"
)

func isNewID(version resource.Version, previous resource.Version) bool {
	return version.ID > previous.ID
}

func isNewVersion(version resource.Version, previous resource.Version) bool {
	return version.ID == previous.ID && version.Version > previous.Version
}

func main() {
	var input struct {
		Source  resource.Source   `json:"source"`
		Version *resource.Version `json:"version,omitempty"`
	}

	err := json.NewDecoder(os.Stdin).Decode(&input)
	if err != nil {
		resource.Fatal("decoding input", err)
	}

	basicAuth := &bitbucket.BasicAuth{
		Username: input.Source.Username,
		Password: input.Source.Password,
	}

	ctx := context.WithValue(context.Background(), bitbucket.BitbucketAuth, basicAuth)

	client, err := server.NewClient(ctx, input.Source.Endpoint)
	if err != nil {
		resource.Fatal("spawning Bitbucket client", err)
	}

	repository, err := client.Repository(input.Source.Owner + "/" + input.Source.Repository)
	if err != nil {
		resource.Fatal("retrieving repository info", err)
	}

	pullRequests, err := repository.PullRequests()
	if err != nil {
		resource.Fatal("retrieving list of pull requests", err)
	}

	versions := []resource.Version{}

	if input.Version != nil {
		for _, pullRequest := range pullRequests {
			version := resource.Version{
				ID:      strconv.Itoa(pullRequest.GetID()),
				Version: strconv.Itoa(pullRequest.GetVersion()),
			}

			if isNewID(version, *input.Version) || isNewVersion(version, *input.Version) {
				versions = append(versions, version)
			}
		}
	} else if len(pullRequests) > 0 {
		pullRequest := pullRequests[len(pullRequests)-1]

		version := resource.Version{
			ID:      strconv.Itoa(pullRequest.GetID()),
			Version: strconv.Itoa(pullRequest.GetVersion()),
		}

		versions = append(versions, version)
	}

	err = json.NewEncoder(os.Stdout).Encode(versions)
	if err != nil {
		resource.Fatal("encoding output", err)
	}
}
