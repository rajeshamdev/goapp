

resource "aws_ecs_task_definition" "bowbow_task" {

  family = "bowbow-task-definition"

  # launch type
  requires_compatibilities = ["EC2"]

  # OS, Architecture, Network mode
  runtime_platform {
    operating_system_family = "LINUX"
    cpu_architecture        = "X86_64"
  }
  # I could not get public IP without bridge net
  network_mode = "bridge"

  # Task size:
  # these limits should be lower than EC2 instance's. If these are equal or greater than EC2 limits,
  # the the creating tasks bound to fail.
  cpu    = 1024
  memory = 1024

  task_role_arn      = "arn:aws:iam::${var.aws_account}:role/ecsTaskExecutionRole"
  execution_role_arn = "arn:aws:iam::${var.aws_account}:role/ecsTaskExecutionRole"

  # TODO (Raj): Fix this
  #task_role_arn      = aws_iam_role.bowbow_ecs_instance_role.arn
  #execution_role_arn = aws_iam_role.bowbow_ecs_instance_role.arn

  container_definitions = jsonencode([{
    name      = "goapp"
    image     = var.goapp_image
    cpu       = 256
    memory    = 512
    essential = true
    portMappings = [
      {
        containerPort = 8080
        hostPort      = 8080
      }
    ]
    "environment" : [
      {
        "name" : "GCP_APIKEY",
        "value" : var.gcp_apikey
      }
    ]

  }])
}
