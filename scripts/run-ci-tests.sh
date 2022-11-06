#!/bin/bash

set -e

docker-compose -f "docker-compose.yml" up -d --build ci
docker-compose -f "docker-compose.yml" run --rm ci
