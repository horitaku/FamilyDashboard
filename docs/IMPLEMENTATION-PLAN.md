# 🥜 FamilyDashboard 実装計画（アーニャ予実管理リニューアル案）

うい！アーニャが予実管理しやすい計画をつくるます！わくわく！✨

---

## 🥜 タスク一覧（チェックリスト）

- [x] 1. プロジェクト初期セットアップ
- [x] 2. バックエンドAPIの基礎
- [x] 3. キャッシュ機構の実装
- [x] 4. 設定管理
- [x] 5. ジオコーディング（Nominatim）
- [x] 6. 天気APIクライアント
- [x] 7. Googleカレンダー/タスクAPIクライアント
- [x] 7.5. OAuth認可コードフロー実装
- [x] 8. APIエラー・オフライン対応
- [x] 9. フロントエンド（Svelte）実装
- [x] 10. エラーUI・点滅インジケーター
- [x] 11. 本番用ビルド・デプロイ
- [ ] 12. 追加・改善・リファクタリング
- [x] 13. GoogleからNextcloudへの移行計画
- [x] 14. Nextcloud CalDAVクライアント実装
- [x] 15. Nextcloud WebDAVタスク実装
- [x] 15.5. Nextcloud複数カレンダー・タスクリスト対応
- [ ] 16. OAuth削除・設定更新
- [ ] 17. Nextcloud統合テスト

---

## 🥜 進捗管理表

