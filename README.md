# Kafka 消费演示样例
本示例程序主要用于展示，如何基于[物联使能应用托管](https://cloud.tencent.com/document/product/1081/50036)
接收和处理[物联网开发平台-数据同步](https://cloud.tencent.com/document/product/1081/47616)，功能转发的数据。


运行本样例，会启动一个 Kafka 消费端和 HTTP 服务。
服务启动后，从指定的 Kafka Topic中消费消息，并记录最近一次消费的消息。
通过访问 HTTP 服务，可以查看最近消费的消息。

## 上手指南
### 前提条件
#### 安装本地 Docker 编译环境

* [Windows](https://docs.docker.com/windows/started)
* [OS X](https://docs.docker.com/mac/started/)
* [Linux](https://docs.docker.com/linux/started/)


#### 准备镜像仓库
在腾讯云上创建一个私有的镜像仓库，假设对应的命名空间为 my-namespace（可自由定义），镜像名称为 kafka-consumer-logger 。
[腾讯云镜像仓库操作快速入门](https://cloud.tencent.com/document/product/1141/50332)

#### 配置 CKafka 转发
参见 [配置指导](https://cloud.tencent.com/document/product/1081/47616) 。

### 使用方法
#### 镜像准备
1. 打开 Makefile

```Makefile
TAG ?= latest
IMAGE_NAME ?= kafka-consumer-logger

# 配置 IMAGE_REPO 设置镜像仓库地址
IMAGE_REPO=ccr.ccs.tencentyun.com/my-namespace

# ...

```

2. 编译并推送

```shell
make && make push
```

#### 实际运行
1. [新建服务](https://cloud.tencent.com/document/product/1081/50044)，此步骤弹窗中的“镜像仓库”选项，选择“绑定已有腾讯云镜像仓库”，并选定前面推送的“my-namespace/kafka-consumer-logger”。
2. [部署服务](https://cloud.tencent.com/document/product/1081/50045)，此步骤弹窗中，对环境变量进行配置。
* `KAFKA_URL` - 填写需要消费的 Kafka 实例的“内网IP与端口”，例如：10.x.x.x:9092，请根据实际情况填写。
* `KAFKA_TOPIC` - 填写需要消费的 Kafka 实例的 topic 。
* `KAFKA_GROUP_ID` - 填写消费的 Kafka 消息时使用的消费组名称。
3. 打开“服务配置”中的“默认公网访问地址”，刷新页面查看消息，也可通过“日志”查看历史记录信息。
