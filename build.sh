#!/bin/bash

set -ex

build() {
    cd $1
    make build
    cd - > /dev/null
}

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
