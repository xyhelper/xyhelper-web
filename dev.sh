#!/bin/bash
set -e

export AUTH_SECRET_KEY=1234568
export KFURL=https://www.lidong.xin/hero.jpeg
export SHOW_PLUS_BTN=https://xyhelper.cn/plus/
export AD_MESSAGE="出少量PLUS独享会员，支持聊天记录漫游，免梯子"
# export BASE_URI=http://localhost:8081
# export ACCESS_TOKEN="hellow world"
cool-tools run main.go