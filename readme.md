# 概要

Misskeyの絵文字を扱うためのツール群です。

# 使い方

```shell
go build
```

## JSON作成ツール

一括で絵文字をインポートできるようにJSONファイルを作成します。

```nano
./cfg/config.yaml`  
```
を編集
```shell
./emojiTool -makejson 画像のあるディレクトリ`
```

完了すると画像のあるディレクトリにjsonファイルが出来るので、それと画像をzipに固めてインポートしてください。

## 命名規則検知ツール

Misskeyで絵文字として登録する際に、不整合が出る可能性のある命名規則を検知するツールです。

```shell
./emojiTool -checkName 画像のあるディレクトリ`  
```
で動作します。