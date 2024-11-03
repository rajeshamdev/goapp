
# IAM Role
resource "aws_iam_role" "eks_node_group_role" {
  name               = "eksNodeGroupRole"
  assume_role_policy = <<EOF
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Principal": {
          "Service": ["ec2.amazonaws.com"]
        },
        "Action": "sts:AssumeRole"
      }
    ]
  }
  EOF
}

resource "aws_iam_role_policy_attachment" "eks_worker_node_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy"
  role       = aws_iam_role.eks_node_group_role.name
}

resource "aws_iam_role_policy_attachment" "eks_nodegroup_cni_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  role       = aws_iam_role.eks_node_group_role.name
}

resource "aws_iam_role_policy_attachment" "eks_ec2_container_registry_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly"
  role       = aws_iam_role.eks_node_group_role.name
}

# EKS Node Group
resource "aws_eks_node_group" "k8nodes" {
  cluster_name    = aws_eks_cluster.my_cluster.name
  node_group_name = "k8nodes"
  node_role_arn   = aws_iam_role.eks_node_group_role.arn

  subnet_ids = aws_subnet.public_subnet[*].id

  scaling_config {
    desired_size = 2
    max_size     = 2
    min_size     = 2
  }

  instance_types = [var.aws_k8node_ec2_type]
  ami_type       = "AL2_x86_64"

  remote_access {
     ec2_ssh_key = aws_key_pair.k8node_group_keys.key_name
  }

  # Tags for the EC2 instances
  tags = {
    Name = "k8nodes-ec2-instance"
  }


  # Ensure that IAM Role permissions are created before and deleted after EKS Node Group handling.
  # Otherwise, EKS will not be able to properly delete EC2 Instances and Elastic Network Interfaces.
  depends_on = [
    aws_iam_role_policy_attachment.eks_worker_node_policy,
    aws_iam_role_policy_attachment.eks_nodegroup_cni_policy,
    aws_iam_role_policy_attachment.eks_ec2_container_registry_policy,
  ]
}


resource "aws_key_pair" "k8node_group_keys" {
  key_name   = "k8node-group-keys"
  public_key = file("~/.ssh/id_rsa.pub")  # Replace with the path to your public key
}

