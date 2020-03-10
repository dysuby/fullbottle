#!/bin/bash

set -ex

build() {
    cd $1
    make build
    cd - > /dev/null
}


if [[ $1 = "-a" ]]
then
    build api
    build user
    build auth
    build bottle
    build share
    build upload
else
    for var in "$@"
    do
        build $var
    done
fi
