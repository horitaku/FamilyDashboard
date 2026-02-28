# 🔑 Nextcloud CalDAV/WebDAV セットアップガイド

Nextcloud のカレンダーとタスクを FamilyDashboard で使うための設定手順を説明するのです。

---

## 📝 前提条件

- Nextcloud サーバーにアクセスできること
- カレンダーアプリとタスクアプリが有効になっていること
- Nextcloud のユーザーアカウントを持っていること

---

## 🔐 ステップ 1: アプリパスワードの生成（推奨）

セキュリティのため、メインパスワードではなく**アプリパスワード**を使用することを推奨するます！

### 1-1. Nextcloud でアプリパスワードを作成

1. Nextcloud にログインするます
2. 右上のプロフィール画像 → **「設定」** をクリック
3. 左側メニュー → **「セキュリティ」** をクリック
4. 下にスクロールして **「アプリパスワード」** セクションを見つける
5. **「アプリ名」** に `FamilyDashboard` と入力
6. **「新しいアプリパスワード」** ボタンをクリック
7. 生成されたパスワードをコピーして安全に保管するのです 🔐

**注意**: 一度しか表示されないので、必ずコピーしてから閉じるます！

### 1-2. メインパスワードを使う場合（非推奨）

アプリパスワードが使えない場合は、Nextcloud のメインパスワードを使用できますが、セキュリティリスクが高くなるます。

---

## 📅 ステップ 2: カレンダー名を確認

### 2-1. Nextcloud Web UI で確認

1. Nextcloud にログインするます
2. 左上のメニュー → **「カレンダー」** をクリック
3. 左側のカレンダーリストを確認
4. 使用したいカレンダーの**名前**をメモするます
   - 例: `family`（家族カレンダー）
   - 例: `work`（仕事カレンダー）

### 2-2. CalDAV URL から確認（上級者向け）

カレンダーの詳細設定で **CalDAV URL** を確認できます：
```
https://your-nextcloud-server.com/remote.php/dav/calendars/username/calendar-name/
```

最後の `calendar-name` がカレンダー名なのです。

### 2-3. 複数カレンダーを使う場合

複数のカレンダーを同時に表示したい場合は、すべてのカレンダー名をメモするます。
後で `settings.json` に配列で設定するます！

---

## ✅ ステップ 3: タスクリスト名を確認

### 3-1. Nextcloud Web UI で確認

1. Nextcloud にログインするます
2. 左上のメニュー → **「タスク」** をクリック
3. 左側のタスクリストを確認
4. 使用したいタスクリストの**名前**をメモするます
   - 例: `tasks`（個人タスク）
   - 例: `family-todo`（家族の TODO）

### 3-2. WebDAV/CalDAV URL から確認（上級者向け）

タスクリストの CalDAV URL を確認できます：
```
https://your-nextcloud-server.com/remote.php/dav/calendars/username/task-list-name/
```

最後の `task-list-name` がタスクリスト名なのです。

### 3-3. 複数タスクリストを使う場合

複数のタスクリストを同時に表示したい場合は、すべてのタスクリスト名をメモするます。
後で `settings.json` に配列で設定するます！

---

## ⚙️ ステップ 4: settings.json を設定

### 4-1. 基本設定（1つずつのカレンダーとタスクリスト）

`data/settings.json` に以下を設定するます：

```json
{
  "refreshIntervals": {
    "weatherSec": 300,
    "calendarSec": 300,
    "tasksSec": 300
  },
  "location": {
    "cityName": "姫路市",
    "country": "JP"
  },
  "nextcloud": {
    "serverUrl": "https://your-nextcloud-server.com",
    "username": "your-username",
    "password": "your-app-password-or-main-password",
    "calendarNames": ["family"],
    "taskListNames": ["tasks"]
  },
  "weather": {
    "provider": "open-meteo",
    "apiKey": "",
    "baseUrl": "https://api.open-meteo.com"
  }
}
```

**設定項目の説明**:
- `serverUrl`: Nextcloud サーバーの URL（https://で始まる）
- `username`: Nextcloud のユーザー名
- `password`: アプリパスワード（または メインパスワード）
- `calendarNames`: カレンダー名の配列（複数指定可能）
- `taskListNames`: タスクリスト名の配列（複数指定可能）

### 4-2. 複数カレンダー・タスクリストを使う場合

複数のカレンダーやタスクリストを同時に表示したい場合：

```json
{
  "nextcloud": {
    "serverUrl": "https://your-nextcloud-server.com",
    "username": "your-username",
    "password": "your-app-password",
    "calendarNames": ["family", "work", "personal"],
    "taskListNames": ["tasks", "family-todo", "shopping"]
  }
}
```

アプリは**すべてのカレンダーとタスクリストからデータを取得**して統合表示するます！✨

---

## 🧪 ステップ 5: 接続テスト

### 5-1. バックエンドを起動

```bash
# ビルド
go build -o familydashboard ./cmd/server

# 起動
./familydashboard
```

### 5-2. API をテスト

```bash
# ステータス確認
curl http://localhost:8080/api/status

# カレンダー取得
curl http://localhost:8080/api/calendar

# タスク取得
curl http://localhost:8080/api/tasks
```

