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
