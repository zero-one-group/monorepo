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

The compose file must **not** pin the container DNS to `127.0.0.1`: nothing
listens there inside the container, so pritunl's public-IP detection fails and
generated client profiles contain `remote None <port> <proto>`, which clients
can never resolve (stuck on "Connecting"). Docker's embedded resolver
(`127.0.0.11`) forwards to the host DNS automatically, and pritunl then
auto-detects the public IP for the `remote` line in profiles.

After first boot, in the web UI (Server → Settings) **remove the DNS server**
(leave it empty) and **add the VPC CIDR prefix to the server routes** (e.g.
`172.31.0.0/16`) so VPN clients can reach your VPC.

On Ubuntu 26.04 the legacy iptables kernel modules are no longer loaded by default
(the kernel defaults to nftables), which breaks Pritunl's NAT and filter rules.
This role loads them on the target host and persists the list in
`/etc/modules-load.d/pritunl-iptables.conf` so they survive a reboot.
