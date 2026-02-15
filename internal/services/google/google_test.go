package google

import (
	"context"
	"testing"
	"time"

	"github.com/rihow/FamilyDashboard/internal/cache"
	"github.com/rihow/FamilyDashboard/internal/config"
	"github.com/rihow/FamilyDashboard/internal/models"
)

// TestNewClient は NewClient 関数で ClientStructが正しく初期化されるかをテストするのです。
func TestNewClient(t *testing.T) {
	fc := &cache.FileCache{} // ダミーキャッシュ
	cfg := &config.Config{}   // ダミー設定

	client := NewClient(fc, cfg)

	if client == nil {
		t.Errorf("NewClient() はクライアントを返すべきですが、nil が返されました")
	}
	if client.cache != fc {
		t.Errorf("NewClient() のキャッシュが一致しません")
	}
	if client.config != cfg {
		t.Errorf("NewClient() の設定が一致しません")
	}
}

// TestIsTokenValid はトークン有効性判定をテストするのです。
func TestIsTokenValid(t *testing.T) {
	client := &Client{
		accessToken:    "",
		tokenExpiresAt: time.Now(),
	}

	// トークンが無い場合は無効
	if client.IsTokenValid() {
		t.Errorf("トークンが無い場合は IsTokenValid() は false を返すべきです")
	}

	// トークンがあり、有効期限内の場合は有効
	client.accessToken = "test-token"
	client.tokenExpiresAt = time.Now().Add(1 * time.Hour)
	if !client.IsTokenValid() {
		t.Errorf("トークンが有効な場合、IsTokenValid() は true を返すべきです")
	}

	// トークンがあるが、有効期限切れの場合は無効
	client.tokenExpiresAt = time.Now().Add(-1 * time.Hour)
	if client.IsTokenValid() {
		t.Errorf("トークンの有効期限が切れている場合、IsTokenValid() は false を返すべきです")
	}
}

// TestGenerateDummyCalendarEvents はダミーカレンダーイベント生成をテストするのです。
func TestGenerateDummyCalendarEvents(t *testing.T) {
	client := &Client{}

	resp := client.generateDummyCalendarEvents()

	if resp == nil {
		t.Errorf("ダミーイベント生成がnilを返しました")
	}
	if len(resp.Days) != 7 {
		t.Errorf("ダミーイベント生成は7日分を返すべきですが、%d 日分が返されました", len(resp.Days))
	}

	// 各日付のイベント数をチェック
	for i, day := range resp.Days {
		if day.Date == "" {
			t.Errorf("日付 %d の Date が空です", i)
		}
		// 偶数日は終日イベントを持つ
		if i%2 == 0 {
			if len(day.AllDay) != 1 {
				t.Errorf("日付 %d は1つの終日イベントを持つべきですが、%d 個が返されました", i, len(day.AllDay))
			}
		}
		// すべての日付は時間帯付きイベントを持つ
		if len(day.Timed) != 1 {
			t.Errorf("日付 %d は1つの時間帯付きイベントを持つべきですが、%d 個が返されました", i, len(day.Timed))
		}
	}
}

// TestSortTaskItems はタスクソート機能をテストするのです。
// ソート順: 1) 期限（昇順、期限なしは最後） 2) 優先度（降順） 3) createdAt（昇順）
func TestSortTaskItems(t *testing.T) {
	client := &Client{}
	location, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(location)

	// ソート前の順序
	resp := &models.TasksResponse{
		Items: []models.TaskItem{
			{
				ID:        "task-1",
				Title:     "期限あり・優先度1・createdAt古い",
				DueDate:   stringPtr(now.Format("2006-01-02")),
				Priority:  1,
				CreatedAt: now.Add(-24 * time.Hour),
			},
			{
				ID:        "task-2",
				Title:     "期限あり・優先度2・createdAt新しい",
				DueDate:   stringPtr(now.Format("2006-01-02")),
				Priority:  2,
				CreatedAt: now,
			},
			{
				ID:        "task-3",
				Title:     "期限なし",
				DueDate:   nil,
				Priority:  1,
				CreatedAt: now,
			},
			{
				ID:        "task-4",
				Title:     "期限あり・優先度3・createdAt古い",
				DueDate:   stringPtr(now.AddDate(0, 0, 1).Format("2006-01-02")),
				Priority:  3,
				CreatedAt: now.Add(-48 * time.Hour),
			},
		},
	}

	// ソート実施
	client.sortTaskItems(resp)

	// ソート後の順序確認
	// 期待順序: task-1 (期限あり、優先度1), task-2 (期限あり、優先度2), task-4 (期限あり、優先度3), task-3 (期限なし)
	expectedOrder := []string{"task-1", "task-2", "task-4", "task-3"}
	for i, item := range resp.Items {
		if item.ID != expectedOrder[i] {
			t.Errorf("ソート順序 %d 番目が一致しません。期待: %s, 実際: %s", i, expectedOrder[i], item.ID)
		}
	}
}

