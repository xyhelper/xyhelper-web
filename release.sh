#!/bin/bash
set -e

docker login -u xyhelper

cd frontend
pnpm run build

cd ..
# docker build -t xyhelper/xyhelper-web:latest .
# docker push xyhelper/xyhelper-web:latest

docker buildx build -f Dockerfile.release --platform linux/amd64,linux/arm64,linux/arm/v7 -t xyhelper/xyhelper-web:latest --push .