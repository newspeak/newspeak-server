#!/bin/sh

# pass username to be created as first parameter

export ADMINTOKEN=$(curl 'http://localhost:8080/management/token?grant_type=password&username=admin&password=admin' | cut -f 2 -d , | cut -f 2 -d : | cut -f 2 -d \")
echo admintoken = $ADMINTOKEN
echo

curl -H "Authorization: Bearer $ADMINTOKEN" \
     -X POST "http://localhost:8080/newspeak/newspeak/users" \
     -d '{ "username":"'$1'", "password":"'$1'", "email":"'$1'@example.com" }'

echo