// TestGenerateDummyTasks はダミータスク生成をテストするのです。
func TestGenerateDummyTasks(t *testing.T) {
	client := &Client{}

	resp := client.generateDummyTasks()

	if resp == nil {
		t.Errorf("ダミータスク生成がnilを返しました")
	}
	if len(resp.Items) == 0 {
		t.Errorf("ダミータスク生成はタスクを返すべきですが、0個が返されました")
	}

	// ソート順を確認（期限→優先度→createdAt）
	for i := 0; i < len(resp.Items)-1; i++ {
		curr := resp.Items[i]
		next := resp.Items[i+1]

		// 期限の比較
		if (curr.DueDate != nil) && (next.DueDate == nil) {
			// 期限を持つタスクが、期限をもたないタスクより先に来ている（OK）
			continue
		}
		if (curr.DueDate == nil) && (next.DueDate != nil) {
			t.Errorf("期限なしタスクが期限ありタスクより先に来ています。ソートが不正です")
		}

		// 両方が期限を持つ場合、期限を比較
		if curr.DueDate != nil && next.DueDate != nil {
			if *curr.DueDate > *next.DueDate {
				t.Errorf("期限ソートが不正です。期待: %s <= %s", *curr.DueDate, *next.DueDate)
			}
		}
	}
}

// TestConvertTasksResponse はGoogle Tasks APIレスポンス変換をテストするのです。
func TestConvertTasksResponse(t *testing.T) {
	client := &Client{}

	gTasksResp := GoogleTasksResponse{
		Items: []GoogleTask{
			{
				ID:      "task-1",
				Title:   "Test Task 1",
				Notes:   "Test Notes",
				Status:  "needsAction",
				DueDate: "2026-02-16",
				Updated: time.Now().Format(time.RFC3339),
			},
			{
				ID:     "task-2",
				Title:  "Test Task 2",
				Status: "completed",
			},
		},
	}

	resp, err := client.convertTasksResponse(gTasksResp)

	if err != nil {
		t.Errorf("変換エラー: %v", err)
	}
	if resp == nil {
		t.Errorf("変換がnilを返しました")
	}
	if len(resp.Items) != 2 {
		t.Errorf("期待: 2個のタスク、実際: %d", len(resp.Items))
	}

	// 最初のタスクをチェック
	if resp.Items[0].Title != "Test Task 1" {
		t.Errorf("タスクタイトルが一致しません。期待: Test Task 1, 実際: %s", resp.Items[0].Title)
	}
	if resp.Items[0].DueDate == nil || *resp.Items[0].DueDate != "2026-02-16" {
		t.Errorf("期限が一致しません。期待: 2026-02-16, 実際: %v", resp.Items[0].DueDate)
	}
}

// TestParseGoogleDateTime は Google Calendar の DateTime/Date パースをテストするのです。
func TestParseGoogleDateTime(t *testing.T) {
	client := &Client{}
	location, _ := time.LoadLocation("Asia/Tokyo")

	// RFC3339 形式のテスト
	dateTime := "2026-02-16T14:30:00+09:00"
	t1, isAllDay, err := client.parseGoogleDateTime(dateTime, "", location)
	if err != nil {
		t.Errorf("DateTime パースエラー: %v", err)
	}
	if isAllDay {
		t.Errorf("DateTime は終日ではないべきです")
	}
	if t1.Year() != 2026 || t1.Month() != 2 || t1.Day() != 16 {
		t.Errorf("DateTime のパースが不正です: %v", t1)
	}

	// Date 形式のテスト
	date := "2026-02-16"
	t2, isAllDay2, err := client.parseGoogleDateTime("", date, location)
	if err != nil {
		t.Errorf("Date パースエラー: %v", err)
	}
	if !isAllDay2 {
		t.Errorf("Date は終日イベントべきです")
	}
	if t2.Year() != 2026 || t2.Month() != 2 || t2.Day() != 16 {
		t.Errorf("Date のパースが不正です: %v", t2)
	}
}