エラーが出る場合は、トラブルシューティングセクションを確認するます！

### 5-3. フロントエンドで確認

ブラウザで http://localhost:8080 にアクセスして、カレンダーとタスクが表示されるか確認するます。

---

## 🔐 セキュリティチェックリスト

- [ ] **アプリパスワードを使用している**（メインパスワードは使わない）
- [ ] `settings.json` が `.gitignore` に含まれている
- [ ] `settings.json` をリポジトリにコミットしていない
- [ ] パスワードを他の人に共有していない
- [ ] HTTPS 接続を使用している（HTTP は使わない）

**もしパスワードを誤ってコミットした場合:**
1. Nextcloud でアプリパスワードを削除
2. 新しいアプリパスワードを生成
3. git history から削除（`git filter-repo` など）

---

## ⚠️ トラブルシューティング

### 🚫 エラー: 「認証エラー」「401 Unauthorized」

#### 原因:
- ユーザー名またはパスワードが間違っている
- アプリパスワードが無効になっている

#### 解決方法:
1. Nextcloud にログインできるか確認
2. アプリパスワードが正しいか確認
3. 新しいアプリパスワードを生成して試す
4. `settings.json` の `username` と `password` を再確認

### 🚫 エラー: 「カレンダーが見つかりません」

#### 原因:
- カレンダー名が間違っている
- カレンダーが存在しない
- 大文字小文字が一致していない

#### 解決方法:
1. Nextcloud Web UI でカレンダー名を再確認
2. `calendarNames` の綴りを確認（大文字小文字も正確に）
3. カレンダーが実際に存在するか確認
4. 複数カレンダーの場合、1つずつ試して問題のカレンダーを特定

### 🚫 エラー: 「タスクが取得できません」

#### 原因:
- タスクリスト名が間違っている
- タスクアプリが有効になっていない

#### 解決方法:
1. Nextcloud Web UI でタスクリスト名を再確認
2. `taskListNames` の綴りを確認
3. Nextcloud のアプリ管理で **Tasks** アプリが有効か確認
4. タスクリストに少なくとも1つのタスクがあるか確認（テスト用に作成）

### 🔴 エラー: 「SSL 証明書エラー」

#### 原因:
- 自己署名証明書を使用している
- 証明書が期限切れ

#### 解決方法（非推奨）:
開発環境でのみ、SSL 検証をスキップできますが**本番では絶対に使わない**でください：

```go
// 開発環境のみ: SSL検証をスキップ（本番では使用禁止）
http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{
    InsecureSkipVerify: true, // 本番では false にすること！
}
```

**推奨**: 正規の SSL 証明書（Let's Encrypt など）を使用するます。

### 📊 エラー: 「データが空です」

#### 原因:
- カレンダーやタスクリストにデータがない
- 日付範囲外

#### 解決方法:
1. Nextcloud Web UI でカレンダーにイベントを作成
2. タスクリストにタスクを作成
3. 今日から7日以内のイベントを作成（アプリは7日分を取得）
4. API レスポンスをログで確認

### 🌐 エラー: 「サーバーに接続できません」

#### 原因:
- サーバー URL が間違っている
- ネットワーク接続の問題
- ファイアウォール/ポート制限

#### 解決方法:
1. `serverUrl` が正しいか確認（`https://` で始まっているか）
2. ブラウザで Nextcloud にアクセスできるか確認
3. サーバーが外部からアクセス可能か確認（ファイアウォール設定）
4. ポート番号が必要な場合は追加（例: `https://server.com:8443`）

---

## 🏠 ローカル Nextcloud サーバーの場合

ローカルネットワークで Nextcloud を動かしている場合:

```json
{
  "nextcloud": {
    "serverUrl": "http://192.168.1.100:8080",
    "username": "admin",
    "password": "your-app-password"
  }
}
```

**注意**:
- HTTP は暗号化されないため、ローカルネットワークでのみ使用してください
- 外部アクセスの場合は HTTPS を使用してください

---

## 📚 参考リンク

- [Nextcloud Calendar ドキュメント](https://docs.nextcloud.com/server/latest/user_manual/en/groupware/calendar.html)
- [Nextcloud Tasks ドキュメント](https://docs.nextcloud.com/server/latest/user_manual/en/groupware/tasks.html)
- [CalDAV 仕様](https://datatracker.ietf.org/doc/html/rfc4791)
- [WebDAV 仕様](https://datatracker.ietf.org/doc/html/rfc4918)

---

## 💡 Tips

### カレンダーの色を設定する

Nextcloud のカレンダー設定で色を変更すると、FamilyDashboard でも同じ色が表示されるます！

1. Nextcloud → カレンダー
2. カレンダー名の横の **⋮** をクリック
3. **「編集」** → 色を選択
4. 保存するます

### 共有カレンダーを使う

家族で共有するカレンダーを作成できるます：

1. Nextcloud → カレンダー
2. 新しいカレンダーを作成
3. カレンダー名の横の **⋮** → **「共有」**
4. 家族のユーザー名を入力して共有
5. 各ユーザーの `settings.json` で同じカレンダー名を設定

これで家族全員が同じイベントを見られるのです！✨

---

がんばってせっていしてくださいなのですー！🥜✨
