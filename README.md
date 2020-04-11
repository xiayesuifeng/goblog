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
- [x] 信息修改
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
- [X] 备份功能
- [X] 还原功能
- [X] 插件管理中心
- [X] 插件机制

## 插件列表
> 插件可通过登录后在后台的插件中心中安装，也可以 [手动安装](https://gitlab.com/xiayesuifeng/goblog-plugins/blob/master/README.md#%E6%8F%92%E4%BB%B6%E5%88%97%E8%A1%A8)
> 注：插件只支持 linux-amd64 ，其他平台请自行编译

- [X] 友链
- [X] 评论
- [ ] 打赏 
- [X] RSS订阅

## 数据库支持

- `MySQL`
- `PortgreSQL`
- `Sqlite3`
- `SQL Server`

其中 `Sqlite3` 与 `SQL Server` 不内置驱动，请自行导入驱动然后重新编译，导入驱动方法如下：

修改 [sql-driver/driver.go](https://gitlab.com/xiayesuifeng/goblog/-/blob/master/sql-driver/driver.go)

Sqlite3 在最后一行加入
```import _ "github.com/jinzhu/gorm/dialects/sqlite"```

而 `SQL Server` 则加入 
```import _ "github.com/jinzhu/gorm/dialects/mssql"```

最后根据 [编译安装](#编译安装) 重新编译既可

## 快速部署

> 下载
```
wget https://gitlab.com/xiayesuifeng/goblog/builds/artifacts/2.4.0/download?job=build-goblog -O goblog.zip
unzip goblog.zip
cd goblog
```
> 配置

方法一（推荐）

根据 [配置](#配置) 进行手动配置

方法二

```
# 生成配置文件(数据库自行创建)
./goblog -i
```

> 后端启动

```
./goblog -p 20181
```

> caddy配置实例

```
your {
    root /your/path/goblog/web
    gzip
    
    rewrite {
        if {path} not_match ^/api
        to {path} {path} /
    }
    proxy /api localhost:20181 {
        transparent
    }
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
go build -ldflags "-s -w" -trimpath gitlab.com/xiayesuifeng/goblog
```

### 配置

```bash
cp config.default.json config.json
```

> config.json详解

```
{
  "mode":"debug",           //调试模式，正式部署改为release
  "name": "goblog",         //个人博客名
  "password": "0925e15d0ae6af196e6295923d76af02b4a3420f",   //登录密码,当前为admin
  "useCategory": true,      //使用分类功能,不使用分类改为false
  "dataDir":"data",         //数据存放路径,当前为goblog运行路径下的data下，用于存储article和图片等
  "database":{              //数据库信息
    "driver":"mysql",       //数据库驱动(支持mysql, portgres, sqlite3, mssql(SQL Server))
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
  },
  "tls": {                    // autotls 配置，可不需要任何 web 服务器直接运行并支持 HTTPS (证书自动申请)
    "enable": false,          // 启用 autotls
    "domain": ["example.com"] // 绑定的域名
  }
}
```

## 裸奔功能 (无需 web 服务器直接监听80与443，自动申请证书)
1. 修改配置文件中的 tls 下的 enable 为 true，如
```json
"tls": {
    "enable": true,
    "domain": ["example.com"]
}
```
2. 修改 domain 中的 `example.con` 为你自己的域名（支持绑定多个），如
```json
"tls": {
    "enable": true,
    "domain": ["example1.com","example2.com"]
}
```
3. 设置 `GOBLOG_WEB_PATH` 变量为前端所在路径，然后启动 `goblog`，如
```
env GOBLOG_WEB_PATH=./web ./goblog
```

> 注：如果需要修改 http 监听端口，可添加 `-autotls-use-custom-http-port` 启动，然后配合 `-p` 参数指定想要的端口既可

## 备份功能 (不支持 `SQL Server` 数据备份)
为保证数据完整性，请确保 `goblog` 未在运行，然后使用 `-b` 参数进行启动，如
```bash
./goblog -b
```
> 备份文件所在位置将在备份成功后提示

## 恢复功能 (不支持 `SQL Server` 数据恢复)

为保证数据完整性，请确保 `goblog` 未在运行，然后
1. 确保待恢复的数据库已存在，如不存在请自行去使用如 `CREATE DATABASE goblog` 这类的去创建
2. 然后使用 `-r` + 备份文件路径 参数进行启动，如
```bash
./goblog -r data/backup/Backup-GoBlog-20200401214718.zip
```
3. 在正式开始恢复前将询问是否使用原备份中的配置文件的 `dataDir` 和 `database` 配置，如果不使用请确保 `config.json` 已配置
> 恢复成功后既可重新启动 `goblog`

## 加密密码生成

```
echo -n yourpassword | openssl dgst -md5 -binary | openssl dgst -sha1
```

## 配合 `systemd` 使用
```
[Unit]
Description=GoBlog Service

[Service]
ExecStart=/path/to/goblog -pid-file /tmp/goblog.pid
ExecReload=/bin/kill -HUP $MAINPID
PIDFile=/tmp/goblog.pid
```

## License

goblog is licensed under [GPLv3](LICENSE).