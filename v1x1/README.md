# v1x1
v1后迭代的第一个具有标准格式的api

## 模块

### 版本特性
-[x] 错误处理（见err与err_code两个包）

v1x1 将采用全新的返回方法，需要前端配合获取数据。

例：
```
{
    ErrCod : "-90",
    Desc : "This is an example error description",
    Data : "[{},{},{}]"
}
```
#### cid
客户端的唯一标识符，用于和验证码进行key value映射

### 业务逻辑
-[ ] 账户的登陆与注册
-[ ] 基础bbs论坛
-[ ] 数据库的管理

