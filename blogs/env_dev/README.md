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

## Go

golang 私仓配置

```bash
# 配置go mod 私有仓库
go env -w GOPRIVATE=gitlab.jhlfund.com
# 配置不使用代理
go env -w GONOPROXY=gitlab.jhlfund.com
# 配置不验证包(无用)
go env -w GONOSUBDB=gitlab.jhlfund.com
# 配置不加密访问
go env -w GOINSECURE=gitlab.jhlfund.com
```

go mod 命令

- tidy: add missing and remove unused modules

go 依赖包保存位置

- `/root/.cache/go-build`
- `/root/gowork/pkg/mod ---> GOPATH/pkg/mod`

todo 添加描述问题，关于 打 tag 引用内容不变问题

kratos 在 service 文件夹下添加了对应 go 文件之后，需要在 service.go 中进行注册`NewxxxSerive`，然后还需要去 `server/http.go` 文件中进行绑定。然后 运行 `make generate` 去生成 wire_gen.go 文件。这样才能通过 http 访问。
