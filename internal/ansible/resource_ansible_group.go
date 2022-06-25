package ansible

import (
	"context"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/database"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/inventory"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/rs/zerolog/log"
	"time"
)

func ansibleGroupResourceQuery() *schema.Resource {
	return &schema.Resource{
		CreateContext: ansibleGroupResourceQueryCreate,
		ReadContext:   ansibleGroupResourceQueryRead,
		UpdateContext: ansibleGroupResourceQueryUpdate,
		DeleteContext: ansibleGroupResourceQueryDelete,
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Second),
			Update: schema.DefaultTimeout(10 * time.Second),
			Delete: schema.DefaultTimeout(10 * time.Second),
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"inventory": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
		},
	}
}

func ansibleGroupResourceQueryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(providerConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	name := d.Get("name").(string)
	inventoryRef := d.Get("inventory").(string)

	conf.Mutex.Lock()
	i := inventory.LoadFromID(inventoryRef)
	db, err := i.GetAndLoadDatabase()
	if err != nil {
		return diag.Errorf("failed to load database '%s': %s", inventoryRef, err.Error())
	}
	g := database.NewGroup(name)
	if err := db.AddGroup(*g); err != nil {
		return diag.Errorf("failed to add group '%s': %s", name, err.Error())
	}

	// Save and export database
	if err := commitAndExport(db, i.GetDatabasePath()); err != nil {
		return diag.FromErr(err)
	}
	conf.Mutex.Unlock()

	d.SetId(g.GetID())
	d.MarkNewResource()
	return ansibleGroupResourceQueryRead(ctx, d, meta)
}

func ansibleGroupResourceQueryRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conf := meta.(providerConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	inventoryRef := d.Get("inventory").(string)

	conf.Mutex.Lock()
	i := inventory.LoadFromID(inventoryRef)
	db, err := i.GetAndLoadDatabase()
	conf.Mutex.Unlock()
	if err != nil {
		return diag.Errorf("failed to load database '%s': %s", inventoryRef, err.Error())
	}

	g := db.Group(d.Id())
	if g == nil {
		return diag.Errorf("unable to find group '%s'", d.Id())
	}

	_ = d.Set("name", g.GetName())

	return diags
}

func ansibleGroupResourceQueryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(providerConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	id := d.Id()
	name := d.Get("name").(string)
	inventoryRef := d.Get("inventory").(string)

	conf.Mutex.Lock()
	i := inventory.LoadFromID(inventoryRef)
	db, err := i.GetAndLoadDatabase()
	if err != nil {
		log.Error().Err(err).Msg("failed to load database")
		return diag.Errorf("failed to load database '%s': %s", inventoryRef, err.Error())
	}

	g := db.Group(d.Id())
	if g == nil {
		return diag.Errorf("unable to group with id '%s'", id)
	}

	if d.HasChange("name") {
		g.SetName(name)
		db.UpdateGroup(*g)

		// Save and export database
		if err := commitAndExport(db, i.GetDatabasePath()); err != nil {
			return diag.FromErr(err)
		}
	}
	conf.Mutex.Unlock()

	return ansibleGroupResourceQueryRead(ctx, d, meta)
}

func ansibleGroupResourceQueryDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conf := meta.(providerConfiguration)
	_, cancel := context.WithCancel(ctx)
	defer cancel()

	inventoryRef := d.Get("inventory").(string)

	log.Debug().Str("id", d.Id()).Str("inventory", inventoryRef).Msg("deleting group")

	conf.Mutex.Lock()
	i := inventory.LoadFromID(inventoryRef)
	db, err := i.GetAndLoadDatabase()
	if err != nil {
		log.Error().Err(err).Msg("failed to load database")
		return diag.Errorf("failed to load database '%s': %s", inventoryRef, err.Error())
	}

	id := d.Id()
	g := db.Group(id)
	if g == nil {
		log.Error().Err(err).Msg("cannot find group so unable to remove, but continuing anyway")
	} else {
		// if we find the group we remove it, if we can't find it we can skip removing it
		if err := db.RemoveGroup(*g); err != nil {
			return diag.Errorf("unable to delete group '%s': %s", g.GetName(), err.Error())
		}
	}

	// Save and export database
	if err := commitAndExport(db, i.GetDatabasePath()); err != nil {
		return diag.FromErr(err)
	}
	conf.Mutex.Unlock()

	return diags
}
