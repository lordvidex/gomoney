#!/bin/bash

# this scripts builds and tags with docker

# it takes one argument, the tag version

TAG=latest

if [ $# -eq 1 ]; then
    TAG=$1
fi

docker build -t lordvidex/gomoney-api:$TAG -f deployments/Dockerfile-api --target production .
docker build -t lordvidex/gomoney-central:$TAG -f deployments/Dockerfile-server --target production .
docker build -t lordvidex/gomoney-telegram:$TAG -f deployments/Dockerfile-telegram --target production .

echo "Login to docker hub"
docker login -u lordvidex

docker push lordvidex/gomoney-api:$TAG
docker push lordvidex/gomoney-central:$TAG
docker push lordvidex/gomoney-telegram:$TAG