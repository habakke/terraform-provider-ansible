package main

import (
	"github.com/habakke/terraform-ansible-provider/internal/ansible"
	"github.com/habakke/terraform-ansible-provider/internal/util"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	// create a logger an log version string
	logger, _ := util.CreateSubLogger("")
	logger.Info().Msgf("%s", util.GetVersionString())

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ansible.Provider,
	})
}
