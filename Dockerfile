#FROM golang:1.17.10 AS builder
#
## ENV GOPROXY      https://goproxy.io
#
#RUN mkdir /app
#ADD . /app/
#WORKDIR /app
#RUN go build -o wechatbot .
#
#FROM centos:centos7
#RUN mkdir /app
#WORKDIR /app
#COPY --from=builder /app/ .
#RUN chmod +x wechatbot && cp config.dev.json config.json && yum -y install vim net-tools telnet wget curl && yum clean all
#
#CMD ./wechatbot

# wechatbot/Dockerfile

# 下面是第二阶段的镜像构建，和之前保持一致
FROM registry.cn-beijing.aliyuncs.com/zjcbit/wechat-chatgpt:base

# 和上个阶段一样设置工作目录
WORKDIR /app

# 而是从上一个阶段构建的 builder容器中拉取
COPY wechatbot /app/
ADD supervisord.conf /etc/supervisord.conf
ADD config.dev.json /app/config.dev.json
RUN cp config.dev.json config.json

# 通过 Supervisor 管理服务
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]