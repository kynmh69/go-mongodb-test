# Go MongoDB User Management API

Echo フレームワークを使用したユーザー管理 REST API です。MongoDB をデータベースとして使用し、完全な CRUD 操作を提供します。

## 機能

- ✅ **ユーザー追加** (POST /api/v1/users)
- ✅ **ユーザー更新** (PUT /api/v1/users/:id)
- ✅ **ユーザー編集** (個別フィールド更新対応)
- ✅ **ユーザー削除** (DELETE /api/v1/users/:id)
- ✅ **パスワードのハッシュ化** (bcrypt)
- ✅ **複数の検索方法** (ID/ユーザーID/メール)

## 技術スタック

- **Go**: 1.24
- **Echo**: v4 (Web フレームワーク)
- **MongoDB**: Latest (データベース)
- **bcrypt**: パスワード暗号化

## セットアップ

### 1. MongoDB の起動

```bash
docker-compose up -d
```

### 2. 依存関係のインストール

```bash
go mod tidy
```

### 3. 環境変数の設定

```bash
cp .env.example .env
# 必要に応じて .env ファイルを編集
```

### 4. アプリケーションの起動

```bash
go run main.go
```

## API エンドポイント

### 基本 URL
```
http://localhost:8080/api/v1
```

### エンドポイント一覧

| メソッド | エンドポイント | 説明 |
|---------|---------------|------|
| POST | `/users` | ユーザー作成 |
| GET | `/users` | 全ユーザー取得 |
| GET | `/users/:id` | ID でユーザー取得 |
| GET | `/users/search?user_id=xxx` | ユーザーID で検索 |
| GET | `/users/search/email?email=xxx` | メールアドレスで検索 |
| PUT | `/users/:id` | ユーザー更新 |
| DELETE | `/users/:id` | ユーザー削除 |
| GET | `/health` | ヘルスチェック |

### リクエスト例

#### ユーザー作成
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "user_id": "user123",
    "email": "user@example.com",
    "password": "password123"
  }'
```

#### ユーザー更新
```bash
curl -X PUT http://localhost:8080/api/v1/users/60f7b1b8e4b0c7a8e4b0c7a8 \
  -H "Content-Type: application/json" \
  -d '{
    "email": "newemail@example.com"
  }'
```

## ユーザーモデル

```go
type User struct {
    ID        primitive.ObjectID `json:"id"`
    UserID    string             `json:"user_id"`
    Email     string             `json:"email"`
    Password  string             `json:"-"` // レスポンスには含まれない
    CreatedAt time.Time          `json:"created_at"`
    UpdatedAt time.Time          `json:"updated_at"`
}
```

## 開発

### テスト実行
```bash
go test ./...
```

### フォーマット
```bash
go fmt ./...
```

### 依存関係の更新
```bash
go mod tidy
```