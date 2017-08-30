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

起動する。
```
./simbase/bin/start
./leaner/learner.py &
gom run hsproom.go
```
