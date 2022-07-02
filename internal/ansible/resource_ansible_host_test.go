package ansible

import (
	"fmt"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/database"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/inventory"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAnsibleHost_Basic(t *testing.T) {
	resourceName := "ansible_host.k3s-master-1"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAnsiblePreCheck(t, resourceName) },
		Providers:    testAnsibleProviders,
		CheckDestroy: testAnsibleHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAnsibleHostBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAnsibleHostExists("ansible_host.k3s-master-1"),
					resource.TestCheckResourceAttr("ansible_host.k3s-master-1", "name", "k3s-master-1"),
					resource.TestCheckResourceAttrSet("ansible_host.k3s-master-1", "group"),
					resource.TestCheckResourceAttrSet("ansible_host.k3s-master-1", "inventory"),
					resource.TestCheckResourceAttr("ansible_host.k3s-master-1", "variables.name", "k3s-master-1"),
					resource.TestCheckResourceAttr("ansible_host.k3s-master-1", "variables.role", "master"),
				),
			},
		},
	})
}

func TestAnsibleHost_Update(t *testing.T) {
	resourceName := "ansible_host.k3s-master-1"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAnsiblePreCheck(t, resourceName) },
		Providers:    testAnsibleProviders,
		CheckDestroy: testAnsibleHostDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAnsibleHostBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAnsibleHostExists("ansible_host.k3s-master-1"),
					resource.TestCheckResourceAttr("ansible_host.k3s-master-1", "name", "k3s-master-1"),
					resource.TestCheckResourceAttrSet("ansible_host.k3s-master-1", "group"),
					resource.TestCheckResourceAttrSet("ansible_host.k3s-master-1", "inventory"),
					resource.TestCheckResourceAttr("ansible_host.k3s-master-1", "variables.name", "k3s-master-1"),
					resource.TestCheckResourceAttr("ansible_host.k3s-master-1", "variables.role", "master"),
				),
			},
			{
				Config: testAnsibleHostUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAnsibleHostExists("ansible_host.k3s-master-1"),
					resource.TestCheckResourceAttr("ansible_host.k3s-master-1", "name", "k3s-master-1-edit"),
					resource.TestCheckResourceAttrSet("ansible_host.k3s-master-1", "group"),
					resource.TestCheckResourceAttrSet("ansible_host.k3s-master-1", "inventory"),
					resource.TestCheckResourceAttr("ansible_host.k3s-master-1", "variables.name", "k3s-master-1-edit"),
					resource.TestCheckResourceAttr("ansible_host.k3s-master-1", "variables.role", "master"),
				),
			},
		},
	})
}

func hostExists(hostID string, rootPath string, inventoryRef string, groupID string) bool {
	i, err := inventory.Load(rootPath, inventoryRef)
	if err != nil {
		return false
	}
	db := database.NewDatabase(i.GetInventoryPath())
	if !db.Exists() {
		return false
	}

	_ = db.Load()
	g, e, err := db.FindEntryByID(hostID)
	if err != nil {
		return false
	}

	return (e != nil) && (g.GetID() == groupID)
}

func testAnsibleHostDestroy(s *terraform.State) error {
	var id *string
	var groupID string
	var inventoryRef string
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "ansible_host" && rs.Primary.Attributes["name"] == "k3s-master-1" || rs.Primary.Attributes["name"] == "k3s-master-1-edit" {
			id = &rs.Primary.ID
			groupID = rs.Primary.Attributes["group"]
			inventoryRef = rs.Primary.Attributes["inventory"]
		}
	}

	if id == nil {
		return fmt.Errorf("unable to find host 'k3s-master-1'")
	}

	if hostExists(*id, "/tmp", inventoryRef, groupID) {
		return fmt.Errorf("host '%s' still exists", *id)
	}

	return nil
}

func testAnsibleHostBasic() string {
	return `
provider "ansible" {
  path = "/tmp"
}

resource "ansible_inventory" "cluster" {
  group_vars = <<-EOT
    ---
    k3s_version: v1.19.5+k3s1
    ansible_user: ubuntu
    systemd_dir: /etc/systemd/system
    master_ip: "{{ hostvars[groups['master'][0]]['ansible_host'] | default(groups['master'][0]) }}"
    extra_server_args: ""
    extra_agent_args: ""
  EOT
}

resource "ansible_group" "master" {
  depends_on = [ansible_inventory.cluster]
  name = "master"
  inventory = ansible_inventory.cluster.id
}

resource "ansible_host" "k3s-master-1" {
  depends_on = [ansible_group.master]
  name = "k3s-master-1"
  inventory = ansible_inventory.cluster.id
  group = ansible_group.master.id
  variables = {
    name = "k3s-master-1"
    role = "master"
  }
}
`
}

func testAnsibleHostUpdate() string {
	return `
provider "ansible" {
  path = "/tmp"
}

resource "ansible_inventory" "cluster" {
  group_vars = <<-EOT
    ---
    k3s_version: v1.19.5+k3s1
    ansible_user: ubuntu
    systemd_dir: /etc/systemd/system
    master_ip: "{{ hostvars[groups['master'][0]]['ansible_host'] | default(groups['master'][0]) }}"
    extra_server_args: ""
    extra_agent_args: ""
  EOT
}

resource "ansible_group" "master" {
  depends_on = [ansible_inventory.cluster]
  name = "master2"
  inventory = ansible_inventory.cluster.id
}

resource "ansible_host" "k3s-master-1" {
  depends_on = [ansible_group.master]
  name = "k3s-master-1-edit"
  inventory = ansible_inventory.cluster.id
  group = ansible_group.master.id
  variables = {
    name = "k3s-master-1-edit"
    role = "master"
  }
}
`
}

func testAnsibleHostExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource ID is set")
		}

		if !hostExists(rs.Primary.ID, "/tmp", rs.Primary.Attributes["inventory"], rs.Primary.Attributes["group"]) {
			return fmt.Errorf("group '%s' does not exist", rs.Primary.ID)
		}
		return nil
	}
}
