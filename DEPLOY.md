# 🥜 FamilyDashboard デプロイガイド

うい！Raspberry Pi 5 でサーバーをデプロイする手順を説明するます！わくわく！✨

---

## 📦 必要なもの

- Raspberry Pi 5（ARM64）
- Docker と Docker Compose がインストールされていること
- インターネット接続

---

## 🚀 デプロイ手順

### 1. リポジトリをクローン

```bash
git clone https://github.com/rihow/FamilyDashboard.git
cd FamilyDashboard
```

### 2. 設定ファイルを準備

```bash
# サンプル設定ファイルをコピー
cp data/settings.example.json data/settings.json

# 設定を編集（都市名、APIキーなど）
nano data/settings.json
```

### 3. Google OAuth 認証設定（必要な場合）

Google カレンダー・タスクを使用する場合は、[GOOGLE_SETUP.md](docs/GOOGLE_SETUP.md) を参照して OAuth クライアントを作成してください。

### 4. Docker イメージをビルド

```bash
# ARM64 向けにビルド（Raspberry Pi 5 で実行）
docker-compose build
```

### 5. サーバーを起動

```bash
# バックグラウンドで起動
docker-compose up -d

# ログを確認
docker-compose logs -f
```

### 6. ブラウザでアクセス

```
http://localhost:8080
```

または、Raspberry Pi の IP アドレスで他の端末からアクセス:

```
http://<RaspberryPi のIP>:8080
```

---

## 🔄 更新・再デプロイ

コードを更新したら、以下のコマンドで再ビルド・再起動するます：

```bash
# リポジトリを更新
git pull

# フロントエンドを再ビルド（必要な場合）
cd frontend
npm install
npm run build
cd ..

# Docker イメージを再ビルド
docker-compose build

# サーバーを再起動
docker-compose down
docker-compose up -d
```

---

## 🛑 サーバーを停止

```bash
# コンテナを停止
docker-compose stop

# コンテナを停止して削除
docker-compose down
```

---

## 📊 ログ確認

```bash
# リアルタイムでログを表示
docker-compose logs -f

# 最新100行を表示
docker-compose logs --tail=100
```

---

## 🔧 トラブルシューティング

### フロントエンドが表示されない

```bash
# フロントエンドを手動でビルド
cd frontend
npm run build
cd ..

# Docker イメージを再ビルド
docker-compose build
docker-compose up -d
```

### Google OAuth が動作しない

1. `data/settings.json` で `clientId` と `clientSecret` が正しく設定されているか確認
2. `/auth/login` にアクセスして OAuth フローを開始
3. `data/tokens.json` にトークンが保存されているか確認

### ポート 8080 が使用中

`docker-compose.yml` の `ports` セクションを編集して別のポートを使用:

```yaml
ports:
  - "8888:8080"  # ホストの8888番ポートにマッピング
```

---

## 🎨 FHD表示確認（1920x1080）

Raspberry Pi Zero 2 W などのキオスク端末で表示確認:

1. フルスクリーンブラウザ（Chromium など）で開く
2. 2m 距離から見やすいか確認
3. タイポグラフィ・余白が適切か確認

---

## 📝 データの永続化

`data/` ディレクトリはボリュームマウントで永続化されています:

- `data/settings.json`: 設定ファイル
- `data/tokens.json`: OAuth トークン
- `data/cache/`: キャッシュファイル

コンテナを削除してもデータは保持されるます！✨

---

がんばるます！💪
