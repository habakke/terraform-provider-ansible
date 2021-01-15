package ansible

import (
	"fmt"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/inventory"
	"github.com/habakke/terraform-ansible-provider/internal/util"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"sync"
)

type ProviderConfiguration struct {
	Path        string
	Mutex       *sync.Mutex
	LogFile     string
	LogLevels   map[string]string
	Inventories inventory.InventoryMap
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
			"log_enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Should logging be enabled or not",
			},
			"log_levels": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Which log levels should be enabled ERROR, WARN, INFO, DEBUG, TRACE",
			},
			"log_file": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "terraform-provider-ansible.log",
				Description: "Name of file where log will be stored",
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

		// look to see what logging we should be outputting according to the provider configuration
		logLevels := make(map[string]string)
		for logger, level := range d.Get("log_levels").(map[string]interface{}) {
			levelAsString, ok := level.(string)
			if ok {
				logLevels[logger] = levelAsString
			} else {
				return nil, fmt.Errorf("Invalid logging level %v for %v. Be sure to use a string.", level, logger)
			}
		}

		// configure logging
		// NOTE: if enable is false here, the configuration will squash all output
		util.ConfigureLogger(
			d.Get("log_enable").(bool),
			d.Get("log_file").(string),
			logLevels,
		)

		path := d.Get("path").(string)
		var mut sync.Mutex
		conf := ProviderConfiguration{
			Path:        path,
			Mutex:       &mut,
			LogFile:     d.Get("log_file").(string),
			LogLevels:   logLevels,
			Inventories: inventory.NewInventoryMap(),
		}
		return conf, nil
	}
}
