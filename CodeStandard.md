# cyf golang代码规范

## v1x1 

### 错误处理
* 关于异常操作

获取的error必须及时处理（在下一行代码处理），如果是最外层函数，例如http的action函数，则调用err.Check进行异常捕捉操作。

业务逻辑中为了代码的美观性，错误处理必须遵守这样的规则
```go
a, b, e := foo( p1, p2 ); err.Check( e ); glg.Log( a ); glg.Log( b )
``` 
错误处理与变量的打印必须写在*同一行上*，否则会影响对业务逻辑代码的阅读。

及时处理也有助于避免隐藏变量问题（Accidental Variable Shadowing）

例如orm一类的非最外层一律不得异常处理，必须逐级函数返回error。*除去http业务处理的action函数之外，其内部所有调用的函数必须返回error。*

不允许任何形式的内部函数直接panic error。这样做是为了让业务逻辑的错误处理更加精确（虽然我也懒得这么做，大部分情况都是直接error由Check函数直接抛出交给recover处理）


### 命名规范
#### 常量名
遵守C/C++的宏名定义规则。

#### 局部变量
* 一般

变量少于5个时可以直接用单词首字母命名，例如account可命名为a。

* error

局部error一律命名为e，形如

```golang
var e error
```
避免与包err冲突。

