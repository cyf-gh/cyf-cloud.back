# Security
可以扫描请求body或param进行安全拦截

## SQL注入

搜索文本或是任何的字符串字句都会可能有sql注入的风险。

```
.Where("account_id = ?",id)
```
这种语句没有风险，id为非string类型

