#!/bin/bash

set -ex

export PublicIP=$(curl -s ip.sb)  # get host ip

declare -A service_map=(["api"]="api" ["user"]="user-service" ["auth"]="auth-service")

services=""
for var in "$@"
do
    services="${services} ${service_map[$var]}"
done

if [[ $services || $1 = "-a" ]]
then
  docker-compose build $services
fi

docker-compose up -d $services
