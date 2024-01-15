#!/bin/sh
# get permissions
#[ "$(id -u)" -ne 0 ] && exec sudo "$0"

set -e

docker image build -f Dockerfile -t ascii-docker .
docker container run -p 8080:8080 -d --name webserver ascii-docker

echo "you can visit http://localhost:8080"
