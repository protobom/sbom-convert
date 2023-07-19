#!/bin/sh 
EVAL=$(eval echo "$@")
echo "{ \"calling\": \"$EVAL\" }"
exec $EVAL
