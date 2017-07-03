package main

import (
	"context"
	"encoding/json"
	"errors"
	resource "github.com/SHyx0rmZ/bitbucket-pull-request-resource"
	"github.com/SHyx0rmZ/go-bitbucket/bitbucket"
	"github.com/SHyx0rmZ/go-bitbucket/server"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type params struct {
	Path string `json:"path"`
}

func main() {
	if len(os.Args) != 2 {
		resource.Fatal("parsing args", errors.New("too fex arguments"))
	}

	var input struct {
		Source resource.Source `json:"source"`
		Params params          `json:"params"`
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

	var output struct {
		Version  resource.Version
		Metadata map[string]string
	}

	output.Metadata = make(map[string]string)

	for _, key := range []string{"from", "to", "title", "description", "reviewers"} {
		path := filepath.Join(os.Args[1], input.Params.Path, key)

		bytes, err := ioutil.ReadFile(path)
		if err != nil {
			resource.Fatal("reading metadata files", err)
		}

		output.Metadata[key] = string(bytes)
	}

	repository, err := client.Repository(input.Source.Owner + "/" + input.Source.Repository)
	if err != nil {
		resource.Fatal("retrieving repository info", err)
	}

	pullRequest, err := repository.CreatePullRequest(
		nil,
		output.Metadata["from"],
		output.Metadata["to"],
		output.Metadata["title"],
		output.Metadata["description"],
		strings.Split(output.Metadata["reviewers"], "\n")...,
	)
	if err != nil {
		resource.Fatal("creating pull request", err)
	}

	output.Version.ID = strconv.Itoa(pullRequest.GetID())
	output.Version.Version = strconv.Itoa(pullRequest.GetVersion())

	err = json.NewEncoder(os.Stdout).Encode(output)
	if err != nil {
		resource.Fatal("encoding output", err)
	}
}
