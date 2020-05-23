#!/bin/bash

set -e
source scripts/aws-tools.sh

# Check dependencies
source scripts/depcheck.sh
depcheck jq
depcheck aws

dev_prefix="longtrail-development"

DEV_BACKEND_API_ID=$(getStackOutput $dev_prefix BackendAPIId)
DEV_USER_POOL_ID=$(getStackOutput $dev_prefix UserPoolId)
DEV_WEB_CLIENT_ID=$(getStackOutput $dev_prefix WebClientId)

echo "let config;

if (process.env.NODE_ENV === 'development') {
  config = {
    Auth: {
      region: 'us-east-1',
      userPoolId: '$DEV_USER_POOL_ID',
      userPoolWebClientId: '$DEV_WEB_CLIENT_ID',
    },
    apiUrl: 'https://$DEV_BACKEND_API_ID.execute-api.us-east-1.amazonaws.com/api',
  };
}

export default config;
" > web/config.js
