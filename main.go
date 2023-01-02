package main

import (
	"context"
	"fmt"
	"github.com/habakke/terraform-ansible-provider/internal/ansible"
	"github.com/habakke/terraform-ansible-provider/internal/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/rs/zerolog/log"
	"os"
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
	ctx := context.Background()
	logger := util.NewTerraformLogger()
	path, err := os.Getwd()
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to initialize provider: %s", err.Error())
	}

	logger.Infof(ctx, "%s", getVersionString("terraform-ansible-provider"))
	logger.Infof(ctx, "%s", path)
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return ansible.New()
		},
	})
}
