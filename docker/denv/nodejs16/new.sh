#!/bin/bash

IMAGE_NAME=envjs16
IMAGE_TAG=v1

read -p "请输入容器名称: " container_name
if docker ps -a --format '{{.Names}}' | grep -q $container_name; then
  read -p "容器$container_name 已存在, 按回车键将进行删除..."
  docker rm -f $container_name
fi
read -p "请输入容器$container_name 的端口映射配置(空格分隔): " ports_map
port_map_array=($ports_map)

# 构建容器生成命令
# == begin ==
docker_run_cmd="docker run -itd --name $container_name --privileged=true"
for port_map in ${port_map_array[@]}; do
  docker_run_cmd="$docker_run_cmd -p $port_map"
done
docker_run_cmd="$docker_run_cmd -v /root/.ssh:/root/.ssh"
docker_run_cmd="$docker_run_cmd -v /root/share:/share"
docker_run_cmd="$docker_run_cmd $IMAGE_NAME:$IMAGE_TAG /bin/bash"
# == end ==

echo -e "将执行的命令如下所示\n\n\t$docker_run_cmd\n"
read -p "按回车键开始执行该命令创建容器..."
$docker_run_cmd

echo -e "安装yarn\n"
docker exec -it $container_name /bin/bash -c 'export PATH=$PATH:$NODE_BIN && npm install -g yarn'

echo -e "容器$container_name 创建完成，进入容器命令为: docker exec -it $container_name /bin/bash"
# docker exec -it tenvjs /bin/bash
# docker rm -f tenvjs
# npm install -g yarn