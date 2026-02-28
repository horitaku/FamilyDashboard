package nextcloud

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
	"github.com/rihow/FamilyDashboard/internal/models"
)

// GetCalendarEvents はNextcloud CalDAVからカレンダーイベントを取得するます。
// 複数のカレンダーから今日から7日分のイベントを取得し、終日/時間帯別に分類して返すのです。
func (c *Client) GetCalendarEvents(ctx context.Context) (*models.CalendarResponse, error) {
	cacheKey := "nextcloud_calendar_events_all"
	ttl := c.config.GetRefreshInterval("calendar")

	// キャッシュを確認するます
	entry, ok, stale, err := c.cache.Read(cacheKey, ttl)
	if ok && !stale && err == nil {
		fmt.Println("📦 カレンダーキャッシュヒット!")
		var resp models.CalendarResponse
		if err := json.Unmarshal(entry.Payload, &resp); err == nil {
			return &resp, nil
		}
		fmt.Printf("⚠️ キャッシュデータのパース失敗: %v\n", err)
	}

	// 複数カレンダー名を取得するます
	calendarNames := c.config.GetCalendarNames()
	if len(calendarNames) == 0 {
		return nil, fmt.Errorf("カレンダー名が設定されていません")
	}

	fmt.Printf("🌐 Nextcloud CalDAV から %d 個のカレンダーを取得するます...\n", len(calendarNames))

	// 今日から7日分の範囲を設定（Asia/Tokyo）
	loc, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(loc)
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endDate := startDate.AddDate(0, 0, 7)

	// 全カレンダーからイベントを収集するます
	allEvents := []eventWithDate{}
	var fetchErrors []error

	for _, calendarName := range calendarNames {
		fmt.Printf("  📅 カレンダー '%s' からイベント取得中...\n", calendarName)

		// CalDAVクエリを実行するます
		calendarPath := c.getCalendarPath(calendarName)
		calendarColor, colorErr := c.getCalendarColor(ctx, calendarPath)
		if colorErr != nil {
			fmt.Printf("⚠️ カレンダー '%s' の色取得失敗: %v\n", calendarName, colorErr)
		}

		query := &caldav.CalendarQuery{
			CompRequest: caldav.CalendarCompRequest{
				Name: "VCALENDAR",
				Comps: []caldav.CalendarCompRequest{
					{
						Name:  "VEVENT",
						Props: []string{"UID", "SUMMARY", "DTSTART", "DTEND", "DESCRIPTION", "COLOR"},
					},
				},
			},
			CompFilter: caldav.CompFilter{
				Name: "VCALENDAR",
				Comps: []caldav.CompFilter{
					{
						Name: "VEVENT",
						Start: startDate,
						End:   endDate,
					},
				},
			},
		}

		calendarObjects, err := c.caldavClient.QueryCalendar(ctx, calendarPath, query)
		if err != nil {
			// エラーを記録するが続行するます（部分的成功を許容）
			fmt.Printf("❌ カレンダー '%s' のCalDAVクエリエラー: %v\n", calendarName, err)
			fetchErrors = append(fetchErrors, fmt.Errorf("calendar '%s': %w", calendarName, err))
			continue
		}

		// iCalendarオブジェクトをパースして構造化するます
		for _, obj := range calendarObjects {
			parsedEvents := parseCalendarObject(obj.Data, startDate, endDate, calendarName, calendarColor)
			allEvents = append(allEvents, parsedEvents...)
		}

		fmt.Printf("  ✅ カレンダー '%s' から %d 件のイベント取得\n", calendarName, len(calendarObjects))
	}

	// すべてのカレンダー取得に失敗した場合
	if len(allEvents) == 0 && len(fetchErrors) > 0 {
		// エラー時はキャッシュから返す試みをするます
		fmt.Println("❌ すべてのカレンダー取得に失敗しました")
		entry, ok, _, readErr := c.cache.Read(cacheKey, 0)
		if ok && readErr == nil {
			fmt.Println("📦 期限切れキャッシュを返すます")
			var resp models.CalendarResponse
			if unmarshalErr := json.Unmarshal(entry.Payload, &resp); unmarshalErr == nil {
				return &resp, fmt.Errorf("全カレンダー取得失敗（キャッシュ返却）: %d エラー", len(fetchErrors))
			}
		}
		return nil, fmt.Errorf("全カレンダー取得失敗: %d エラー", len(fetchErrors))
	}

	// 日付ごとにイベントを分類するます
	response := convertToCalendarResponse(allEvents, startDate, endDate)

	// キャッシュに保存するます
	meta := map[string]string{"source": "nextcloud_calendar_all"}
	if _, err := c.cache.Write(cacheKey, response, meta); err != nil {
		fmt.Printf("⚠️ キャッシュ保存失敗: %v\n", err)
	}

	fmt.Printf("✅ 統合カレンダーイベント取得成功: %d日分、合計 %d イベント\n", len(response.Days), len(allEvents))
	if len(fetchErrors) > 0 {
		fmt.Printf("⚠️ 一部のカレンダーで取得エラーがありました: %d 件\n", len(fetchErrors))
	}

	return response, nil
}

