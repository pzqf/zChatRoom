# zChatRoom

**注意事项**
1. 仅支持在linux下运行,win下客户端界面有问题。
2. go版本：1.18

## 简介
这是一个在线聊天室，通信基于TCP/IP协议。有简单的字符界面，目的在于学习。

## 使用

### 编译

进入代码目录，新建bin目录，分别对服务器和客户端进行编译
```
cd zChatRoom
mkdir -p bin
go build  -o ./bin/server ./ChatServer/main.go
go build  -o ./bin/client ./ChatClient/main.go
```
以上过程可以使用make linux命令进行编译
```
cd zChatRoom
make linux
```

### 运行

通信使用的是9106端口。

#### 服务端
将敏感词文件keyword.txt,及分词文件dictionary.txt放于bin目录下。然后直接运行
```
./bin/server
```

#### 客户端
客户端运行需指定服务器地址，默认地址为127.0.0.1
```
./bin/client -a 192.168.50.100
```

## 性能
未做大规模性能测试，预计在线一万人，应该没有压力，通信库经历过12万客户端压力测试。

## 代码解读
整个项目代码量较少，阅读简单。目录及各模块如下：

#### zChatRoom 项目框架
1. bin 运行程序目录

2. ChatClient 客户端代码  

    2.1. cui 界面及界面管理模块  
    2.2. handler 通信响应处理模块  
    2.3. model 数据交互模块  

3. ChatServer 服务端代码

    3.1. gm GM命令处理模块  
    3.2. handler 通信响应处理模块  
    3.3. player 玩家模块  
    3.4. playerMgr 玩家管理模块  
    3.5. room 聊到室及聊天室管理模块  
    3.6. segmenter 分词模块  

5. proto 通信协议模块

#### 代码流程是
1. 初始化分词处理器，初始化敏感词过滤器，初始化玩家，房间管理器。
2. 初始化网络，注册协议处理函数，初始化网络包处理模式及包的最大大小。
3. 开启服务监听。
4. 优雅关闭模块监听指令。

#### 关健算法
1. 网络通信模块，服务器单协程监听，建立新的连接后，创建Session模块,每个连接一个Session，Session异步处理数据的接收和发送。
收到数据后（已做粘包处理），会进行解包，封装成NetPacket数据包。包里会根据设置的通信协议方式进行解包，目前支持，json,gob,二进制，默认是JSON。
自定义的可使用二进制模式，解包后，就会根据协议号调用注册的数据处理函数。传递数据给业务层，此处采用协程池处理。handler协程，
未采用content,建议在调用新的协程处理时使用上下文。此处不强制。业务层发送数据，可直接调用Session.Send接口，导步发送。
2. 敏感词过滤采用了DEA算法，详见[zKeyWordFilter](https://github.com/pzqf/zUtil/tree/develop/zKeyWordFilter)。

## 库说明
[zEngine](https://github.com/pzqf/zEngine) 该库是我自己写的，包含通信zNet，优雅关闭zSignal，具体代码可点击链接查看，将来会进行扩展

[zUtil](https://github.com/pzqf/zUtil) 该库是我自己写的常用工具集合，包含有敏感词过滤，时间处理，队列，高性能map等。

[sego](https://github.com/huichen/sego) 大牛写的分词处理器。

[ants](https://github.com/panjf2000/ants) 潘神写的协程池库。

## 未开发及不足之处
以下内容未开发，以后会增加。
1. 未处理日志。
2. 未处理数据存储，服务器重启就会丢失所有数据。
3. 分词库未做具体考量，不知是否合理，性能是否达标。
4. 业务层许多功能未做，如断线重连等。
5. 代码简单易懂，未做过多注释。
6. 心跳模块已支持，未开启。
7. 客户端界面粗糙，快捷健支持不足。