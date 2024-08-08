
# create a application load balancer:
# 1. LB listens on :80 and forwards the traffic to goapp containers listening on :8080
# 2. Creates DNS name(A record) wrapping EC2 instances under it.
#

resource "aws_lb" "bowbow_lb" {

  name               = "bowbow-lb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.bowbow_lb_sg.id]
  subnets            = [aws_subnet.public_subnet.0.id, aws_subnet.public_subnet.1.id]

  enable_deletion_protection       = false
  enable_cross_zone_load_balancing = true
}

resource "aws_lb_target_group" "bowbow_target_group" {

  name        = "bowbow-ecs-target"
  target_type = "instance" # This is the target type for EC2 instances

  port     = 8080
  protocol = "HTTP"
  vpc_id   = aws_vpc.bowbow_vpc.id


  health_check {
    path                = "/v1/api/health"
    interval            = 30
    timeout             = 5
    healthy_threshold   = 3
    unhealthy_threshold = 3
  }
}

resource "aws_lb_listener" "http" {
  load_balancer_arn = aws_lb.bowbow_lb.arn
  port              = 80
  protocol          = "HTTP"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.bowbow_target_group.arn
  }
}
