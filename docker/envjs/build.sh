#!/bin/bash

IMAGE_NAME=envjs
IMAGE_TAG=v0.5

docker build -t $IMAGE_NAME:$IMAGE_TAG .