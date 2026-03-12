#!/bin/bash

# 构建GoFlow项目

echo "Building GoFlow project..."

# 创建build目录
mkdir -p build

# 构建项目
go build -o build/GoFlow_exp ./src

if [ $? -eq 0 ]; then
    echo "Build successful! Executable is in build/GoFlow_exp"
    # 复制控制脚本到build目录
    cp control.sh build/
    chmod +x build/control.sh
    echo "Control script copied to build/control.sh"
else
    echo "Build failed!"
    exit 1
fi
