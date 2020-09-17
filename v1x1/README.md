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

- [ ] 更高级的账户功能
更高级的账户功能目前拥有两个结构体（见 orm/account.go）（2020.9.17）
* AccountEx - 头像图片，个人简介markdown文本等数据
* AccountPermission - 用户的各种权限

### 论坛
- [ ] 基础bbs论坛

### Cookie
注：为了保证cookie的正确携带，服务器和nginx服务器应当保证 Access-Control-Allow-Credentials 头为 true。
某些情况下还需要保证 Access-Control-Allow-Origin 不为 *。
#### 字段
* cid - 客户端的唯一标识符，用于和验证码进行key value映射（security模块）
* atk - 服务器颁发的账户令牌，用于在服务器端进行账户的信息访问
    * atk=cyhaoshuaicyfhaoshuaicyfhaoshuaicyfhaoshuai
   
        这是一个特殊令牌，用于访问测试账户。


