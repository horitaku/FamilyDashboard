# FamilyDashboard

家族向けダッシュボード（Go/Gin + Svelte 静的）プロジェクトです。
Raspberry Pi 5 がバックエンド、Pi Zero 2 W が表示専用キオスク端末を想定しています。

## できること
- 天気・カレンダー・タスクの情報を1画面に固定レイアウトで表示
- バックエンドが外部APIをキャッシュし、フロントはAPI経由で表示
- オフライン時は直近キャッシュを表示（エラー状態はヘッダーで通知予定）

## アーキテクチャ
- バックエンド: Go + Gin（REST API / 静的ファイル配信）
- フロントエンド: Svelte + Vite（静的ビルド）
- ストレージ: ./data 配下（設定・トークン・キャッシュ）

## 要件
- Go 1.22+（PATH に go があること）
- Node.js 18+（npm を使用）

## 主要ディレクトリ
- cmd/server: Gin サーバーエントリ
- internal: 各サービス・設定・キャッシュ・DTO
- frontend: Svelte アプリ
- data: 設定・キャッシュ・トークン保存

## フォルダ構成

```text
.
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── cache/
│   ├── config/
│   ├── http/
│   ├── models/
│   └── services/
│       ├── geocode/
│       ├── google/
│       └── weather/
├── frontend/
│   ├── dist/
│   └── src/
├── data/
│   ├── cache/
│   └── settings.json
├── docs/
└── README.md
```

## セットアップ
1) 依存関係の準備

```bash
go mod download
```

```bash
cd frontend
npm install
```

2) 設定ファイルを確認

`data/settings.json` を環境に合わせて更新してください。

## 開発起動

### バックエンド

```bash
go run cmd/server/main.go
```

### フロントエンド（開発）

```bash
cd frontend
npm run dev
```

### フロントエンド（ビルド）

```bash
cd frontend
npm run build
```

ビルド後は `frontend/dist` をバックエンドが配信します。

## Docker（Raspberry Pi 5 デプロイ）

Raspberry Pi 5 上では Docker で起動する想定です。
データ永続化のために `data` をボリュームでマウントします。

### 想定構成
- コンテナ: Go/Gin サーバー
- 静的ファイル: `frontend/dist` をコンテナへ含める
- データ永続化: `./data` をコンテナの `/app/data` にマウント

### 予定コマンド（後で整備）

```bash
docker compose up -d --build
```

Dockerfile と docker-compose.yml はステップ11で作成します。

## API（予定）
- GET /api/status
- GET /api/calendar
- GET /api/tasks
- GET /api/weather

## タイムゾーン
すべての計算と表示は Asia/Tokyo を使用します。

## ライセンス
後で決める予定です。
