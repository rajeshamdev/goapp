
# provision ASG (Auto Scaling Group)

resource "aws_launch_template" "bowbow_ec2_launch_template" {

  name_prefix   = "bowbow-ec2-launch-template"
  image_id      = var.aws_ec2_ami_id # Replace with ECS-optimized AMI ID
  instance_type = var.aws_ec2_type   # Choose appropriate instance type
  key_name      = "ecs"              # Replace with your key pair for ssh access

  iam_instance_profile {
    arn = aws_iam_instance_profile.bowbow_ecs_instance_profile.arn
    #name = "ecsInstanceRole" # gave it manually
  }

  network_interfaces {
    associate_public_ip_address = true
    security_groups             = [aws_security_group.bowbow_ecs_sg.id]
  }

  # It is required to pass ECS cluster name, so AWS can register EC2 instance
  # as node of ECS cluster. I wasted 2 days on missing this simple point.
  user_data = base64encode(<<-EOF
      #!/bin/bash
      echo ECS_CLUSTER=${aws_ecs_cluster.bowbow_ecs_cluster.name} >> /etc/ecs/ecs.config;
    EOF
  )

}

resource "aws_autoscaling_group" "bowbow_asg" {

  name = "bowbow-asg"

  launch_template {
    id      = aws_launch_template.bowbow_ec2_launch_template.id
    version = "$Latest"
  }

  min_size            = 2
  max_size            = 2
  desired_capacity    = 2
  vpc_zone_identifier = [aws_subnet.public_subnet.*.id[0]]

  tag {
    key                 = "AmazonECSManaged"
    value               = true
    propagate_at_launch = true
  }


  # Auto Scaling policies
  lifecycle {
    create_before_destroy = true
  }

}

# Auto Scaling Policy to Scale Up
resource "aws_autoscaling_policy" "bowbow_scale_up" {

  name                   = "bowbow-scale-up"
  scaling_adjustment     = 1
  adjustment_type        = "ChangeInCapacity"
  cooldown               = 300
  autoscaling_group_name = aws_autoscaling_group.bowbow_asg.name
}

# Auto Scaling Policy to Scale Down
resource "aws_autoscaling_policy" "bowbow_scale_down" {

  name                   = "bowbow-scale-down"
  scaling_adjustment     = -1
  adjustment_type        = "ChangeInCapacity"
  cooldown               = 300
  autoscaling_group_name = aws_autoscaling_group.bowbow_asg.name
}

resource "aws_ecs_capacity_provider" "bowbow_capcity_provider" {

  name = "bowbow-capacity-provider"

  auto_scaling_group_provider {
    auto_scaling_group_arn         = aws_autoscaling_group.bowbow_asg.arn
    managed_termination_protection = "DISABLED"

    managed_scaling {
      status          = "ENABLED"
      target_capacity = 80 # Target capacity as a percentage
    }
  }
}

resource "aws_ecs_cluster_capacity_providers" "example" {

  cluster_name = aws_ecs_cluster.bowbow_ecs_cluster.name

  capacity_providers = [aws_ecs_capacity_provider.bowbow_capcity_provider.name]

  default_capacity_provider_strategy {
    capacity_provider = aws_ecs_capacity_provider.bowbow_capcity_provider.name
    weight            = 1
  }
}
