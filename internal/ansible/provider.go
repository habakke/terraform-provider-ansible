package ansible

import (
	"context"
	"github.com/habakke/terraform-ansible-provider/internal/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"sync"
)

type providerConfiguration struct {
	Path  string
	Mutex *sync.Mutex
}

// Provider represents a terraform provider definition
func Provider() *schema.Provider {
	return New()
}

// New represents a terraform provider definition
func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"path": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("INVENTORY_PATH", nil),
				ValidateFunc: validation.NoZeroValues,
				Description:  "Path to where the ansible inventory files are stored",
			},
			"log_caller": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Include calling function in log entries",
				Default:     false,
			},
		},
		DataSourcesMap: map[string]*schema.Resource{},
		ResourcesMap: map[string]*schema.Resource{
			"ansible_inventory": ansibleInventoryResourceQuery(),
			"ansible_group":     ansibleGroupResourceQuery(),
			"ansible_host":      ansibleHostResourceQuery(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	// load provider config vars
	path := util.ResourceToString(d, "path")

	var mut sync.Mutex
	conf := providerConfiguration{
		Path:  path,
		Mutex: &mut,
	}
	return conf, diags
}
