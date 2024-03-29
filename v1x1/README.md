# v1x1
v1后迭代的第一个具有标准格式的api

## 模块

### 版本特性
- [x] 错误处理（见err与err_code两个包）

v1x1 将采用全新的返回方法，需要前端配合获取数据。

例：
```
{
    ErrCod : "-90",
    Desc : "This is an example error description",
    Data : "[{},{},{}]"
}
```
一个不携带数据的返回her形如：
```
{
    ErrCod : "0",
    Desc : "ok",
    Data : ""
}
```

### 账户
- [x] 账户的登陆与注册
- [ ] 账户找回密码

- [x] 更高级的账户功能
更高级的账户功能目前拥有两个结构体（见 orm/account.go）（2020.9.17）
* AccountEx - 头像图片，个人简介markdown文本等数据

### 论坛
- [ ] 基础bbs论坛

### Cookie
注：为了保证cookie的正确携带，服务器和nginx服务器应当保证 Access-Control-Allow-Credentials 头为 true。
某些情况下还需要保证 Access-Control-Allow-Origin 不为 *。
#### 字段
* cid - 客户端的唯一标识符，用于和验证码进行key value映射（security模块）
* atk - 服务器颁发的账户令牌，用于在服务器端进行账户的信息访问
    * atk=cyfhaoshuaicyfhaoshuaicyfhaoshuaicyfhaoshuai
   
        这是一个特殊令牌，用于访问测试账户。

- [x] atk 持久化？
* 2020.10.4 似乎redis自带持久化

### 安全
- [x] 验证码

- [ ] 桶功能

用于限流防止过度访问

### 中间件
- [x] 日志输出调用接口名

- [x] 统计每个接口调用所需时间

以上两个功能见middleware/helper LogUsedTime中间件

- [x] 携带Cookie中间件

- [ ] 统计调用源中间件 

- [ ] 统计调用次数（redis）

### 高频率访问API

见 http/post_freq.go

#### View
为文章阅读次数

其cache设计为：
```golang
$post_view${{post_id}} : {{int64}}
```

表示Key为某个文章的id，而值为观看过文章的用户id列表
