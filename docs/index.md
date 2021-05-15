# Ansible provider for terraform
This document describes how to build and use the ansible terraform provider.
The provider lets you define ansible resources in terraform code, which generates
an inventory to be consumed by ansible during provisioning.

## Example usage

```terraform
terraform {
  required_providers {
    ansible = "~> 1.0.9"
  }
}

provider "ansible" {
  path            = "/data/ansible/inventory"
  log_enable      = false
  log_file        = "terraform-provider-ansible.log"
  log_levels      = {
    _default = "debug"
    _capturelog = ""
  }
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
}
```
