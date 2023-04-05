#!/bin/bash
set -e

docker login -u xyhelper

cd frontend
pnpm run build

cd ..
docker build -t xyhelper/xyhelper-web:latest .
docker push xyhelper/xyhelper-web:latest