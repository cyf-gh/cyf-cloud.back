# cyf-cloud.back

cyf-cloud的服务器端，使用Golang编写。

## 基本

### API 版本
* v1（停止维护）
* v1x1 

### 依赖
#### 数据库
* SQLite3

### 配置文件 
#### server.cfg
```
[server]
    udp=":2333"			# UDP监听端口号
    tcp=":2334"			# TCP监听端口号
    log="127.0.0.1:2335" # 日志服务
[fresh_interval]
    log="2000"			# 写日志频率
    udp="1"				# UDP暂歇频率，单位毫秒
```
