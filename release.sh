#!/bin/bash
set -e
# set -x

# 生成版本号 如 1.2.3
if [ -z "$1" ]; then
    echo "请输入版本号 如 1.2.3"
    exit 1
fi
version=$1

# 获取 Major Minor Patch
major=$(echo $version | cut -d. -f1)
minor=$(echo $version | cut -d. -f2)
patch=$(echo $version | cut -d. -f3)

echo "major: $major"
echo "minor: $minor"
echo "patch: $patch"

docker login -u xyhelper

cd frontend
pnpm run build

cd ..
# docker build -t xyhelper/xyhelper-web:latest .
# docker push xyhelper/xyhelper-web:latest

docker buildx build -f Dockerfile.release --build-arg VERSION=v$version --platform linux/amd64,linux/arm64,linux/arm/v7 -t xyhelper/xyhelper-web:latest --push .
docker buildx build -f Dockerfile.release --build-arg VERSION=v$version --platform linux/amd64,linux/arm64,linux/arm/v7 -t xyhelper/xyhelper-web:$major --push .
docker buildx build -f Dockerfile.release --build-arg VERSION=v$version --platform linux/amd64,linux/arm64,linux/arm/v7 -t xyhelper/xyhelper-web:$major.$minor --push .
docker buildx build -f Dockerfile.release --build-arg VERSION=v$version --platform linux/amd64,linux/arm64,linux/arm/v7 -t xyhelper/xyhelper-web:$major.$minor.$patch --push .