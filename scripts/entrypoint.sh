#!/bin/sh

git config --global --add safe.directory $PWD

# shellcheck disable=SC2068
exec $NAME $@
