
# create VPC
resource "aws_vpc" "ohio_vpc" {

  cidr_block = var.aws_vpc_cidr

  tags = {
    Name = "${var.company_name}-VPC"
  }
}

# create two public subneits
resource "aws_subnet" "public_subnet" {

  count             = length(var.aws_public_subnet_cidr)
  vpc_id            = aws_vpc.ohio_vpc.id
  cidr_block        = element(var.aws_public_subnet_cidr, count.index)
  availability_zone = element(var.aws_azs, count.index)

  tags = {
    Name = "${var.company_name}-Public-Subnet-${count.index + 1}"
  }

  map_public_ip_on_launch = true # needed for EKS

}

# create two private subnets
resource "aws_subnet" "private_subnet" {

  count             = length(var.aws_private_subnet_cidr)
  vpc_id            = aws_vpc.ohio_vpc.id
  cidr_block        = element(var.aws_private_subnet_cidr, count.index)
  availability_zone = element(var.aws_azs, count.index)

  tags = {
    Name = "${var.company_name}-Private-Subnet-${count.index + 1}"
  }
}

# create internet gateway
resource "aws_internet_gateway" "igw" {

  vpc_id = aws_vpc.ohio_vpc.id

  tags = {
    Name = "${var.company_name}-igw"
  }
}

# create route table and attach to internet gateway
resource "aws_route_table" "rt" {

  vpc_id = aws_vpc.ohio_vpc.id

  route {
    cidr_block = var.aws_default_route
    gateway_id = aws_internet_gateway.igw.id
  }

  tags = {
    Name = "${var.company_name}-Default-Route"
  }
}


# create route table associations
resource "aws_route_table_association" "rt_association" {

  count          = length(var.aws_public_subnet_cidr)
  subnet_id      = element(aws_subnet.public_subnet.*.id, count.index)
  route_table_id = aws_route_table.rt.id
}
#
#
# create security group
resource "aws_security_group" "app_security_group" {
  name        = "AWS EKS Security Group"
  description = "Security group for goapp, Prometheus, Grafana and Loki"
  vpc_id      = aws_vpc.ohio_vpc.id

  tags = {
    Name = "EKS-Security-Group"
  }
}

# Allow http traffic
resource "aws_vpc_security_group_ingress_rule" "http" {
  security_group_id = aws_security_group.app_security_group.id
  ip_protocol       = "tcp"
  from_port         = 80
  to_port           = 80
  cidr_ipv4         = "0.0.0.0/0"
  description       = "allow http access"
}

# Allow https traffic
resource "aws_vpc_security_group_ingress_rule" "https" {
  security_group_id = aws_security_group.app_security_group.id
  ip_protocol       = "tcp"
  from_port         = 443
  to_port           = 443
  cidr_ipv4         = "0.0.0.0/0"
  description       = "allow https access"
}

# allow kubelet metrics (port 10250)
resource "aws_vpc_security_group_ingress_rule" "kubelet_metrics" {
  security_group_id = aws_security_group.app_security_group.id
  ip_protocol       = "tcp"
  from_port         = 10250
  to_port           = 10250
  cidr_ipv4         = "0.0.0.0/0"
  description       = "allow access to kubelet metrics endpoint"
}

# add the rule that accepts traffic to bowbow-app
resource "aws_vpc_security_group_ingress_rule" "goapp" {
  security_group_id = aws_security_group.app_security_group.id
  ip_protocol       = "tcp"
  from_port         = 8080
  to_port           = 8080
  cidr_ipv4         = "0.0.0.0/0"
  description       = "allow access to goapp endpoint"
}


# Allow prometheus UI
resource "aws_vpc_security_group_ingress_rule" "prometheus_ui" {
  security_group_id = aws_security_group.app_security_group.id
  ip_protocol       = "tcp"
  from_port         = 9090
  to_port           = 9090
  cidr_ipv4         = "0.0.0.0/0"
  description      = "allow access to Prometheus web UI"
}

# Allow Grafana UI
resource "aws_vpc_security_group_ingress_rule" "grafana_ui" {
  security_group_id = aws_security_group.app_security_group.id
  ip_protocol       = "tcp"
  from_port         = 3000
  to_port           = 3000
  cidr_ipv4         = "0.0.0.0/0"
  description      = "Allow access to Grafana web UI"
}

# Allow Loki API
resource "aws_vpc_security_group_ingress_rule" "loki_api" {
  security_group_id = aws_security_group.app_security_group.id
  ip_protocol       = "tcp"
  from_port         = 3100
  to_port           = 3100
  cidr_ipv4         = "0.0.0.0/0"
  description      = "Allow access to Loki API"
}

# add the rule that accepts ssh connections
resource "aws_vpc_security_group_ingress_rule" "allow_ssh" {
  security_group_id = aws_security_group.app_security_group.id
  ip_protocol       = "tcp"
  from_port         = 22
  to_port           = 22
  cidr_ipv4         = "0.0.0.0/0"
}

# NOTE: inbound connections to security group are stateful
# (so there is no need of rule for these outbound connections).
# But, connections initiated by resources in security group requires
# outbound rule in case they want to communicate outside
# (for example git clone or wget etc).

resource "aws_vpc_security_group_egress_rule" "allow_https" {
  security_group_id = aws_security_group.app_security_group.id
  ip_protocol       = "tcp"
  from_port         = 443
  to_port           = 443
  cidr_ipv4         = "0.0.0.0/0"
}