// getCalendarColor はカレンダーコレクションの色（calendar-color）を取得するます。
// Nextcloud は calendar-color を #RRGGBB または #RRGGBBAA で返すことがあるのです。
func (c *Client) getCalendarColor(ctx context.Context, calendarPath string) (string, error) {
	requestBody := `<?xml version="1.0" encoding="UTF-8"?>
<d:propfind xmlns:d="DAV:" xmlns:a="http://apple.com/ns/ical/">
	<d:prop>
		<a:calendar-color/>
	</d:prop>
</d:propfind>`

	targetURL, err := c.resolveDAVURL(calendarPath)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "PROPFIND", targetURL, bytes.NewBufferString(requestBody))
	if err != nil {
		return "", fmt.Errorf("calendar color request作成失敗: %w", err)
	}
	req.Header.Set("Depth", "0")
	req.Header.Set("Content-Type", "application/xml; charset=utf-8")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("calendar color取得失敗: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusMultiStatus && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("calendar color取得HTTPエラー: %d", resp.StatusCode)
	}

	var result struct {
		Responses []struct {
			PropStats []struct {
				Status string `xml:"status"`
				Prop   struct {
					CalendarColor string `xml:"http://apple.com/ns/ical/ calendar-color"`
				} `xml:"prop"`
			} `xml:"propstat"`
		} `xml:"response"`
	}

	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("calendar colorレスポンス解析失敗: %w", err)
	}

	for _, response := range result.Responses {
		for _, propStat := range response.PropStats {
			if !strings.Contains(propStat.Status, "200") {
				continue
			}

			if normalized, ok := normalizeHexColor(propStat.Prop.CalendarColor); ok {
				return normalized, nil
			}
		}
	}

	return "", nil
}

// resolveDAVURL は base URL と DAVパスを結合するます。
func (c *Client) resolveDAVURL(davPath string) (string, error) {
	base, err := url.Parse(c.config.Nextcloud.ServerURL)
	if err != nil {
		return "", fmt.Errorf("ServerURL解析失敗: %w", err)
	}

	ref, err := url.Parse(davPath)
	if err != nil {
		return "", fmt.Errorf("DAVパス解析失敗: %w", err)
	}

	return base.ResolveReference(ref).String(), nil
}

// eventWithDate はイベントと日付情報を保持する内部構造体なのです。
type eventWithDate struct {
	event  models.Event
	date   time.Time
	allDay bool
}

// parseCalendarObject はiCalendarデータをパースしてイベントリストに変換するます。
func parseCalendarObject(cal *ical.Calendar, startDate, endDate time.Time, calendarName, calendarColor string) []eventWithDate {
	events := []eventWithDate{}

	if cal == nil {
		return events
	}

	loc, _ := time.LoadLocation("Asia/Tokyo")

	for _, comp := range cal.Children {
		if comp.Name != "VEVENT" {
			continue
		}

		// イベント情報を抽出するます
		uid := comp.Props.Get("UID")
		summary := comp.Props.Get("SUMMARY")
		dtStart := comp.Props.Get("DTSTART")
		dtEnd := comp.Props.Get("DTEND")
		description := comp.Props.Get("DESCRIPTION")
		color := comp.Props.Get("COLOR")

		if uid == nil || summary == nil || dtStart == nil {
			continue
		}

		// 開始日時をパースするます
		startTime, isAllDay := parseDateTime(dtStart.Value, loc)
		if startTime.IsZero() {
			continue
		}

		// 期間外のイベントはスキップ
		if startTime.Before(startDate) || startTime.After(endDate) {
			continue
		}

		// 終了日時をパースするます
		endTime := startTime.Add(1 * time.Hour) // デフォルト1時間
		if dtEnd != nil {
			parsedEnd, _ := parseDateTime(dtEnd.Value, loc)
			if !parsedEnd.IsZero() {
				endTime = parsedEnd
			}
		}

		// 色を決定（優先順位: イベント色 > カレンダー色 > デフォルト）
		colorValue := "#3788d8"
		if normalized, ok := normalizeHexColor(calendarColor); ok {
			colorValue = normalized
		}
		if color != nil && color.Value != "" {
			if normalized, ok := normalizeHexColor(color.Value); ok {
				colorValue = normalized
			}
		}

		// Eventオブジェクトを作成
		event := models.Event{
			ID:       uid.Value,
			Title:    summary.Value,
			Start:    startTime.Format(time.RFC3339),
			End:      endTime.Format(time.RFC3339),
			Color:    colorValue,
				Calendar: calendarName,
			Desc:     "",
		}
		if description != nil {
			event.Desc = description.Value
		}

		events = append(events, eventWithDate{
			event:  event,
			date:   startTime,
			allDay: isAllDay,
		})
	}

	return events
}

