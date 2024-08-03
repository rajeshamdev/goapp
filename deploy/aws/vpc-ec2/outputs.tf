output "public_subnets" {
  value = aws_subnet.public_subnet.*.id
}

output "private_subnets" {
  value     = aws_subnet.private_subnet
  sensitive = true
}

output "security_group_id" {
  value = aws_security_group.app_security_group.id
}

output "ec2_public_ip" {
  value = aws_instance.web_app.public_ip
}

output "webapp_state" {
  value = aws_instance.web_app.instance_state
}

output "ec2_public_ip_sensitive" {
  value     = aws_instance.web_app.public_ip
  sensitive = true
}

