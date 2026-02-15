# 🔑 Google API セットアップガイド

Google Calendar API と Google Tasks API の認証設定手順を説明するのです。

---

## 📝 ステップ 1: Google Cloud プロジェクト作成

1. [Google Cloud Console](https://console.cloud.google.com/) にアクセスするますよー
2. 上部の **「プロジェクトを選択」** → **「新しいプロジェクト」** をクリック
3. プロジェクト名を入力（例: `FamilyDashboard`）
4. **「作成」** ボタンをクリック

---

## 🔌 ステップ 2: Google APIs を有効化

### Google Calendar API

1. Cloud Console の左側メニュー → **「APIとサービス」** → **「ライブラリ」**
2. 検索欄に **「Google Calendar API」** と入力
3. 検索結果をクリック
4. **「有効にする」** ボタンをクリック

### Google Tasks API

1. 同じく左側メニュー → **「ライブラリ」**
2. 検索欄に **「Google Tasks API」** と入力
3. **「有効にする」** ボタンをクリック

---

## ⚙️ ステップ 3: OAuth 2.0 クライアント認証情報を作成

### 3-1. 同意画面を設定（初回のみ）

1. **「APIとサービス」** → **「OAuth同意画面」** をクリック
2. **「ユーザーの種類」** から **「外部」** を選択
3. **「作成」** をクリック
4. 以下の基本情報を入力:
   - **アプリ名**: `FamilyDashboard`
   - **ユーザーサポートメール**: （自分のGoogleメールアドレス）
   - **デベロッパー連絡先情報**: （自分のメールアドレス）
5. **「保存して続行」** をクリック
6. **スコープ**画面 → **「保存して続行」** をクリック（デフォルトのままでOK）
7. **テストユーザー**画面 → **「テストユーザーを追加」** から自分のGoogleアカウントを追加
8. **「保存して完了」** をクリック

### 3-2. クライアント認証情報を作成

1. **「APIとサービス」** → **「認証情報」** をクリック
2. **「+ 認証情報を作成」** ボタンをクリック
3. **「OAuth クライアントID」** を選択
4. **アプリケーションの種類** を選択:
   - **ローカル開発**: `デスクトップアプリケーション`
   - **サーバー運用**: `ウェブアプリケーション`

#### ウェブアプリケーション（サーバー）を選択した場合:

5. **「承認済みの JavaScript 生成元」** に以下を追加:
   ```
   http://localhost:8080
   http://localhost:3000
   ```

6. **「承認済みのリダイレクト URI」** に以下を追加:
   ```
   http://localhost:8080/auth/callback
   http://localhost:8080/auth/oauth2callback
   ```

   本番環境の場合:
   ```
   https://example.com/auth/callback
   https://example.com/auth/oauth2callback
   ```

7. **「作成」** をクリック

---

## 🔑 ステップ 4: ClientID と ClientSecret を取得

作成後、以下の画面が表示されるのです：

```
クライアント ID: 123456789-abcdef...apps.googleusercontent.com
クライアント シークレット: GOCSPX-xxxxxxxxxxxxx
```

**この情報を安全に保管するのです！** 🔐

`data/settings.json` にコピー:

```json
{
  "google": {
    "clientId": "123456789-abcdef...apps.googleusercontent.com",
    "clientSecret": "GOCSPX-xxxxxxxxxxxxx",
    "redirectUri": "http://localhost:8080/auth/callback",
    "calendarId": "",
    "taskListId": "@default"
  }
}
```

---

## 📅 ステップ 5: Calendar ID を取得

1. [Google Calendar](https://calendar.google.com/) を開くます
2. 左側の **「その他のカレンダー」** セクションで、共有カレンダーを右クリック
3. **「設定」** をクリック
4. **「カレンダーの統合」** セクションで **「カレンダーID」** を見つける
   ```
   例: xyz123@group.calendar.google.com
   ```
5. ID をコピーして `data/settings.json` に追加:
   ```json
   "calendarId": "xyz123@group.calendar.google.com"
   ```

---

## ✅ ステップ 6: Tasks List ID を確認

Google Tasks はデフォルトのリストを使用する場合、以下を設定するのです：

```json
"taskListId": "@default"
```

特定のタスクリストを使用する場合は、以下のコマンドでリストを取得:

```bash
curl -X GET "https://www.googleapis.com/tasks/v1/users/@me/lists" \
  -H "Authorization: Bearer $ACCESS_TOKEN"
```

結果から `id` フィールドをコピー。

---

## 🧪 ローカルテスト用（簡易方法）

開発中は環境変数で設定することもできるのです：

```bash
# .env ファイルを作成（git には含まれません）
export GOOGLE_CLIENT_ID="your-client-id"
export GOOGLE_CLIENT_SECRET="your-client-secret"
export GOOGLE_CALENDAR_ID="your-calendar-id"
export GOOGLE_TASKS_LIST_ID="@default"
```

その後、`main.go` で環境変数から読み込むるよう実装可能なのです。

---

## 🔐 セキュリティチェックリスト

- [ ] ClientSecret を他の人に共有していない？
- [ ] settings.json が `.gitignore` に含まれている？
- [ ] settings.json をリポジトリにコミットしなかった？
- [ ] 公開リポジトリにアップロードする前に確認した？

**もし ClientSecret を誤ってコミットした場合:**
1. 無効化: Google Cloud Console で既存の認証情報を削除
2. 新しい認証情報を生成
3. git history から削除（`git filter-branch` など）

---

## 📚 参考リンク

- [Google Calendar API ドキュメント](https://developers.google.com/calendar/api/guides/overview)
- [Google Tasks API ドキュメント](https://developers.google.com/tasks/overview)
- [OAuth 2.0 認可フロー](https://developers.google.com/identity/protocols/oauth2)
