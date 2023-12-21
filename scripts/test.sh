#!/bin/bash

#  tparse needs -json flag to parse the output.json file
set -o pipefail && TEST_OPTIONS="-json" task test | tee output.json | tparse -follow
success=$?

set -e
NO_COLOR=true tparse -format markdown -slow 10 -file output.json > $GITHUB_STEP_SUMMARY

exit $success
