[![Build Status](https://drone.io/github.com/mjhd-devlion/hsproom/status.png)](https://drone.io/github.com/mjhd-devlion/hsproom/latest)

# HSP部屋
HSP製のプログラムを、hsp3dish.jsを使って手軽に実行できるWebサイト。
重い開発の最中にあります。

# インストール

以下のコマンドを実行する。
```
git clone https://github.com/mjhd-devlion/hsproom.git
cd hsproom

go get "github.com/mattn/gom"
gom install
```

Google+API、TwitterAPI、TwitterBot用のアクセストークンを用意し、config/config.go.exampleを編集し、config/config.goとして保存する。

HSP部屋を起動する。
```
go run hsproom.go --nodaemonize
```
