#!/bin/bash
set -e
# set -x
gf docker main.go -p -t xyhelper/xyhelper-web:latest
# 修改镜像标签为当前日期时间
time=$(date "+%Y%m%d%H%M%S")
docker tag xyhelper/xyhelper-web:latest xyhelper/xyhelper-web:$time
# 推送镜像到docker hub
docker push xyhelper/xyhelper-web:latest
docker push xyhelper/xyhelper-web:$time

