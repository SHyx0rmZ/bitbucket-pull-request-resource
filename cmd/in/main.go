package main

import (
	"context"
	"encoding/json"
	"errors"
	resource "github.com/SHyx0rmZ/bitbucket-pull-request-resource"
	"github.com/SHyx0rmZ/go-bitbucket/bitbucket"
	"github.com/SHyx0rmZ/go-bitbucket/server"
	"github.com/concourse/atc"
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
		Metadata []atc.MetadataField
	}

	output.Version = input.Version
	output.Metadata = []atc.MetadataField{
		{Name: "state", Value: pullRequest.GetState()},
		{Name: "author", Value: pullRequest.GetAuthorName()},
		{Name: "from", Value: pullRequest.GetFromRef()},
		{Name: "to", Value: pullRequest.GetToRef()},
		{Name: "title", Value: pullRequest.GetTitle()},
		{Name: "description", Value: pullRequest.GetDescription()},
	}

	for _, field := range output.Metadata {
		path := filepath.Join(os.Args[1], field.Name)
		mode := os.FileMode(0755)

		err = os.MkdirAll(filepath.Dir(path), os.FileMode(0755))
		if err != nil {
			resource.Fatal("creating metadata directory", err)
		}

		err = ioutil.WriteFile(path, []byte(field.Value), mode)
		if err != nil {
			resource.Fatal("writing metadata files", err)
		}
	}

	err = json.NewEncoder(os.Stdout).Encode(output)
	if err != nil {
		resource.Fatal("encoding output", err)
	}
}
