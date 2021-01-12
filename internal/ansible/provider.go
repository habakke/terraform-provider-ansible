package ansible

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

type ProviderMetadata struct {
	Path string
}

func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("INVENTORY_PATH", nil),
				Description: "Path to where the ansible inventory files are stored",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ResourcesMap: map[string]*schema.Resource{
			"ansible_inventory": ansibleInventoryResourceQuery(),
			"ansible_group":     ansibleGroupResourceQuery(),
			"ansible_host":      ansibleHostResourceQuery(),
		},
	}
	p.ConfigureFunc = providerConfigure(p)
	return p
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		path := d.Get("path").(string)
		meta := ProviderMetadata{
			Path: path,
		}
		return meta, nil
	}
}
