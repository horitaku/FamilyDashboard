# data ディレクトリ

このディレクトリには、アプリケーション実行時の設定・キャッシュ・トークンが保存されるのです。

## 📋 ファイル説明

### settings.json (秘密情報を含む - gitには含めない)

アプリケーション設定ファイル。ClientSecret などの秘密情報を含むため、**git で版管理されません**。

**セットアップ方法:**

1. `settings.example.json` をコピーして `settings.json` にリネームするます:
   ```bash
   cp data/settings.example.json data/settings.json
   ```

2. `settings.json` を編集して、以下の情報を入力するのです:
   - `google.clientId`: Google OAuth 2.0 Client ID
   - `google.clientSecret`: Google OAuth 2.0 Client Secret
   - `google.redirectUri`: OAuth リダイレクトURI（ローカル開発は `http://localhost:8080/auth/callback`）
   - `google.calendarId`: Google Calendar の共有カレンダーID
   - `google.taskListId`: Google Tasks のタスクリストID（`@default` でも動作）

### settings.example.json (テンプレート - gitに含まれる)

`settings.json` のテンプレート。秘密情報は**プレースホルダーのみ**で、git で版管理されます。

新しい開発者がリポジトリをクローンした時の参考になるのです。

### cache/ (実行時生成キャッシュ)

天気 API、Google Calendar、Google Tasks のキャッシュが保存されるのです。
自動で生成されるため、git には含まれません。

### tokens.json (将来: トークン永続化用)

Google OAuth トークンをファイルに保存する場合に使用（実装予定）。
秘密情報なので、git には含まれません。

---

## 🔐 セキュリティ上の注意

⚠️ **settings.json に秘密情報を含める場合:**

- `.gitignore` に `data/settings.json` が含まれていることを確認
- リポジトリを公開する前に、秘密情報が含まれていないか確認
- ClientSecret を誤ってコミットした場合は、**無効化して新しい認証情報を生成**するのです！

---

## 📝 Google API 設定手順

Google OAuth 2.0 の ClientID/ClientSecret を取得する手順は、
リポジトリの [docs/GOOGLE_SETUP.md](../docs/GOOGLE_SETUP.md) を参照するますね。
