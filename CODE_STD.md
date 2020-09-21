# cyf golang代码规范

## v1x1 

### 错误处理
* 关于异常操作

获取的error必须及时处理（在下一行代码处理），如果是最外层函数，例如http的action函数，则调用err.Check进行异常捕捉操作。

及时处理也有助于避免隐藏变量问题（Accidental Variable Shadowing）

例如orm一类的非最外层一律不得异常处理，必须逐级函数返回error。

### 命名规范
#### 常量名
形如_xxxXxxxx

以下划线+小写字母开头。

例如：_valueApple

#### 局部变量
* 一般

变量少于5个时可以直接用单词首字母命名，例如account可命名为a。

* error

局部error一律命名为e，形如

```golang
var e error
```
避免与包err冲突。

