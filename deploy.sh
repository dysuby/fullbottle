#!/bin/bash

set -ex

declare -A service_map=(["api"]="api" ["user"]="user-service" ["auth"]="auth-service")

./build.sh $@

services=""
for var in "$@"
do
    services="${services} ${service_map[$var]}"
done

docker-compose build $services

docker-compose up -d $services
