package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rihow/FamilyDashboard/internal/config"
	"github.com/rihow/FamilyDashboard/internal/models"
	"github.com/rihow/FamilyDashboard/internal/services/weather"
)

// ============================================================================
// /api/status ハンドラー
// ============================================================================

// GetStatus は /api/status のGETハンドラーなのです。
// 現在の状態・エラー・各ソースの最終更新時刻を返すもなのです。
func GetStatus(ctx *gin.Context) {
	// 現在時刻を取得する（Asia/Tokyo）
	now := time.Now().UTC().Format(time.RFC3339)

	// ダミーのエラーリスト（いまは空）
	errors := []models.ErrorInfo{}

	// ダミーの最終更新時刻
	lastUpdated := models.LastUpdatedTimes{
		Weather:  time.Now().UTC().Add(-5 * time.Minute).Format(time.RFC3339),
		Calendar: time.Now().UTC().Add(-3 * time.Minute).Format(time.RFC3339),
		Tasks:    time.Now().UTC().Add(-2 * time.Minute).Format(time.RFC3339),
	}

	response := models.StatusResponse{
		OK:          true,
		Now:         now,
		Errors:      errors,
		LastUpdated: lastUpdated,
	}

	ctx.JSON(http.StatusOK, response)
}

// ============================================================================
// /api/calendar ハンドラー
// ============================================================================

// GetCalendar は /api/calendar のGETハンドラーなのです。
// 最大7日分のイベントを返すもなのです（いまはダミーデータ）。
func GetCalendar(ctx *gin.Context) {
	// ダミーイベント
	dummyEvent := models.Event{
		ID:       "dummy-event-1",
		Title:    "チームミーティング",
		Start:    time.Now().UTC().Format(time.RFC3339),
		End:      time.Now().UTC().Add(1 * time.Hour).Format(time.RFC3339),
		Color:    "#4285F4",
		Calendar: "仕事",
		Desc:     "定期ミーティング",
	}

	dummyAllDayEvent := models.Event{
		ID:       "dummy-event-2",
		Title:    "家族のおでかけ",
		Start:    time.Now().UTC().Format("2006-01-02"),
		End:      time.Now().UTC().Format("2006-01-02"),
		Color:    "#EA4335",
		Calendar: "家族",
		Desc:     "全日予定",
	}

	today := time.Now().UTC().Format("2006-01-02")

	response := models.CalendarResponse{
		Days: []models.CalendarDay{
			{
				Date:   today,
				AllDay: []models.Event{dummyAllDayEvent},
				Timed:  []models.Event{dummyEvent},
			},
			{
				Date:   time.Now().UTC().AddDate(0, 0, 1).Format("2006-01-02"),
				AllDay: []models.Event{},
				Timed:  []models.Event{},
			},
		},
	}

	ctx.JSON(http.StatusOK, response)
}

// ============================================================================
// /api/tasks ハンドラー
// ============================================================================

// GetTasks は /api/tasks のGETハンドラーなのです。
// サーバー側ソート済みのタスクリストを返すもなのです（いまはダミーデータ）。
func GetTasks(ctx *gin.Context) {
	// ダミーのタスク
	today := time.Now().UTC().Format("2006-01-02")
	tomorrow := time.Now().UTC().AddDate(0, 0, 1).Format("2006-01-02")

	task1 := models.TaskItem{
		ID:        "task-1",
		Title:     "シュッポとミサイルで遊ぶ",
		Notes:     "2時間ほど遊ぶます",
		Status:    "needsAction",
		DueDate:   &today,
		Priority:  1, // 最高優先度
		CreatedAt: time.Now().UTC().Add(-24 * time.Hour),
	}

	task2 := models.TaskItem{
		ID:        "task-2",
		Title:     "スパイ任務をこなす",
		Notes:     "情報収集タイム",
		Status:    "needsAction",
		DueDate:   &tomorrow,
		Priority:  2,
		CreatedAt: time.Now().UTC().Add(-12 * time.Hour),
	}

	task3 := models.TaskItem{
		ID:        "task-3",
		Title:     "テレパシーの練習",
		Notes:     "毎日する~",
		Status:    "needsAction",
		DueDate:   nil, // 期限なし
		Priority:  3,
		CreatedAt: time.Now().UTC().Add(-1 * time.Hour),
	}

	response := models.TasksResponse{
		Items: []models.TaskItem{task1, task2, task3},
	}

	ctx.JSON(http.StatusOK, response)
}

// ============================================================================
// /api/weather ハンドラー
// ============================================================================

// GetWeather は /api/weather のGETハンドラーなのです。
// 現在の天候・今日の気温・降水確率・警報を返すもなのです。
// 設定から都市名を取得して、weather クライアントで Open-Meteo API から
// 気象庁データを含む最新の天気情報を取得するます。
func GetWeather(ctx *gin.Context) {
	// コンテキストから設定と weather クライアントを取得するます
	cfgRaw, exists := ctx.Get("config")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "設定が見つからないです",
		})
		return
	}
	cfg := cfgRaw.(*config.Config)

	weatherRaw, exists := ctx.Get("weather")
	if !exists {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "天気クライアントが見つかりません",
		})
		return
	}
	weatherClient := weatherRaw.(*weather.Client)

	// 設定から都市名と国を取得するます
	cityName := cfg.Location.CityName
	if cityName == "" {
		cityName = "姫路市" // デフォルト都市
	}
	country := cfg.Location.Country
	if country == "" {
		country = "JP" // デフォルト国コード
	}

	// 天気データを取得するます（キャッシュから または API から）
	weatherRsp, err := weatherClient.GetWeather(ctx, cityName, country)
	if err != nil {
		// エラーが発生した場合、ログに出力してダミーデータを返します
		fmt.Printf("❌ 天気データ取得エラー: %v\n", err)
		ctx.JSON(http.StatusOK, &models.WeatherResponse{
			Location: cityName,
			Current: models.CurrentWeather{
				Temperature: 0,
				Condition:   "データ取得失敗",
				Icon:        "04u",
				Humidity:    0,
				WindSpeed:   0,
			},
			Today: models.TodayWeather{
				MaxTemp: 0,
				MinTemp: 0,
				Summary: "データ取得失敗",
			},
			PrecipSlots: []models.PrecipSlot{},
			Alerts:      []models.WeatherAlert{},
		})
		return
	}

	ctx.JSON(http.StatusOK, weatherRsp)
}
