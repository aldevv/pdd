# (to be changed) create uploads folder manually

# Installation

## Download the google_cloud_credentials.json file from GCP

## Create a .env file with these variables: 
```bash
JWT_SECRET=
SQS_QUEUE_URL=
SQS_ENDPOINT_URL=
SENDGRID_API_KEY=
PDD_EMAIL=
MONGO_URL=mongodb://<user>:<pass>@localhost:27017/photos?authSource=admin
```
## Run the mongo database
```bash
docker compose up -d
```
## Install the project's dependencies
```bash
go mod tidy
```

## Run the project with
```bash
go run cmd/photos_api/main.go -s google # or local
```
