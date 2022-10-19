#!/bin/bash
# do Y:@"
# let $default="http://localhost:8080"

# upload file
curl -s -X POST $default/upload \
	-F "uploads=@/home/kanon/my_file.md"

# create user
curl -s -X POST $default/create_user \
	-d @- <<EOF
{
    "name":"John Doe",
    "password":"my weak password",
    "email":"example3@example.com",
    "phone":"+573105236382"
}
EOF

# get user
curl -s -X GET $default/get_user

# get users
curl -s -X GET $default/get_users
