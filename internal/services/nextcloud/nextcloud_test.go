package nextcloud

import (
	"testing"
	"time"

	"github.com/rihow/FamilyDashboard/internal/cache"
	"github.com/rihow/FamilyDashboard/internal/config"
	"github.com/rihow/FamilyDashboard/internal/models"
)

// TestNewClient はクライアント初期化のテストなのです。
func TestNewClient(t *testing.T) {
	// テスト用設定を作成
	cfg := &config.Config{
		Nextcloud: config.Nextcloud{
			ServerURL:     "https://nextcloud.example.com",
			Username:      "testuser",
			Password:      "testpass",
			CalendarNames: []string{"family"},
			TaskListNames: []string{"tasks"},
		},
	}

	fc := cache.New("./test_cache")
	client, err := NewClient(fc, cfg)

	if err != nil {
		t.Fatalf("NewClient エラー: %v", err)
	}

	if client == nil {
		t.Fatal("クライアントが nil なのです")
	}

	if client.config.Nextcloud.Username != "testuser" {
		t.Errorf("Username が一致しません: got %s, want testuser", client.config.Nextcloud.Username)
	}
}

// TestNewClientWithInvalidConfig は不正な設定でエラーになるかのテストなのです。
func TestNewClientWithInvalidConfig(t *testing.T) {
	tests := []struct {
		name   string
		config *config.Config
	}{
		{
			name:   "nil config",
			config: nil,
		},
		{
			name: "empty ServerURL",
			config: &config.Config{
				Nextcloud: config.Nextcloud{
					Username: "user",
					Password: "pass",
				},
			},
		},
		{
			name: "empty Username",
			config: &config.Config{
				Nextcloud: config.Nextcloud{
					ServerURL: "https://example.com",
					Password:  "pass",
				},
			},
		},
		{
			name: "empty Password",
			config: &config.Config{
				Nextcloud: config.Nextcloud{
					ServerURL: "https://example.com",
					Username:  "user",
				},
			},
		},
	}

	fc := cache.New("./test_cache")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewClient(fc, tt.config)
			if err == nil {
				t.Errorf("%s: エラーが期待されましたが nil でした", tt.name)
			}
		})
	}
}

// TestParseDateTime は日時パース機能のテストなのです。
func TestParseDateTime(t *testing.T) {
	loc, _ := time.LoadLocation("Asia/Tokyo")

	tests := []struct {
		name      string
		input     string
		wantAllDay bool
		wantError bool
	}{
		{
			name:      "終日イベント",
			input:     "20260228",
			wantAllDay: true,
			wantError: false,
		},
		{
			name:      "時間指定イベント",
			input:     "20260228T143000",
			wantAllDay: false,
			wantError: false,
		},
		{
			name:      "UTC時間",
			input:     "20260228T143000Z",
			wantAllDay: false,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, isAllDay := parseDateTime(tt.input, loc)
			
			if tt.wantError && !result.IsZero() {
				t.Errorf("エラーが期待されましたがパース成功しました: %v", result)
			}
			
			if !tt.wantError && result.IsZero() {
				t.Errorf("パース失敗しました: %s", tt.input)
			}
			
			if !tt.wantError && isAllDay != tt.wantAllDay {
				t.Errorf("allDay 不一致: got %v, want %v", isAllDay, tt.wantAllDay)
			}
		})
	}
}

// TestSortTasks はタスクソート機能のテストなのです。
func TestSortTasks(t *testing.T) {
	// テストタスクを作成
	dueDate1 := "2026-03-01"
	dueDate2 := "2026-03-05"
	
	tasks := []models.TaskItem{
		{
			ID:        "task1",
			Title:     "期限なし・優先度低",
			Priority:  1,
			DueDate:   nil,
			CreatedAt: time.Date(2026, 2, 20, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        "task2",
			Title:     "期限3/5・優先度高",
			Priority:  3,
			DueDate:   &dueDate2,
			CreatedAt: time.Date(2026, 2, 21, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        "task3",
			Title:     "期限3/1・優先度中",
			Priority:  2,
			DueDate:   &dueDate1,
			CreatedAt: time.Date(2026, 2, 19, 10, 0, 0, 0, time.UTC),
		},
		{
			ID:        "task4",
			Title:     "期限3/1・優先度高",
			Priority:  3,
			DueDate:   &dueDate1,
			CreatedAt: time.Date(2026, 2, 22, 10, 0, 0, 0, time.UTC),
		},
	}

	// ソート実行
	sortTasks(tasks)

	// 期待される順序:
	// 1. task4: 期限3/1・優先度高・作成2/22
	// 2. task3: 期限3/1・優先度中・作成2/19
	// 3. task2: 期限3/5・優先度高・作成2/21
	// 4. task1: 期限なし・優先度低・作成2/20

	expectedOrder := []string{"task4", "task3", "task2", "task1"}

	for i, expectedID := range expectedOrder {
		if tasks[i].ID != expectedID {
			t.Errorf("ソート順序エラー [%d]: got %s, want %s", i, tasks[i].ID, expectedID)
		}
	}
}

// TestGetCalendarPath はカレンダーパス生成のテストなのです。
func TestGetCalendarPath(t *testing.T) {
	cfg := &config.Config{
		Nextcloud: config.Nextcloud{
			ServerURL:     "https://nextcloud.example.com",
			Username:      "testuser",
			Password:      "testpass",
			CalendarNames: []string{"family"},
		},
	}

	fc := cache.New("./test_cache")
	client, _ := NewClient(fc, cfg)

	path := client.getCalendarPath("family")
	expected := "/remote.php/dav/calendars/testuser/family/"

	if path != expected {
		t.Errorf("カレンダーパス不一致: got %s, want %s", path, expected)
	}
}

// TestGetTasksPath はタスクパス生成のテストなのです。
func TestGetTasksPath(t *testing.T) {
	cfg := &config.Config{
		Nextcloud: config.Nextcloud{
			ServerURL:     "https://nextcloud.example.com",
			Username:      "testuser",
			Password:      "testpass",
			TaskListNames: []string{"tasks"},
		},
	}

	fc := cache.New("./test_cache")
	client, _ := NewClient(fc, cfg)

	path := client.getTasksPath("tasks")
	expected := "/remote.php/dav/calendars/testuser/tasks/"

	if path != expected {
		t.Errorf("タスクパス不一致: got %s, want %s", path, expected)
	}
}

// TestParsePriority は優先度パースのテストなのです。
func TestParsePriority(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"1", 1},
		{"5", 5},
		{"9", 9},
		{"invalid", 0},
		{"", 0},
	}

	for _, tt := range tests {
		result := parsePriority(tt.input)
		if result != tt.expected {
			t.Errorf("parsePriority(%s): got %d, want %d", tt.input, result, tt.expected)
		}
	}
}

