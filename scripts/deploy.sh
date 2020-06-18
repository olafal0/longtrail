#!/bin/bash

set -e
source scripts/aws-tools.sh
test -f scripts/env.sh && source scripts/env.sh

# Check dependencies
source scripts/depcheck.sh
depcheck jq
depcheck aws

prefix="longtrail-$DEPLOY_ENV"

deployStack deployment-prereq.yml $prefix-deployments
deploymentBucket=$(getStackOutput $prefix-deployments DeploymentBucket)

deployLambdaCode $deploymentBucket longtrail-api
LAMBDA_S3_VERSION=$(getLatestVersion $deploymentBucket "lambdas/longtrail-api.zip")
echo "version: $LAMBDA_S3_VERSION"

deployStack longtrail.yml $prefix "
    --capabilities CAPABILITY_IAM
    --parameter-overrides
        APILambdaKey=lambdas/longtrail-api.zip
        APILambdaVersion=$LAMBDA_S3_VERSION
        DeploymentBucket=$deploymentBucket
        SenderEmailArn=$SENDER_EMAIL_ARN
        WebhookURL=$WEBHOOK_URL"
