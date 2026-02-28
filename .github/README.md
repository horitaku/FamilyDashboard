# 🥜 GitHub Actions CI/CD パイプライン

このフォルダには、FamilyDashboard のCI/CDワークフロー設定が格納されるます！✨

## 📋 ワークフロー一覧

### `ci.yml` - 継続的インテグレーション（自動テスト）

**トリガー:**
- `push` イベント（main/develop ブランチへのプッシュ）
- `pull_request` イベント（main/develop ブランチへのPR）

**実行内容:**

#### 1️⃣ バックエンド（Go）テスト・ビルド
- Go ユニットテスト実行（`go test -v -race -coverprofile=coverage.out ./...`）
- テストカバレッジレポート（Codecov へアップロード）
- Go バイナリのビルド確認（`go build`）
- 依存関係のキャッシュ活用（高速化）

#### 2️⃣ フロントエンド（Svelte）ビルド
- npm 依存関係のインストール（`npm ci`）
- Svelte 静的ビルド実行（`npm run build`）
- ビルドアーティファクト確認（`./frontend/build` 存在確認）

#### 3️⃣ Docker イメージ ビルド確認
- Dockerfile のビルド確認（プッシュはしない）
- ビルドキャッシュ活用（Github Actions キャッシュ）

#### 4️⃣ CI サマリー
- すべてのジョブ結果を集約
- 1つでも失敗していれば全体失敗扱い

**例：PR時のCI実行画面**
```
✓ Go Backend Tests & Build (3m 45s)
✓ Svelte Frontend Build (1m 20s)
✓ Docker Build Verification (2m 10s)
✓ CI Summary (5s)
```

---

### `build-release.yml` - リリースビルド

**トリガー:**
- `push` イベント（`v*` タグ作成時）
  - 例: `git tag v1.0.0 && git push origin v1.0.0`

**実行内容:**

1. CI パイプラインすべて実行（テスト・ビルド確認）
2. Go バイナリをアーティファクトとしてアップロード
3. GitHub Releases を自動作成

**使用方法:**
```bash
# リリース用タグを作成＆プッシュ
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

⚠️ **Docker レジストリプッシュについて:**
本番運用時に Docker イメージをレジストリ（Docker Hub, GHCR など）にプッシュする場合：
- `build-release.yml` の `push: false` を `push: true` に変更
- GitHub Secrets へレジストリ認証情報を登録（例: `DOCKER_USERNAME`, `DOCKER_PASSWORD`）

---

## 🔧 セットアップ手順

### 前提条件
- GitHub リポジトリに `.github/workflows/` フォルダが存在
- Go 1.21以上がインストール済み（ローカル開発）
- Node.js 18以上がインストール済み（ローカル開発）

### ワークフロー設定確認

ワークフロー実行状況は GitHub リポジトリの **Actions** タブで確認できるます！

```
Repository → Actions → ci.yml をクリック
  ↓
各 push/PR マージ時のテスト実行ログを表示
```

### Codecov 統計情報の有効化（オプション）

テストカバレッジをCodecovで可視化したい場合：

1. https://codecov.io にアクセス
2. GitHub 連携（ログイン）
3. リポジトリを有効化
4. CI実行後、自動的にカバレッジデータが送信される

現在は `continue-on-error: true` で Codecov 失敗は無視します。

---

## 📊 CI パイプライン実行時間の目安

| ステップ | 時間 |
|---------|------|
| Go テスト              | 1-2分 |
| npm インストール          | 30-60秒 |
| Svelte ビルド          | 30-60秒 |
| Docker ビルド           | 1-2分 |
| **全体**                | **5-7分** |

※キャッシュがある場合はより高速化

---

## 💡 ベストプラクティス

### 1. ローカルテストは必須
**CI 実行前にローカルで検証すみ！**
```bash
# バックエンド
go test ./...
go build -o ./bin/familydashboard ./cmd/server

# フロントエンド
cd frontend
npm run build
```

### 2. コミットメッセージを明確に
```bash
git commit -m "fix: weather API timeout issue"
git commit -m "feat: add calendar color support"
```

### 3. PR時は早期フィードバック
- PR作成直後にCIが実行
- 赤❌が出ていないか確認してからマージ

### 4. リリースタグは慎重に
- テストが全て PASS してから `v*` タグを作成
- タグ作成後は自動的に Release ワークフローが実行される

---

## 🚨 トラブルシューティング

### ❌ Go テストが失敗する場合
```bash
# ローカルで再現
go test -v ./...

# 依存関係のリセット
go mod tidy
go mod download
```

### ❌ npm ビルドが失敗する場合
```bash
cd frontend
npm ci  # package-lock.json から正確なバージョンをインストール
npm run build
```

### ❌ Docker ビルドが失敗する場合
```bash
# ローカルで Docker ビルドをテスト
docker build -t familydashboard:test .
```

### 📌 キャッシュをクリアしたい場合
GitHub リポジトリ → Settings → Actions → Caches → キャッシュ削除

---

## 📖 参考リンク

- [GitHub Actions ドキュメント](https://docs.github.com/en/actions)
- [Codecov ドキュメント](https://docs.codecov.io/)
- [docker/build-push-action](https://github.com/docker/build-push-action)

---

作成日: 2026-03-01
更新日: 2026-03-01
