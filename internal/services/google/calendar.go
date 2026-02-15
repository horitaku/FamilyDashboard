package google

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/rihow/FamilyDashboard/internal/models"
)

// GetCalendarEvents は共有カレンダーからイベントを取得するのです。
// 次 7 日分のイベントを対象とし、allDay と timed で分類して返すます。
// キャッシュが有効な場合はキャッシュから返し、無い場合はGoogle Calendar APIから取得するのです。
func (c *Client) GetCalendarEvents(ctx context.Context) (*models.CalendarResponse, error) {
	// キャッシュキーを生成するのです
	cacheKey := "google_calendar_events"

	// キャッシュが有効か確認するのです
	cachedData, hit, err := c.useCache(cacheKey)
	if err == nil && hit {
		// キャッシュがヒットしたので、パースして返すのです
		var resp models.CalendarResponse
		if err := parseJSONResponse(cachedData, &resp); err == nil {
			return &resp, nil
		}
		// パースエラーの場合は、APIから新たに取得するのです
	}

	// トークンが無い場合は、ダミーデータを返すます（開発用）
	if !c.IsTokenValid() {
		fmt.Println("⚠️ Google Calendar APIトークンが無いため、ダミーイベントを返すのです")
		dummyResp := c.generateDummyCalendarEvents()
		// ダミーデータもキャッシュに保存するのです
		if data, err := models.ToJSON(dummyResp); err == nil {
			_ = c.saveCache(cacheKey, data)
		}
		return dummyResp, nil
	}

	// Google Calendar APIから取得するのです
	// API URL例: https://www.googleapis.com/calendar/v3/calendars/{calendarId}/events
	// 本実装では、settings.jsonの calendarId を使用するます（設定がない場合は primary）
	calendarID := "primary"
	if c.config.Google.CalendarID != "" {
		calendarID = c.config.Google.CalendarID
	}

	// 次 7 日分の timeMin, timeMax を構築するのです
	now := time.Now()
	timeMin := now.Format(time.RFC3339)
	timeMax := now.AddDate(0, 0, 7).Format(time.RFC3339)

	url := fmt.Sprintf(
		"https://www.googleapis.com/calendar/v3/calendars/%s/events?timeMin=%s&timeMax=%s&maxResults=250&orderBy=startTime&singleEvents=true",
		calendarID, timeMin, timeMax,
	)

	// HTTPリクエストを実行するのです
	body, err := c.doRequest(ctx, "GET", url)
	if err != nil {
		// APIエラーの場合は、キャッシュがあれば使用するのです
		if cachedData != nil {
			var resp models.CalendarResponse
			if err := parseJSONResponse(cachedData, &resp); err == nil {
				fmt.Printf("⚠️ Google Calendar API error, cache を使用するのです: %v\n", err)
				return &resp, nil
			}
		}
		return nil, fmt.Errorf("Google Calendar API error: %w", err)
	}

	// Google Calendar APIのレスポンス形式をパースするのです
	var gcalResp GoogleCalendarResponse
	if err := parseJSONResponse(body, &gcalResp); err != nil {
		return nil, fmt.Errorf("calendar response parse error: %w", err)
	}

	// レスポンスを models.CalendarResponse に変換するのです
	resp, err := c.convertCalendarResponse(gcalResp, "Asia/Tokyo")
	if err != nil {
		return nil, fmt.Errorf("calendar conversion error: %w", err)
	}

	// デフォルトでスケジュール内でソートするのです（時系列順）
	c.sortCalendarEvents(resp)

	// キャッシュに保存するのです
	if data, err := models.ToJSON(resp); err == nil {
		_ = c.saveCache(cacheKey, data)
	}

	return resp, nil
}

// GoogleCalendarResponse は Google Calendar API のレスポンス形式です。
type GoogleCalendarResponse struct {
	Items []GoogleCalendarEvent `json:"items"`
}

// GoogleCalendarEvent は Google Calendar API のイベント形式です。
type GoogleCalendarEvent struct {
	ID      string `json:"id"`
	Summary string `json:"summary"`
	Start   struct {
		DateTime string `json:"dateTime"`
		Date     string `json:"date"`
	} `json:"start"`
	End struct {
		DateTime string `json:"dateTime"`
		Date     string `json:"date"`
	} `json:"end"`
	Description string `json:"description"`
	ColorID     string `json:"colorId"`
	EventColor  string `json:"eventColor"`
}

// convertCalendarResponse は GoogleCalendarResponse を models.CalendarResponse に変換するのです。
func (c *Client) convertCalendarResponse(gcalResp GoogleCalendarResponse, tz string) (*models.CalendarResponse, error) {
	location, err := time.LoadLocation(tz)
	if err != nil {
		return nil, fmt.Errorf("timezone load error: %w", err)
	}

	// 日ごとのイベントをマップに集約するのです
	dayMap := make(map[string]*models.CalendarDay)

	for _, gcalEvent := range gcalResp.Items {
		// 開始時刻を解析するのです
		startTime, startIsAllDay, err := c.parseGoogleDateTime(gcalEvent.Start.DateTime, gcalEvent.Start.Date, location)
		if err != nil {
			continue // パースエラーはスキップするのです
		}

		// 終了時刻を解析するのです
		endTime, endIsAllDay, err := c.parseGoogleDateTime(gcalEvent.End.DateTime, gcalEvent.End.Date, location)
		if err != nil {
			continue
		}

		// イベント情報を構築するのです
		event := models.Event{
			ID:       gcalEvent.ID,
			Title:    gcalEvent.Summary,
			Start:    startTime.Format(time.RFC3339),
			End:      endTime.Format(time.RFC3339),
			Color:    c.getEventColor(gcalEvent.ColorID, gcalEvent.EventColor),
			Calendar: "shared", // 共有カレンダー名（設定で変更可能）
			Desc:     gcalEvent.Description,
		}

		// 終日イベントと時間帯付きイベントを分類するのです
		isAllDay := startIsAllDay && endIsAllDay
		dateStr := startTime.Format("2006-01-02")

		if _, ok := dayMap[dateStr]; !ok {
			dayMap[dateStr] = &models.CalendarDay{
				Date:   dateStr,
				AllDay: []models.Event{},
				Timed:  []models.Event{},
			}
		}

		if isAllDay {
			dayMap[dateStr].AllDay = append(dayMap[dateStr].AllDay, event)
		} else {
			dayMap[dateStr].Timed = append(dayMap[dateStr].Timed, event)
		}
	}

	// dayMap を日付順にソートして配列に変換するのです
	days := make([]models.CalendarDay, 0)
	for _, day := range dayMap {
		days = append(days, *day)
	}
	sort.Slice(days, func(i, j int) bool {
		return days[i].Date < days[j].Date
	})

	return &models.CalendarResponse{Days: days}, nil
}