// TestMultipleCalendarPaths は複数カレンダー名でパス生成のテストなのです。
func TestMultipleCalendarPaths(t *testing.T) {
	cfg := &config.Config{
		Nextcloud: config.Nextcloud{
			ServerURL:     "https://nextcloud.example.com",
			Username:      "testuser",
			Password:      "testpass",
			CalendarNames: []string{"family", "work", "personal"},
		},
	}

	fc := cache.New("./test_cache")
	client, _ := NewClient(fc, cfg)

	// 複数カレンダー名をテスト
	tests := []struct {
		calendarName string
		expectedPath string
	}{
		{"family", "/remote.php/dav/calendars/testuser/family/"},
		{"work", "/remote.php/dav/calendars/testuser/work/"},
		{"personal", "/remote.php/dav/calendars/testuser/personal/"},
	}

	for _, tt := range tests {
		path := client.getCalendarPath(tt.calendarName)
		if path != tt.expectedPath {
			t.Errorf("カレンダーパス不一致 (%s): got %s, want %s", tt.calendarName, path, tt.expectedPath)
		}
	}
}

// TestMultipleTaskListPaths は複数タスクリスト名でパス生成のテストなのです。
func TestMultipleTaskListPaths(t *testing.T) {
	cfg := &config.Config{
		Nextcloud: config.Nextcloud{
			ServerURL:     "https://nextcloud.example.com",
			Username:      "testuser",
			Password:      "testpass",
			TaskListNames: []string{"tasks", "personal", "work"},
		},
	}

	fc := cache.New("./test_cache")
	client, _ := NewClient(fc, cfg)

	// 複数タスクリスト名をテスト
	tests := []struct {
		taskListName string
		expectedPath string
	}{
		{"tasks", "/remote.php/dav/calendars/testuser/tasks/"},
		{"personal", "/remote.php/dav/calendars/testuser/personal/"},
		{"work", "/remote.php/dav/calendars/testuser/work/"},
	}

	for _, tt := range tests {
		path := client.getTasksPath(tt.taskListName)
		if path != tt.expectedPath {
			t.Errorf("タスクパス不一致 (%s): got %s, want %s", tt.taskListName, path, tt.expectedPath)
		}
	}
}

// TestGetCalendarNames は設定から複数カレンダー名を取得するテストなのです。
func TestGetCalendarNames(t *testing.T) {
	cfg := &config.Config{
		Nextcloud: config.Nextcloud{
			CalendarNames: []string{"family", "work", "personal"},
		},
	}

	names := cfg.GetCalendarNames()
	expected := []string{"family", "work", "personal"}

	if len(names) != len(expected) {
		t.Errorf("カレンダー名数不一致: got %d, want %d", len(names), len(expected))
		return
	}

	for i, name := range names {
		if name != expected[i] {
			t.Errorf("カレンダー名不一致 [%d]: got %s, want %s", i, name, expected[i])
		}
	}
}

// TestGetTaskListNames は設定から複数タスクリスト名を取得するテストなのです。
func TestGetTaskListNames(t *testing.T) {
	cfg := &config.Config{
		Nextcloud: config.Nextcloud{
			TaskListNames: []string{"tasks", "personal", "shopping"},
		},
	}

	names := cfg.GetTaskListNames()
	expected := []string{"tasks", "personal", "shopping"}

	if len(names) != len(expected) {
		t.Errorf("タスクリスト名数不一致: got %d, want %d", len(names), len(expected))
		return
	}

	for i, name := range names {
		if name != expected[i] {
			t.Errorf("タスクリスト名不一致 [%d]: got %s, want %s", i, name, expected[i])
		}
	}
}
