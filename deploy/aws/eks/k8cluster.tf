
# IAM Role
resource "aws_iam_role" "eks_cluster_role" {
  name               = "eksClusterRole"
  assume_role_policy = <<EOF
  {
    "Version": "2012-10-17",
    "Statement": [
      {
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "eks.amazonaws.com"
          ]
        },
        "Action": "sts:AssumeRole"
      }
    ]
  }
  EOF
}

resource "aws_iam_role_policy_attachment" "eks_cluster_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKSClusterPolicy"
  role       = aws_iam_role.eks_cluster_role.name
}


resource "aws_iam_role_policy_attachment" "eks_cni_policy" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy"
  role       = aws_iam_role.eks_cluster_role.name
}

# EKS Cluster
resource "aws_eks_cluster" "my_cluster" {
  name     = "k8cluster"
  role_arn = aws_iam_role.eks_cluster_role.arn

  vpc_config {
    subnet_ids              = aws_subnet.public_subnet[*].id
    security_group_ids      = [aws_security_group.app_security_group.id]
    endpoint_public_access  = true
    endpoint_private_access = true
  }

  # Specify Kubernetes version
  version = "1.31" # Change as needed

  # Upgrade settings
  lifecycle {
    create_before_destroy = true
  }

  depends_on = [
    aws_iam_role_policy_attachment.eks_cluster_policy
  ]
}

# EKS Add-ons
resource "aws_eks_addon" "vpc_cni" {
  cluster_name = aws_eks_cluster.my_cluster.name
  addon_name   = "vpc-cni"
  # addon_version     = "v1.18.3-eksbuild.2"
}

resource "aws_eks_addon" "kube_proxy" {
  cluster_name = aws_eks_cluster.my_cluster.name
  addon_name   = "kube-proxy"
}

resource "aws_eks_addon" "coredns" {
  cluster_name = aws_eks_cluster.my_cluster.name
  addon_name   = "coredns"
}

/*
resource "aws_eks_addon" "pod_identity" {
  cluster_name = aws_eks_cluster.my_cluster.name
  addon_name   = "eks-pod-identity-agent"
}
*/
