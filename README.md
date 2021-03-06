# 概述
sparrow是一个分布式服务器框架，采用golang编写，只实现了网络和路由部分，其他组件根据自己需求自行添加就是

sparrow支持分布式和单机部署，分布式需要消息中间件(nats)支持，配置文件gate_config.json中的Distributed字段决定是否采用中间件分布式架构,。

sparrow的基本思想是，对同步性以及消息顺序要求不高的请求，采用消息中间件(nats)转发到微服务处理。 对同步性和消息顺序要求高的请求，建议采用grpc的方式处理，[grpc负载均衡方法](https://blog.csdn.net/weixin_43733451/article/details/84262506)

sparrow支持tcp和websocket两种传输协议,相关配置在gate_config.json配置文件里面，
demo里有本地网关直接返回处理，mq微服务处理，grpc处理三种方式

# 消息格式
## tcp:
包头5个字节，4个字节的消息总长度，1个字节的tag长度

包头后面紧跟tag内容和data内容

**| 消息总长度(4字节) | tag长度(1字节) | tag | data |**

## websocket:
由于websocket自身有消息长度检查，所以包头不包含消息总长度，所以组成为
包头就1个字节的tag长度，后面紧跟tag内容和data内容

**|tag长度(1字节) | tag | data |**


# 配置文件说明
ServerID 网关id,分布式识别网关用的

TcpPort tcp监听端口

WebsocketPort websocket监听端口

MaxConnection 最大连接数

WirteQueLen 子线程的写入缓冲队列大小

MaxMsgLen 消息最大长度

BigEndian 是否采用大端序

Distributed 是否采用中间件分布式

PublisherNum 发布队列大小

PublishAddr 发布消息中间件地址

SubcriberNum 订阅队列大小

SubcribAddr 订阅消息中间件地址


# 安装sparrow
go get github.com/qianlidongfeng/sparrow/gate

或

git clone https://github.com/qianlidongfeng/sparrow.git 克隆到gopath相关目录下


# 安装消息中间件(nats)
[nats安装方法](https://www.nats.io/documentation/managing_the_server/installing/)

go get github.com/nats-io/gnatsd

cd github.com/nats-io/gnatsd

go install

nats是一款消息中间件，由golang编写,去中心化，支持分布式横向扩展，能非常方便实现微服务负载均衡。速度极快，官网测试是目前市面上最快的。缺点是消息无法持久化。

# 运行demo

在安装好sparrow和消息中间件的前提下

cd example

安装网关服务器
go install gate.go

安装mq微服务
go install mqServer.go

安装grpc服务
go install grpcServer.go

安装测试客户端
go install sparrowClient.go

## 启动相关程序
按默认配置启动消息中间件(启动命令gnatsd),默认为4222端口

复制example/gate_config.json到bin目录，启动gate,mqServer,grpcServer

### tcp

启动客户端sparrowClient,你会看到打印出服务器返回的消息,返回消息的tag和内容:

tag:NORMAL

data:hello I am sparrow

tag:GRPC

data:hello I am sparrow

tag:MQ

data:hello I am sparrow


### websocket

浏览器运行 example/websocketClient.html,你会看到打印和上面相同的效果

