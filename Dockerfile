# 🥜 FamilyDashboard Dockerfile（Raspberry Pi 5 向け）
# マルチステージビルドで効率的にイメージを作成するます！

# ========================================
# ステージ1: フロントエンドビルド（Node.js）
# ========================================
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# package.json と package-lock.json をコピーして依存関係をインストール
COPY frontend/package*.json ./
RUN npm ci --only=production

# フロントエンドのソースをコピーしてビルド
COPY frontend/ ./
RUN npm run build

# ========================================
# ステージ2: バックエンドビルド（Go）
# ========================================
FROM golang:1.23-alpine AS backend-builder

WORKDIR /app

# Go モジュールファイルをコピーして依存関係をダウンロード
COPY go.mod go.sum ./
RUN go mod download

# バックエンドのソースをコピーしてビルド
COPY cmd/ ./cmd/
COPY internal/ ./internal/

# バイナリをビルド（静的リンクでコンパイル）
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -installsuffix cgo -o server ./cmd/server

# ========================================
# ステージ3: 本番用最小イメージ（Alpine Linux）
# ========================================
FROM alpine:latest

# タイムゾーンとca-certificatesをインストール（HTTPS通信に必要）
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Asia/Tokyo タイムゾーンを設定
ENV TZ=Asia/Tokyo

# ビルド済みバイナリをコピー
COPY --from=backend-builder /app/server ./server

# フロントエンドのビルド成果物をコピー
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# data ディレクトリを作成（ボリュームマウントで永続化）
RUN mkdir -p /app/data/cache

# ポート8080を公開
EXPOSE 8080

# 非rootユーザーで実行（セキュリティ向上）
RUN addgroup -g 1000 familydashboard && \
    adduser -D -u 1000 -G familydashboard familydashboard && \
    chown -R familydashboard:familydashboard /app

USER familydashboard

# サーバーを起動するます！
CMD ["./server"]
