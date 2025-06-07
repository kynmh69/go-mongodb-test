# Go MongoDB User Management System

[![Go Unit Tests](https://github.com/kynmh69/go-mongodb-test/actions/workflows/go-tests.yml/badge.svg)](https://github.com/kynmh69/go-mongodb-test/actions/workflows/go-tests.yml)

Echo フレームワークを使用したユーザー管理 REST API と Next.js フロントエンドの統合システムです。MongoDB をデータベースとして使用し、完全な CRUD 操作を提供します。

## 機能

- ✅ **ユーザー追加** (POST /api/v1/users)
- ✅ **ユーザー更新** (PUT /api/v1/users/:id)
- ✅ **ユーザー編集** (個別フィールド更新対応)
- ✅ **ユーザー削除** (DELETE /api/v1/users/:id)
- ✅ **パスワードのハッシュ化** (bcrypt)
- ✅ **複数の検索方法** (ID/ユーザーID/メール)

## 技術スタック

### バックエンド
- **Go**: 1.24
- **Echo**: v4 (Web フレームワーク)
- **MongoDB**: Latest (データベース)
- **bcrypt**: パスワード暗号化

### フロントエンド
- **Next.js**: 15.1.3 (React フレームワーク)
- **TypeScript**: 型安全性
- **Tailwind CSS**: スタイリング
- **shadcn/ui**: UI コンポーネントライブラリ
- **Radix UI**: アクセシブルなプリミティブ

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

### 4. バックエンドの起動

```bash
go run main.go
```

### 5. フロントエンドの起動

```bash
cd frontend
npm install
npm run dev
```

アプリケーションは以下のURLでアクセスできます：
- **フロントエンド**: http://localhost:3000
- **バックエンド API**: http://localhost:8080

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

## プロジェクト構成

```
go-mongodb-test/
├── database/           # データベース接続
├── handlers/           # HTTP ハンドラー
├── models/            # データモデル
├── services/          # ビジネスロジック
├── frontend/          # Next.js フロントエンド
│   ├── src/
│   │   ├── app/       # Next.js app directory
│   │   ├── components/ # React コンポーネント
│   │   └── lib/       # ユーティリティと API クライアント
│   └── README.md      # フロントエンド専用README
├── main.go            # メインアプリケーション
├── docker-compose.yml # MongoDB コンテナ設定
└── README.md          # このファイル
```

## 開発

### バックエンド
```bash
# テスト実行
go test ./...

# フォーマット
go fmt ./...

# 依存関係の更新
go mod tidy
```

### フロントエンド
```bash
cd frontend

# 開発サーバー起動
npm run dev

# ビルド
npm run build

# 本番サーバー起動
npm start

# リント
npm run lint
```