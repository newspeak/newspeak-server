#!/bin/sh

# display all users

export ADMINTOKEN=$(curl 'http://localhost:8080/management/token?grant_type=password&username=admin&password=admin' | cut -f 2 -d , | cut -f 2 -d : | cut -f 2 -d \")
echo admintoken = $ADMINTOKEN
echo

curl -H "Authorization: Bearer $ADMINTOKEN" \
     -X GET "http://localhost:8080/newspeak/newspeak/users"

echo
