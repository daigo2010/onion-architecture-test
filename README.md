# onion

Go 言語でオニオンアーキテクチャを意識して構築した商品・ユーザ管理 API サーバ。
機能ごとの vertical slicing + Module パターンで、新規機能の追加コストを抑える設計になっています。

## 特徴

- 標準ライブラリのみで動作(Go 1.22 の `net/http` PathValue を活用)
- データソースは in-memory(`Repository` インタフェース経由で差し替え可能)
- 各機能 (`product`, `user`) が `internal/<feature>/` 配下にオニオン一式を持ち、自身のルート登録まで完結

## 必要環境

- Go 1.22 以上

## 起動

```bash
go run .
# => listening on :8080
```

`ADDR` 環境変数で待ち受けアドレスを変更できます。

```bash
ADDR=:9000 go run .
```

## ディレクトリ構成

```
.
├── main.go                          機能を組み立てて HTTP サーバ起動
└── internal/
    ├── app/
    │   └── module.go                Module IF / Deps / ModuleFactory
    ├── shared/                      共有インフラ
    │   ├── httpx/io.go              JSON I/O ヘルパ
    │   └── system/system.go         Clock / IDGenerator 実装
    ├── product/                     商品機能(オニオン一式)
    │   ├── product.go               Entity + Repository IF
    │   ├── usecase.go               UseCase
    │   ├── inmemory_repository.go   Repository 実装
    │   ├── handler.go               HTTP ハンドラ
    │   └── module.go                組み立て + ルート登録
    └── user/                        ユーザ機能(同じ形)
        └── ...
```

各 feature パッケージ内では、外側のファイルから内側のファイルにのみ依存する形(`module.go` → `handler.go` → `usecase.go` → `product.go`)で書かれています。

## 設計の要点

### Module パターンによる feature 分離

各機能は `app.Module` インタフェースを実装し、自身のルート登録を `RegisterRoutes` で引き受けます。

```go
type Module interface {
    RegisterRoutes(mux *http.ServeMux)
}

type ModuleFactory func(Deps) (Module, error)
```

`main.go` は `moduleFactories` を回して各モジュールを生成・登録するだけです。

### 新規機能の追加手順

1. `internal/<feature>/` に 5 ファイル(`<feature>.go`, `usecase.go`, `inmemory_repository.go`, `handler.go`, `module.go`)を追加
2. `main.go` の `moduleFactories` に **1 行追加**

`router.go` のような集約点が存在しないため、機能追加で他機能のファイルを触る必要がありません。

### データソースの差し替え

各 feature の `inmemory_repository.go` を別の実装(SQLite / PostgreSQL / Redis 等)に置き換えるだけで切替できます。`Repository` インタフェースは feature 内に閉じているため、`usecase.go` / `handler.go` には影響しません。

## API

### Product

| Method | Path             | 説明     |
|--------|------------------|----------|
| POST   | `/products`      | 登録     |
| GET    | `/products`      | 一覧     |
| GET    | `/products/{id}` | 1 件取得 |
| PUT    | `/products/{id}` | 更新     |
| DELETE | `/products/{id}` | 削除     |

リクエストボディ:

```json
{ "name": "りんご", "price": 150, "stock": 10 }
```

### User

| Method | Path          | 説明     |
|--------|---------------|----------|
| POST   | `/users`      | 登録     |
| GET    | `/users`      | 一覧     |
| GET    | `/users/{id}` | 1 件取得 |
| PUT    | `/users/{id}` | 更新     |
| DELETE | `/users/{id}` | 削除     |

リクエストボディ:

```json
{ "name": "山田太郎", "email": "taro@example.com" }
```

## curl での動作確認

### Product

```bash
# 登録
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"りんご","price":150,"stock":10}'

# 一覧
curl http://localhost:8080/products

# 1 件取得
curl http://localhost:8080/products/<id>

# 更新
curl -X PUT http://localhost:8080/products/<id> \
  -H "Content-Type: application/json" \
  -d '{"name":"青りんご","price":180,"stock":8}'

# 削除
curl -X DELETE http://localhost:8080/products/<id>
```

### User

```bash
# 登録
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"山田太郎","email":"taro@example.com"}'

# (一覧 / 取得 / 更新 / 削除 は Product と同形式)
```

## レスポンスとステータスコード

| コード | 意味                                                  |
|--------|-------------------------------------------------------|
| 200    | 取得・更新成功                                         |
| 201    | 登録成功                                               |
| 204    | 削除成功(ボディなし)                                |
| 400    | バリデーションエラー(空文字、負数、不正な email 等) |
| 404    | 該当 ID のデータなし                                   |
| 500    | その他のエラー                                         |

エラーレスポンスは `{"error":"..."}` 形式で返却されます。

## 注意事項

- 永続化は **in-memory** のため、プロセス再起動でデータは消えます。
- ID は 16 バイトのランダム値を hex で表現した文字列です。
- 時刻は UTC で保持されます。
