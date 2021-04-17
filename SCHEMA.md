## BUGS
/v1x1/post in post.go


当输入不存在的文章id时

```p, e = orm.GetPostById( id ); err.Assert( e )```

并未返回正确的错误代码，直到

``` a, e := orm.GetAccount( p.OwnerId ); err.Assert( e ) ```

才进行了返回。

##