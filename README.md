# UEC SSH プロキシ

SSH の ProxyCommand として利用すると、sol/CED/IED などの学内サーバーに簡単に接続できます。
使用には、あらかじめ[公開鍵認証の設定](https://www.cc.uec.ac.jp/ug/ja/remote/ssh/index.html#remote-ssh-keypair)が必要です。

## インストール

実行ファイルをダウンロードして、適当な場所に配置してください。

SSH の設定ファイル(Windows: `%USERPROFILE%\.ssh\config`(`C:\Users\ユーザー\.ssh\config`), macOS/Linux: `~/.ssh/config`)に以下のように設定を追加してください。

```ssh
Host uec-*
  User あなたのUECアカウント
  ProxyCommand 実行ファイルの配置先パス -host %h -user %r -port %p
```

## 使用方法

SSH の接続は通常通り行います。例えば、以下のように接続できます。

```bash
ssh uec-sol
ssh uec-ced
ssh uec-ied
```

VPN が必要な CED や IED の環境であっても、VPN なしで接続することができます。

### ホスト名の書き換え

ssh コマンドに指定したホスト名に以下の文字列が含まれる場合は、接続先を以下の通りに変更します。

- `sol` → `sol.edu.cc.uec.ac.jp`
- `ced` → `orange[01-30].ced.edu.cc.uec.ac.jp`
- `ied` → `[ab][11-78].ied.inf.uec.ac.jp`
- `ssh` → `ssh.cc.uec.ac.jp`

CED に接続する場合は、CED Monitor の表示に応じて利用可能でユーザーが最も少ない端末に接続します。
IED に接続する場合は、すでに起動している利用可能な端末に優先して接続し、
起動中の端末がすべて利用中である場合には停止中の端末を起動してから接続します。

## SSH の設定例

### `.ssh`フォルダ内に`uec-ssh-proxy`という名前で実行ファイルを配置した場合

```ssh
# Windows
Host uec-*
  User z2000000
  ProxyCommand C:/Users/username/.ssh/uec-ssh-proxy -host %h -user %r -port %p
```

```ssh
# macOS/Linux
Host uec-*
  User z2000000
  ProxyCommand ~/.ssh/uec-ssh-proxy -host %h -user %r -port %p
```

### ホストを明示的に設定する場合

```ssh
Host uec-sol
Host uec-ced
Host uec-ied
Host uec-*
  User あなたのUECアカウント
  ProxyCommand /path/to/uec-ssh-proxy -host %h -user %r -port %p
```

`uec-sol`, `uec-ced`, `uec-ied`を明示的に設定することで、VSCode の Remote SSH 拡張機能などで接続先を選択しやすくなります。

### カスタム認証鍵を使用する場合

`id_rsa`や`id_ed25519`などのデフォルトのファイル名の認証鍵ではなく、カスタムの認証鍵を使用する場合は、SSH の設定ファイルに以下のように追加してください。

```ssh
Host uec-*
  User あなたのUECアカウント
  IdentityFile ~/.ssh/id_ed25519_custom
  ProxyCommand /path/to/uec-ssh-proxy -host %h -user %r -port %p

Host ssh.cc.uec.ac.jp
  User あなたのUECアカウント
  IdentityFile ~/.ssh/id_ed25519_custom
```

特に、ホスト `ssh.cc.uec.ac.jp` に対して鍵の指定がない場合にはすべての接続に失敗します。

あるいは、

```
Host *uec*
  User あなたのUECアカウント
  IdentityFile ~/.ssh/id_ed25519_custom
  ProxyCommand /path/to/uec-ssh-proxy -host %h -user %r -port %p
```

のようにすることで、設定を 1 つにまとめることもできます。

### 別のホスト名パターンを使用する場合

```ssh
Host *-uec
  User あなたのUECアカウント
  ProxyCommand /path/to/uec-ssh-proxy -host %h -user %r -port %p
```

`Host`のパターンを変更することで、ホスト名の接続先を変更できます。例えば、`*-uec`とすることで、`sol-uec`, `ced-uec`, `ied-uec`などのホスト名に対応できます。

技術的には、`sol`/`ced`/`ied`/`ssh`が含まれて`.uec.ac.jp`で終わらない任意のホスト名についてホスト名の書き換え処理を行っているため、
パターンの変更により、ホスト名の接続先を変更できます。

このことから、以下のようにホスト名を明示的に指定することもできます。

```ssh
Host sol
  User あなたのUECアカウント
  ProxyCommand /path/to/uec-ssh-proxy -host %h -user %r -port %p

Host ied
  User あなたのUECアカウント
  ProxyCommand /path/to/uec-ssh-proxy -host %h -user %r -port %p

Host ced
  User あなたのUECアカウント
  ProxyCommand /path/to/uec-ssh-proxy -host %h -user %r -port %p

Host ssh
  User あなたのUECアカウント
  ProxyCommand /path/to/uec-ssh-proxy -host %h -user %r -port %p
```

# 開発

```sh
$ go mod download
$ go build
```

`ProxyCommand`でコマンドを指定するときに、`-verbose`オプションをつけるとログを詳細に表示します。
