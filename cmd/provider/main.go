package main

import (
	"github.com/habakke/terraform-ansible-provider/internal/ansible"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: ansible.Provider,
	})
}
