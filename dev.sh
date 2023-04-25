#!/bin/bash
set -e

export AUTH_SECRET_KEY=1234568
export KFURL=https://www.lidong.xin/hero.jpeg
export BASE_URI=http://localhost:8081
export ACCESS_TOKEN="hellow world"
cool-tools run main.go