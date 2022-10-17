#!/bin/bash

# upload file
curl -s -X POST http://localhost:8080/upload \
	-F "uploads=@/home/kanon/my_file.md" \
	-H "Content-Type: multipart/form-data"

# create user
curl -s -X POST http://localhost:8080/create_user
