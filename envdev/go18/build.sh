#!/bin/bash

IMAGE_NAME=devgo18
IMAGE_TAG=v1

# 删除原先的容器
# == begin ==
container_ids=$(docker ps -aq)
for id in $container_ids; do
  container_image=$(docker inspect --format='{{.Config.Image}}' $id)
  if [ "$container_image" == "$IMAGE_NAME:$IMAGE_TAG" ]; then
    echo "删除容器 $id"
    docker rm -f $id
  fi
done
# == end ==

# 构建镜像
docker build -t $IMAGE_NAME:$IMAGE_TAG .

# 清理tag=none的镜像
# == begin ==
image_ids=$(docker images -f "dangling=true" -q)
num=$(docker images -f "dangling=true" -q | wc -l)
if [[ $num == 0 ]]; then
  echo "没有无效镜像需要清除"
else
  echo "清除无效镜像个数为$num"
  docker rmi -f $image_ids 
fi
# == end ==
