#!/bin/bash
set -e

# 获取脚本所在目录全路径
SCRIPT_PATH=$(cd `dirname $0`; pwd)

export AUTH_SECRET_KEY=111111
# export KFURL=https://www.lidong.xin/hero.jpeg
export SHOW_PLUS_BTN=https://xyhelper.cn/plus/
export AD_MESSAGE="出少量PLUS独享会员，支持聊天记录漫游，免梯子"
# export BASE_URI=http://localhost:8081
# export ACCESS_TOKEN="hellow world"

cd $SCRIPT_PATH/frontend
pnpm run dev &
cd $SCRIPT_PATH
cool-tools run main.go 

