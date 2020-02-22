#!/bin/bash

set -ex

build() {
    cd $1
    make build
    cd - > /dev/null
}


if [[ $1 -eq "-b" ]]
then
  exit 0
fi

if [[ $# -eq 0 ]]
then
    build api
    build user
    build auth
else
    for var in "$@"
    do
        build $var
    done
fi
