#!/bin/bash

# 拉取镜像
# docker pull mysql:8

# 运行
docker run -d --name mydb \
  -p 3306:3306 \
  -e MYSQL_ROOT_PASSWORD=root mysql:8

# docker exec –it mydb /bin/bash
# mysql –u root –p <enter> root
# ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'root';
# LAN conn: 192.168.15.42:3306 root@root