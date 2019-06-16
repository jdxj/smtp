# MailHole

![](https://img.shields.io/badge/go-1.12-brightgreen.svg)&ensp;
![](https://img.shields.io/badge/build-passing-brightgreen.svg)&ensp;
![](https://img.shields.io/badge/Powered%20by-Jdxj-orange.svg)

## 简介

- MailHole 使用 Go 语言编写.
- MailHole 是一个玩具, 用于接收一次性邮件. 如果你不想用自己的私人邮箱接收一些注册类邮件, 你可以使用 MailHole.
- MailHole 尽量不使用第三方库.
- MailHole 的[演示地址](http://test.aaronkir.xyz:8025/mail).

## 快速开始

1. 克隆 & 构建

```
$ git clone https://github.com/jdxj/smtp.git
$ cd smtp
$ go build -o server *.go
```

2. 部署

你可能需要一个公网服务器并且正确配置 DNS 才能接收邮件.

```
$ ./server
```

## 特性

1. 接受轰炸并予以反击.

## TODO

1. 优化前端显示.
2. 不同邮件供应商解码兼容性测试.
3. 过滤某些邮件地址.
4. **接受**轰炸.
5. 优化存储后端.
6. 推送.
7. 防攻击?
    - ![](https://img.shields.io/badge/-%E2%88%9A-brightgreen.svg) 由于每个连接的 handler 数量与连接的数量对应, 且都在新 goroutine 跑, 如果出现大量连接可能会把服务器内存搞垮, 所以需要 goroutine 池.
    - 超时后清除 goroutine.
    - 发现 MailHole 接收了 http 的 get 方法, 参数是个 url, 可能要下载东西?
    - 也接收到了在 http 常用 mime.
    - 对方发送了 ping 命令.
    - 对方发送我的 ip 给我.
    - 对方发送相对路径给我: `get /../../../../../../../../../../../../../../../../../../etc/passwd http/1.1`.
        - 严重的一个命令: `[SMTPServer]2019/06/07 18:32:50 receiver.go:82: Unresolved command: get, data: Cmd: get, Param: /../../../../../../../../../../../../../../../../../../etc/passwd      http/1.1`.
    - 对方连续创建连接并不发送数据.
    - 对方尝试弱口令.
    - **对方仍不放弃攻击**.
8. 解析 smtp 协议的命令及其参数.
    - ![](https://img.shields.io/badge/-%E2%88%9A-brightgreen.svg) reset
    - 在一些奇怪命令后携带了类似脚本写法的参数.
    - auth login.
9. smtp 命令出现乱码, 不清楚对方发送非 utf-8 编码, 还是二进制数据.
10. ![](https://img.shields.io/badge/-%E2%88%9A-brightgreen.svg) 日志到达指定大小后新建日志文件.
11. 统计对方 ehlo 后的内容.
    - 建立黑名单, 并利用反击手段.
12. 有些邮件确实由内容, 需要输出到文件.
13. *假装被黑.*
    - 对所有命令返回 "ok".
14. ![](https://img.shields.io/badge/-%E2%88%9A-brightgreen.svg) 记录对方 ip.
15. 优雅退出.

## FAQ

Q: MailHole 能发邮件吗?

A: 不能.

Q: MailHole 打算支持发送邮件吗?

A 会考虑.

Q: MailHole 的其他用途?

A: 嘿嘿...

## License

- [MIT](https://opensource.org/licenses/MIT)
