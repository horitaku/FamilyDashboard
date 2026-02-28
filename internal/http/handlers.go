package http

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rihow/FamilyDashboard/internal/cache"
	"github.com/rihow/FamilyDashboard/internal/config"
	"github.com/rihow/FamilyDashboard/internal/models"
	"github.com/rihow/FamilyDashboard/internal/services/nextcloud"
	"github.com/rihow/FamilyDashboard/internal/services/weather"
	"github.com/rihow/FamilyDashboard/internal/status"
)

// ============================================================================
// /api/status ハンドラー
// ============================================================================

// GetStatus は /api/status のGETハンドラーなのです。
// 現在の状態・エラー・各ソースの最終更新時刻を返すもなのです。
func GetStatus(ctx *gin.Context) {
	// 現在時刻を取得する（Asia/Tokyo）
	now := status.NowRFC3339()

	// エラーリストを取得するのです（記録が無い場合は空）
	errors := []models.ErrorInfo{}
	if store := getErrorStore(ctx); store != nil {
		errors = store.List()
	}

	// キャッシュの最終更新時刻を集計するのです
	lastUpdated := models.LastUpdatedTimes{}
	if fc := getCache(ctx); fc != nil {
		cfg := getConfig(ctx)
		cityName := "姫路市"
		country := "JP"
		if cfg != nil {
			if cfg.Location.CityName != "" {
				cityName = cfg.Location.CityName
			}
			if cfg.Location.Country != "" {
				country = cfg.Location.Country
			}
		}

		lastUpdated.Weather = readFetchedAt(fc, fmt.Sprintf("weather:%s:%s", country, cityName))
		lastUpdated.Calendar = readFetchedAt(fc, "nextcloud_calendar_events")
		lastUpdated.Tasks = readFetchedAt(fc, "nextcloud_tasks_items")
	}

	response := models.StatusResponse{
		OK:          len(errors) == 0,
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
// Nextcloud CalDAV からイベントを取得し、最大7日分を返すます。
// クライアントが無い場合はダミーデータを返すのです。
func GetCalendar(ctx *gin.Context) {
	// コンテキストから Nextcloud クライアントと設定を取得するます
	nextcloudRaw, exists := ctx.Get("nextcloud")
	if !exists {
		// Nextcloud クライアントが無い場合はダミーデータを返す
		fmt.Println("⚠️ Nextcloud クライアントが見つかりません。ダミーデータを返すのです")
		dummyResp := &models.CalendarResponse{
			Days: []models.CalendarDay{
				{
					Date:   time.Now().Format("2006-01-02"),
					AllDay: []models.Event{},
					Timed:  []models.Event{},
				},
			},
		}
		ctx.JSON(http.StatusOK, dummyResp)
		return
	}
	nextcloudClient := nextcloudRaw.(*nextcloud.Client)

	// Nextcloud CalDAV からイベントを取得するます
	calendarResp, err := nextcloudClient.GetCalendarEvents(ctx)
	if err != nil {
		fmt.Printf("❌ カレンダーデータ取得エラー: %v\n", err)
		setSourceError(ctx, "calendar", err)
		if calendarResp != nil {
			ctx.JSON(http.StatusOK, calendarResp)
			return
		}
		ctx.JSON(http.StatusOK, &models.CalendarResponse{
			Days: []models.CalendarDay{},
		})
		return
	}

	clearSourceError(ctx, "calendar")
	ctx.JSON(http.StatusOK, calendarResp)
}

// ============================================================================
// /api/tasks ハンドラー
// ============================================================================

// GetTasks は /api/tasks のGETハンドラーなのです。
// Nextcloud WebDAV からタスクを取得し、サーバー側ソート済みのタスクリストを返すます。
// クライアントが無い場合はダミーデータを返すのです。
func GetTasks(ctx *gin.Context) {
	// コンテキストから Nextcloud クライアントを取得するます
	nextcloudRaw, exists := ctx.Get("nextcloud")
	if !exists {
		// Nextcloud クライアントが無い場合はダミーデータを返す
		fmt.Println("⚠️ Nextcloud クライアントが見つかりません。ダミーデータを返すのです")
		dummyResp := &models.TasksResponse{
			Items: []models.TaskItem{},
		}
		ctx.JSON(http.StatusOK, dummyResp)
		return
	}
	nextcloudClient := nextcloudRaw.(*nextcloud.Client)

	// Nextcloud WebDAV からタスクを取得するます
	tasksResp, err := nextcloudClient.GetTaskItems(ctx)
	if err != nil {
		fmt.Printf("❌ タスクデータ取得エラー: %v\n", err)
		setSourceError(ctx, "tasks", err)
		if tasksResp != nil {
			ctx.JSON(http.StatusOK, tasksResp)
			return
		}
		ctx.JSON(http.StatusOK, &models.TasksResponse{
			Items: []models.TaskItem{},
		})
		return
	}

	clearSourceError(ctx, "tasks")
	ctx.JSON(http.StatusOK, tasksResp)
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
		// エラーが発生した場合、ログに出力してキャッシュを優先するます
		fmt.Printf("❌ 天気データ取得エラー: %v\n", err)
		setSourceError(ctx, "weather", err)
		if weatherRsp != nil {
			ctx.JSON(http.StatusOK, weatherRsp)
			return
		}
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

	clearSourceError(ctx, "weather")
	ctx.JSON(http.StatusOK, weatherRsp)
}

func getCache(ctx *gin.Context) *cache.FileCache {
	cacheRaw, exists := ctx.Get("cache")
	if !exists {
		return nil
	}
	fc, ok := cacheRaw.(*cache.FileCache)
	if !ok {
		return nil
	}
	return fc
}

func getConfig(ctx *gin.Context) *config.Config {
	cfgRaw, exists := ctx.Get("config")
	if !exists {
		return nil
	}
	cfg, ok := cfgRaw.(*config.Config)
	if !ok {
		return nil
	}
	return cfg
}

func getErrorStore(ctx *gin.Context) *status.ErrorStore {
	storeRaw, exists := ctx.Get("errorStore")
	if !exists {
		return nil
	}
	store, ok := storeRaw.(*status.ErrorStore)
	if !ok {
		return nil
	}
	return store
}

func setSourceError(ctx *gin.Context, source string, err error) {
	if err == nil {
		return
	}
	if store := getErrorStore(ctx); store != nil {
		store.Set(source, err.Error())
	}
}

func clearSourceError(ctx *gin.Context, source string) {
	if store := getErrorStore(ctx); store != nil {
		store.Clear(source)
	}
}

func readFetchedAt(fc *cache.FileCache, cacheKey string) string {
	entry, exists, _, err := fc.Read(cacheKey, 0)
	if err != nil || !exists {
		return ""
	}
	return entry.FetchedAt
}

// ============================================================================
// /auth/login ハンドラー（Google OAuth - 将来削除予定）
// ============================================================================

/*
// AuthLogin は Google OAuth ログインへのリダイレクトリンクを生成するのです。
// ブラウザこのエンドポイントにアクセスするとGoogle ログイン画面に遷移するます。
// NOTE: NextcloudではBasic認証を使うため、このハンドラーは不要になるます（Step16で削除予定）
func AuthLogin(ctx *gin.Context) {
	cfg := ctx.MustGet("config").(*config.Config)

	if cfg.Google.ClientID == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Google OAuth設定が未設定のです",
		})
		return
	}

	// OAuth 認可画面へのURL生成
	// スコープ: Google Calendar と Google Tasks の読み取り権限
	scopes := "https://www.googleapis.com/auth/calendar.readonly https://www.googleapis.com/auth/tasks.readonly"

	authURL := fmt.Sprintf(
		"https://accounts.google.com/o/oauth2/v2/auth?"+
			"client_id=%s&"+
			"redirect_uri=%s&"+
			"response_type=code&"+
			"scope=%s&"+
			"access_type=offline",
		cfg.Google.ClientID,
		cfg.Google.RedirectUri,
		scopes,
	)

	// リダイレクト
	ctx.Redirect(http.StatusTemporaryRedirect, authURL)
}

// ============================================================================
// /auth/callback ハンドラー（Google OAuth - 将来削除予定）
// ============================================================================

// AuthCallback は Google OAuth のコールバックハンドラーなのです。
// クエリから "code" パラメータを受け取り、トークン取得を実行するます。
// NOTE: NextcloudではBasic認証を使うため、このハンドラーは不要になるます（Step16で削除予定）
func AuthCallback(ctx *gin.Context) {
	// ユーザーが認可をキャンセルした場合
	if err := ctx.Query("error"); err != "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": fmt.Sprintf("OAuth キャンセル: %s", err),
		})
		return
	}

	// 認可コードを取得
	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "認可コードが未設定のです",
		})
		return
	}

	// Google クライアントを取得
	googleClient := ctx.MustGet("google").(*google.Client)

	// OAuth認可コードフロー実行（トークン取得）
	tokenResp, err := googleClient.OAuthAuthorizationCodeFlow(ctx, code)
	if err != nil {
		fmt.Printf("❌ OAuthエラー: %v\n", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": fmt.Sprintf("OAuth フロー失敗: %v", err),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":       "認可成功したます！✨",
		"access_token":  tokenResp.AccessToken[:20] + "...",
		"expires_in":    tokenResp.ExpiresIn,
		"refresh_token": "保存済み",
	})
}
*/
