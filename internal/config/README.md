# config 📋

settings.json の読み込みとバリデーション機能を実装するますね。

## 機能一覧 🎯

### Config 構造体
- `RefreshIntervals`: データソース別（天気/カレンダー/タスク）の更新間隔を設定
- `Location`: ロケーション情報（都市名、国コード）
- `Google`: Google API の認証・設定（clientId, clientSecret等）
- `Weather`: 天気API の設定（プロバイダ、APIキー等）

### 主要な関数・メソッド

#### LoadConfig(filepath string) (*Config, error)
- settings.json ファイルからJSON設定を読み込み
- 自動的にバリデーションを実行
- エラー時は詳細なエラーメッセージを返す

```go
cfg, err := config.LoadConfig("./data/settings.json")
if err != nil {
    log.Fatal(err)
}
```

#### Validate() error
- 設定値のバリデーション
- チェック内容：
  - 更新間隔は正の数か？（weatherSec, calendarSec, tasksSec > 0）
  - ロケーション情報は必須か？（cityName, country が空でないか？）
- バリデーション失敗時は詳細なエラーメッセージを返す

#### GetRefreshInterval(source string) time.Duration
- "weather", "calendar", "tasks" に対応した更新間隔を Duration で取得
- source が無効な場合は既定値（5分）を返す

```go
weatherInterval := cfg.GetRefreshInterval("weather")
// → time.Duration（秒）
```

#### GetLocationString() string
- ロケーション情報を "CityName, Country" 形式の文字列で返す
- cityName が空の場合は "Unknown" を返す

```go
locStr := cfg.GetLocationString()
// → "姫路市, JP"
```

#### LoadedAt() time.Time
- 設定の読み込み時刻を返す（デバッグ用途）

## ファイル構成

- `config.go`: Config構造体、LoadConfig、Validate の実装
- `config_test.go`: ユニットテスト

## テスト事項 ✅

### 実装済みテスト
- ✅ 正しい設定ファイル読み込み
- ✅ ファイル不在エラー
- ✅ 無効なJSON パース エラー
- ✅ 無効な更新間隔バリデーション
- ✅ 無効なロケーション情報バリデーション
- ✅ GetRefreshInterval メソッド実装テスト
- ✅ GetLocationString メソッド実装テスト

テスト実行コマンド:
```bash
go test -v ./internal/config
```

## 使用例（main.go での統合）

```go
// 設定読み込み
cfg, err := config.LoadConfig("./data/settings.json")
if err != nil {
    log.Fatal(err)
}

// Gin ミドルウェアで設定をコンテキストに保存
router.Use(func(ctx *gin.Context) {
    ctx.Set("config", cfg)
    ctx.Next()
})

// ハンドラー内で設定を取得
func MyHandler(ctx *gin.Context) {
    cfg, _ := ctx.Get("config")
    appConfig := cfg.(*config.Config)
    interval := appConfig.GetRefreshInterval("weather")
}
```

## 今後の改善 🚀

- [ ] 設定値の変更API（PUT /api/config）
- [ ] 環境変数オーバーライド（例：WEATHER_SEC=600）
- [ ] ホットリロード（ファイル変更時に自動再読み込み）
- [ ] 設定ファイルの暗号化（OAuthトークンなど機密情報向け）
