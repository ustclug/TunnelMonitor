FROM smartentry/alpine:3.4-0.3.5

MAINTAINER Yifan Gao <docker@yfgao.com>

ADD . /opt/tunnelmonitor

ENV ASSETS_DIR=/opt/tunnelmonitor/.docker

RUN smartentry.sh build
