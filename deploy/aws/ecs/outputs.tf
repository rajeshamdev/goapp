output "public_subnets" {
  value = aws_subnet.public_subnet.*.id
}

output "lb_dns_name" {
  value = aws_lb.bowbow_lb.dns_name
}
