# goblog

> 个人博客goblog后端

[![pipeline status](https://gitlab.com/xiayesuifeng/goblog/badges/master/pipeline.svg)](https://gitlab.com/xiayesuifeng/goblog/commits/master)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/xiayesuifeng/goblog)](https://goreportcard.com/report/gitlab.com/xiayesuifeng/goblog)
[![GoDoc](https://godoc.org/gitlab.com/xiayesuifeng/goblog?status.svg)](https://godoc.org/gitlab.com/xiayesuifeng/goblog)
[![Sourcegraph](https://sourcegraph.com/gitlab.com/xiayesuifeng/goblog/-/badge.svg)](https://sourcegraph.com/gitlab.com/xiayesuifeng/goblog)

## 博客预览
[夏叶随风](https://blog.firerain.me)

## 前端

[goblog-web](https://gitlab.com/xiayesuifeng/goblog-web.git)

## 重构计划

现已使用 [gin](https://github.com/gin-gonic/gin) , [gorm](https://github.com/jinzhu/gorm) , [gin-sessions](https://github.com/gin-contrib/sessions) 进行后端重构，前端使用react重构，使用axios和后端 API 进行交互。

- [X] 一键部署
- [ ] 信息修改(后端已完成)
- [X] 架构搭建
- [X] 登录
- [X] 获取article
- [X] 新建article
- [X] 编辑article
- [X] 删除article
- [X] 获取分类
- [X] 新建分类
- [X] 编辑分类
- [X] 删除分类
- [ ] 插件机制

## 插件列表
- [ ]评论
- [ ]打赏 
- [ ]RSS订阅

## 快速部署

> 下载
```
wget https://gitlab.com/xiayesuifeng/goblog-web/builds/artifacts/master/download?job=build-web -O web.zip
unzip web.zip
wget https://gitlab.com/xiayesuifeng/goblog/builds/artifacts/2.1.1/download?job=build-goblog -O goblog.zip
unzip goblog.zip
mv build/goblog ./
```
> 配置
```
# 生成配置文件(数据库自行创建)
./goblog -i
```

> 后端启动

```
goblog -p 20181
```

> caddy配置实例

```
your {
    root /your/path/goblog-web
    gzip
    
    rewrite {
        if {path} not_match ^/api
        to {path} {path} /
    }
    proxy /api localhost:20181
}
```

## 编译安装

### 编译
> 前端
```
git clone https://gitlab.com/xiayesuifeng/goblog-web.git goblog
cd goblog
npm install
npm build
```
> 后端
```
go get gitlab.com/xiayesuifeng/goblog
go build -ldflags "-s -w" gitlab.com/xiayesuifeng/goblog
```

### 配置

```bash
cp $(go env GOPATH)/src/gitlab.com/xiayesuifeng/goblog/config.json ./
```

> config.json详解

```
{
  "mode":"debug",           //调试模式，正式部署改为release
  "name": "goblog",         //个人博客名
  "password": "0925e15d0ae6af196e6295923d76af02b4a3420f",   //登录密码,当前为admin
  "useCategory": true,      //使用分类功能,不使用分类改为false
  "dataDir":"data",         //数据存放路径,当前为goblog运行路径下的data下，用于存储article和图片等
  "database":{              //数据库信息(暂只支持mysql,敬请期待别的数据库支持)
    "driver":"mysql",       //数据库驱动(支持mysql与portgres)
    "username":"root",      //数据库用户名
    "password":"",          //数据库密码
    "dbname":"goblog",      //数据库名(需要手动创建)
    "address":"127.0.0.1",  //数据库地址，当前为本地
    "port":"3306"           //数据库端口
  },
  "smtp":{                  //邮箱配置,用于当article有新评论时发送邮件通知(尚未支持,无需配置,敬请期待)
    "username":"",
    "password": "",
    "host": ""
  }
}
```

## 加密密码生成

```
echo -n yourpassword | openssl dgst -md5 -binary | openssl dgst -sha1
```

## License

goblog is licensed under [GPLv3](LICENSE).