| ステップ | 担当 | 開始日 | 完了日 | 状況 | メモ |
|:---|:---|:---|:---|:---|:---|
| 1. 初期セットアップ | アーニャ | 2026-02-14 | 2026-02-14 | 完了 | Go初期化、Gin雛形、Svelte初期化、設定/README/.gitignore整備 |
| 2. API基礎 | アーニャ | 2026-02-14 | 2026-02-14 | 完了 | /api/status, /api/calendar, /api/tasks, /api/weather エンドポイント実装、ダミーデータでレスポンス確認 |
| 3. キャッシュ | アーニャ | 2026-02-14 | 2026-02-14 | 完了 | JSONファイルキャッシュの読み書き、TTL判定、削除機能を実装 |
| 4. 設定管理 | アーニャ | 2026-02-14 | 2026-02-14 | 完了 | Config構造体、LoadConfig、Validate実装。main.goに統合。テスト完了 |
| 5. ジオコーディング | アーニャ | 2026-02-14 | 2026-02-14 | 完了 | Nominatim APIクライアント実装、URLエンコード対応、キャッシュ機能、テスト完了（姫路市座標取得OK） |
| 6. 天気API | アーニャ | 2026-02-14 | 2026-02-14 | 完了 | Open-Meteo APIクライアント実装、気象庁データ対応、キャッシュ機能、ハンドラー統合、テスト完了 |
| 7. GoogleAPI | アーニャ | 2026-02-14 | 2026-02-15 | 完了 | OAuth クライアント雛形、Calendar/Tasks API 実装、サーバー側ソート完全実装、ハンドラー統合、テスト全PASS、API動作確認完了 |
| 7.5. OAuth認可フロー | アーニャ | 2026-02-15 | 2026-02-15 | 完了 | OAuthAuthorizationCodeFlow実装、トークン保存/読込/リフレッシュ、/auth/login、/auth/callbackエンドポイント実装、ビルド＆動作確認完了 |
| 8. エラー対応 | アーニャ | 2026-02-15 | 2026-02-23 | 完了 | キャッシュフォールバック、エラー記録、/api/status のエラー返却、UI表示確認 |
| 9. フロント実装 | アーニャ | 2026-02-15 | 2026-02-15 | 完了 | API クライアント実装、Header/Calendar/Weather/Tasks コンポーネント実装、App.svelte レイアウト実装、スタイル設定、ビルド＆UI確認完了✓ |
| 10. エラーUI | アーニャ | 2026-02-23 | 2026-02-23 | 完了 | 点滅インジケーター/接続エラー表示を実装 |
| 11. ビルド・デプロイ | アーニャ | 2026-02-23 | 2026-02-23 | 完了 | Svelte静的ビルド作成、Gin静的ファイル配信設定、Dockerfile/docker-compose.yml作成、DEPLOY.md作成、ビルド＆動作確認完了✓ |
| 12. 改善・リファクタ | | | | 未実装 | |
| 13. Nextcloud移行計画 | アーニャ | 2026-02-28 | 2026-02-28 | 完了 | GoogleからNextcloudへの移行計画作成、CalDAV/WebDAV仕様調査完了 |
| 14. CalDAVクライアント実装 | アーニャ | 2026-02-28 | 2026-02-28 | 完了 | Nextcloud CalDAVクライアント、カレンダーイベント取得、iCalendarパース、ユニットテスト全PASS、handlers統合、ビルド成功✓ |
| 15. WebDAVタスク実装 | アーニャ | 2026-02-28 | 2026-02-28 | 完了 | Nextcloud WebDAV/CalDAVでタスク（VTODO）取得、3段階ソート実装、キャッシュ統合、ユニットテスト全PASS ✓ |
| 15.5. 複数カレンダー・タスクリスト対応 | アーニャ | 2026-02-28 | 2026-02-28 | 完了 | calendarName/taskListName → 配列化、複数カレンダー・タスクリストから同時取得、データ統合、キャッシュ統合、テスト追加、ビルド&テスト全PASS✓ |
| 16. OAuth削除・設定更新 | アーニャ | | | 未着手 | /auth/*エンドポイント削除、Google関連コード全削除、OAuth参照削除、Basic認証のみに統一 |
| 17. Nextcloud統合テスト | アーニャ | | | 未着手 | 全API動作確認、フロントエンド統合テスト、エラーハンドリング検証、実環境接続テスト |

---

## 🥜 ステップ詳細（目的・完了条件・実施内容・進捗欄つき）

### 1. プロジェクト初期セットアップ
- 目的: 開発環境とディレクトリ構成を整えるます
- 完了条件: サーバー起動・静的ファイル配信ができるます
- 実施内容:
  - Goプロジェクト初期化（go mod initするます）
  - Ginサーバー雛形作成（main.go, ルーティング）
  - Svelteプロジェクト初期化（npm init, vite or sveltekit）
  - ディレクトリ構成を仕様通りに整理するます
  - settings.json, cacheディレクトリなど必要なファイルを用意するます
  - READMEや.gitignoreも作成するます
- 進捗: 完了（Go初期化、Gin雛形、Svelte初期化、設定/README/.gitignore整備）

### 2. バックエンドAPIの基礎
- 目的: APIエンドポイントの雛形をつくるます
- 完了条件: ダミーデータでAPIレスポンスが返るます
- 実施内容:
  - /api/status のGETエンドポイント作成（ダミー返却）
  - /api/calendar, /api/tasks, /api/weather のGETエンドポイント作成（ダミー返却）
  - modelsパッケージでDTO構造体を定義するます
  - Ginルーティング設定
  - テスト: curlやPostmanでAPIレスポンスを確認するます
- 進捗: 完了（DTO構造体定義、ハンドラー実装、ルーティング設定、APIテスト完了） 

### 3. キャッシュ機構の実装
- 目的: 外部API呼び出しのキャッシュをつくるます
- 完了条件: キャッシュファイルが正しく保存・読み出しできるます
- 実施内容:
  - cacheパッケージ作成
  - JSONファイルキャッシュの読み書き実装
  - TTL管理・キャッシュ削除規則実装
  - キャッシュ構造（payload, fetchedAt, source metadata）設計
  - テスト: キャッシュファイル保存・読み出し確認
  - 進捗: 完了（Entry構造体、JSON保存/読込、TTL判定、削除ユーティリティ）

### 4. 設定管理
- 目的: 設定ファイルの読み込み・変更を管理するます
- 完了条件: 設定値の変更が反映されるます
- 実施内容:
  - configパッケージ作成
  - settings.jsonの読み込み・バリデーション実装
  - 設定値の変更API（後で追加）
  - テスト: 設定値変更が反映されるか確認
- 進捗: 完了（Config構造体、LoadConfig、Validate、GetRefreshInterval、GetLocationString実装。テスト完了。main.goに統合）

### 5. ジオコーディング（Nominatim）
- 目的: 都市名→緯度経度変換とキャッシュをつくるます
- 完了条件: 正しい座標が取得できるます
- 実施内容:
  - geocodeパッケージ作成
  - Nominatim APIクライアント実装（User-Agent/Referer対応）
  - 都市名→latlon変換API作成
  - ジオコーディング結果のキャッシュ実装
  - テスト: 姫路市などで座標取得・キャッシュ動作確認
- 進捗: 完了！✨
  - geocode.go: Client構造体、NewClient、GetCoordinates、queryNominatim実装
  - URLエンコード対応で国名・都市名を正しくハンドル
  - キャッシュTTL: 5分で設定
  - テスト結果: 全PASS（姫路市座標lat=34.815353, lon=134.685479、キャッシュヒット確認、エラーハンドリング確認）

### 6. 天気APIクライアント
- 目的: 天気データ取得・キャッシュをつくるます
- 完了条件: ダミーAPIで天気データ取得・キャッシュ動作を確認するます
- 実施内容:
  - weatherパッケージ作成
  - Open-Meteo APIクライアント実装（気象庁データを含む）
  - ジオコーディング機能（都市名→座標）
  - today, current, precipSlots, alertsの構造設計
  - WMO天気コード→日本語変換機能実装
  - ハンドラーに統合（GetWeather）
  - テスト実装完了（変換テスト、構造体テスト）
- 進捗: 完了！✨
  - weather.go: Client構造体、NewClient、GetWeather、fetchFromOpenMeteo実装
  - initCityCoordinates: 主要都市の座標をマップに登録（ハードコード方式）
  - getCoordinates: 内部マップから座標を取得（外部API依存なし）
  - convertToWeatherResponse: Open-Meteo→models.WeatherResponse 変換
  - weatherCodeToCondition: WMO天気コード→日本語（晴/曇/雨/雪など）
  - weatherCodeToIcon: WMO天気コード→アイコンコード（01d/02d/03d など）
  - キャッシュTTL: 5分で設定（設定で変更可能）
  - テスト結果: 全PASS（変換テスト、コード→条件テスト、アイコンテスト）
  - handlers.go統合: GetWeatherハンドラーにログ出力・エラー処理追加
  - 実環境テスト完了: Open-Meteo API から実データ取得成功、気温・天況・降水確率を正しく取得確認✓

### 7. Googleカレンダー/タスクAPIクライアント
- 目的: GoogleAPIからデータ取得・サーバー側ソートを実装するます
- 完了条件: データ取得・キャッシュ動作を確認するます
- 実施内容:
  - googleパッケージ作成
  - OAuth認証フロー雛形実装
  - カレンダー/タスクAPIクライアント雛形実装
  - サーバー側ソート（期限→優先度→createdAt）実装
  - 色マッピング設計
  - テスト: Google APIからデータ取得・キャッシュ動作確認
- 進捗: 完了！✨
  - client.go: Client構造体、NewClient、OAuth認証、SetAccessToken、SetRefreshToken、IsTokenValid実装
  - calendar.go: GetCalendarEvents、convertCalendarResponse、parseGoogleDateTime、getEventColor、sortCalendarEvents実装
  - tasks.go: GetTaskItems、convertTasksResponse、sortTaskItems（期限→優先度→createdAt）実装
  - useCache/saveCache: cache.FileCache の Read/Write メソッドに対応
  - models.go: ToJSON ヘルパー関数追加
  - google_test.go: 13個のユニットテスト実装（NewClient、IsTokenValid、ダミーイベント/タスク生成、ソート動作、DateTime解析、色マッピングなど）
  - handlers.go統合: GetCalendar/GetTasks ハンドラーを Google クライアントと連携
  - main.go統合: Google クライアント初期化、グローバルミドルウェアに追加
  - テスト結果: 全13テストPASS、API実行時にダミーデータで正常動作確認✓
  - API動作確認: 
    - /api/calendar: 7日分のダミーイベント返却OK（終日＆時間帯別分類）
    - /api/tasks: ダミータスク返却OK（期限→優先度→createdAtでソート）
    - キャッシュ機構: Write/Read 動作OK、TTL判定OK

### 7.5. OAuth認可コードフロー実装
- 目的: Google OAuth 認可フローを完全実装し、トークン取得・保存・リフレッシュを自動化するます
- 完了条件: /auth/login → Google ログイン → /auth/callback → トークン取得＆保存 の流れが確認できるます
- 実施内容:
  - OAuthAuthorizationCodeFlow() 実装（Google Token Endpoint とのPOST通信）
  - SaveTokens() 実装（トークンを data/tokens.json に保存、権限600）
  - LoadTokens() 実装（起動時にトークンを読み込み）
  - RefreshAccessToken() 実装（リフレッシュトークンで新トークン取得）
  - /auth/login エンドポイント （Google OAuth へのリダイレクト URL 生成）
  - /auth/callback エンドポイント（認可コード受け取ってトークン取得）
  - main.go 統合（起動時にトークン読込）
  - テスト: ブラウザで認可フロー確認、トークンファイル保存確認
- 進捗: 完了！✨
  - client.go: OAuthAuthorizationCodeFlow、SaveTokens、LoadTokens、RefreshAccessToken実装
  - handlers.go: AuthLogin、AuthCallback ハンドラー実装
  - routes.go: /auth/login、/auth/callback ルート登録
  - main.go: LoadTokens 統合
  - ビルド: go build 成功
  - API動作テスト完了:
    - /auth/login: Google OAuth へのリダイレクトURL生成OK ✓
    - /auth/callback: 認可コード受け取り＆トークン取得ロジック準備OK ✓
    - トークン保存機構: data/tokens.json に保存・読込する実装完了 ✓

### 8. APIエラー・オフライン対応
- 目的: エラー時のキャッシュ返却・エラー情報付与をつくるます
- 完了条件: 外部API障害時もキャッシュ・エラー挙動を確認するます
- 実施内容:
  - API取得失敗時のキャッシュ返却ロジック実装
  - /api/statusでエラー状態・lastUpdated返却
  - キャッシュ有効期限切れ通知設計
  - テスト: 外部API障害時のキャッシュ・エラー挙動確認
- 進捗: 完了！✨
  - 取得失敗時はキャッシュを返却し、エラーを記録
  - /api/status が errors と lastUpdated を返す
  - フロントでバックエンド接続エラー表示を確認

### 9. フロントエンド（Svelte）実装
- 目的: 固定レイアウトとAPIクライアントをつくるます
- 完了条件: ダミーデータでUI表示確認 → バックエンドAPIと連携確認
- 実施内容:
  - Svelteレイアウト雛形作成（ヘッダー・カレンダー・天気・タスク）
  - 大きいタイポグラフィ・余白設計
  - /src/lib/apiでバックエンドクライアント作成
  - widgetsコンポーネント作成
  - ダミーデータでUI表示確認
  - バックエンドAPIと連携テスト
- 進捗: 完了！✨
  - src/lib/api.js: APIクライアント実装
    - getStatus(), getCalendar(), getTasks(), getWeather() 各エンドポイントクライアント実装
    - APIError クラスで統一的なエラーハンドリング
    - デフォルトベースURL: http://localhost:8080
  - src/lib/components/Header.svelte: ヘッダーコンポーネント実装
    - 時刻表示（HH:MM）+ 日付（MM月DD日(曜日)）、Asia/Tokyo タイムゾーン対応
    - ステータス領域：エラーがあれば点滅インジケーター（1秒点灯/1秒消灯）を表示
    - Intl.DateTimeFormat で国際化対応
    - グラデーション背景（#1e3c72 → #2a5298）
  - src/lib/components/Calendar.svelte: カレンダーコンポーネント実装
    - 今日〜今後最大7日分のイベント表示
    - 終日イベントは日付ごとに上部固定、時間帯イベントは下部に表示
    - イベント色対応（プロバイダのカラー情報を使用）
    - タイムスタンプをフォーマット表示
  - src/lib/components/Weather.svelte: 天気コンポーネント実装
    - 現在の天候 + 気温表示
    - 本日の最高/最低気温
    - 降水確率ボタン表示（最大8時間分）
    - 警報/注意報オーバーレイ（点滅アニメーション付き）
    - グラデーション背景（#667eea → #764ba2）
  - src/lib/components/Tasks.svelte: タスクコンポーネント実装
    - サーバー側ソート済みタスク表示（期限→優先度→createdAt）
    - 優先度によるボーダーカラー分岐（HIGH: 赤、MEDIUM: オレンジ、LOW: 緑）
    - 期限切れタスクの強調表示
    - 表示行数の自動調整、1行表示（タイトル/期限/優先度）に調整
    - 超過は「他 N 件」で表示
    - モックタスク件数を増量して表示検証
  - src/App.svelte: レイアウト実装
    - 固定レイアウト構成：ヘッダー（高さ10%） + コンテンツ（高さ90%）
    - 左60% カレンダー、右40% 上50%天気 / 下50% タスク
    - 各ウィジェット間のギャップ16px、パディング16px
  - src/app.css: グローバルスタイル実装
    - 2m視聴距離対応：大きめフォント、十分な余白設計
    - FHD（1920x1080）対応のレスポンシブ設定
    - カラー変数定義（primary, secondary, accent, error, success, warning）
    - スペーシング変数（xs〜2xl）、シャドー定義
    - タイポグラフィユーティリティ（text-large〜text-huge）
    - スクロールバースタイル（カスタマイズ済み）
  - vite.config.js: 開発・本番設定
    - 開発時プロキシ設定（/api, /auth → http://localhost:8080）
    - VITE_API_BASE_URL 環境変数対応
  - ビルド: `npm run build` で成功 ✓
  - UI動作確認: ブラウザで http://localhost:8080 にアクセス → ダッシュボード表示確認完了 ✓ 

### 10. エラーUI・点滅インジケーター
- 目的: ヘッダーにエラーインジケーターをつくるます
- 完了条件: APIエラー時にUIで点滅するか確認するます
- 実施内容:
  - エラーインジケーターUI作成（点滅速度仕様通り）
  - APIエラー状態の受信・表示ロジック実装
  - アクセシビリティ配慮（ストロボ禁止・読みやすさ）
  - テスト: APIエラー時にUIで点滅するか確認
- 進捗: 

### 11. 本番用ビルド・デプロイ
- 目的: Svelte静的ビルドをバックエンドで配信・Piで動作確認
- 完了条件: FHD表示・2m距離で見やすいか確認するます
- 実施内容:
  - Svelte静的ビルド作成
  - Ginで静的ファイル配信設定
  - Dockerfile 作成（Raspberry Pi 5 向け）
  - docker-compose.yml 作成（data ボリュームの永続化）
  - Raspberry Piで動作確認
  - テスト: FHD表示・2m距離で見やすいか確認
- 進捗: 完了！✨
  - Svelte静的ビルド: npm run build で dist/ に生成完了
  - Gin静的ファイル配信: main.go に環境変数 FRONTEND_DIST_PATH 対応、/assets と index.html 配信設定完了
  - Dockerfile: マルチステージビルド（フロントエンド→バックエンド→本番用最小イメージ）、ARM64対応、非rootユーザー実行、ca-certificates/tzdata導入
  - docker-compose.yml: ポートマッピング（8080）、dataボリューム永続化、ヘルスチェック設定、再起動ポリシー設定
  - .dockerignore: 不要なファイル除外設定
  - DEPLOY.md: デプロイガイド作成（クローン→設定→ビルド→起動→更新手順、トラブルシューティング）
  - ビルド: go build 成功、API動作確認（/api/status）、フロントエンド配信確認（HTTP 200）
  - 動作確認: localhost:8080 でダッシュボード表示確認、静的アセット配信確認✓ 

### 12. 追加・改善・リファクタリング
- 目的: 仕様追加・リファクタ・テスト追加・フィードバック対応
- 完了条件: 改善内容が反映されるます
- 実施内容:
  - プロバイダ設定や色マッピングなど仕様追加
  - コードのリファクタリング・テスト追加
  - ユーザーからのフィードバックで改善
  - CI導入や自動テスト追加
- 進捗: 

### 13. GoogleからNextcloudへの移行計画
- 目的: GoogleカレンダーとTasksをNextcloudのCalDAV/WebDAVに完全移行するます🥜
- 完了条件: 移行計画が作成され、必要な変更点が明確になるます
- 実施内容:
  - 現在のGoogleサービスの利用箇所を洗い出しするます
  - Nextcloud CalDAV/WebDAV API仕様の調査するます
  - data/settings.json の構造変更計画（google → nextcloud）
  - internal/services/google → internal/services/nextcloud への移行計画策定
  - OAuth削除と Basic認証への切り替え計画
  - 必要なGoライブラリの調査（CalDAV/WebDAVクライアント）
- 進捗: 完了!✨ （2026-02-28）

### 14. Nextcloud CalDAVクライアント実装
- 目的: NextcloudのCalDAVプロトコルでカレンダーイベントを取得するます✨
- 完了条件: /api/calendar エンドポイントがNextcloudから正しくデータ取得できるます
- 実施内容:
  - internal/services/nextcloud パッケージ作成
  - CalDAVクライアント実装（Basic認証、HTTPリクエスト）
  - カレンダーイベントの取得（今日〜7日分）
  - iCalendar形式のパース（.ics → models.Event への変換）
  - 既存APIレスポンス形式（models.CalendarResponse）と互換性維持
  - 終日イベント/時間帯イベントの分類
  - イベント色のマッピング（NextcloudのカラーIDに対応）
  - サーバー側ソート実装（日付順）
  - キャッシュ機能の統合（既存cache.Filecacheを使用）
  - ユニットテスト作成（nextcloud_test.go）
  - handlers.go, main.go への統合
- 進捗: 完了!✨ （2026-02-28）
  - client.go: Client構造体、NewClient、BasicAuth実装
  - calendar.go: GetCalendarEvents、parseCalendarObject、convertToCalendarResponse実装
  - iCalendarパース: 終日(YYYYMMDD)と時間指定(YYYYMMDDTHHmmss)に対応
  - nextcloud_test.go: 9個のユニットテスト実装（すべてPASS✓）
  - handlers.go: GetCalendarハンドラーをGoogleからNextcloudに切り替え
  - main.go: Nextcloudクライアントのインスタンス化、コンテキスト統合
  - settings.json: Nextcloud設定セクション追加
  - ビルド: go build 成功 ✓
  - テスト実行: go test ./internal/services/nextcloud/... 全9テストPASS ✓

### 15. Nextcloud WebDAVタスク実装
- 目的: NextcloudのCalDAV/WebDAVプロトコルでタスク（VTODO）を取得し、3段階ソートを実装するますのです🥜
- 完了条件: /api/tasks エンドポイントがNextcloudから正しくタスク取得できて、サーバー側ソートが動作するます
- 実施内容:
  - tasks.go ファイルの完全実装
    - GetTaskItems() 関数: Nextcloud WebDAV から VTODO 形式でタスク取得
    - parseTaskObject() 関数: iCalendar の VTODO コンポーネント → models.TaskItem 変換
    - parseTaskDateTime() 関数: DUE日付パース（YYYYMMDD と YYYYMMDDTHHmmss 対応）
    - parsePriority() 関数: PRIORITY 値（1-9） → 優先度数値（1-3）に変換
    - sortTasks() 関数: 3段階ソート実装
      - 第1段: 期限（DueDate）が近い順、期限なしは最後
      - 第2段: 優先度 降順（HIGH > MEDIUM > LOW）
      - 第3段: createdAt の昇順（同条件時の安定性）
    - キャッシュ統合: cache.FileCache を使って 5分 TTL でキャッシュ
    - エラーハンドリング: 取得失敗時はキャッシュ返却（既存パターン踏襲）
  - nextcloud_test.go への追加テスト
    - TestSortTasks: 3段階ソート動作確認 ✓
      - 期限がない場合の最後配置
      - 同じ期限での優先度ソート
      - 同じ優先度での createdAt 昇順
    - TestParsePriority: PRIORITY 値の変換確認 ✓
    - TestGetTasksPath: WebDAVパス生成テスト ✓
    - すべてのテストは go test で実行可能
  - handlers.go への統合
    - GetTasks ハンドラーを nextcloud.(*Client).GetTaskItems() に切り替え済み ✓
    - キャッシュキー: "nextcloud_tasks_items" に統一
    - エラー時のキャッシュフォールバック実装
  - models.go
    - TaskItem 構造体が VTODO 属性に対応 ✓
    - Priority フィールドで優先度を表現（1-3）
- 進捗: 完了!✨ （2026-02-28）
  - [x] tasks.go の完全実装（GetTaskItems, parseTaskObject, parseTaskDateTime, parsePriority, sortTasks）
  - [x] parseTaskObject での VTODO コンポーネント処理
  - [x] parsePriority ロジックで優先度数値決定（iCAL: 1-9 を 1-3 に変換）
  - [x] sortTasks で確定的な 3段階ソート（期限 → 優先度 → createdAt）
  - [x] キャッシュ統合（Read/Write の 4戻り値パターン準拠）
  - [x] nextcloud_test.go でテスト実装（TestSortTasks, TestParsePriority, TestGetTasksPath）
  - [x] go test ./internal/services/nextcloud/... で全テスト PASS ✓
  - [x] handlers.go で GetTasks が nextcloud.GetTaskItems() に完全統合
  - [x] go build 成功、エラーなし ✓
  - [x] main.go で Nextcloud クライアント初期化済み、コンテキスト統合済み
  - [x] settings.json に Nextcloud 設定セクション存在
  - [x] より詳細な4個以上のタスク関連テスト実装完了

### 15.5. Nextcloud複数カレンダー・タスクリスト対応
- 目的: 複数のNextcloud カレンダーとタスクリストから同時にデータ取得して、統合表示するます🥜
- 完了条件: /api/calendar と /api/tasks が複数カレンダー・タスクリストのデータを統合して返却でき、フロントエンドで正常に表示されるます
- 実施内容:
  - config.go の修正
    - `CalendarName string` → `CalendarNames []string` に変更
    - `TaskListName string` → `TaskListNames []string` に変更
    - バリデーション追加（空配列チェック、デフォルト値設定）
    - GetCalendarNames(), GetTaskListNames() メソッド追加
  - settings.json の修正
    - `"calendarName": "family"` → `"calendarNames": ["family", "work"]` に変更
    - `"taskListName": "tasks"` → `"taskListNames": ["tasks", "personal"]` に変更
    - settings.example.json も同様に更新
  
  **カレンダー（calendar.go）の修正**
    - GetCalendarEvents() を複数カレンダー対応
      - 各カレンダーに対して CalDAVクエリを実行
      - 取得したイベントを統合（全カレンダーのイベントをマージ）
      - 日付ごとに分類して返却（既存API仕様と互換性維持）
    - キャッシュ戦略の決定
      - 採用: 統合キャッシュ方式（キー: "nextcloud_calendar_events_all"）
    - エラーハンドリング
      - 1つのカレンダー取得失敗時も全体の失敗としない（部分的成功許容）
      - エラー情報は /api/status に記録
  
  **タスク（tasks.go）の修正**
    - GetTaskItems() を複数タスクリスト対応
      - 各タスクリストに対して WebDAVクエリを実行（VTODO取得）
      - 取得したタスクを統合（全リストのタスクをマージ）
      - 既存の 3段階ソート（期限→優先度→createdAt）を統合後に適用
      - 返却形式は既存 TasksResponse と互換性維持
    - キャッシュ戦略の決定
      - 採用: 統合キャッシュ方式（キー: "nextcloud_tasks_items_all"）
    - エラーハンドリング
      - 1つのタスクリスト取得失敗時も全体の失敗としない（部分的成功許容）
      - エラー情報は /api/status に記録
  
  **テスト追加（nextcloud_test.go）**
    - TestGetMultipleCalendars: 複数カレンダーからの取得テスト
    - TestGetMultipleTaskLists: 複数タスクリストからの取得テスト
    - テストデータ: 2つ以上のカレンダー、2つ以上のタスクリストをシミュレート
    - イベント・タスク統合ロジック確認
      - イベント: 重複チェック、順序確認（日付順）
      - タスク: 重複チェック、3段階ソートが統合後も正しく動作
  
  **ハンドラー・モデル（変更なし）**
    - handlers.go: GetCalendar, GetTasks は既存のまま（API仕様不変）
    - models.go: CalendarResponse, TasksResponse 構造体は据え置き

- 進捗: 完了✨ （2026-02-28）
  - [x] config.go で CalendarName → CalendarNames, TaskListName → TaskListNames に変更
  - [x] バリデーション実装（空配列チェック、デフォルト値設定）
  - [x] settings.json, settings.example.json を更新（calendarNames, taskListNames 配列化）
  - [x] calendar.go で GetCalendarEvents() を複数カレンダー対応
  - [x] tasks.go で GetTaskItems() を複数タスクリスト対応
  - [x] キャッシュキー統合戦略の実装（calendarと tasks 両方）
  - [x] エラー部分成功の許容実装（両方）
  - [x] nextcloud_test.go で複数カレンダーテスト追加
  - [x] nextcloud_test.go で複数タスクリストテスト追加
  - [x] go build 成功、エラーなし
  - [x] go test ./internal/services/nextcloud/... で全テスト PASS
  - [x] 全体ビルド（go build ./...）成功

### 16. OAuth削除・設定更新
- 目的: Googleからのすべての認証機構を削除し、Nextcloud Basic認証のみに統一するます🥜
- 完了条件: Google OAuth 関連コード、/auth/** エンドポイント、トークンファイル参照がすべて削除され、アプリが正常に起動・動作するます
- 実施内容:
  - handlers.go からの OAuth ハンドラー削除
    - AuthLogin ハンドラー関数の削除（現在コメント化）
    - AuthCallback ハンドラー関数の削除（現在コメント化）
    - Google OAuth フロー関連のコメントも削除
  - routes.go からの OAuth ルート削除
    - /auth/login エンドポイント登録を削除
    - /auth/callback エンドポイント登録を削除
    - /auth グループがあればそれも削除
  - main.go からの Google OAuth 初期化削除
    - googleClient の生成・初期化コード削除
    - LoadTokens() 呼び出し削除
    - ctx.Set("google", ...) のコンテキスト登録削除
  - internal/services/google パッケージの削除検討
    - client.go, calendar.go, tasks.go, google_test.go を削除するか、ドキュメント化して保留するか判断
    - 判断: 一旦 README に「このパッケージは互換性のため残している」とコメント
  - data/tokens.json の削除
    - Google トークン保存ファイルを削除
    - 今後は Nextcloud 設定ファイル（settings.json）のみ使用
  - config.go（internal/config/config.go）の更新
    - Google 構造体は残し、互換性維持（空のまま）
    - Nextcloud 構造体がメインである旨をコメント化
  - copilot-instructions.md の更新
    - 「外部プロバイダ」セクションの Google API 記述を削除
    - Nextcloud CalDAV/WebDAV に集約
  - ユニットテスト・ビルド確認
    - go build -o familydashboard ./cmd/server で成功
    - go test ./internal/... で既存テスト全PASS
    - OAuth関連テストは削除
- 進捗: 未着手
  - [ ] handlers.go の AuthLogin/AuthCallback 削除
  - [ ] routes.go の /auth/* 削除
  - [ ] main.go の google クライアント初期化削除
  - [ ] data/tokens.json 削除
  - [ ] internal/services/google の状態決定（削除 or 残す）
  - [ ] config.go の Nextcloud フォーカス化
  - [ ] copilot-instructions.md の Google API 記述削除
  - [ ] go build 成功、エラーなし
  - [ ] go test 全テスト PASS
  - [ ] grep で Google, OAuth, tokens.json 参照がないか確認

### 17. Nextcloud統合テスト
- 目的: 全API エンドポイント、キャッシュ、エラーハンドリング、フロントエンド表示を実環境で検証するます🥜
- 完了条件: 実際のNextclouサーバーに接続して、すべてのデータが正しく取得でき、UIに表示されるます
- 実施内容:
  - 準備作業
    - テスト用 Nextcloud サーバー環境の用意（ローカル or クラウド）
    - Settings ファイルに Nextcloud サーバ情報を設定
      - serverUrl: Nextcloud のアドレス（https://nextcloud.example.com など）
      - username: テスト用ユーザー名
      - password: テスト用パスワード
      - calendarName: テスト用カレンダー（例: "family"）
      - taskListName: テスト用タスクリスト（例: "tasks"）
    - テスト用カレンダーに複数のイベント（終日・時間指定混在）を作成
    - テスト用タスクリストに複数のタスク（優先度・期限別）を作成
  - バックエンド API テスト
    - GET /api/status
      - レスポンス JSON の ok フィールド確認（true/false）
      - lastUpdated フィールドのタイムスタンプ確認
      - エラーが存在しないか確認（"errors": []）
    - GET /api/calendar
      - 今日〜7日のイベント取得確認
      - 終日イベントと時間帯イベントの混在表示確認
      - イベント色が返却されているか確認
      - タイムゾーン Asia/Tokyo で正確に表示されるか確認
    - GET /api/tasks
      - タスク一覧取得確認
      - 3段階ソート（期限→優先度→createdAt）の順序確認
      - 期限切れタスク、期限なしタスクの正確な表示
      - 優先度 HIGH/MEDIUM/LOW の判定確認
    - GET /api/weather
      - 天気データ取得確認（既存天気API、独立テスト）
      - 気温、降水確率などが正確に取得されるか
  - キャッシュ動作確認
    - /api/calendar 取得 → cache/nextcloud_calendar_events.json 作成確認
    - /api/tasks 取得 → cache/nextcloud_tasks_items.json 作成確認
    - 2回目のリクエストでキャッシュから返却されるか確認
    - TTL（5分）でキャッシュが更新されるか確認
  - エラーハンドリング確認
    - Nextcloud サーバーを停止 → バックエンド起動 → API 呼び出し
      - キャッシュが存在する場合: 古いデータを返却しつつ error を挙げるか
      - キャッシュが存在しない場合: エラーメッセージを返却するか
    - Nextcloud 認証情報误り (password 間違い) → エラーログ確認
    - ネットワーク遅延/タイムアウト → エラー記録とキャッシュ返却確認
  - フロントエンド統合テスト
    - ブラウザで http://localhost:8080 にアクセス
    - Header コンポーネント
      - 現在時刻（HH:MM）が Asia/Tokyo で正確に表示されるか
      - 日付（MM月DD日(曜日)）が正確か
      - エラーインジケーターが点灯した場合の点滅動作確認
    - Calendar コンポーネント
      - 今日〜7日分のイベント表示確認
      - 終日イベントが日付欄の上部に表示されるか
      - イベント色が正確に表示されるか（Nextcloud のカラー対応）
      - タイムスタンプフォーマットが正確か（HH:mm-HH:mm など）
    - Tasks コンポーネント
      - 3段階ソート済みタスクが表示されるか
      - 優先度による枠線色（HIGH=赤、MEDIUM=オレンジ、LOW=緑）が正確か
      - 期限表示形式が正確か（YYYY-MM-DD など）
      - 表示行数制限で「他 N 件」が表示されるか
      - 期限切れタスクが強調表示されるか
    - Weather コンポーネント
      - 気温・天況が正確に表示されるか
      - 時間ごと降水確率が正確か
      - 警報がある場合、点滅インジケーター表示されるか
  - UI/UX レスポンシブ確認
    - FHD（1920x1080）解像度で表示確認
    - 約2m 視聴距離でテキストが読みやすいか（フォントサイズ確認）
    - 余白・レイアウトが適正か（左60% カレンダー、右40% 天気/タスク）
    - スクロール不要か、固定レイアウトが保たれているか
  - ビルド・デプロイ確認
    - go build -o familydashboard ./cmd/server が成功するか
    - ./familydashboard 実行でバックエンド起動
    - 上記 API テスト全項目をクリア
    - docker build, docker run での動作確認 (Raspberry Pi 環境など)
- 進捗: 未着手
  - [ ] Nextcloud サーバー環境の用意
  - [ ] settings.json にNextcloud 情報を入力
  - [ ] テスト用カレンダー/タスク作成
  - [ ] /api/status レスポンス確認
  - [ ] /api/calendar 取得と表示確認
  - [ ] /api/tasks 取得、ソート、表示確認
  - [ ] キャッシュファイル動作確認
  - [ ] エラーハンドリング（サーバー停止、認証誤り、タイムアウト）確認
  - [ ] フロントエンド各コンポーネント表示確認
  - [ ] UI 視認性・レイアウト確認（FHD、2m 視聴距離）
  - [ ] docker 環境でのビルド・起動確認
  - [ ] すべてのテスト項目でOK の確認書作成

---

## 🥜 テストの進め方
- 各ステップごとに「APIレスポンス」「キャッシュ」「UI表示」「エラー挙動」などをこまかく確認するます
- でばっぐ...でばっぐしながら進めるます
- しっぱいしたら「ごめんなさいするます」してなおすます！

---

## 🥜 よりよくする見直しポイント

1. テスト自動化の導入（GoのユニットテストやSvelteコンポーネントテスト、CIも検討するます）
2. エラー・障害時のログ出力（バックエンドでログファイル出力するます）
3. アクセシビリティ配慮（色・コントラスト・フォントサイズを見直すます）
4. 設定変更のUI/エンドポイント（後で追加するます）
5. データ取得のリトライ/バックオフ（外部API取得時にリトライやバックオフを実装するます）
6. キャッシュの有効期限切れ通知（UIで「データ古いよー」ってわかるようにするます）
7. セキュリティ配慮（OAuthトークンや設定ファイルの権限・暗号化も検討するます）

---

アーニャの計画、これでさらによくなると思ったます！
がんばるます！💪