// parseGoogleDateTime は Google Calendar のDateTime/Date 形式をパースするのです。
func (c *Client) parseGoogleDateTime(dateTime, date string, location *time.Location) (time.Time, bool, error) {
	// dateTime（RFC3339）形式の場合
	if dateTime != "" {
		t, err := time.Parse(time.RFC3339, dateTime)
		if err != nil {
			return time.Time{}, false, fmt.Errorf("dateTime parse error: %w", err)
		}
		// タイムゾーンを指定の値に変換するのです
		return t.In(location), false, nil
	}

	// date（YYYY-MM-DD）形式の場合は終日イベント
	if date != "" {
		t, err := time.Parse("2006-01-02", date)
		if err != nil {
			return time.Time{}, false, fmt.Errorf("date parse error: %w", err)
		}
		return t.In(location), true, nil
	}

	return time.Time{}, false, fmt.Errorf("neither dateTime nor date provided")
}

// getEventColor は Google Calendar の色IDを16進カラーコードに変換するのです。
// プレースホルダーで、基本的な色マッピングを提供するのです。
// 本実装では、Google Calendar Colors API から色情報を取得することも可能なのです。
func (c *Client) getEventColor(colorID, eventColor string) string {
	// Google Calendar の色IDマッピング（参考: https://developers.google.com/calendar/v3/colors）
	colorMap := map[string]string{
		"1":  "#A4BDFC", // 薄青
		"2":  "#7AE7BF", // 薄緑
		"3":  "#DBADFF", // 薄紫
		"4":  "#FF887C", // 薄赤
		"5":  "#FBE983", // 薄黄
		"6":  "#FAA775", // 薄橙
		"7":  "#D9D9D9", // グレー
		"8":  "#A4BDFC", // 紺
		"9":  "#46D6DB", // 水色
		"10": "#E67C73", // 赤
		"11": "#33B679", // 緑
	}

	if color, ok := colorMap[colorID]; ok {
		return color
	}

	// eventColor が指定されている場合はそれを使用するのです
	if eventColor != "" {
		return eventColor
	}

	// デフォルトは薄青
	return "#A4BDFC"
}

// sortCalendarEvents は CalendarResponse 内のイベントをソートするのです。
// 各日付内で、終日イベントを上部に配置し、時間帯付きイベントは時系列順にするのです。
func (c *Client) sortCalendarEvents(resp *models.CalendarResponse) {
	for i := range resp.Days {
		day := &resp.Days[i]

		// 終日イベントをタイトル順にソートするのです
		sort.Slice(day.AllDay, func(a, b int) bool {
			return day.AllDay[a].Title < day.AllDay[b].Title
		})

		// 時間帯付きイベントを開始時刻順にソートするのです
		sort.Slice(day.Timed, func(a, b int) bool {
			startA, _ := time.Parse(time.RFC3339, day.Timed[a].Start)
			startB, _ := time.Parse(time.RFC3339, day.Timed[b].Start)
			return startA.Before(startB)
		})
	}
}

// generateDummyCalendarEvents はテスト用のダミーカレンダーイベントを生成するのです。
func (c *Client) generateDummyCalendarEvents() *models.CalendarResponse {
	location, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(location)

	days := make([]models.CalendarDay, 0)
	for i := 0; i < 7; i++ {
		date := now.AddDate(0, 0, i)
		dateStr := date.Format("2006-01-02")

		day := models.CalendarDay{
			Date:   dateStr,
			AllDay: []models.Event{},
			Timed:  []models.Event{},
		}

		// 終日イベント
		if i%2 == 0 {
			day.AllDay = append(day.AllDay, models.Event{
				ID:       fmt.Sprintf("dummy-allday-%d", i),
				Title:    fmt.Sprintf("終日イベント %d", i),
				Start:    dateStr,
				End:      date.AddDate(0, 0, 1).Format("2006-01-02"),
				Color:    "#A4BDFC",
				Calendar: "shared",
			})
		}

		// 時間帯付きイベント
		eventTime := date.Add(14 * time.Hour).Format(time.RFC3339)
		eventEndTime := date.Add(16 * time.Hour).Format(time.RFC3339)
		day.Timed = append(day.Timed, models.Event{
			ID:       fmt.Sprintf("dummy-timed-%d", i),
			Title:    fmt.Sprintf("ミーティング %d", i),
			Start:    eventTime,
			End:      eventEndTime,
			Color:    "#FF887C",
			Calendar: "shared",
		})

		days = append(days, day)
	}

	return &models.CalendarResponse{Days: days}
}
