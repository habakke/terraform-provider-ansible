package main

import (
	"fmt"
	"github.com/habakke/terraform-ansible-provider/internal/ansible"
	"github.com/habakke/terraform-ansible-provider/internal/util"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

var (
	version   string // build version number
	commit    string // sha1 revision used to build the program
	buildTime string // when the executable was built
	buildBy   string
)

func getVersionString(name string) string {
	return fmt.Sprintf("%s %s (%s at %s by %s)", name, version, commit, buildTime, buildBy)
}
func main() {
	// create a logger an log version string
	logger, _ := util.CreateSubLogger("")
	logger.Info().Msgf("%s", getVersionString("terraform-ansible-provider"))

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ansible.Provider,
	})
}
