# Deploy Golang backend server and React GUI App on AWS

Backend App is written in golang:
 1) collects YouTube channel metrics/stats
 2) walks through list of videos of a channel
 3) finds each video's metrics (views, likes etc)
 4) Gets comments, and the user who commented
 5) does sentiment analysis (VADER - Valence Aware Dictionary and sEntiment Reasoner) on the comments.

The app uses youtube data APIs (golang). App provides REST APIs by utilizing Golang gin router framework.

Also, wrote a simple React GUI app.

Both apps (backend and frontend) are built as docker images, tested standalone, and also tested by deploying
on AWS with API Gateway + Lambda, ECS with EC2 and Fargate. Used Terraform to deploy on AWS.

