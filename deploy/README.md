# Multi-container Multi-zone App Deployment on AWS

Technologies used:
 - VPC (subnets, security groups, internet gateway, routes)
 - IAM (Roles)
 - EC2 Instances (with ASG - Auto Scaling Group)
 - Serverless compute (Fargate and Lambda)
 - Container deployment with ECS and EKS
 - Terraform (or drop in compatible OpenTofu)
 - Dockers and containers

There are two apps - a backend REST API server and the front-end React app.

Backend app is written in golang utilizing gin framework (github.com/gin-gonic/gin for REST API routes)
and YouTube Data APIs:
 1) To collect YouTube channel metrics/stats
 2) Walks through the list of videos on a channel
 3) Find the video's metrics (views, likes)
 4) Gets comments, and the user who commented
 5) Does sentiment analysis (VADER - Valence Aware Dictionary and sEntiment Reasoner) on the comments.

Note: YouTube Data APIs are subjective to API Rate Limits.

Front-end is a very simple React GUI app.

Docker images built for both apps (backend and frontend):
 - Tested standalone
 - Tested by deploying on AWS with API Gateway + Lambda
 - Tested by deploying on ECS:
     1) with EC2 (Auto Scaling Group) infrastructure
     2) Fargate serverless Infrastructure

Used Terraform to deploy on AWS.

