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

## 📊 Capacity Planning

How many users you can create, and how many of them can stay connected at the
same time, are two different limits:

- **Concurrent clients** is bound by the Pritunl node (CPU and network).
- **Total users** is bound by MongoDB (memory).

Most of the numbers below come from the official
[Pritunl scaling guide](https://docs.pritunl.com/kb/vpn/system/scaling);
`t3.micro` and `t3.small` are our own estimates.

### Concurrent clients per node

| Size class | AWS         | Google Cloud   | Resources                   | Max concurrent clients |
| ---------- | ----------- | -------------- | --------------------------- | ---------------------- |
| micro      | `t3.micro`  | `e2-micro`     | 2 vCPU (burstable), low net | ~60 (estimate)¹        |
| small      | `t3.small`  | `e2-small`     | 2 vCPU (burstable), low net | ~125 (estimate)¹       |
| medium     | `t3.medium` | `n1-highcpu-2` | 2 vCPU, low to moderate net | 250                    |
| large      | `c5.large`  | `n1-highcpu-4` | 2–4 vCPU, up to 10 Gigabit  | 1000                   |
| xlarge     | `c5.xlarge` | `n1-highcpu-8` | 4–8 vCPU, up to 10 Gigabit  | 2000                   |

Pritunl doesn't list anything below `t3.medium`, so these two are scaled down
from the documented 250 — worth checking under real load before you rely on
them. Both sizes suit dev environments more than production: their CPU is
burstable, so heavy OpenVPN traffic tends to eat the credit balance and slow
things down well before you hit the client count.

The GCP column is matched to the AWS row by max clients, not by vCPU count.
"Max clients" means concurrent connections, not the number of user accounts.

### Total users (MongoDB)

| AWS         | Google Cloud   | Memory     | Max users |
| ----------- | -------------- | ---------- | --------- |
| `r3.large`  | `n1-highmem-2` | 13–15 GB   | 20000     |
| `r3.xlarge` | `n1-highmem-4` | 26–30.5 GB | 40000     |

### What this means for this setup

This role keeps everything on one VM. The image bundles Pritunl, OpenVPN and
MongoDB in a single container, with data living in `./pritunl` and `./mongodb`
next to `vpn.yml` — see [`files/vpn.yml`](files/vpn.yml). There's no separate
database host and no second node.

So treat the client table as a ceiling rather than a target: the database, the
web console and OpenVPN all share the same vCPUs, memory and NIC. Pick a size
with some headroom over the client count you actually expect.

The MongoDB table assumes a dedicated database host, so it doesn't map onto this
setup directly. Here the user count is rarely what you run out of first — on a
single VM you'll feel the concurrent-client limit long before the database
becomes the bottleneck.

Pritunl's own answer for larger deployments is a separate MongoDB, several nodes
behind DNS names, and a low **Max Clients** per server so clients spill over to
another node instead of piling onto one. That's all worth reading if you outgrow
this, but none of it is what this role sets up — it would mean reworking
`vpn.yml` and running MongoDB somewhere else.

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
