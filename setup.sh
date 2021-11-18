#!/bin/sh

DEPLOY_SERVICE_ACCOUNT_NAME=coppe-deploy
PROJECT_ID=`yq eval '.GCP_PROJECT_ID' .env.yaml`
DEPLOY_SERVICE_ACCOUNT_EMAIL="$DEPLOY_SERVICE_ACCOUNT_NAME@$PROJECT_ID.iam.gserviceaccount.com"
BUCKET_NAME=$PROJECT_ID-coppe-tfstate-bucket

# create service account to deploy
gcloud iam service-accounts create $DEPLOY_SERVICE_ACCOUNT_NAME --display-name='For deploying Coppe'

# add policy binding for owner
gcloud projects add-iam-policy-binding $PROJECT_ID --member="serviceAccount:$DEPLOY_SERVICE_ACCOUNT_EMAIL" --role='roles/owner'

# create storage bucket for managing tfstate
gsutil mb -p $PROJECT_ID -c STANDARD -l US-CENTRAL1 gs://$BUCKET_NAME

# set the bucket's versioning on
gsutil versioning set on gs://$BUCKET_NAME
