# GCP Logger

### Basic Idea:
1. Listen for messages (over HTTP)
1. Send messages to Google Cloud Logging (via golang api)

### App Configuration

1. Get service account token from GCP (with access to send logs to cloud logging)
1. `export GOOGLE_APPLICATION_CREDENTIALS=<PATH-TO-SERVICE-ACCOUNT-KEY-JSON>` #optional-if-running-in-gcp
1. `export PROJECT_ID=<GCP-PROJECT-ID>`
1. `export PORT=<PORT>` #defaults 8080
1. `export QUEUE_MAX_SIZE=<QUEUE_MAX_SIZE>` #defaults 10
1. `export LOGGER_NAME=<LOGGER_NAME>` #defaults cute-animal-name-logger

### K8s Stuff

1. Image is available at [oskoss/gcp-logger](https://hub.docker.com/repository/docker/oskoss/gcp-logger)
1. Namespace/Deployment/Service work with internal access only
1. Istio is optional