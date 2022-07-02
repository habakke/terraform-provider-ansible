//go:build integration

package ansible

import (
	"fmt"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/database"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/inventory"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"testing"
)

func TestAnsibleGroup_Basic(t *testing.T) {
	resourceName := "ansible_group.master"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAnsiblePreCheck(t, resourceName) },
		Providers:    testAnsibleProviders,
		CheckDestroy: testAnsibleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAnsibleGroupBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAnsibleGroupExists("ansible_group.master"),
					resource.TestCheckResourceAttr("ansible_group.master", "name", "master"),
					resource.TestCheckResourceAttrSet("ansible_group.master", "inventory"),
				),
			},
		},
	})
}

func TestAnsibleGroup_Update(t *testing.T) {
	resourceName := "ansible_group.master"
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAnsiblePreCheck(t, resourceName) },
		Providers:    testAnsibleProviders,
		CheckDestroy: testAnsibleGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAnsibleGroupBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAnsibleGroupExists("ansible_group.master"),
					resource.TestCheckResourceAttr("ansible_group.master", "name", "master"),
					resource.TestCheckResourceAttrSet("ansible_group.master", "inventory"),
				),
			},
			{
				Config: testAnsibleGroupUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAnsibleGroupExists("ansible_group.master"),
					resource.TestCheckResourceAttr("ansible_group.master", "name", "master2"),
					resource.TestCheckResourceAttrSet("ansible_group.master", "inventory"),
				),
			},
		},
	})
}

func groupExists(groupID string, rootPath string, inventoryRef string) bool {
	i, err := inventory.Load(rootPath, inventoryRef)
	if err != nil {
		return false
	}

	db := database.NewDatabase(i.GetInventoryPath())
	if !db.Exists() {
		return false
	}

	_ = db.Load()
	g := db.Group(groupID)
	return g != nil
}

func testAnsibleGroupDestroy(s *terraform.State) error {
	var gid *string = nil
	var inventoryRef string
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "ansible_group" && rs.Primary.Attributes["name"] == "master" || rs.Primary.Attributes["name"] == "master2" {
			gid = &rs.Primary.ID
			inventoryRef = rs.Primary.Attributes["inventory"]
		}
	}

	if gid == nil {
		return fmt.Errorf("Unable to find group 'master'")
	}

	if groupExists(*gid, "/tmp", inventoryRef) {
		return fmt.Errorf("group '%s' still exists", *gid)
	}

	return nil
}

func testAnsibleGroupBasic() string {
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
`
}

func testAnsibleGroupUpdate() string {
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
`
}

func testAnsibleGroupExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource ID is set")
		}
		if !groupExists(rs.Primary.ID, "/tmp", rs.Primary.Attributes["inventory"]) {
			return fmt.Errorf("group '%s' does not exist", rs.Primary.ID)

		}
		return nil
	}
}
