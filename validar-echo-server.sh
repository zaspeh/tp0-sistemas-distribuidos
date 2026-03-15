#!/bin/bash

NETWORK=$(docker network ls --format '{{.Name}}' | grep testing_net)
SERVER="server"

PORT=$(grep -i port server/config.ini | cut -d '=' -f2 | tr -d ' ')

MESSAGE="hola"

RESPONSE=$(docker run --rm --network "$NETWORK" busybox \
sh -c "echo $MESSAGE | nc -w 1 $SERVER $PORT")

if [ "$RESPONSE" = "$MESSAGE" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi