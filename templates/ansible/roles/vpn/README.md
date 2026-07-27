# VPN Pritunl Setup

## 📘 Overview

This Ansible role automates the setup of **Pritunl VPN**.  
It uses **Ansible Galaxy**, so make sure Ansible is already installed on your system.  
You can follow the official [Ansible installation guide](https://docs.ansible.com/ansible/latest/installation_guide/intro_installation.html).

## ⚙️ Prerequisites

Before running this role, ensure that:

1. You have configured your **inventory file (`inventory.ini`)** correctly.
2. The **host** defined in your inventory matches the one used in `initial-setup.yml`.
3. You have properly set up your **SSH key** to access the target host.

## 🚀 Usage

Run the following command to execute the playbook with the VPN setup tag:

```sh
ansible-playbook -i inventory.ini initial-setup.yml --tags vpn
```

## 📝 Notes

On Ubuntu 26.04 the legacy iptables kernel modules are no longer loaded by default
(the kernel defaults to nftables), which breaks Pritunl's NAT and filter rules.
This role loads them on the target host and persists the list in
`/etc/modules-load.d/pritunl-iptables.conf` so they survive a reboot.
