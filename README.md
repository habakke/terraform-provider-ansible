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
  variables = {
    name = "k3s-master-1"
    role = "master"
  }
}
```

## How to build
To build an test the plugin locally first create a `~/.terraformrc` file

```shell
provider_installation {

  dev_overrides {
    "habakke/ansible" = "/Users/habakke/.terraform.d/plugins"
  }
  direct {}
}
```

Then build and install the plugin locally using

```shell
make install
```

## Running tests
To run the internal unit tests run test `test` make target

```shell
make test
```

To run terraform acceptance tests, the `TF_ACC` env variable must be set to true before making the
`test` make target, or the `testacc` make target can be used

```shell
make testacc
```

## TODO

* Add support for Ansible group variables (https://docs.ansible.com/ansible/latest/user_guide/intro_inventory.html#assigning-a-variable-to-many-machines-group-variables)
* Add proper docs as seen in other community providers (https://github.com/paultyng/terraform-provider-unifi/tree/main/docs)
* Upgrade to Terraform Plugin SDK v2 (https://www.terraform.io/docs/extend/guides/v2-upgrade-guide.html)
