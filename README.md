# 适配[gin-gonic/gin](https://github.com/gin-gonic/gin)的session管理

参考以下项目，因为改动非常大，所以并非基于某一个clone的

1. <https://github.com/gorilla/sessions>
1. <https://github.com/martini-contrib/sessions>


gorilla/sessions依赖于<https://github.com/gorilla/context>，后者内部依赖一个加锁的map，不是很中意。在1.7之后，内建了[context](https://golang.org/pkg/context/)模块，可以在一定程度上优化**gorilla/context**的问题。

之所以**一定程度**是因为context库并不会改变现有的http.Request，而是返回一个新的对象，这导致一个很严重的问题，除非直接修改传入的http.Request对象，否则就无法链式的调用下去，参见如下代码

``` go
// 注意传入的next
func middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        userContext:=context.WithValue(context.Background(),"user","张三")
        ageContext:=context.WithValue(userContext,"age",18)
        // 这里必须递归调用
        next.ServeHTTP(rw, r.WithContext(ageContext))
    })
}
```

在上面的代码中，通过context包可以在http.Request对象上附加信息，但是由于会生成新的http.Request对象，所以链式调用，后续的handle并不会读取到新添加的数据，在很多场景无法使用或者会导致代码很难看。

此问题参考

1. <http://www.flysnow.org/2017/07/29/go-classic-libs-gorilla-context.html>
1. <https://stackoverflow.com/questions/40199880/how-to-use-golang-1-7-context-with-http-request-for-authentication?rq=1>

martini-contrib/sessions也直接依赖gorilla/sessions，所以也需要优化

---

之所以基于martini-contrib/sessions来修改，是因为

1. 存在redis的store，因个人喜好，替换了一个新redis库[redis.v5](https://gopkg.in/redis.v5)
1. 适配了gin.Context对象

由于**gin.Context自带context**，可以直接附加数据，所以完全可以绕开<https://github.com/gorilla/context>和<https://golang.org/pkg/context/>

另外**增加了Store的delete**方法，用于删除整个cookie，而不是cookie里面的某个key。