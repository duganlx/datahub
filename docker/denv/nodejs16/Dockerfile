# FROM ubuntu:jammy-20230308
FROM jfrog.jhlfund.com/docker/eam/ubuntu:22.04
LABEL author=duganlx email=lvxiang@jhlfund.com

SHELL ["/bin/bash", "-c"]

ENV DOWNLOAD_DIR=/download
ENV WORKSPACE_DIR=/workspace
ENV SHARE_DIR=/share

RUN mkdir ${DOWNLOAD_DIR} ${WORKSPACE_DIR}
RUN apt-get update
RUN apt-get install -y vim net-tools iputils-ping wget git

# 安装: nodejs & npm -> https://nodejs.org/dist/
# == begin ==
RUN cd ${DOWNLOAD_DIR} && wget https://nodejs.org/dist/v16.20.2/node-v16.20.2-linux-x64.tar.gz
RUN cd ${DOWNLOAD_DIR} && tar -xvf node-v16.20.2-linux-x64.tar.gz
ENV NODE_BIN='/download/node-v16.20.2-linux-x64/bin'
# 不推荐用软链的方式, 因为后续如果用 npm 全局安装的命令 仍然需要软链
# RUN ln -s ${DOWNLOAD_DIR}/node-v16.20.2-linux-x64/bin/node /usr/local/bin/node
# RUN ln -s ${DOWNLOAD_DIR}/node-v16.20.2-linux-x64/bin/npm /usr/local/bin/npm
RUN echo 'export PATH=$PATH:'${NODE_BIN} >> /root/.bashrc
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
# todo 删除时仍然是单字节删除，所以汉字需要两次删除键
# == begin ==
RUN cd /root && touch .inputrc
RUN echo 'set meta-flag on' >> /root/.inputrc
RUN echo 'set convert-meta off' >> /root/.inputrc
RUN echo 'set input-meta on' >> /root/.inputrc
RUN echo 'set output-meta on' >> /root/.inputrc
# == end ==

WORKDIR ${WORKSPACE_DIR}