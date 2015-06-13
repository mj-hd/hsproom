[![Build Status](https://drone.io/github.com/mjhd-devlion/hsproom/status.png)](https://drone.io/github.com/mjhd-devlion/hsproom/latest)

# HSP部屋
HSP製のプログラムを、hsp3dish.jsを使って手軽に実行できるWebサイト。
重い開発の最中にあります。

# インストール

以下のコマンドを実行する。
```
git clone https://github.com/mjhd-devlion/hsproom.git
cd hsproom
go get "github.com/go-sql-driver/mysql" "github.com/gorilla/context" "github.com/gorilla/sessions" "github.com/lestrrat/go-ngram" "github.com/microcosm-cc/bluemonday" "github.com/mrjones/oauth" "github.com/russross/blackfriday" "golang.org/x/oauth2" "golang.org/x/oauth2/google"
```

mysqlサーバに、専用ユーザと専用データベースを作る。そして、config/mysql.sqlの内容をSQLサーバに実行する。以下は例。
```
mysql -u root

CREATE DATABASE hsproom;
GRANT ALL ON hsproom.* TO 'hsproom'@'localhost' IDENTIFIED BY 'password';
USE hsproom;
source config/mysql.sql
```

Google+API、TwitterAPI、TwitterBot用のアクセストークンを用意し、config/config.go.exampleを編集し、config/config.goとして保存する。

HSP部屋を起動する。
```
go run hsproom.go --nodaemonize
```
