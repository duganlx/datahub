#!/bin/bash
#
# Author: lvx
# Date: 2023-11-25
# Description: 搭建开发环境的 Main 程序

SCRIPT_DIR=$(dirname "$(readlink -f "$0")")
TMP_DIR=$SCRIPT_DIR/tmp
BASIC_IMAGE_NAME=devbasic
BASIC_IMAGE_TAG=v1.0.0

cat <<EOF
操作引导
  0 [生成basic镜像]
  1 [生成basic容器]
  2 [生成go镜像] - Deprecated 
  3 [生成golang容器]
  4 [生成nodejs镜像] - Deprecated 
  5 [生成nodejs容器]
  6 [生成mysql容器]
  7 [生成python容器]
EOF
read -p "选择进行的操作: " opt

case $opt in
  0)
    # 有三种情况: 1. 无镜像; 2. 有镜像无容器; 3. 有镜像有容器
    container_ids=$(docker ps -q --filter ancestor="$BASIC_IMAGE_NAME:$BASIC_IMAGE_TAG")
    if [ -n "$container_ids" ]; then
      echo -e "镜像$BASIC_IMAGE_NAME:$BASIC_IMAGE_TAG 存在如下容器:\n$container_ids"
      read -p "按回车键将进行删除..."
      for id in $container_ids; do
        docker rm -f $id
      done
    fi

    docker build -t $BASIC_IMAGE_NAME:$BASIC_IMAGE_TAG -f Dockerfile .
  ;;
  1)
    read -p "请输入容器名称: " container_name
    if docker ps -a --format '{{.Names}}' | grep -q $container_name; then
      read -p "容器$container_name 已存在, 按回车键将进行删除..."
      docker rm -f $container_name
    fi

    # 数据安全性考虑 /workspace 目录进行挂载
    container_wsdir=$TMP_DIR/workspace/$container_name
    mkdir -p $container_wsdir
    docker_run_cmd="docker run -itd --name $container_name --privileged=true"
    docker_run_cmd="$docker_run_cmd -v $container_wsdir:/workspace"
    docker_run_cmd="$docker_run_cmd $BASIC_IMAGE_NAME:$BASIC_IMAGE_TAG /bin/bash"
    echo -e "创建容器的命令如下:\n\n\t$docker_run_cmd\n"
    read -p "按回车键开始执行该命令创建容器..."
    $docker_run_cmd
  ;;
  3)
    # 在basic镜像上生成 golang 开发容器
    # 生成basic容器
    # == begin ==
    read -p "请输入容器名称: " container_name
    if docker ps -a --format '{{.Names}}' | grep -q $container_name; then
      read -p "容器$container_name 已存在, 按回车键将进行删除..."
      docker rm -f $container_name
    fi
    read -p "请输入容器$container_name 的端口映射配置(空格分隔): " ports_map
    port_map_array=($ports_map)
    
    container_wsdir=$TMP_DIR/workspace/$container_name
    mkdir -p $container_wsdir
    
    docker_run_cmd="docker run -itd --name $container_name --privileged=true"
    for port_map in ${port_map_array[@]}; do
      docker_run_cmd="$docker_run_cmd -p $port_map"
    done
    docker_run_cmd="$docker_run_cmd -v $container_wsdir:/workspace"
    docker_run_cmd="$docker_run_cmd -v /root/.ssh:/root/.ssh"
    docker_run_cmd="$docker_run_cmd $BASIC_IMAGE_NAME:$BASIC_IMAGE_TAG /bin/bash"

    echo -e "将执行的命令如下所示\n\n\t$docker_run_cmd\n"
    read -p "按回车键开始执行该命令创建容器..."
    $docker_run_cmd
    # == end ==

    # golang 版本
    echo -e "可选go版本:\n  0 [go1.18.10]\n  1 [go1.20.11]" 
    read -p "版本选择: " goopt
    gozip=""
    case $goopt in
      0)
        gozip="go1.18.10.linux-amd64.tar.gz"
      ;;
      1) 
        gozip="go1.20.11.linux-amd64.tar.gz"
      ;;
      *)
        echo "输入无效"
        exit 2
       ;;
    esac

    # 资源准备
    download_dir=$TMP_DIR/download/go
    mkdir -p $download_dir
    if [ ! -e "$download_dir/$gozip" ]; then
      wget -P "$download_dir" https://golang.google.cn/dl/$gozip
    fi
    if [ ! -e "$download_dir/protoc-22.2-linux-x86_64.zip" ]; then 
      wget -P "$download_dir" https://github.com/protocolbuffers/protobuf/releases/download/v22.2/protoc-22.2-linux-x86_64.zip
    fi

    cat <<EOT > "$download_dir/inrun.sh"
