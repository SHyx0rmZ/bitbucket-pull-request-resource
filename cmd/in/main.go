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
)

type params struct{}

func main() {
	if len(os.Args) != 2 {
		resource.Fatal("parsing args", errors.New("too few arguments"))
	}

	var input struct {
		Source  resource.Source  `json:"source"`
		Params  params           `json:"params"`
		Version resource.Version `json:"version"`
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

	version, err := strconv.Atoi(input.Version.ID)
	if err != nil {
		resource.Fatal("parsing version", errors.New("malformed version: "+input.Version.ID))
	}

	pullRequest, err := repository.PullRequest(version)
	if err != nil {
		resource.Fatal("retrieving pull request", err)
	}

	var output struct {
		Version  resource.Version
		Metadata map[string]string
	}

	output.Version = input.Version
	output.Metadata = make(map[string]string)
	output.Metadata["state"] = pullRequest.GetState()
	output.Metadata["author"] = pullRequest.GetAuthorName()
	output.Metadata["from"] = pullRequest.GetFromRef()
	output.Metadata["to"] = pullRequest.GetToRef()
	output.Metadata["title"] = pullRequest.GetTitle()
	output.Metadata["description"] = pullRequest.GetDescription()

	for key, value := range output.Metadata {
		path := filepath.Join(os.Args[1], key)
		mode := os.FileMode(0755)

		err = os.MkdirAll(filepath.Dir(path), os.FileMode(0755))
		if err != nil {
			resource.Fatal("creating metadata directory", err)
		}

		err = ioutil.WriteFile(path, []byte(value), mode)
		if err != nil {
			resource.Fatal("writing metadata files", err)
		}
	}

	err = json.NewEncoder(os.Stdout).Encode(output)
	if err != nil {
		resource.Fatal("encoding output", err)
	}
}
