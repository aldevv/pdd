# Installation

## if using GCP
Download the google_cloud_credentials.json file from the GCP IAM

## Create a .env file with these variables: 
```bash
JWT_SECRET=
SQS_QUEUE_URL=http://localhost:9324/queue/default
SQS_ENDPOINT_URL=http://localhost:9324
SENDGRID_API_KEY=
PDD_EMAIL=
MONGO_USER=
MONGO_PASS=
MONGO_URL=mongodb://<user>:<pass>@localhost:27017/photos?authSource=admin
RECOVERY_PAGE=
```
## Run the containers
```bash
make local-up
```
## Install the project's dependencies
```bash
go mod tidy
```

## Run the project with
```bash
# uses GCP
make backend
# or
# uses local filesystem
make backend-local
```
