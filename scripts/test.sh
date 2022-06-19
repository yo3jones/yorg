#!/bin/bash

if [ -z "$1" ]
then
  path="./..."
else
  path=$1
fi

if [[ ! $path = ./* ]]
then
  path="./$path"
fi

clear \
&& printf "go test -v -coverprofile cover.out $path \\\\\n" \
&& printf "&& go tool cover -func cover.out\n\n" \
&& go test -v -coverprofile cover.out $path \
&& go tool cover -func cover.out
