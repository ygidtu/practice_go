@2019.08.08 今天研究别人的作图库，看了大量前端框架的东西（其实就是vue），突然搞明白前端框架到底是个什么玩意了。后续有空可以基于vue将这个单页面app重写了。另外还可以考虑学一下如何打包electron app呢

---

以gin写了一个简单的数据下载页面，通过命令行参数指定目录和用户名密码。安全方面应该凑凑合合，没有加太多特别复杂的东西。

前端用jQuery的datatables和vue简单写了一下，就是个单页面的应用。自己用一下给其他人分享什么数据刚刚好

使用了

- [goptions](https://github.com/voxelbrain/goptions)：构建命令行
- [gin](https://github.com/gin-gonic/gin)：比起beego这种重型框架，还是喜欢这种轻量级的
- [gin-sessions](https://github.com/gin-contrib/sessions)：gin的session插件，用户登录操作
- [archiver](https://github.com/mholt/archiver)：压缩文件，压缩文件夹，方便下载文件夹。不过当然不是边压缩边下载，那么复杂干嘛。
- [go-rice](https://github.com/GeertJohan/go.rice)：打包静态文件，便于编译成统一的二进制


---

```bash
$ ./server/server_linux_amd64 -h
Usage: server_linux_amd64 [global options]

Global options:
            --host           host (default: 127.0.0.1)
            --port           port (default: 5000)
            --dir            File directory (default: server)
            --user           Username (default: admin)
            --passwd         Password (default: admin)
            --disable-delete Disable delete button
        -h, --help           Show this help

```
