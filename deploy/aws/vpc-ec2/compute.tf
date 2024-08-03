
resource "aws_instance" "web_app" {
  ami                         = var.aws_ec2_ami_id # "ami-0649bea3443ede307"
  instance_type               = var.aws_ec2_type
  availability_zone           = var.aws_azs[0]
  subnet_id                   = aws_subnet.public_subnet.0.id
  vpc_security_group_ids      = [aws_security_group.app_security_group.id]
  associate_public_ip_address = true
  key_name                    = aws_key_pair.webapp_keypair.key_name
  user_data                   = file("./webApp.sh")

  tags = {
    Name = "Web Application"
  }
}


resource "aws_key_pair" "webapp_keypair" {
  key_name   = "webapp-ssh-keypair"
  public_key = file("~/.ssh/id_rsa.pub") # Path to your public key
}
