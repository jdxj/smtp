# MailHole

![](https://img.shields.io/badge/go-1.12-brightgreen.svg)&ensp;
![](https://img.shields.io/badge/build-passing-brightgreen.svg)&ensp;
![](https://img.shields.io/badge/Powered%20by-Jdxj-orange.svg)

## 简介

- MailHole 使用 Go 语言编写.
- MailHole 是一个玩具, 用于接收一次性邮件. 如果你不想用自己的私人邮箱接收一些注册类邮件, 你可以使用 MailHole.
- MailHole 尽量不使用第三方库.
- MailHole 的[演示地址](http://test.aaronkir.xyz:8025/mail)

## 快速开始

1. 克隆 & 构建

```
$ git clone https://github.com/jdxj/smtp.git
$ cd smtp
$ go build -o server *.go
```

2. 部署

你可能需要一个公网服务器才能接收邮件.

```
$ ./server
```

## FAQ

Q: MailHole 能发邮件吗?

A: 不能.

Q: MailHole 打算支持发送邮件吗?

A 会考虑.

Q: MailHole 的其他用途?

A: 嘿嘿...

## License

- [MIT](https://opensource.org/licenses/MIT)
