
# Create Service and attach ASG

resource "aws_ecs_service" "bowbow_service" {

  name            = "bowbow-service"
  cluster         = aws_ecs_cluster.bowbow_ecs_cluster.id
  task_definition = aws_ecs_task_definition.bowbow_task.arn
  desired_count   = 2
  launch_type     = "EC2"


  force_new_deployment = true

  load_balancer {
    target_group_arn = aws_lb_target_group.bowbow_target_group.arn
    container_name   = "goapp"
    container_port   = 8080
  }

  depends_on = [aws_autoscaling_group.bowbow_asg]
}
