#!/bin/bash
# 配置kratos框架开发环境
# 该脚本可携带一个参数
# @param1: 表示是否下载 protoc压缩包, 如果不需要下载携带 skipdownproto

# 环境变量配置 /root/.bashrc
# == begin ==
export PATH=$PATH:/usr/local/go/bin
export GOPATH=/workspace/go
export GOBIN=$GOPATH/bin
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=sum.golang.google.cn
export PATH=$PATH:$GOBIN
# == end ==

DOWNLOAD_DIR=/workspace
DOWNLOAD_DIR=/download

# 安装: protoc 
# == begin ==
cd $DOWNLOAD_DIR 
if [ "$1" == "skipdownproto" ]; then
  echo "默认你已下载好protoc-22.2-linux-x86_64.zip文件"
else
  wget https://github.com/protocolbuffers/protobuf/releases/download/v22.2/protoc-22.2-linux-x86_64.zip
fi 
unzip -d $DOWNLOAD_DIR/tmp protoc-22.2-linux-x86_64.zip
mv $DOWNLOAD_DIR/tmp/bin/protoc $WORKSPACE_DIR/go/bin/
rm -rf $DOWNLOAD_DIR/tmp
# == end ==

# 安装: kratos 依赖
# == begin ==
go install github.com/go-kratos/kratos/cmd/kratos/v2@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@latest
go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@latest
go install github.com/google/gnostic/cmd/protoc-gen-openapi@latest
go install github.com/google/wire/cmd/wire@latest
go install github.com/envoyproxy/protoc-gen-validate@latest
# == end ==

# 配置: jhl私仓配置
# == begin ==
go env -w GOPRIVATE=gitlab.jhlfund.com
go env -w GONOPROXY=gitlab.jhlfund.com
# go env -w GONOSUBDB=gitlab.jhlfund.com
go env -w GOINSECURE=gitlab.jhlfund.com
# == end ==