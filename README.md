# 概述
sparrow是一个分布式服务器框架，采用golang编写

sparrow采用消息中间件(nats)转发网关和微服务之间的消息，消息采用protobuf传输。 至于网关和客户端之间自定义即可，建议消息data部分用protobuf格式

sparrow支持分布式和单机部署，分布式需要消息中间件支持，gate模块的RigsterLocalMsg是注册本地消息，由网关直接处理。RigsterMqMsg注册的消息则由消息中间件转发给微服务处理。 如果只需要单机部署，统统注册成本地消息即可。

sparrow支持tcp和websocket两种传输协议,相关配置在gate_config.json配置文件里面

# 消息格式
包头5个字节，4个字节的消息总长度，1个字节的tag长度

包头后面紧跟tag内容和data内容

**| 消息总长度(4字节) | tag长度(1字节) | tag | data |**
# 安装sparrow
go install github.com/qianlidongfeng/sparrow
# 安装消息中间件(nats)
[nats安装方法](https://www.nats.io/documentation/managing_the_server/installing/)

# 运行demo

cd example/chatroom

安装聊天服务器
go install chatroom.go

安装微服务
go install microserver.go

安装测试客户端
go install sparrow_client.go

## 启动相关程序
按默认配置启动消息中间件(启动命令gnatsd),默认为4222端口

复制example/chatroom/gate_config.json到bin目录，启动chatroom,microserver

启动客户端测试demo,sparrow_client,你会看到打印出服务器返回的消息,返回消息的tag和内容，分别为respone, I am sparrow
