#!/bin/bash

set -ex

ip=$(curl -s ip.sb)
export public_ip=$ip  # get host ip

declare -A service_map=(["api"]="api" ["user"]="user-service" ["auth"]="auth-service"
["bottle"]="bottle-service" ["share"]="share-service" ["upload"]="upload-service")

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
