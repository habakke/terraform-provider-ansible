package ansible

import (
	"context"
	"fmt"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/database"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/inventory"
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
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
		},
	}
}

func ansibleHostResourceQueryCreate(d *schema.ResourceData, meta interface{}) error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	name := d.Get("name").(string)
	groupId := d.Get("group").(string)
	inventoryRef := d.Get("inventory").(string)

	i := inventory.LoadFromId(inventoryRef)
	db := database.NewDatabase(i.GetDatabasePath())

	if err := db.Load(); err != nil {
		return fmt.Errorf("failed to load groups from temporary database: %e", err)
	}

	g := db.Group(groupId)
	if g == nil {
		return fmt.Errorf("unable to find group '%s'", groupId)
	}

	h := database.NewHost(name)
	g.UpdateEntity(h)
	db.UpdateGroup(*g)

	// Save and export database
	if err := commitAndExport(db, i.GetDatabasePath()); err != nil {
		return err
	}

	d.SetId(h.GetId())
	d.MarkNewResource()
	return ansibleHostResourceQueryRead(d, meta)
}

func ansibleHostResourceQueryRead(d *schema.ResourceData, meta interface{}) error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	inventoryRef := d.Get("inventory").(string)
	i := inventory.LoadFromId(inventoryRef)
	db := database.NewDatabase(i.GetDatabasePath())

	if err := db.Load(); err != nil {
		return fmt.Errorf("failed to load groups from temporary database: %e", err)
	}

	id := d.Id()
	g, entry, err := db.FindEntryById(id)
	if err != nil {
		return fmt.Errorf("unable to find entry '%s': %e", id, err)
	}

	_ = d.Set("name", (*entry).GetName())
	_ = d.Set("group", g.GetId())

	return nil
}

func ansibleHostResourceQueryUpdate(d *schema.ResourceData, meta interface{}) error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	name := d.Get("name").(string)
	groupId := d.Get("group").(string)

	inventoryRef := d.Get("inventory").(string)
	i := inventory.LoadFromId(inventoryRef)
	db := database.NewDatabase(i.GetDatabasePath())

	if err := db.Load(); err != nil {
		return fmt.Errorf("failed to load groups from temporary database: %e", err)
	}

	g, entry, err := db.FindEntryById(d.Id())
	if err != nil {
		return fmt.Errorf("unable to find entry '%s': %e", d.Id(), err)
	}

	// check if name has changed
	if d.HasChange("name") {
		(*entry).SetName(name)
		db.UpdateGroup(*g)
	}

	// check if group has changed
	if d.HasChange("group") {
		// remove host from old group
		if err := g.RemoveEntity(*entry); err != nil {
			return fmt.Errorf("failed remove entry from group '%s': %e", g.GetId(), err)
		}
		db.UpdateGroup(*g)

		// load new group
		ng := db.Group(groupId)
		if ng == nil {
			return fmt.Errorf("failed to locate group '%s': %e", groupId, err)
		}

		// update name and add entity to new group
		ng.UpdateEntity(*entry)
	}

	if d.HasChanges("name", "group") {
		// Save and export database
		if err := commitAndExport(db, i.GetDatabasePath()); err != nil {
			return err
		}
	}

	return ansibleHostResourceQueryRead(d, meta)
}

func ansibleHostResourceQueryDelete(d *schema.ResourceData, meta interface{}) error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	inventoryRef := d.Get("inventory").(string)
	i := inventory.LoadFromId(inventoryRef)
	db := database.NewDatabase(i.GetDatabasePath())

	if err := db.Load(); err != nil {
		return fmt.Errorf("failed to load groups from temporary database: %e", err)
	}

	id := d.Id()
	g, entry, err := db.FindEntryById(id)
	if err != nil {
		return fmt.Errorf("unable to find entry with id '%s': %e", id, err)
	}

	// remove entry from group
	if err := g.RemoveEntity(*entry); err != nil {
		return fmt.Errorf("unable to remove entry from group with id: %e", err)
	}

	// update group
	db.UpdateGroup(*g)

	// Save and export database
	if err := commitAndExport(db, i.GetDatabasePath()); err != nil {
		return err
	}

	return nil
}
