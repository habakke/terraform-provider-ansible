package ansible

import (
	"context"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/inventory"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rs/zerolog/log"
	"time"
)

func ansibleInventoryResourceQuery() *schema.Resource {
	return &schema.Resource{
		CreateContext: ansibleInventoryResourceQueryCreate,
		ReadContext:   ansibleInventoryResourceQueryRead,
		UpdateContext: ansibleInventoryResourceQueryUpdate,
		DeleteContext: ansibleInventoryResourceQueryDelete,
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

func ansibleInventoryResourceQueryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(providerConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	groupVars := d.Get("group_vars").(string)

	conf.Mutex.Lock()
	i := inventory.NewInventory(conf.Path)
	log.Debug().Str("id", i.GetID()).Msg("created new inventory")
	if err := i.Commit(groupVars); err != nil {
		return diag.Errorf("failed to commit inventory: %s", err.Error())
	}
	conf.Mutex.Unlock()

	d.SetId(i.GetID())
	d.MarkNewResource()
	return ansibleInventoryResourceQueryRead(ctx, d, meta)
}

func ansibleInventoryResourceQueryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conf := meta.(providerConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	id := d.Id()

	conf.Mutex.Lock()
	i := inventory.LoadFromID(id)
	groupVars, err := i.Load()
	if err != nil {
		log.Error().Err(err).Msg("failed to load inventory")
		return diag.Errorf("failed to load inventory '%s': %s", id, err.Error())
	}
	conf.Mutex.Unlock()

	_ = d.Set("group_vars", groupVars)

	return diags
}

func ansibleInventoryResourceQueryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(providerConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	id := d.Id()
	groupVars := d.Get("group_vars").(string)

	i := inventory.LoadFromID(id)
	if d.HasChange("group_vars") {
		conf.Mutex.Lock()
		if err := i.Commit(groupVars); err != nil {
			log.Error().Err(err).Msg("failed to update inventory")
			return diag.Errorf("failed to update inventory: %s", err.Error())
		}
		conf.Mutex.Unlock()
	}

	return ansibleInventoryResourceQueryRead(ctx, d, meta)
}

func ansibleInventoryResourceQueryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conf := meta.(providerConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	id := d.Id()

	conf.Mutex.Lock()
	i := inventory.LoadFromID(id)
	if err := i.Delete(); err != nil {
		log.Error().Err(err).Msg("failed to delete inventory")
		return diag.Errorf("failed to delete inventory: %s", err.Error())
	}
	conf.Mutex.Unlock()
	return diags
}