#!/bin/bash

tar -C /usr/local -zxf /download/$gozip
mkdir -p /root/go/bin /root/go/pkg

# 环境变量配置
# == begin ==
echo 'export PATH=\$PATH:/usr/local/go/bin' >> /root/.bashrc
echo 'export GOBIN=/root/go/bin' >> /root/.bashrc
echo 'export GOPROXY=https://goproxy.cn,direct' >> /root/.bashrc
echo 'export GOSUMDB=sum.golang.google.cn' >> /root/.bashrc
echo 'export PATH=\$PATH:\$GOBIN' >> /root/.bashrc

export PATH=\$PATH:/usr/local/go/bin
export GOBIN=/root/go/bin
export GOPROXY=https://goproxy.cn,direct
export PATH=\$PATH:\$GOBIN
# == end ==

unzip -d /download/tmp /download/protoc-22.2-linux-x86_64.zip
mv /download/tmp/bin/protoc /root/go/bin
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

    docker cp $download_dir/$gozip $container_name:/download
    docker cp $download_dir/inrun.sh $container_name:/download
    docker cp $download_dir/protoc-22.2-linux-x86_64.zip $container_name:/download

    echo "进行容器内的配置..."
    docker exec -it $container_name /bin/bash -c 'chmod 750 /download/inrun.sh'
    docker exec -it $container_name /bin/bash -c 'bash /download/inrun.sh'
  ;;
  5)
    # 在basic镜像上生成 nodejs 开发容器
    # 生成basic容器
    # == begin ==
    read -p "请输入容器名称: " container_name
    if docker ps -a --format '{{.Names}}' | grep -q $container_name; then
      read -p "容器$container_name 已存在, 按回车键将进行删除..."
      docker rm -f $container_name
    fi
    read -p "请输入容器$container_name 的端口映射配置(空格分隔): " ports_map
    port_map_array=($ports_map)

    container_wsdir=$TMP_DIR/workspace/$container_name
    mkdir -p $container_wsdir

    docker_run_cmd="docker run -itd --name $container_name --privileged=true"
    for port_map in ${port_map_array[@]}; do
      docker_run_cmd="$docker_run_cmd -p $port_map"
    done
    docker_run_cmd="$docker_run_cmd -v $container_wsdir:/workspace"
    docker_run_cmd="$docker_run_cmd -v /root/.ssh:/root/.ssh"
    docker_run_cmd="$docker_run_cmd $BASIC_IMAGE_NAME:$BASIC_IMAGE_TAG /bin/bash"

    echo -e "将执行的命令如下所示\n\n\t$docker_run_cmd\n"
    read -p "按回车键开始执行该命令创建容器..."
    $docker_run_cmd
    # == end ==

    # nodejs 版本 https://nodejs.org/dist/
    echo -e "可选nodejs版本:\n  0 [v16.20.2]\n  1 [v18.9.1]" 
    read -p "版本选择: " nodeopt
    nodearr=""
    nodezip=""
    nodeunzip=""
    case $nodeopt in
      0)
        nodearr="https://nodejs.org/dist/v16.20.2/node-v16.20.2-linux-x64.tar.gz"
        nodezip="node-v16.20.2-linux-x64.tar.gz"
        nodeunzip="node-v16.20.2-linux-x64"
      ;;
      1) 
        nodearr="https://nodejs.org/dist/v18.9.1/node-v18.9.1-linux-x64.tar.gz"
        nodezip="node-v18.9.1-linux-x64.tar.gz"
        nodeunzip="node-v18.9.1-linux-x64"
      ;;
      *)
        echo "输入无效"
        exit 2
       ;;
    esac

    # 资源准备
    download_dir=$TMP_DIR/download/nodejs
    mkdir -p $download_dir
    if [ ! -e "$download_dir/$nodezip" ]; then
      wget -P "$download_dir" $nodearr
    fi

    docker cp $download_dir/$nodezip $container_name:/download
    docker exec -it $container_name /bin/bash -c "tar -C /usr/local -zxf /download/$nodezip"
    # 不推荐用软链的方式, 因为后续如果用 npm 全局安装的命令 仍然需要软链 (ln -s src dist)
    docker exec -it $container_name /bin/bash -c "echo 'export PATH=\$PATH:/usr/local/$nodeunzip/bin' >> /root/.bashrc"

    echo -e "安装yarn\n"
    docker exec -it $container_name /bin/bash -c "export PATH=\$PATH:/usr/local/$nodeunzip/bin && npm install -g yarn"
  ;;
  6)
    read -p "请输入容器名称: " container_name
    if docker ps -a --format '{{.Names}}' | grep -q $container_name; then
      read -p "容器$container_name 已存在, 按回车键将进行删除..."
      docker rm -f $container_name
    fi

    if ! docker images --format "{{.Repository}}:{{.Tag}}" | grep -q "mysql:8"; then
      docker pull mysql:8
    fi

    docker run -d --name $container_name -p 3306:3306 -e MYSQL_ROOT_PASSWORD=root mysql:8
    echo -e "mysql容器$container_name 创建完成, 还需要进行如下配置:"
    echo -e "\n\tdocker exec -it mydb /bin/bash"
    echo -e "\tmysql -u root -p <enter> root"
    echo -e "\tALTER USER 'root'@'%' IDENTIFIED WITH mysql_native_password BY 'root';\n"
  ;;
  7)
    # 生成 python 开发环境 - miniconda
    # 生成basic容器
    # == begin ==
    read -p "请输入容器名称: " container_name
    if docker ps -a --format '{{.Names}}' | grep -q $container_name; then
      read -p "容器$container_name 已存在, 按回车键将进行删除..."
      docker rm -f $container_name
    fi
    read -p "请输入容器$container_name 的端口映射配置(空格分隔): " ports_map
    port_map_array=($ports_map)

    container_wsdir=$TMP_DIR/workspace/$container_name
    mkdir -p $container_wsdir

    docker_run_cmd="docker run -itd --name $container_name --privileged=true"
    for port_map in ${port_map_array[@]}; do
      docker_run_cmd="$docker_run_cmd -p $port_map"
    done
    docker_run_cmd="$docker_run_cmd -v $container_wsdir:/workspace"
    docker_run_cmd="$docker_run_cmd -v /root/.ssh:/root/.ssh"
    docker_run_cmd="$docker_run_cmd $BASIC_IMAGE_NAME:$BASIC_IMAGE_TAG /bin/bash"

    echo -e "将执行的命令如下所示\n\n\t$docker_run_cmd\n"
    read -p "按回车键开始执行该命令创建容器..."
    $docker_run_cmd
    # == end ==

    # 资源准备
    download_dir=$TMP_DIR/download/py
    condash="Miniconda3-latest-Linux-x86_64.sh"
    mkdir -p $download_dir
    if [ ! -e "$download_dir/$condash" ]; then
      wget -P "$download_dir" "https://repo.anaconda.com/miniconda/$condash"
    fi

    docker cp $download_dir/$condash $container_name:/download
    docker exec -it $container_name /bin/bash -c "chmod 750 /download/$condash"
    docker exec -it $container_name /bin/bash -c "bash /download/$condash -b -p /usr/local/miniconda"
    docker exec -it $container_name /bin/bash -c "echo 'export PATH=\$PATH:/usr/local/miniconda/bin' >> /root/.bashrc"
  ;;
  *)
    echo "输入无效"
    exit 1
  ;;
esac

# 清理tag=none的镜像
# == begin ==
found=false
create_images_opts=(0 2 4)
for num in "${create_images_opts[@]}"
do
  if [ "$num" -eq $opt ]; then
    found=true
    break
  fi
done
if [ "$found" = true ]; then
  image_ids=$(docker images -f "dangling=true" -q)
  num=$(docker images -f "dangling=true" -q | wc -l)
  if [ $num -ne 0 ]; then
    echo -e "\n清除无效镜像个数为$num"
    docker rmi -f $image_ids 
  fi
fi
# == end ==

# 进入容器命令
# == begin ==
found=false
create_container_opts=(1 3 5 7)
for num in "${create_container_opts[@]}"
do
  if [ "$num" -eq $opt ]; then
    found=true
    break
  fi
done
if [ "$found" = true ]; then
  echo -e "容器$container_name 创建完成，进入容器命令为: docker exec -it $container_name /bin/bash"
fi
# == end ==