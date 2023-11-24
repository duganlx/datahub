#!/bin/bash
# 
# Author: lvx
# Description: 微服务测试 - go 开发环境

# common
image_name=devgo18
image_tag=v1
container_name=ncgodev
script_path="$(dirname "$(readlink -f "$0")")"
microsrv_dir="$(dirname "${script_path}")"
src_dir=$microsrv_dir/src

# 创建容器
# == begin ==
if docker ps -a --format '{{.Names}}' | grep -q $container_name; then
  read -p "go开发环境已存在, 按回车键将进行删除并重新构建..."
  docker rm -f $container_name
fi

read -p "请输入容器$container_name 的端口映射配置(空格分隔): " ports_map
port_map_array=($ports_map)

docker_run_cmd="docker run -itd --name $container_name --privileged=true"
for port_map in ${port_map_array[@]}; do
  docker_run_cmd="$docker_run_cmd -p $port_map"
done
docker_run_cmd="$docker_run_cmd -v /root/.ssh:/root/.ssh"
docker_run_cmd="$docker_run_cmd -v $src_dir:/workspace/go/src"
docker_run_cmd="$docker_run_cmd $image_name:$image_tag /bin/bash"

echo -e "将执行的命令如下所示\n\n\t$docker_run_cmd\n"
read -p "按回车键开始执行该命令创建容器..."
$docker_run_cmd

echo -e "容器$container_name 创建完成，开始进行初始化..."
# == end ==

if [ ! -e "$microsrv_dir/tmp/protoc-22.2-linux-x86_64.zip" ]; then
  wget -P "$microsrv_dir/tmp" https://github.com/protocolbuffers/protobuf/releases/download/v22.2/protoc-22.2-linux-x86_64.zip
fi

docker cp $microsrv_dir/tmp/protoc-22.2-linux-x86_64.zip $container_name:/download

# inrun.sh 创建&执行
# == begin ==
inrun_filename=inrun.sh

if [ ! -e "$microsrv_dir/tmp/$inrun_filename" ]; then
  touch "$microsrv_dir/tmp/$inrun_filename"
fi

cat <<EOT > "$microsrv_dir/tmp/$inrun_filename"
export PATH=\$PATH:/usr/local/go/bin
export GOPATH=/workspace/go
export GOBIN=\$GOPATH/bin
export GOPROXY=https://goproxy.cn,direct
export GOSUMDB=sum.golang.google.cn
export PATH=\$PATH:\$GOBIN

unzip -d /download/tmp /download/protoc-22.2-linux-x86_64.zip
mv /download/tmp/bin/protoc /workspace/go/bin/
rm -rf /download/tmp

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
EOT

docker cp $microsrv_dir/tmp/$inrun_filename $container_name:/download
docker exec -it $container_name /bin/bash -c 'chmod 750 /download/'$inrun_filename
docker exec -it $container_name /bin/bash -c 'bash /download/'$inrun_filename
# == end ==

echo -e "容器$container_name 初始化完成，进入容器命令为: docker exec -it $container_name /bin/bash"
