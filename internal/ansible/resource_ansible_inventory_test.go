package ansible

import (
	"fmt"
	"github.com/habakke/terraform-ansible-provider/internal/ansible/inventory"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"testing"
)

const TestGroupVarsData = `---
k3s_version: v1.19.5+k3s1
ansible_user: ubuntu
systemd_dir: /etc/systemd/system
master_ip: "{{ hostvars[groups['master'][0]]['ansible_host'] | default(groups['master'][0]) }}"
extra_server_args: ""
extra_agent_args: ""
`

const TestGroupVarsData2 = `---
k3s_version: v1.19.5+k3s1
ansible_user: ubuntu
systemd_dir: /etc/systemd/system
master_ip: "{{ hostvars[groups['master'][0]]['ansible_host'] | default(groups['master'][0]) }}"
extra_server_args: "--no-deploy traefik --flannel-backend wireguard"
extra_agent_args: ""
`

func TestAnsibleInventory_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAnsiblePreCheck(t) },
		Providers:    testAnsibleProviders,
		CheckDestroy: testAnsibleInventoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAnsibleInventoryBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAnsibleInventoryExists("ansible_inventory.cluster"),
					resource.TestCheckResourceAttrSet("ansible_inventory.cluster", "id"),
					resource.TestCheckResourceAttr("ansible_inventory.cluster", "group_vars", TestGroupVarsData),
				),
			},
		},
	})
}

func TestAnsibleInventory_Update(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAnsiblePreCheck(t) },
		Providers:    testAnsibleProviders,
		CheckDestroy: testAnsibleInventoryDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAnsibleInventoryBasic(),
				Check: resource.ComposeTestCheckFunc(
					testAnsibleInventoryExists("ansible_inventory.cluster"),
					resource.TestCheckResourceAttrSet("ansible_inventory.cluster", "id"),
					resource.TestCheckResourceAttr("ansible_inventory.cluster", "group_vars", TestGroupVarsData),
				),
			},
			{
				Config: testAnsibleInventoryUpdate(),
				Check: resource.ComposeTestCheckFunc(
					testAnsibleInventoryExists("ansible_inventory.cluster"),
					resource.TestCheckResourceAttrSet("ansible_inventory.cluster", "id"),
					resource.TestCheckResourceAttr("ansible_inventory.cluster", "group_vars", TestGroupVarsData2),
				),
			},
		},
	})
}

func testAnsibleInventoryDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ansible_inventory" {
			continue
		}

		if inventory.Exists(rs.Primary.ID) {
			return fmt.Errorf("inventory '%s' still exists", rs.Primary.ID)
		}
	}

	return nil
}

func testAnsibleInventoryBasic() string {
	return fmt.Sprintf(`
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
`)
}

func testAnsibleInventoryUpdate() string {
	return fmt.Sprintf(`
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
    extra_server_args: "--no-deploy traefik --flannel-backend wireguard"
    extra_agent_args: ""
  EOT
}
`)
}

func testAnsibleInventoryExists(resource string) resource.TestCheckFunc {
	return func(state *terraform.State) error {
		rs, ok := state.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("not found: %s", resource)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("no resource ID is set")
		}

		if !inventory.Exists(rs.Primary.ID) {
			return fmt.Errorf("inventory '%s' does not exist", rs.Primary.ID)

		}

		return nil
	}
}
