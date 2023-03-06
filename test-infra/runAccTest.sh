#!/bin/bash

echoLog () {
  echo "===> $1 <==="
}

if [ $# -eq 0 ]; then
  echo "At least 1 argument is required"
  exit 1
fi

if [ "$1" == "ALL" ]; then
  echoLog "Running acceptance tests for all resources"
  cd .. &&
    TF_ACC=1 go test -v ./... -timeout=1200m
else
  resources=" "
  for var in "$@"; do
    resources="$resources$var "
  done
  echoLog "Running acceptance tests for $# resources [$resources]"

  test_regex=""
  first=true
  for var in "$@"; do
    if [ $first ]; then
      test_regex="TestAccZIA${var}_basic$"
      first=false
    else
      test_regex="$test_regex/TestAccZIA${var}_basic$"
    fi
  done
  echoLog "Go test regex: $test_regex"
  cd .. &&
    TF_ACC=1 go test -v ./aviatrix -run $test_regex -timeout=1800s
fi
