#!/bin/bash
# nacos文档: https://nacos.io/zh-cn/docs/quick-start.html

script_dir=$(dirname "$(readlink -f "$0")")
microsrv_dir="$(dirname "${script_path}")"
nacos_conf_dir=$microsrv_dir/tmp/nacos
log_dir=$microsrv_dir/tmp/nacos/logs/
conf_dir=$microsrv_dir/tmp/nacos/conf/

echo -e "操作引导:\n0 [拉取镜像]\n1 [初始化]\n2 [创建容器]"
read -p "选择进行的操作: " op

case $op in
  0)
    # 拉去镜像
    docker pull nacos/nacos-server
    ;;
  1)
    # 初始化 nacos
    mkdir -p $log_dir $conf_dir
    docker run -d --name tmp nacos/nacos-server
    docker cp tmp:/home/nacos/logs/ $nacos_conf_dir
    docker cp tmp:/home/nacos/conf/ $nacos_conf_dir
    docker rm -f tmp

    # 后续操作
    # [info] 在mysql的数据库 nacosconf 执行脚本 conf/mysql-schema.sql
    # [info] 在 conf/application.properties 中修改数据库配置 如下
    # {
    #   spring.sql.init.platform: mysql, 
    #   db.num: 1, 
    #   db.url.0: {MYSQL_SERVICE_HOST: 192.168.15.42, MYSQL_SERVICE_PORT: 3306, MYSQL_SERVICE_DB_NAME: nacosconf}, 
    #   db.user.0: root, 
    #   db.password.0: root
    # }
    ;;
  2)
    # 创建容器
    # --restart=always
    docker run -d --name nacos --privileged=true \
      -p 8848:8848 -p 9848:9848 -p 9849:9849 \
      -e JVM_XMS=256m -e JVM_XMX=256m -e MODE=standalone \
      -v $log_dir:/home/nacos/logs/ \
      -v $conf_dir:/home/nacos/conf/ \
      nacos/nacos-server
    
    # nacos访问: http://192.168.15.42:8848/nacos/index.html
    ;;
  *)
    echo "输入无效"
    exit 1
    ;;
esac


