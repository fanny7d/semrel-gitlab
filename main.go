/*
go-semrel-gitlab provides tools to automate parts of release process on Gitlab CI

More documentation can be found at https://juhani.gitlab.io/go-semrel-gitlab/
*/
package main

import (
	"github.com/fanny7d/semrel-gitlab/cmd"
)

func main() {
	cmd.Execute()
}
