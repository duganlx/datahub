# FROM ubuntu:jammy-20230308
FROM jfrog.jhlfund.com/docker/eam/ubuntu:22.04
LABEL author=duganlx email=lvxiang@jhlfund.com

SHELL ["/bin/bash", "-c"]

ENV DOWNLOAD_DIR=/download
ENV WORKSPACE_DIR=/workspace
ENV ROOT_DEFAULT_PASSWD=root

RUN mkdir ${DOWNLOAD_DIR} ${WORKSPACE_DIR}
RUN apt-get update
RUN apt-get install -y git