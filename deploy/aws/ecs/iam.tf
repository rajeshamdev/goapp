
# This has few bugs. TBF (Raj)

resource "aws_iam_role" "bowbow_ecs_instance_role" {

  name = "bowbow-ecs-instance-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Principal = {
          Service = "ec2.amazonaws.com"
        }
        Action = "sts:AssumeRole"
      },
    ]
  })
}

resource "aws_iam_role_policy_attachment" "bowbow_ecs_instance_policy" {

  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role"
  role       = aws_iam_role.bowbow_ecs_instance_role.name

}

resource "aws_iam_instance_profile" "bowbow_ecs_instance_profile" {

  name = "bowbow-ecs-instance-profile"
  role = aws_iam_role.bowbow_ecs_instance_role.name

}
