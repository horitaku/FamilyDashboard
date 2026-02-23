# google

Google API（Calendar/Tasks）のクライアントを実装するます。

## 概要

`google` パッケージは、Google Calendar API と Google Tasks API からデータを取得し、
キャッシュ・ソート・エラーハンドリング機能を提供するのです。

## ファイル構成

- `client.go`: Google API 基本クライアント（OAuth認証、キャッシュ管理）
- `calendar.go`: Google Calendar API 実装（イベント取得、変換、ソート）
- `tasks.go`: Google Tasks API 実装（タスク取得、サーバー側ソート）
- `google_test.go`: ユニットテスト（13個、全PASS）

## 主要機能

### Client

- **OAuth認証フロー**: SetAccessToken/SetRefreshToken でトークン管理
- **トークン有効性判定**: IsTokenValid() で有効期限を確認
- **キャッシュ管理**: useCache/saveCache で FileCache と統合

### Calendar

- **イベント取得**: GetCalendarEvents() で次7日分のイベントを取得
- **イベント分類**: 終日イベント（AllDay）と時間帯付きイベント（Timed）を分類
- **色マッピング**: Google Calendar の色ID から16進カラーコードに変換
- **ソート**: 日付順→ 終日イベントはタイトル順 → 時間帯イベントは時系列順

### Tasks

- **タスク取得**: GetTaskItems() でダッシュボード用タスクを取得
- **サーバー側ソート**: 期限（昇順） → 優先度（降順） → createdAt（昇順）
- **期限なし対応**: 期限を持たないタスクは最後に配置

## キャッシュ機構

- TTL: 設定の refreshIntervals に従う（既定: 5分）
- ソース: "google" として meta に記録
- 失敗時: キャッシュが有効なら使用、無ければエラー返却

## トークン無し時の動作

トークンが無い場合はダミーデータを返すます：

- Calendar: 7日分のダミーイベント
- Tasks: ダミータスク4件（ソート規則適用済）

これにより、Google 認証を設定する前も UI テストが可能なのです。

## 今後の実装

- [ ] OAuth認可コードフロー完全実装（ブラウザロジン）
- [ ] リフレッシュトークンの自動更新
- [ ] トークン有効期限切れ時の自動延長
- [ ] 複数カレンダー対応
- [ ] カスタム優先度フィールド対応（Google Tasks）
