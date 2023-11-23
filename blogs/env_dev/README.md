# 运维 develop

windows + wsl2 + docker desktop

不能直接在 wsl2 中进行开发，会出现 LAN 中其他机器无法访问本机 wsl2 中的服务，目前采用 docker 容器内启动服务提供服务 W

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
