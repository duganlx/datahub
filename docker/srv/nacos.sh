#!/bin/bash
# nacos文档: https://nacos.io/zh-cn/docs/quick-start.html

log_dir=/root/github/nacos/logs/
conf_dir=/root/github/nacos/conf/

# 启动容器
# == begin ==
docker run -d --name tmp -p 8848:8848 nacos/nacos-server
docker cp tmp:/home/nacos/logs/ /root/github/nacos/logs/
docker cp tmp:/home/nacos/conf/ /root/github/nacos/conf/
docker rm -f tmp
docker run -d --name nacos \
  -p 8848:8848 -p 9848:9848 -p 9849:9849 --privileged=true \
  -e JVM_XMS=256m -e JVM_XMX=256m -e MODE=standalone \
  -v /root/github/nacos/logs/:/home/nacos/logs/ \
  -v /root/github/nacos/conf/:/home/nacos/conf/ \
  --restart=always nacos/nacos-server
# == end ==

# [info] 在mysql的数据库 nacosconf 执行脚本 conf/mysql-schema.sql
# [info] 在 conf/application.properties 中修改数据库配置 {spring.sql.init.platform: mysql, db.num: 1, db.url.0: {MYSQL_SERVICE_HOST: 192.168.15.42, MYSQL_SERVICE_PORT: 3306, MYSQL_SERVICE_DB_NAME: nacosconf}, db.user.0: root, db.password.0: root}
# [info] http://192.168.15.42:8848/nacos/index.html