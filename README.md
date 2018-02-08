# goblog

> goblog后端

## goblog 前端

[goblog-vue](https://github.com/1377195627/goblog-vue.git)

## Build

> bash

``` bash
# download
go get github.com/1377195627/goblog

cd $(go env GOPATH)/src/github.com/1377195627/goblog

# build
go build -o goblog app/main.go

```

> fish

``` fish
# download
go get -t github.com/1377195627/goblog

cd (go env GOPATH)/src/github.com/1377195627/goblog

# build
go build -o goblog app/main.go

```

## Run

goblog [[-p] port] [[-S] ip]

-p default 8080
-S default 127.0.0.1

``` bash
./goblog
```

``` bash
./goblog -p 8080 -S 127.0.0.1
```

# Build goblog-vue

``` bash
# download
git clone https://github.com/1377195627/goblog-vue.git

#install dependencies
npm install

# build
npm build
```

## Caddy

```
your.domain {
    root [your goblog path]/dist
    # Optional
    gzip
    tls {
	    dns cloudflare
    }

    # 8080 replace with your goblog listen port
    proxy /api localhost:8080
}
```

## License

goblog is licensed under [GPLv3](LICENSE).