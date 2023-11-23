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

命令 docker

```bash
docker images
docker exec -it $CONTAINER_NAME /bin/bash
docker build -t $IMAGE_NAME:$IMAGE_TAG .
```

## jhl

jfrog: https://jfrog.jhlfund.com/ui/packages

## vscode

插件:

- `Docker`; `Remote - Containers`;
- `Prettier - Code formatter`
