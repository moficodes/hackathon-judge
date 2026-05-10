#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -euo pipefail

echo "=========================================="
echo " Starting GOOGLE_CLOUD Infrastructure Provisioning"
echo "=========================================="

# 1. Source the .env file
if [ -f .env ]; then
  echo "Sourcing .env file..."
  # Export variables from .env ignoring comments
  export $(grep -v '^#' .env | xargs)
else
  echo "Error: .env file not found at the root of the project."
  echo "Please create one using the provided template."
  exit 1
fi

# Ensure critical variables are set
if [ -z "${GOOGLE_CLOUD_PROJECT_ID:-}" ] || [ "$GOOGLE_CLOUD_PROJECT_ID" == "your-project-id-here" ]; then
  echo "Error: GOOGLE_CLOUD_PROJECT_ID is not set correctly in .env"
  exit 1
fi

echo "Project ID: $GOOGLE_CLOUD_PROJECT_ID"
echo "Region: $GOOGLE_CLOUD_REGION"

# 2. Make sure the user is authenticated
echo "Checking GOOGLE_CLOUD authentication..."
if ! gcloud auth print-access-token &> /dev/null; then
  echo "Error: You are not authenticated with Google Cloud."
  echo "Please run 'gcloud auth login' and 'gcloud auth application-default login' first."
  exit 1
fi

# Set the active project
gcloud config set project "$GOOGLE_CLOUD_PROJECT_ID"

# 3. Enable necessary Google Cloud Services
echo "Enabling required Google Cloud APIs..."
gcloud services enable \
  container.googleapis.com \
  artifactregistry.googleapis.com \
  cloudbuild.googleapis.com \
  pubsub.googleapis.com

# 4. Create an Artifact Registry Docker repo
echo "Checking Artifact Registry repository: $ARTIFACT_REPO_NAME..."
if ! gcloud artifacts repositories describe "$ARTIFACT_REPO_NAME" --location="$ARTIFACT_REGISTRY_LOCATION" &> /dev/null; then
  echo "Creating Artifact Registry repository in location: $ARTIFACT_REGISTRY_LOCATION..."
  gcloud artifacts repositories create "$ARTIFACT_REPO_NAME" \
    --repository-format=docker \
    --location="$ARTIFACT_REGISTRY_LOCATION" \
    --description="Docker repository for Hackathon Judge services"
else
  echo "Artifact Registry repository '$ARTIFACT_REPO_NAME' already exists in $ARTIFACT_REGISTRY_LOCATION."
fi

# 5. Create GKE Autopilot Cluster
echo "Checking GKE Autopilot cluster: $CLUSTER_NAME..."
if ! gcloud container clusters describe "$CLUSTER_NAME" --region="$GOOGLE_CLOUD_REGION" &> /dev/null; then
  echo "Creating GKE Autopilot cluster (this may take 5-10 minutes)..."
  gcloud container clusters create-auto "$CLUSTER_NAME" \
    --region="$GOOGLE_CLOUD_REGION" \
    --project="$GOOGLE_CLOUD_PROJECT_ID"
else
  echo "GKE Autopilot cluster '$CLUSTER_NAME' already exists."
fi

# 6. Create Pub/Sub Topics and Subscriptions
echo "Configuring Pub/Sub Topics and Subscriptions..."

# Helper function to create topic
create_topic() {
  local topic=$1
  if ! gcloud pubsub topics describe "$topic" &> /dev/null; then
    echo "Creating topic: $topic"
    gcloud pubsub topics create "$topic"
  else
    echo "Topic '$topic' already exists."
  fi
}

# Helper function to create subscription
create_subscription() {
  local sub=$1
  local topic=$2
  if ! gcloud pubsub subscriptions describe "$sub" &> /dev/null; then
    echo "Creating subscription: $sub for topic: $topic"
    gcloud pubsub subscriptions create "$sub" --topic="$topic"
  else
    echo "Subscription '$sub' already exists."
  fi
}

# Topics
create_topic "$PUBSUB_TOPIC_JUDGING_TASKS"
create_topic "$PUBSUB_TOPIC_JUDGING_RESULTS"

# Subscriptions
create_subscription "$PUBSUB_SUB_AGENT" "$PUBSUB_TOPIC_JUDGING_TASKS"
create_subscription "$PUBSUB_SUB_BACKEND" "$PUBSUB_TOPIC_JUDGING_RESULTS"

# 7. Run Cloud Build
echo "Triggering Cloud Build for all services..."
gcloud builds submit --config cloudbuild.yaml . \
  --substitutions=_REGION="$ARTIFACT_REGISTRY_LOCATION",_REPO="$ARTIFACT_REPO_NAME"

echo "=========================================="
echo " Infrastructure Setup & Build Complete!"
echo "=========================================="
