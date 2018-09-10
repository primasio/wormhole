#!/bin/bash
docker build -t primasio/wormhole dist/
echo "$DOCKER_PASSWORD" | docker login -u "primasio" --password-stdin
docker push primasio/wormhole
