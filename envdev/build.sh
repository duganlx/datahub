#!/bin/bash
#
# Author: lvx
# Date: 2023-11-24
# Description: 搭建开发环境的 Main 程序

BASIC_IMAGE_NAME=devbasic
BASIC_IMAGE_TAG=v1.0.0

cat <<EOF
操作引导
  0 [生成basic镜像]
  1 [生成basic容器]
EOF
read -p "选择进行的操作: " opt

case $opt in
  0)
    # 有三种情况: 1. 无镜像; 2. 有镜像无容器; 3. 有镜像有容器
    container_ids=$(docker ps -q --filter ancestor="$BASIC_IMAGE_NAME:$BASIC_IMAGE_TAG")
    if [ -n "$container_ids" ]; then
      echo -e "镜像$BASIC_IMAGE_NAME:$BASIC_IMAGE_TAG 存在如下容器:\n$container_ids"
      read -p "按回车键将进行删除..."
      for id in $container_ids; do
        docker rm -f $id
      done
    fi

    docker build -t $BASIC_IMAGE_NAME:$BASIC_IMAGE_TAG -f docker/Dockerfile_basic .
  ;;
  1)
    read -p "请输入容器名称: " container_name
    if docker ps -a --format '{{.Names}}' | grep -q $container_name; then
      read -p "容器$container_name 已存在, 按回车键将进行删除..."
      docker rm -f $container_name
    fi

    # todo 考虑安全性问题: 将/workspace 挂载出来
    docker_run_cmd="docker run -itd --name $container_name --privileged=true"
    docker_run_cmd="$docker_run_cmd $BASIC_IMAGE_NAME:$BASIC_IMAGE_TAG /bin/bash"
    echo -e "创建容器的命令如下:\n\n\t$docker_run_cmd\n"
    read -p "按回车键开始执行该命令创建容器..."
    $docker_run_cmd
  ;;
  *)
    echo "输入无效"
    exit 1
  ;;
esac

echo $xxx