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

func ansibleGroupResourceQuery() *schema.Resource {
	return &schema.Resource{
		Create: ansibleGroupResourceQueryCreate,
		Read:   ansibleGroupResourceQueryRead,
		Update: ansibleGroupResourceQueryUpdate,
		Delete: ansibleGroupResourceQueryDelete,
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

func ansibleGroupResourceQueryCreate(d *schema.ResourceData, meta interface{}) error {
	conf := meta.(providerConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	name := d.Get("name").(string)
	inventoryRef := d.Get("inventory").(string)

	// create a logger for this function
	logger, _ := util.CreateSubLogger("resource_group_create")
	logger.Debug().Str("name", name).Str("inventory", inventoryRef).Msg("invoking creation of group")

	conf.Mutex.Lock()
	i := inventory.LoadFromID(inventoryRef)
	db, err := i.GetAndLoadDatabase()
	if err != nil {
		logger.Error().Err(err).Msg("failed to load database")
		return fmt.Errorf("failed to load database '%s': %s", inventoryRef, err.Error())
	}
	g := database.NewGroup(name)
	if err := db.AddGroup(*g); err != nil {
		return fmt.Errorf("failed to add group '%s': %s", name, err.Error())
	}

	// Save and export database
	if err := commitAndExport(db, i.GetDatabasePath()); err != nil {
		return err
	}
	conf.Mutex.Unlock()

	d.SetId(g.GetID())
	d.MarkNewResource()
	return ansibleGroupResourceQueryRead(d, meta)
}

func ansibleGroupResourceQueryRead(d *schema.ResourceData, meta interface{}) error {
	conf := meta.(providerConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	inventoryRef := d.Get("inventory").(string)

	// create a logger for this function
	logger, _ := util.CreateSubLogger("resource_group_read")
	logger.Debug().Str("id", d.Id()).Str("inventory", inventoryRef).Msg("reading configuration for group")

	conf.Mutex.Lock()
	i := inventory.LoadFromID(inventoryRef)
	db, err := i.GetAndLoadDatabase()
	conf.Mutex.Unlock()
	if err != nil {
		logger.Error().Err(err).Msg("failed to load database")
		return fmt.Errorf("failed to load database '%s': %s", inventoryRef, err.Error())
	}

	g := db.Group(d.Id())
	if g == nil {
		return fmt.Errorf("unable to find group '%s'", d.Id())
	}

	_ = d.Set("name", g.GetName())

	return nil
}

func ansibleGroupResourceQueryUpdate(d *schema.ResourceData, meta interface{}) error {
	conf := meta.(providerConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	id := d.Id()
	name := d.Get("name").(string)
	inventoryRef := d.Get("inventory").(string)

	// create a logger for this function
	logger, _ := util.CreateSubLogger("resource_group_update")
	logger.Debug().Str("id", d.Id()).Str("name", name).Str("inventory", inventoryRef).Msg("updating configuration for group")

	conf.Mutex.Lock()
	i := inventory.LoadFromID(inventoryRef)
	db, err := i.GetAndLoadDatabase()
	if err != nil {
		logger.Error().Err(err).Msg("failed to load database")
		return fmt.Errorf("failed to load database '%s': %s", inventoryRef, err.Error())
	}

	g := db.Group(d.Id())
	if g == nil {
		return fmt.Errorf("unable to group with id '%s'", id)
	}

	if d.HasChange("name") {
		g.SetName(name)
		db.UpdateGroup(*g)

		// Save and export database
		if err := commitAndExport(db, i.GetDatabasePath()); err != nil {
			return err
		}
	}
	conf.Mutex.Unlock()

	return ansibleGroupResourceQueryRead(d, meta)
}

func ansibleGroupResourceQueryDelete(d *schema.ResourceData, meta interface{}) error {
	conf := meta.(providerConfiguration)
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	inventoryRef := d.Get("inventory").(string)

	// create a logger for this function
	logger, _ := util.CreateSubLogger("resource_group_delete")
	logger.Debug().Str("id", d.Id()).Str("inventory", inventoryRef).Msg("deleting group")

	conf.Mutex.Lock()
	i := inventory.LoadFromID(inventoryRef)
	db, err := i.GetAndLoadDatabase()
	if err != nil {
		logger.Error().Err(err).Msg("failed to load database")
		return fmt.Errorf("failed to load database '%s': %s", inventoryRef, err.Error())
	}

	id := d.Id()
	g := db.Group(id)
	if g == nil {
		logger.Error().Err(err).Msg("cannot find group so unable to remove, but continuing anyway")
	} else {
		// if we find the group we remove it, if we can't find it we can skip removing it
		if err := db.RemoveGroup(*g); err != nil {
			return fmt.Errorf("unable to delete group '%s': %s", g.GetName(), err.Error())
		}
	}

	// Save and export database
	if err := commitAndExport(db, i.GetDatabasePath()); err != nil {
		return err
	}
	conf.Mutex.Unlock()

	return nil
}
