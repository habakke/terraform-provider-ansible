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

func ansibleGroupResourceQuery() *schema.Resource {
	return &schema.Resource{
		Create: ansibleGroupResourceQueryCreate,
		Read:   ansibleGroupResourceQueryRead,
		Update: ansibleGroupResourceQueryUpdate,
		Delete: ansibleGroupResourceQueryDelete,
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
		},
	}
}

func ansibleGroupResourceQueryCreate(d *schema.ResourceData, meta interface{}) error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	name := d.Get("name").(string)
	inventoryRef := d.Get("inventory").(string)
	i := inventory.LoadFromId(inventoryRef)
	db := database.NewDatabase(i.GetDatabasePath())

	g := database.NewGroup(name)
	if err := db.AddGroup(*g); err != nil {
		return fmt.Errorf("failed to add group '%s': %e", name, err)
	}

	// Save and export database
	if err := commitAndExport(db, i.GetDatabasePath()); err != nil {
		return err
	}

	d.SetId(g.GetId())
	d.MarkNewResource()
	return ansibleGroupResourceQueryRead(d, meta)
}

func ansibleGroupResourceQueryRead(d *schema.ResourceData, meta interface{}) error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	inventoryRef := d.Get("inventory").(string)
	i := inventory.LoadFromId(inventoryRef)
	db := database.NewDatabase(i.GetDatabasePath())

	if err := db.Load(); err != nil {
		return fmt.Errorf("failed to load groups from temporary database: %e", err)
	}

	g := db.Group(d.Id())
	if g == nil {
		return fmt.Errorf("unable to find group '%s'", d.Id())
	}

	_ = d.Set("name", g.GetName())

	return nil
}

func ansibleGroupResourceQueryUpdate(d *schema.ResourceData, meta interface{}) error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	id := d.Id()
	name := d.Get("name").(string)
	inventoryRef := d.Get("inventory").(string)
	i := inventory.LoadFromId(inventoryRef)
	db := database.NewDatabase(i.GetDatabasePath())

	if err := db.Load(); err != nil {
		return fmt.Errorf("failed to load groups from temporary database: %e", err)
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

	return ansibleGroupResourceQueryRead(d, meta)
}

func ansibleGroupResourceQueryDelete(d *schema.ResourceData, meta interface{}) error {
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	inventoryRef := d.Get("inventory").(string)
	i := inventory.LoadFromId(inventoryRef)
	db := database.NewDatabase(i.GetDatabasePath())

	if err := db.Load(); err != nil {
		return fmt.Errorf("failed to load groups from temporary database: %e", err)
	}

	id := d.Id()
	g := db.Group(id)
	if g == nil {
		return fmt.Errorf("unable to find group with id '%s'", d.Id())
	}

	if err := db.RemoveGroup(*g); err != nil {
		return fmt.Errorf("unable to delete group '%s': %e", g.GetName(), err)
	}

	// Save and export database
	if err := commitAndExport(db, i.GetDatabasePath()); err != nil {
		return err
	}

	return nil
}
