[![Build Status](https://drone.io/github.com/mjhd-devlion/hsproom/status.png)](https://drone.io/github.com/mjhd-devlion/hsproom/latest)

# HSP部屋
HSP製のプログラムを、hsp3dish.jsを使って手軽に実行できるWebサイト。
重い開発の最中にあります。

# インストール

以下のコマンドを実行する。
```
git clone https://github.com/mjhd-devlion/hsproom.git
cd hsproom
make
```

Google+API、TwitterAPI、TwitterBot用のアクセストークンを用意し、config/config.go.exampleを編集し、config/config.goとして保存する。

HSP部屋を起動する。
```
gom run hsproom.go
```

# [cbe0e851445a7a772fa6cfc9ee954253fa379c35](https://github.com/mjhd-devlion/hsproom/commit/cbe0e851445a7a772fa6cfc9ee954253fa379c35)以前のデータベースから移行するには

[移行プログラム](https://gist.github.com/mjhd-devlion/e5d9fc116c0b19e4688bhttps://gist.github.com/mjhd-devlion/e5d9fc116c0b19e4688bhttps://gist.github.com/mjhd-devlion/e5d9fc116c0b19e4688b])をhsproomディレクトリ直下に置き、`go run db_migrate.db`を実行する。
