# data ディレクトリ

このディレクトリには、アプリケーション実行時の設定・キャッシュ・トークンが保存されるのです。

## 📋 ファイル説明

### settings.json (秘密情報を含む - gitには含めない)

アプリケーション設定ファイル。Nextcloud のパスワードなどの秘密情報を含むため、**git で版管理されません**。

**セットアップ方法:**

1. `settings.example.json` をコピーして `settings.json` にリネームするます:
   ```bash
   cp data/settings.example.json data/settings.json
   ```

2. `settings.json` を編集して、以下の情報を入力するのです:
   - `nextcloud.serverUrl`: Nextcloud サーバーの URL（例: `https://nextcloud.example.com`）
   - `nextcloud.username`: Nextcloud のユーザー名
   - `nextcloud.password`: Nextcloud のアプリパスワード（または メインパスワード）
   - `nextcloud.calendarNames`: カレンダー名の配列（例: `["family", "work"]`）
   - `nextcloud.taskListNames`: タスクリスト名の配列（例: `["tasks", "shopping"]`）
   - `location.cityName`: 天気情報を取得する都市名（例: `"姫路市"`）

詳細な設定方法は [docs/NEXTCLOUD_SETUP.md](../docs/NEXTCLOUD_SETUP.md) を参照してください。

### settings.example.json (テンプレート - gitに含まれる)

`settings.json` のテンプレート。秘密情報は**プレースホルダーのみ**で、git で版管理されます。

新しい開発者がリポジトリをクローンした時の参考になるのです。

### cache/ (実行時生成キャッシュ)

天気 API、Nextcloud カレンダー、Nextcloud タスクのキャッシュが保存されるのです。
自動で生成されるため、git には含まれません。

キャッシュファイルの例:
- `weather:JP:姫路市.json`: 天気データのキャッシュ
- `nextcloud_calendar_events.json`: カレンダーイベントのキャッシュ
- `nextcloud_tasks_items.json`: タスクリストのキャッシュ

---

## 🔐 セキュリティ上の注意

⚠️ **settings.json に秘密情報を含める場合:**

- `.gitignore` に `data/settings.json` が含まれていることを確認
- リポジトリを公開する前に、秘密情報が含まれていないか確認
- アプリパスワードを誤ってコミットした場合は、**無効化して新しいパスワードを生成**するのです！

---

## 📝 Nextcloud 設定手順

Nextcloud CalDAV/WebDAV のアプリパスワードやカレンダー名、タスクリスト名を取得する手順は、
リポジトリの [docs/NEXTCLOUD_SETUP.md](../docs/NEXTCLOUD_SETUP.md) を参照するますね。
