
# create VPC
resource "aws_vpc" "bowbow_vpc" {

  cidr_block = var.aws_vpc_cidr

  tags = {
    Name = "${var.company_name}-VPC"
  }
}

# create two public subneits
resource "aws_subnet" "public_subnet" {

  count             = length(var.aws_public_subnet_cidr)
  vpc_id            = aws_vpc.bowbow_vpc.id
  cidr_block        = element(var.aws_public_subnet_cidr, count.index)
  availability_zone = element(var.aws_azs, count.index)

  tags = {
    Name = "${var.company_name}-Public-Subnet-${count.index + 1}"
  }
}

# create two private subnets
resource "aws_subnet" "private_subnet" {

  count             = length(var.aws_private_subnet_cidr)
  vpc_id            = aws_vpc.bowbow_vpc.id
  cidr_block        = element(var.aws_private_subnet_cidr, count.index)
  availability_zone = element(var.aws_azs, count.index)

  tags = {
    Name = "${var.company_name}-Private-Subnet-${count.index + 1}"
  }
}

# create internet gateway
resource "aws_internet_gateway" "igw" {

  vpc_id = aws_vpc.bowbow_vpc.id

  tags = {
    Name = "${var.company_name}-igw"
  }
}

# create route table and attach to internet gateway
resource "aws_route_table" "rt" {

  vpc_id = aws_vpc.bowbow_vpc.id

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

# security group for ecs
resource "aws_security_group" "bowbow_ecs_sg" {

  name   = "bowbow-ecs-security-group"
  vpc_id = aws_vpc.bowbow_vpc.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 8080
    to_port     = 8080
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# security group for lb
resource "aws_security_group" "bowbow_lb_sg" {

  name   = "bowbow-lb-security-group"
  vpc_id = aws_vpc.bowbow_vpc.id

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
