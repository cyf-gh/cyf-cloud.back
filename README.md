# Vt.Server

Vt的服务器端，使用Golang编写。

## 0 开箱即用

### 0.1 必备项

安装[Golang](https://golang.google.cn/)并[配置GoPath](https://studygolang.com/articles/17598)

安装[Git](https://git-scm.com/)

### 0.2 clone本项目

```bash
$ git clone https://github.com/cyf-gh/Vt.Server.git
$ cd Vt.Server
```

### 0.3 安装依赖

```shell
$ . install_dep.sh
```



直接运行update_and_deploy.sh脚本

```shell
$ . update_and_deploy.sh
```

### 0.2 配置

#### vtserver.cfg

```
[server]
    udp=":2333"			# UDP监听端口号
    tcp=":2334"			# TCP监听端口号
    log="127.0.0.1:2335" # 日志服务
[fresh_interval]
    log="2000"			# 写日志频率
    udp="1"				# UDP暂歇频率，单位毫秒
```

