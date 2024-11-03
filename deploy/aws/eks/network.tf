
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
  name        = "Web Security Group"
  description = "Allow http and ssh access"
  vpc_id      = aws_vpc.ohio_vpc.id

  tags = {
    Name = "Web-Security-Group"
  }
}

# add the rule that accepts any http traffic
resource "aws_vpc_security_group_ingress_rule" "allow_http" {
  security_group_id = aws_security_group.app_security_group.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 80
  to_port           = 80
  ip_protocol       = "tcp"
}

# add the rule that accepts any https traffic
resource "aws_vpc_security_group_ingress_rule" "allow_https" {
  security_group_id = aws_security_group.app_security_group.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 443
  to_port           = 443
  ip_protocol       = "tcp"
}

# add the rule that accepts any http traffic
resource "aws_vpc_security_group_ingress_rule" "allow_8080" {
  security_group_id = aws_security_group.app_security_group.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 8080
  to_port           = 8080
  ip_protocol       = "tcp"
}


# add the rule that accepts ssh connections
resource "aws_vpc_security_group_ingress_rule" "allow_ssh" {
  security_group_id = aws_security_group.app_security_group.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 22
  to_port           = 22
  ip_protocol       = "tcp"
}

# NOTE: inbound connections to security group are stateful
# (so there is no need of rule for these outbound connections).
# But, connections initiated by resources in security group requires
# outbound rule in case they want to communicate outside
# (for example git clone or wget etc).

resource "aws_vpc_security_group_egress_rule" "allow_http" {
  security_group_id = aws_security_group.app_security_group.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 80
  to_port           = 80
  ip_protocol       = "tcp"
}

resource "aws_vpc_security_group_egress_rule" "allow_https" {
  security_group_id = aws_security_group.app_security_group.id
  cidr_ipv4         = "0.0.0.0/0"
  from_port         = 443
  to_port           = 443
  ip_protocol       = "tcp"
}