// normalizeHexColor は #RGB/#RRGGBB/#RRGGBBAA/先頭#なし の色を #RRGGBB に正規化するます。
func normalizeHexColor(value string) (string, bool) {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return "", false
	}

	trimmed = strings.TrimPrefix(trimmed, "#")

	if len(trimmed) == 3 {
		trimmed = strings.Repeat(string(trimmed[0]), 2) +
			strings.Repeat(string(trimmed[1]), 2) +
			strings.Repeat(string(trimmed[2]), 2)
	}

	if len(trimmed) == 8 {
		trimmed = trimmed[:6]
	}

	if len(trimmed) != 6 {
		return "", false
	}

	for _, ch := range trimmed {
		if (ch < '0' || ch > '9') && (ch < 'a' || ch > 'f') && (ch < 'A' || ch > 'F') {
			return "", false
		}
	}

	return "#" + strings.ToUpper(trimmed), true
}

// parseDateTime はiCalendar日時文字列をパースするます。
// YYYYMMDD形式（終日）とYYYYMMDDTHHMMSS形式（時間指定）に対応するのです。
func parseDateTime(value string, loc *time.Location) (time.Time, bool) {
	value = strings.TrimSpace(value)

	// 終日イベント（YYYYMMDD形式）
	if len(value) == 8 {
		t, err := time.ParseInLocation("20060102", value, loc)
		if err == nil {
			return t, true // 終日
		}
	}

	// 時間指定イベント（YYYYMMDDTHHMMSSフォーマット）
	if len(value) >= 15 {
		// タイムゾーン指定を取り除く
		value = strings.TrimSuffix(value, "Z")
		t, err := time.ParseInLocation("20060102T150405", value, loc)
		if err == nil {
			return t, false // 時間指定
		}
	}

	// RFC3339形式もサポート
	t, err := time.Parse(time.RFC3339, value)
	if err == nil {
		return t.In(loc), false
	}

	return time.Time{}, false
}

// convertToCalendarResponse はイベントリストを日付ごとに分類して
// CalendarResponseに変換するます。
func convertToCalendarResponse(events []eventWithDate, startDate, endDate time.Time) *models.CalendarResponse {
	// 日付ごとのマップを作成
	dayMap := make(map[string]*models.CalendarDay)

	// 7日分の日付を初期化
	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		dayMap[dateStr] = &models.CalendarDay{
			Date:   dateStr,
			AllDay: []models.Event{},
			Timed:  []models.Event{},
		}
	}

	// イベントを分類
	for _, evt := range events {
		dateStr := evt.date.Format("2006-01-02")
		day, exists := dayMap[dateStr]
		if !exists {
			continue
		}

		if evt.allDay {
			day.AllDay = append(day.AllDay, evt.event)
		} else {
			day.Timed = append(day.Timed, evt.event)
		}
	}

	// CalendarDayリストに変換してソート
	days := []models.CalendarDay{}
	for d := startDate; d.Before(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		if day, exists := dayMap[dateStr]; exists {
			// 時間帯イベントを時刻順にソート
			sort.Slice(day.Timed, func(i, j int) bool {
				return day.Timed[i].Start < day.Timed[j].Start
			})
			days = append(days, *day)
		}
	}

	return &models.CalendarResponse{
		Days: days,
	}
}
