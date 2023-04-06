#!/bin/bash
set -e

# 生成版本号
if [ -z "$1" ]; then
    echo "请输入版本号"
    exit 1
fi
version=$1


docker login -u xyhelper

cd frontend
pnpm run build

cd ..
# docker build -t xyhelper/xyhelper-web:latest .
# docker push xyhelper/xyhelper-web:latest

docker buildx build -f Dockerfile.release --build-arg Version=$version --platform linux/amd64,linux/arm64,linux/arm/v7 -t xyhelper/xyhelper-web:latest --push .