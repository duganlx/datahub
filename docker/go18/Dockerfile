FROM jfrog.jhlfund.com/docker/eam/ubuntu:22.04
LABEL author=duganlx email=lvxiang@jhlfund.com

SHELL ["/bin/bash", "-c"]

ENV DOWNLOAD_DIR=/download
ENV WORKSPACE_DIR=/workspace
ENV SHARE_DIR=/share
ENV LANG=C.UTF-8

RUN mkdir ${DOWNLOAD_DIR} ${WORKSPACE_DIR}
RUN apt-get update
RUN apt-get install -y vim net-tools iputils-ping wget git
RUN apt-get install -y autoconf automake libtool curl make g++ unzip

# 安装: golang -> https://golang.google.cn/dl/
# == begin ==
RUN cd ${DOWNLOAD_DIR} && wget https://golang.google.cn/dl/go1.18.10.linux-amd64.tar.gz
RUN tar -C /usr/local -xzf ${DOWNLOAD_DIR}/go1.18.10.linux-amd64.tar.gz
RUN mkdir -p ${WORKSPACE_DIR}/go/bin ${WORKSPACE_DIR}/go/pkg ${WORKSPACE_DIR}/go/src
RUN echo 'export PATH=$PATH:/usr/local/go/bin' >> /root/.bashrc
RUN echo 'export GOPATH=/workspace/go' >> /root/.bashrc
RUN echo 'export GOBIN=$GOPATH/bin' >> /root/.bashrc
RUN echo 'export GOPROXY=https://goproxy.cn,direct' >> /root/.bashrc
RUN echo 'export GOSUMDB=sum.golang.google.cn' >> /root/.bashrc
RUN echo 'export PATH=$PATH:$GOBIN' >> /root/.bashrc
# == end ==

# 解决: vim中不允许输入中文问题
# == begin ==
RUN cd /root && touch .vimrc
RUN echo 'set fileencodings=utf-8,cp936,gb2312,gbk,gb18030' >> /root/.vimrc
RUN echo 'set termencoding=utf-8' >> /root/.vimrc
RUN echo 'set fileformats=unix' >> /root/.vimrc
RUN echo 'set encoding=prc' >> /root/.vimrc
# == end ==

# 解决: 终端不允许输入中文问题
# == begin ==
RUN cd /root && touch .inputrc
RUN echo 'set meta-flag on' >> /root/.inputrc
RUN echo 'set convert-meta off' >> /root/.inputrc
RUN echo 'set input-meta on' >> /root/.inputrc
RUN echo 'set output-meta on' >> /root/.inputrc
# == end ==

WORKDIR ${WORKSPACE_DIR}