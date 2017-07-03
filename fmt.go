package bitbucket_pull_request_resource

import (
	"fmt"
	"github.com/k0kubun/go-ansi"
	"github.com/mitchellh/colorstring"
	"os"
)

func Fatal(doing string, err error) {
	Sayf(colorstring.Color("[red]error %s: %s\n"), doing, err)
	os.Exit(1)
}

func Sayf(message string, args ...interface{}) {
	fmt.Fprintf(ansi.NewAnsiStderr(), message, args...)
}
