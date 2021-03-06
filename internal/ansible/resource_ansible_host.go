package ansible

import (
	"context"
	"fmt"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/database"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/inventory"
	"github.com/habakke/terraform-ansible-provider/internal/util"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"time"
)

func ansibleHostResourceQuery() *schema.Resource {
	return &schema.Resource{
		Create: ansibleHostResourceQueryCreate,
		Read:   ansibleHostResourceQueryRead,
		Update: ansibleHostResourceQueryUpdate,
		Delete: ansibleHostResourceQueryDelete,
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
			"group": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.NoZeroValues,
			},
			"variables": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					Required:     true,
					ValidateFunc: validation.NoZeroValues,
				},
			},
		},
	}
}

func ansibleHostResourceQueryCreate(d *schema.ResourceData, meta interface{}) error {
	conf := meta.(providerConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	name := d.Get("name").(string)
	groupID := d.Get("group").(string)
	inventoryRef := d.Get("inventory").(string)
	variables := d.Get("variables").(map[string]interface{})

	// create a logger for this function
	logger, _ := util.CreateSubLogger("resource_host_create")
	logger.Debug().Str("name", name).Str("group", groupID).Str("inventory", inventoryRef).Msg("invoking creation of host")

	conf.Mutex.Lock()
	i := inventory.LoadFromID(inventoryRef)
	db, err := i.GetAndLoadDatabase()
	if err != nil {
		logger.Error().Err(err).Msg("failed to load database")
		return fmt.Errorf("failed to load database '%s': %e", inventoryRef, err)
	}

	g := db.Group(groupID)
	if g == nil {
		return fmt.Errorf("unable to find group '%s'", groupID)
	}

	h := database.NewHost(name, variables)
	g.UpdateEntity(h)
	db.UpdateGroup(*g)

	// Save and export database
	if err := commitAndExport(db, i.GetDatabasePath()); err != nil {
		return err
	}
	conf.Mutex.Unlock()

	d.SetId(h.GetID())
	d.MarkNewResource()
	return ansibleHostResourceQueryRead(d, meta)
}

func ansibleHostResourceQueryRead(d *schema.ResourceData, meta interface{}) error {
	conf := meta.(providerConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	inventoryRef := d.Get("inventory").(string)

	// create a logger for this function
	logger, _ := util.CreateSubLogger("resource_host_read")
	logger.Debug().Str("id", d.Id()).Str("inventory", inventoryRef).Msg("reading configuration for host")

	conf.Mutex.Lock()
	i := inventory.LoadFromID(inventoryRef)
	db, err := i.GetAndLoadDatabase()
	conf.Mutex.Unlock()
	if err != nil {
		logger.Error().Err(err).Msg("failed to load database")
		return fmt.Errorf("failed to load database '%s': %e", inventoryRef, err)
	}

	id := d.Id()
	g, entry, err := db.FindEntryByID(id)
	if err != nil {
		return fmt.Errorf("unable to find entry '%s': %e", id, err)
	}

	_ = d.Set("name", entry.GetName())
	_ = d.Set("group", g.GetID())

	h, ok := entry.(*database.Host)
	if ok {
		_ = d.Set("variables", h.GetVariables())
	}
	return nil
}

func ansibleHostResourceQueryUpdate(d *schema.ResourceData, meta interface{}) error {
	conf := meta.(providerConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	name := d.Get("name").(string)
	groupID := d.Get("group").(string)
	inventoryRef := d.Get("inventory").(string)
	variables := d.Get("variables").(map[string]interface{})

	// create a logger for this function
	logger, _ := util.CreateSubLogger("resource_host_update")
	logger.Debug().Str("id", d.Id()).Str("group", groupID).Str("inventory", inventoryRef).Msg("updating configuration for host")

	conf.Mutex.Lock()
	i := inventory.LoadFromID(inventoryRef)
	db, err := i.GetAndLoadDatabase()
	if err != nil {
		logger.Error().Err(err).Msg("failed to load database")
		return fmt.Errorf("failed to load database '%s': %e", inventoryRef, err)
	}

	g, entry, err := db.FindEntryByID(d.Id())
	if err != nil {
		return fmt.Errorf("unable to find entry '%s': %e", d.Id(), err)
	}

	// check if name has changed
	if d.HasChange("name") {
		entry.SetName(name)
		db.UpdateGroup(*g)
	}

	// check if group has changed
	if d.HasChange("group") {
		// remove host from old group
		if err := g.RemoveEntity(entry); err != nil {
			return fmt.Errorf("failed remove entry from group '%s': %e", g.GetID(), err)
		}
		db.UpdateGroup(*g)

		// load new group
		ng := db.Group(groupID)
		if ng == nil {
			return fmt.Errorf("failed to locate group '%s': %e", groupID, err)
		}

		// update name and add entity to new group
		ng.UpdateEntity(entry)
	}

	if d.HasChange("variables") {
		h, ok := entry.(*database.Host)
		if ok {
			for k := range variables {
				h.SetVariable(k, variables[k])
			}
		}
		db.UpdateGroup(*g)
	}

	if d.HasChanges("name", "group", "variables") {
		// Save and export database
		if err := commitAndExport(db, i.GetDatabasePath()); err != nil {
			return err
		}
	}

	conf.Mutex.Unlock()

	return ansibleHostResourceQueryRead(d, meta)
}

func ansibleHostResourceQueryDelete(d *schema.ResourceData, meta interface{}) error {
	conf := meta.(providerConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	inventoryRef := d.Get("inventory").(string)

	// create a logger for this function
	logger, _ := util.CreateSubLogger("resource_host_delete")
	logger.Debug().Str("id", d.Id()).Str("inventory", inventoryRef).Msg("deleting host")

	conf.Mutex.Lock()
	i := inventory.LoadFromID(inventoryRef)
	db, err := i.GetAndLoadDatabase()
	if err != nil {
		logger.Error().Err(err).Msg("failed to load database")
		return fmt.Errorf("failed to load database '%s': %e", inventoryRef, err)
	}

	id := d.Id()
	g, entry, err := db.FindEntryByID(id)
	if err != nil {
		logger.Error().Err(err).Msg("cannot find host so unable to remove, but continuing anyway")
	} else {
		// only remove host from group if we actually find it there. if we dont find it, then everything is ok and we
		// can skip the removing it.

		// remove entry from group
		if err := g.RemoveEntity(entry); err != nil {
			return fmt.Errorf("unable to remove entry from group with id: %e", err)
		}

		// update group
		db.UpdateGroup(*g)
	}

	// Save and export database
	if err := commitAndExport(db, i.GetDatabasePath()); err != nil {
		return err
	}
	conf.Mutex.Unlock()

	return nil
}
