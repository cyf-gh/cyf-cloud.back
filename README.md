# cyf-server.back

cyf-server的服务器端，使用Golang编写。

## 0 开箱即用

### 0.1 安装依赖

```shell
$ . install_dep.sh
```

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
