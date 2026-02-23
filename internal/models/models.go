package models

import (
	"encoding/json"
	"time"
)

// ============================================================================
// APIレスポンス共通構造体
// ============================================================================

// StatusResponse は /api/status のレスポンスです。
// ダッシュボード全体の状態・エラー・最終更新時刻を返すもなのです。
type StatusResponse struct {
	OK          bool             `json:"ok"`          // サーバー稼働状態
	Now         string           `json:"now"`         // 現在時刻（RFC3339）
	Errors      []ErrorInfo      `json:"errors"`      // エラーリスト
	LastUpdated LastUpdatedTimes `json:"lastUpdated"` // 各ソースの最終更新時刻
}

// ErrorInfo はエラー情報を表すのです。
type ErrorInfo struct {
	Source  string `json:"source"`  // エラー源（"weather", "calendar", "tasks" など）
	Message string `json:"message"` // エラーメッセージ
	At      string `json:"at"`      // エラー発生時刻（RFC3339）
}

// LastUpdatedTimes は各データソースの最終更新時刻なのです。
type LastUpdatedTimes struct {
	Weather  string `json:"weather"`  // 天気の最終更新時刻（RFC3339）
	Calendar string `json:"calendar"` // カレンダーの最終更新時刻（RFC3339）
	Tasks    string `json:"tasks"`    // タスクの最終更新時刻（RFC3339）
}

// ============================================================================
// カレンダー関連の構造体
// ============================================================================

// CalendarResponse は /api/calendar のレスポンスなのです。
type CalendarResponse struct {
	Days []CalendarDay `json:"days"` // 日ごとのイベントリスト
}

// CalendarDay は1日分のイベント情報を表すのです。
type CalendarDay struct {
	Date   string  `json:"date"`   // 日付（YYYY-MM-DD）
	AllDay []Event `json:"allDay"` // 終日イベント
	Timed  []Event `json:"timed"`  // 時間帯付きイベント
}

// Event はカレンダーのイベント情報なのです。
type Event struct {
	ID       string `json:"id"`          // Google イベントID
	Title    string `json:"title"`       // イベント名
	Start    string `json:"start"`       // 開始時刻（RFC3339 または YYYY-MM-DD）
	End      string `json:"end"`         // 終了時刻（RFC3339 または YYYY-MM-DD）
	Color    string `json:"color"`       // 色コード（#RRGGBB など）
	Calendar string `json:"calendar"`    // カレンダー名
	Desc     string `json:"description"` // 説明（省略可）
}

// ============================================================================
// タスク関連の構造体
// ============================================================================

// TasksResponse は /api/tasks のレスポンスなのです。
type TasksResponse struct {
	Items []TaskItem `json:"items"` // タスクリスト（サーバー側ソート済）
}

// TaskItem はタスク1件の情報なのです。
type TaskItem struct {
	ID        string    `json:"id"`        // Google タスクID
	Title     string    `json:"title"`     // タスク名
	Notes     string    `json:"notes"`     // 説明
	Status    string    `json:"status"`    // "needsAction" か "completed"
	DueDate   *string   `json:"dueDate"`   // 期限（ISO 8601 形式 YYYY-MM-DD、null 可能）
	Priority  int       `json:"priority"`  // 優先度（1-3、1が最高）
	CreatedAt time.Time `json:"createdAt"` // 作成日時
}

// ============================================================================
// 天気関連の構造体
// ============================================================================

// WeatherResponse は /api/weather のレスポンスなのです。
type WeatherResponse struct {
	Location    string         `json:"location"`    // 場所（都市名など）
	Current     CurrentWeather `json:"current"`     // 現在の天候
	Today       TodayWeather   `json:"today"`       // 今日の天況
	PrecipSlots []PrecipSlot   `json:"precipSlots"` // 時間帯ごとの降水確率
	Alerts      []WeatherAlert `json:"alerts"`      // 注意報・警報
}

// CurrentWeather は現在の天況なのです。
type CurrentWeather struct {
	Temperature float64 `json:"temperature"` // 気温（℃）
	Condition   string  `json:"condition"`   // 天候（"晴" "曇" "雨" など）
	Icon        string  `json:"icon"`        // 天候アイコンコード
	Humidity    int     `json:"humidity"`    // 湿度（%）
	WindSpeed   float64 `json:"windSpeed"`   // 風速（m/s）
}

// TodayWeather は今日の天況なのです。
type TodayWeather struct {
	MaxTemp float64 `json:"maxTemp"` // 最高気温（℃）
	MinTemp float64 `json:"minTemp"` // 最低気温（℃）
	Summary string  `json:"summary"` // 概況
}

// PrecipSlot は時間帯ごとの降水確率なのです。
type PrecipSlot struct {
	Time   string `json:"time"`   // 時刻（HH:00 形式）
	Precip int    `json:"precip"` // 降水確率（%）
}

// WeatherAlert は注意報・警報なのです。
type WeatherAlert struct {
	Title    string `json:"title"`       // 警報名（e.g., "大雨警報"）
	Headline string `json:"headline"`    // 見出し
	Desc     string `json:"description"` // 詳細
	Severity string `json:"severity"`    // 重大度（"注意報" "警報" "特別警報"）
}

// ============================================================================
// ユーティリティ関数
// ============================================================================

// ToJSON は任意の構造体をJSON バイト列に変換するのです。
// キャッシュへの保存や、APIレスポンスの詳細な処理に使用するのです。
func ToJSON(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
