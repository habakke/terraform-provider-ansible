# Ansible provider for terraform
This documents describes how to build and use the ansible terraform provider. 
The provider lets you define ansible resources in terraform code, which generates
a inventory to be consumed by ansible during provisioning.

## How to use

```terraform
terraform {
  required_providers {
    ansible = "~> 0.0.1"
  }
}

provider "ansible" {
  path = "/data/ansible/inventory"
}

resource "ansible_inventory" "cluster" {
  group-vars = <<-EOT
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

## How to build

```shell
make build-dev version=v0.0.1
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

* Make provider available publicly at registry.hashicorp.com (https://www.terraform.io/guides/terraform-provider-development-program.html)