// TestGetEventColor は イベント色マッピングをテストするのです。
func TestGetEventColor(t *testing.T) {
	client := &Client{}

	// 既知の色IDをテスト
	color1 := client.getEventColor("1", "")
	if color1 != "#A4BDFC" {
		t.Errorf("色ID 1 のマッピングが不正です。期待: #A4BDFC, 実際: %s", color1)
	}

	color10 := client.getEventColor("10", "")
	if color10 != "#E67C73" {
		t.Errorf("色ID 10 のマッピングが不正です。期待: #E67C73, 実際: %s", color10)
	}

	// 未知の色IDはデフォルト色を返す
	colorDefault := client.getEventColor("unknown", "")
	if colorDefault != "#A4BDFC" {
		t.Errorf("未知の色IDはデフォルト色を返すべきです。期待: #A4BDFC, 実際: %s", colorDefault)
	}

	// eventColor が指定されている場合はそれを使用
	colorEvent := client.getEventColor("unknown", "#FF0000")
	if colorEvent != "#FF0000" {
		t.Errorf("eventColor が指定されている場合はそれを使用するべきです。期待: #FF0000, 実際: %s", colorEvent)
	}
}

// TestParseJSONResponse は JSON パースをテストするのです。
func TestParseJSONResponse(t *testing.T) {
	type TestStruct struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	jsonData := []byte(`{"name": "Test User", "age": 30}`)
	var result TestStruct

	err := parseJSONResponse(jsonData, &result)
	if err != nil {
		t.Errorf("JSON パースエラー: %v", err)
	}

	if result.Name != "Test User" || result.Age != 30 {
		t.Errorf("パース結果が不正です。Name=%s, Age=%d", result.Name, result.Age)
	}
}

// TestSetAccessToken はアクセストークン設定をテストするのです。
func TestSetAccessToken(t *testing.T) {
	client := &Client{}

	client.SetAccessToken("test-token", 3600) // 1時間有効

	if client.accessToken != "test-token" {
		t.Errorf("Set AccessToken が失敗しました")
	}
	if !client.IsTokenValid() {
		t.Errorf("設定後、トークンは有効であるべきです")
	}
}

// TestContextCancellation は GetCalendarEvents でコンテキストキャンセルをテストするのです。
func TestGetCalendarEvents_NoToken(t *testing.T) {
	fc := &cache.FileCache{} // ダミーキャッシュ
	cfg := &config.Config{
		RefreshIntervals: config.RefreshIntervals{
			CalendarSec: 300,
		},
	}

	client := NewClient(fc, cfg)

	// トークンが無いため、ダミーイベントが返される
	resp, err := client.GetCalendarEvents(context.Background())

	if err != nil {
		t.Errorf("トークンが無い場合にダミーイベント返却で、エラーは返されないべきです: %v", err)
	}
	if resp == nil {
		t.Errorf("レスポンスはnilではないべきです")
	}
	if len(resp.Days) != 7 {
		t.Errorf("期待: 7日分のイベント、実際: %d", len(resp.Days))
	}
}

// TestGetTaskItems_NoToken は GetTaskItems でトークン無し時をテストするのです。
func TestGetTaskItems_NoToken(t *testing.T) {
	fc := &cache.FileCache{} // ダミーキャッシュ
	cfg := &config.Config{
		RefreshIntervals: config.RefreshIntervals{
			TasksSec: 300,
		},
	}

	client := NewClient(fc, cfg)

	// トークンが無いため、ダミータスクが返される
	resp, err := client.GetTaskItems(context.Background())

	if err != nil {
		t.Errorf("トークンが無い場合にダミータスク返却で、エラーは返されないべきです: %v", err)
	}
	if resp == nil {
		t.Errorf("レスポンスはnilではないべきです")
	}
	if len(resp.Items) == 0 {
		t.Errorf("ダミータスクは返されるべきです")
	}
}
