package ansible

import (
	"context"
	"fmt"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/inventory"
	"github.com/habakke/terraform-ansible-provider/internal/util"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"
)

func ansibleInventoryResourceQuery() *schema.Resource {
	return &schema.Resource{
		Create: ansibleInventoryResourceQueryCreate,
		Read:   ansibleInventoryResourceQueryRead,
		Update: ansibleInventoryResourceQueryUpdate,
		Delete: ansibleInventoryResourceQueryDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Second),
			Update: schema.DefaultTimeout(10 * time.Second),
			Delete: schema.DefaultTimeout(10 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"group_vars": {
				Type:         schema.TypeString,
				Required:     true,
				DefaultFunc:  schema.EnvDefaultFunc("GROUP_VARS", nil),
				Description:  "Ansible inventory group vars",
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func ansibleInventoryResourceQueryCreate(d *schema.ResourceData, meta interface{}) error {
	conf := meta.(ProviderConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	groupVars := d.Get("group_vars").(string)

	// create a logger for this function
	logger, _ := util.CreateSubLogger("resource_inventory_create")
	logger.Debug().Str("path", conf.Path).Msg("invoking creation of inventory")

	conf.Mutex.Lock()
	i := inventory.NewInventory(conf.Path)
	if err := conf.Inventories.Add(i); err != nil {
		logger.Error().Err(err).Msg("failed to create inventory")
		return fmt.Errorf("failed to create inventory: %e", err)
	}
	if err := i.Commit(groupVars); err != nil {
		logger.Error().Err(err).Msg("failed to commit inventory")
		return fmt.Errorf("failed to commit inventory: %e", err)
	}
	conf.Mutex.Unlock()

	d.SetId(i.GetId())
	d.MarkNewResource()
	return ansibleInventoryResourceQueryRead(d, meta)
}

func ansibleInventoryResourceQueryRead(d *schema.ResourceData, meta interface{}) error {
	conf := meta.(ProviderConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	id := d.Id()

	// create a logger for this function
	logger, _ := util.CreateSubLogger("resource_inventory_read")
	logger.Debug().Str("id", d.Id()).Msg("reading configuration for inventory")

	conf.Mutex.Lock()
	i, err := conf.Inventories.GetOrCreate(id)
	if err != nil {
		logger.Error().Err(err).Msg("faile to read inventory")
		return fmt.Errorf("failed to read inventory '%s': %e", id, err)
	}

	err, groupVars := i.Load()
	if err != nil {
		logger.Error().Err(err).Msg("failed to load inventory")
		return fmt.Errorf("failed to load inventory '%s': %e", id, err)
	}
	conf.Mutex.Unlock()

	_ = d.Set("group_vars", groupVars)

	return nil
}

func ansibleInventoryResourceQueryUpdate(d *schema.ResourceData, meta interface{}) error {
	conf := meta.(ProviderConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	id := d.Id()
	groupVars := d.Get("group_vars").(string)

	// create a logger for this function
	logger, _ := util.CreateSubLogger("resource_host_update")
	logger.Debug().Str("id", d.Id()).Str("groupVars", groupVars).Msg("updating configuration for inventory")

	i, err := conf.Inventories.GetOrCreate(id)
	if err != nil {
		logger.Error().Err(err).Msg("failed to load inventory")
		return fmt.Errorf("failed to load inventory '%s': %e", id, err)
	}

	if d.HasChange("group_vars") {
		conf.Mutex.Lock()
		if err := i.Commit(groupVars); err != nil {
			logger.Error().Err(err).Msg("failed to update inventory")
			return fmt.Errorf("failed to update inventory: %e", err)
		}
		conf.Mutex.Unlock()
	}

	return ansibleInventoryResourceQueryRead(d, meta)
}

func ansibleInventoryResourceQueryDelete(d *schema.ResourceData, meta interface{}) error {
	conf := meta.(ProviderConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	id := d.Id()

	// create a logger for this function
	logger, _ := util.CreateSubLogger("resource_host_delete")
	logger.Debug().Str("id", d.Id()).Msg("deleting inventory")

	i, err := conf.Inventories.GetOrCreate(id)
	if err != nil {
		logger.Error().Err(err).Msg("failed to load inventory")
		return fmt.Errorf("failed to load inventory '%s': %e", id, err)
	}

	conf.Mutex.Lock()
	if err := i.Delete(); err != nil {
		logger.Error().Err(err).Msg("failed to delete inventory")
		return fmt.Errorf("failed to delete inventory: %e", err)
	}
	if err := conf.Inventories.Delete(id); err != nil {
		logger.Error().Err(err).Msg("failed to remove inventory from internal memory")
		return fmt.Errorf("failed to remove inventory from internal memory: %e", err)
	}
	conf.Mutex.Unlock()
	return nil
}
