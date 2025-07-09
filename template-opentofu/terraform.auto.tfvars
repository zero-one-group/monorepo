instance_type = {
  "bastion" : "t3a.micro"
  "master" : "t3a.small"
  "worker" : "t3a.small"
  "master-dev-staging" : "t3a.micro"
  "worker-dev-staging" : "t3a.micro"
  "monitoring" : "t3a.medium"
  "database" : "t3a.medium"
  "database-dev-staging" : "t3a.small"
  "nginx" : "t3a.medium"
  "nginx-dev-staging" : "t3a.small"
}
project          = "lgtm"
public_key_path  = "~/.ssh/id_rsa.pub"
region           = "ap-southeast-1"
ubuntu_jammy_ami = "ami-0e0ddf453092e1e37"
ubuntu_noble_ami = "ami-0b874c2ac1b5e9957"
user_data_ec2 = {
  "bastion" : "user-data.sh",
  "master" : "user-data-master.sh",
  "master_join" : "user-data-master-join.sh",
  "monitoring" : "user-data-monitoring.sh",
  "worker" : "user-data-worker.sh"
  "nginx" : "user-data-nginx.sh"
}
volume = {
  "bastion" : 10,
  "master" : 15,
  "monitoring" : 20,
  "worker" : 15,
  "database" : 20,
  "nginx" : 20,
  "database-dev-staging" : 25
}
vpc_cidr_block = "10.201.0.0/16"
priv_key       = "~/.ssh/id_rsa"
