# practice GO

以gin写了一个简单的数据下载页面，通过命令行参数指定目录和用户名密码。安全方面应该凑凑合合，没有加太多特别复杂的东西。

前端用jQuery的datatables和vue简单写了一下，就是个单页面的应用。自己用一下给其他人分享什么数据刚刚好

使用了

- [goptions](https://github.com/voxelbrain/goptions)：构建命令行
- [gin](https://github.com/gin-gonic/gin)：比起beego这种重型框架，还是喜欢这种轻量级的
- [gin-sessions](https://github.com/gin-contrib/sessions)：gin的session插件，用户登录操作
- [archiver](https://github.com/mholt/archiver)：压缩文件，压缩文件夹，方便下载文件夹。不过当然不是边压缩边下载，那么复杂干嘛。
- [go-rice](https://github.com/GeertJohan/go.rice)：打包静态文件，便于编译成统一的二进制

基本上不会在这个基础上搞了，反而有个私人在玩的东西，会上心在做一定程度的修改和编写