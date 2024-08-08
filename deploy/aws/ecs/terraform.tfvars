
aws_region   = "us-east-2"
company_name = "bowbow"

# The goal: 3 subnets in 3 AZ. Each subnet with 1024 IP address range.
#
# CIDR /19: 2^(32-19=13) = 8192 IP v4 addresses
# Total 6 subnets - 3 public and 3 private
# Each subnet CIDS /22: 2^(32-22=10) = 1024 IP addresses
#
# Subnet 1: 10.0.0.0/22
# Range:    10.0.0.0 to 10.0.3.255
#
# Subnet 2: 10.0.4.0/22
# Range:    10.0.4.0 to 10.0.7.255
#
# Subnet 3: 10.0.8.0/22
# Range:    10.0.8.0 to 10.0.11.255
#
# Subnet 4: 10.0.12.0/22
# Range:    10.0.12.0 to 10.0.15.255
#
# Subnet 5: 10.0.16.0/22
# Range:    10.0.16.0 to 10.0.19.255
#
# Subnet 6: 10.0.20.0/22
# Range:    10.0.20.0 to 10.0.23.255
#

aws_vpc_cidr = "10.0.0.0/19"

aws_azs = ["us-east-2a", "us-east-2b", "us-east-2c"]

aws_public_subnet_cidr = ["10.0.0.0/22", "10.0.4.0/22", "10.0.8.0/22"]

aws_private_subnet_cidr = ["10.0.12.0/22", "10.0.16.0/22", "10.0.20.0/22"]

aws_default_route = "0.0.0.0/0"


#
# For ECS: ami-071d7370f6f6a5ba1
# amzn2-ami-ecs-kernel-5.10-hvm-2.0.20240802-x86_64-ebs
aws_ec2_ami_id = "ami-071d7370f6f6a5ba1"

# 1 vCPU, 2GB memory
aws_ec2_type = "t2.small"

aws_account = "<aws_account>"
goapp_image = "<aws_account>.dkr.ecr.us-east-2.amazonaws.com/rajeshamdev:goapp-v1.0"
gcp_apikey  = "<GCP APIKEY>"
