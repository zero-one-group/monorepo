
# AWS Infrastructure with OpenTofu

A comprehensive AWS infrastructure setup using OpenTofu (Terraform fork) with modular design, remote state management, and Docker Swarm deployment.


## Project Overview

This project implements Infrastructure as Code (IaC) to provision AWS resources in a modular and reusable manner across multiple environments (development, staging, production).

### Prerequisites

- OpenTofu installed

- AWS CLI configured with appropriate permissions

- AWS account with necessary access

## Infrastructure Deployment Process

### Base VPC Creation

First, deploy the base VPC module which provides foundational networking:
```
tofu init
tofu apply -target=module.vpc-base
```
- Note the generated **vpc_id** and **igw_id** for future reference you may see this on output after apply or check console
- The base VPC uses CIDR block 10.201.0.0/16(i will show later my example tfvars this will have impact because in security groups need to adjust a little bit if different CIDR block)

## State Management Resources

Next, create the S3 bucket and DynamoDB table for state management:

```
# Comment out the backend configuration in main.tf
tofu apply
-target=aws_s3_bucket.state_bucket
-target=aws_s3_bucket_server_side_encryption_configuration.state_bucket_encryption
-target=aws_dynamodb_table.state_lock
```

## Configure Remote Backend

After creating state resources:

- Uncomment the backend configuration in terraform.tf
- Update with your bucket and DynamoDB table names
- Initialize the backend:

```
# Here example backend
terraform {
  backend "s3" {
    bucket         = "lgtm-bucket-states"
    key            = "state/terraform.tfstate"
    region         = "ap-southeast-1"
    encrypt        = true
    dynamodb_table = "lgtm-locks"
  }
}
# do this comment for Initialize the backend
tofu init -reconfigure
```

## Complete Infrastructure Deployment

Before deploy here example for **terraform.auto.tfvars**
```
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
project             = "lgtm"
public_key_path     = "~/.ssh/id_rsa.pub"
region              = "ap-southeast-1"
ubuntu_jammy_ami    = "ami-0e0ddf453092e1e37"
ubuntu_noble_ami    = "ami-0b874c2ac1b5e9957"
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

```

After setup tfvars deploy the infrastructure components:
```
tofu apply
```

- This will apply **VPC** with 4 **Subnets** which is **Development**, **Staging**, **Production**, **Miscellanous(for monitoring, db and other services)**
- On **VPC** there will be elastic ip which is **Bastion** and **Nginx(load balancer)**
- Beside **VPC** this will also apply **EC2**, **ECR**, **S3(bucket for monitoring)**
- The structure will be 3 cluster docker swarm which is **Monitoring**, **Dev/Staging**, **Production** and **Nginx(load balancer)** as standalone docker
