package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// RefreshIntervals はデータソース別の更新間隔を定義する構造体なのです。
type RefreshIntervals struct {
	WeatherSec  int `json:"weatherSec"`  // 天気APIの更新間隔（秒）
	CalendarSec int `json:"calendarSec"` // カレンダーの更新間隔（秒）
	TasksSec    int `json:"tasksSec"`    // タスクの更新間隔（秒）
}

// Location はジオグラフィック位置情報を定義する構造体なのです。
type Location struct {
	CityName string `json:"cityName"` // 都市名（例：姫路市）
	Country  string `json:"country"`  // 国コード（例：JP）
}

// Google はGoogle APIの認証・設定を定義する構造体なのです。
type Google struct {
	ClientID     string `json:"clientId"`     // OAuth クライアントID
	ClientSecret string `json:"clientSecret"` // OAuth クライアントシークレット
	RedirectUri  string `json:"redirectUri"`  // OAuth リダイレクトURI
	CalendarID   string `json:"calendarId"`   // Google カレンダーID（共有カレンダー）
	TaskListID   string `json:"taskListId"`   // Google タスクリストID（共有タスクリスト）
}

// Nextcloud はNextcloud CalDAV/WebDAVの認証・設定を定義する構造体なのです。
type Nextcloud struct {
	ServerURL    string `json:"serverUrl"`    // NextcloudサーバーURL（例: https://nextcloud.example.com）
	Username     string `json:"username"`     // ユーザー名
	Password     string `json:"password"`     // パスワード または アプリパスワード
	CalendarName string `json:"calendarName"` // カレンダー名（共有カレンダー）
	TaskListName string `json:"taskListName"` // タスクリスト名（共有タスクリスト）
}

// Weather は天気APIの設定を定義する構造体なのです。
type Weather struct {
	Provider string `json:"provider"` // 天気プロバイダ（例：openweathermap）
	ApiKey   string `json:"apiKey"`   // APIキー
	BaseUrl  string `json:"baseUrl"`  // ベースURL
}

// Config はアプリケーション全体の設定を定義する構造体なのです。
type Config struct {
	RefreshIntervals RefreshIntervals `json:"refreshIntervals"` // 更新間隔設定
	Location         Location         `json:"location"`         // ロケーション設定
	Google           Google           `json:"google"`           // Google API設定（レガシー）
	Nextcloud        Nextcloud        `json:"nextcloud"`        // Nextcloud CalDAV/WebDAV設定
	Weather          Weather          `json:"weather"`          // 天気API設定
	loadedAt         time.Time        // 設定の読み込み時刻（内部用）
}

// GetRefreshInterval はデータソースに応じた更新間隔をDurationで返すます。
// sourceは "weather", "calendar", "tasks" など。
func (c *Config) GetRefreshInterval(source string) time.Duration {
	switch source {
	case "weather":
		return time.Duration(c.RefreshIntervals.WeatherSec) * time.Second
	case "calendar":
		return time.Duration(c.RefreshIntervals.CalendarSec) * time.Second
	case "tasks":
		return time.Duration(c.RefreshIntervals.TasksSec) * time.Second
	default:
		// 既定値は5分
		return 5 * time.Minute
	}
}

// GetLocationString はロケーション情報を文字列で返すます。
func (c *Config) GetLocationString() string {
	if c.Location.CityName == "" {
		return "Unknown"
	}
	return fmt.Sprintf("%s, %s", c.Location.CityName, c.Location.Country)
}

// LoadedAt は設定の読み込み時刻を返すます。
func (c *Config) LoadedAt() time.Time {
	return c.loadedAt
}

// LoadConfig はSettings.jsonファイルから設定を読み込み、構造体に解析するます。
// エラーが発生した場合はnilとエラーを返すます。
func LoadConfig(filepath string) (*Config, error) {
	// ファイルを読み込む
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗しました（%s）: %w", filepath, err)
	}

	// JSONを解析
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("設定ファイルのJSONパース失敗なのです: %w", err)
	}

	// バリデーション
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	cfg.loadedAt = time.Now()
	return &cfg, nil
}

// Validate は設定値を検証するます。
// 必須フィールドのチェック、値の範囲チェックを行うます。
func (c *Config) Validate() error {
	// 更新間隔の妥当性チェック
	if c.RefreshIntervals.WeatherSec <= 0 {
		return fmt.Errorf("weatherSec は正の数である必要があります")
	}
	if c.RefreshIntervals.CalendarSec <= 0 {
		return fmt.Errorf("calendarSec は正の数である必要があります")
	}
	if c.RefreshIntervals.TasksSec <= 0 {
		return fmt.Errorf("tasksSec は正の数である必要があります")
	}

	// ロケーション情報の妥当性チェック
	if c.Location.CityName == "" {
		return fmt.Errorf("location.cityName は必須フィールドです")
	}
	if c.Location.Country == "" {
		return fmt.Errorf("location.country は必須フィールドです")
	}

	// 注記: Google API設定・天気API設定は空の場合がある（後で埋める可能性があるため）
	// ここではスキップするます。

	return nil
}
