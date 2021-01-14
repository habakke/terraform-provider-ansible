package ansible

import (
	"context"
	"fmt"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/inventory"
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
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"groupvars": {
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
	providerMeta := meta.(ProviderMetadata)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	groupVars := d.Get("groupvars").(string)

	i := inventory.NewInventory(providerMeta.Path)
	if err := i.Commit(groupVars); err != nil {
		return fmt.Errorf("failed to create inventory: %e", err)
	}

	d.SetId(i.GetId())
	d.MarkNewResource()
	return ansibleInventoryResourceQueryRead(d, meta)
}

func ansibleInventoryResourceQueryRead(d *schema.ResourceData, meta interface{}) error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	id := d.Id()
	i := inventory.LoadFromId(id)
	err, groupVars := i.Load()
	if err != nil {
		return fmt.Errorf("failed to load inventory '%s': %e", id, err)
	}

	_ = d.Set("groupvars", groupVars)

	return nil
}

func ansibleInventoryResourceQueryUpdate(d *schema.ResourceData, meta interface{}) error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	id := d.Id()
	groupVars := d.Get("groupvars").(string)
	i := inventory.LoadFromId(id)

	if d.HasChange("groupvars") {
		if err := i.Commit(groupVars); err != nil {
			return fmt.Errorf("failed to update inventory: %e", err)
		}
	}

	return ansibleInventoryResourceQueryRead(d, meta)
}

func ansibleInventoryResourceQueryDelete(d *schema.ResourceData, meta interface{}) error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	id := d.Id()
	i := inventory.LoadFromId(id)
	if err := i.Delete(); err != nil {
		return fmt.Errorf("failed to delete inventory: %e", err)
	}
	return nil
}
