# johma_site_reborn_portfolio

Go と Gin を用いて開発した個人 Web サイトです。

Firebase Authentication と Cloud Firestore を利用し、ユーザー登録・認証、認証済みユーザー向け API などを実装しています。

## 主な技術

- Go
- Gin
- Firebase Authentication
- Cloud Firestore
- HTML Template
- Docker
- Firebase Emulator Suite

## 主な実装

- Gin によるルーティングと API
- Firebase ID Token の検証
- Bearer Token を利用した認証 Middleware
- ユーザー本人のみアクセスできる Owner チェック
- Firebase Authentication と Firestore を組み合わせたユーザー作成
- Firestore からのユーザーデータ取得
- HTML Template を利用したページ描画
- Cookie を利用したテーマ・言語設定
- Docker による実行環境
- Firebase Emulator を利用したローカル開発

## 構成

```text
src/
  main.go                  # サーバー起動・ルーティング
  handlers/
    user_handler.go        # ユーザー作成・取得 API
  middleware/
    auth.go                # 認証・Owner チェック
  models/
    user.go                # リクエスト / User モデル

templates/                 # HTML Template
static/                    # CSS
firestore.rules            # Firestore Security Rules
```

## コードを見る場合

特に以下のファイルでバックエンド側の実装を確認できます。

- `src/main.go`  
  Gin、Firebase、Firestore の初期化とルーティング
- `src/middleware/auth.go`  
  Firebase ID Token の検証と認証 Middleware
- `src/handlers/user_handler.go`  
  Firebase Authentication / Firestore を利用したユーザー処理
- `firestore.rules`  
  認証ユーザー本人だけが自分のデータへアクセスできるルール

## ローカル実行

Firebase CLI と Docker を利用します。

Firebase Emulator を起動します。

```bash
firebase emulators:start --project demo-johma-site --only firestore,auth
```

別ターミナルで Docker イメージを起動します。

```bash
make lrun
```

`http://localhost:8080` からアクセスできます。

実際の Firebase プロジェクトを利用する場合は、`templates/register.html` と実行時のプロジェクト設定を環境に合わせて変更してください。
