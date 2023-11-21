# 运维 develop

windows + wsl2 + docker desktop

不能直接在 wsl2 中进行开发，会出现 LAN 中其他机器无法访问本机 wsl2 中的服务，目前采用 docker 容器内启动服务提供服务

搭建 mysql

```
docker pull mysql:8
docker run -d -p 3306:3306 --name test -e MYSQL_ROOT_PASSWORD=test mysql:8
docker exec –it test /bin/bash
mysql –u root –p <enter> test
ALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'root’;

[info] LAN conn: 192.168.15.42:3306 root@root
```

搭建 nacos

文档: https://nacos.io/zh-cn/docs/quick-start.html

```
docker run -d --name tmp -p 8848:8848 nacos/nacos-server
docker cp tmp:/home/nacos/logs/ /root/github/nacos/logs/
docker cp tmp:/home/nacos/conf/ /root/github/nacos/conf/
docker rm -f tmp
docker run -d --name nacos -p 8848:8848 -p 9848:9848 -p 9849:9849 --privileged=true -e JVM_XMS=256m -e JVM_XMX=256m -e MODE=standalone -v /root/github/nacos/logs/:/home/nacos/logs/ -v /root/github/nacos/conf/:/home/nacos/conf/ --restart=always nacos/nacos-server

[info] 在mysql的数据库 nacosconf 执行脚本 conf/mysql-schema.sql
[info] 在 conf/application.properties 中修改数据库配置 {spring.sql.init.platform: mysql, db.num: 1, db.url.0: {MYSQL_SERVICE_HOST: 192.168.15.42, MYSQL_SERVICE_PORT: 3306, MYSQL_SERVICE_DB_NAME: nacosconf}, db.user.0: root, db.password.0: root}
[info] http://192.168.15.42:8848/nacos/index.html
```

命令 docker

```bash
docker images
docker exec -it $CONTAINER_NAME /bin/bash
docker build -t $IMAGE_NAME:$IMAGE_TAG .
```

# jhl

jfrog: https://jfrog.jhlfund.com/ui/packages
