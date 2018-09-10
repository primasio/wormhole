#!/bin/bash
docker build -t primasio/wormhole dist/
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
docker push primasio/wormhole
