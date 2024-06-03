# safeline-utils

通过模拟登录，调用web服务的API，执行定制化操作，将本地SSL证书更新到雷池WAF中。

## 概述

在内网部署雷池WAF后，发现官方只提供使用[Let's Encrypt 的 HTTP-01 验证方法](https://letsencrypt.org/docs/challenge-types/#http-01-challenge)进行免费证书的更新，大部分内网或者家庭环境是无法提供80端口的，在浏览社区讨论以及GitHub相关issue后，发现一篇博客[从指定路径更新雷池WAF证书](https://www.iots.vip/post/update-safeline-ssl-cert-from-file.html)。仔细浏览源码后，发现跟个人需求有一定出入（个人内网账号未添加OTP，lucky的docker环境不适合用python脚本），故而基于该博客提供的登录与更新证书的逻辑，对代码进行了封装，构建二进制文件，更好的兼容各类第三方工具回调。

## 安装

以下安装过程以 Ubuntu 为例，其他系统请使用对应版本。

### 二进制文件安装

1. 从[release 页面](https://github.com/doratiger/safeline-utils/releases)获取最新版本压缩文件
2. 解压缩，并赋予执行权限
3. 部署至`/usr/local/bin`目录

```bash
## example for ubuntu
# download
wget https://github.com/DoraTiger/safeline-utils/releases/download/v0.1.0/safeline-utils-linux-amd64.tar.gz
# unzip and grant
tar -zxf ./safeline-utils-linux-amd64.tar.gz
chmod +x ./safeline-util
# move
sudo cp ./safeline-util  /usr/local/bin/
```

### 源码安装

1. 准备 go 语言环境，可参考该[博客](https://www.superheaoz.top/2022/10/1036/)的 2.3 节。
2. 编译项目(考虑到服务器本身无网络的情况，提供了vendor目录，支持离线编译，如不需要，请删除makefile中的`-mod=vendor`)，可参考`Makefile`文件自行选择编译目标平台。
3. 部署至`/usr/local/bin`目录

```bash
## example for ubuntu
# download
git clone https://github.com/DoraTiger/safeline-utils.git
cd NEU_IPGW
# build 
make all
# grant
chmod +x ./build/linux-amd64/safeline-utils
# move
sudo cp ./build/linux-amd64/safeline-utils /usr/local/bin/
```


## 使用

前提条件已经部署好了雷池WAF，若开启了两部验证要额外保存TOTP的secret。

### lucky

1. 部署lucky，参考以下compose文件
```YAML
services:
  lucky:
    image: gdy666/lucky
    container_name: lucky
    volumes:
      - /pathto/lucky:/goodluck # 自定义持久化路径，用于保存证书与命令工具
    network_mode: host
    restart: always
```
2. 获取safeline-utils工具，并部署
```BASH
cd /pathto/lucky
# 单独构建证书目录，方便区分
mkdir acme
cd acme
wget https://github.com/DoraTiger/NEU_IPGW/releases/download/v0.1.0/safeline-utils-linux-amd64.tar.gz # 假定系统为amd64架构，其他架构自行更改下载地址
# unzip and grant
tar -zxf ./safeline-utils-linux-amd64.tar.gz
```
3. lucky安全管理界面添加证书并配置`映射路径`为`/goodluck/acme`,`证书改变后触发脚本`为：
```SH
# 自行修改网站地址、用户名、密码、证书路径，若有TOTP，则需额外添加 -o yoursecret,证书ID一般从1开始。
/goodluck/acme/safeline-utils cert update -u https://your.waf.site:9443 -n username -p password -i certID -c /goodluck/acme/certname.crt -k /goodluck/acme/certname.key
```

### local

参考lucky部署第二三步。

## 参考
- [从指定路径更新雷池WAF证书](https://www.iots.vip/post/update-safeline-ssl-cert-from-file.html)

## 存在问题
