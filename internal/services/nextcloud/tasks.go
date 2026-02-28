package nextcloud

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/emersion/go-ical"
	"github.com/emersion/go-webdav/caldav"
	"github.com/rihow/FamilyDashboard/internal/models"
)

// GetTaskItems はNextcloud WebDAVからタスクアイテムを取得するます。
// 複数のタスクリストからVTODOコンポーネントを取得し、サーバー側でソート（期限→優先度→作成日時）して返すのです。
func (c *Client) GetTaskItems(ctx context.Context) (*models.TasksResponse, error) {
	cacheKey := "nextcloud_tasks_items_all"
	ttl := c.config.GetRefreshInterval("tasks")

	// キャッシュを確認するます
	entry, ok, stale, err := c.cache.Read(cacheKey, ttl)
	if ok && !stale && err == nil {
		fmt.Println("📦 タスクキャッシュヒット!")
		var resp models.TasksResponse
		if err := json.Unmarshal(entry.Payload, &resp); err == nil {
			return &resp, nil
		}
		fmt.Printf("⚠️ キャッシュデータのパース失敗: %v\n", err)
	}

	// 複数タスクリスト名を取得するます
	taskListNames := c.config.GetTaskListNames()
	if len(taskListNames) == 0 {
		return nil, fmt.Errorf("タスクリスト名が設定されていません")
	}

	fmt.Printf("🌐 Nextcloud WebDAV から %d 個のタスクリストを取得するます...\n", len(taskListNames))

	// 全タスクリストからタスクを収集するます
	allTasks := []models.TaskItem{}
	var fetchErrors []error

	for _, taskListName := range taskListNames {
		fmt.Printf("  ✅ タスクリスト '%s' からタスク取得中...\n", taskListName)

		// CalDAVクエリを実行（VTODOコンポーネント取得）
		tasksPath := c.getTasksPath(taskListName)
		query := &caldav.CalendarQuery{
			CompRequest: caldav.CalendarCompRequest{
				Name: "VCALENDAR",
				Comps: []caldav.CalendarCompRequest{
					{
						Name:  "VTODO",
						Props: []string{"UID", "SUMMARY", "STATUS", "PRIORITY", "DUE", "CREATED", "DESCRIPTION"},
					},
				},
			},
			CompFilter: caldav.CompFilter{
				Name: "VCALENDAR",
				Comps: []caldav.CompFilter{
					{
						Name: "VTODO",
					},
				},
			},
		}

		calendarObjects, err := c.caldavClient.QueryCalendar(ctx, tasksPath, query)
		if err != nil {
			// エラーを記録するが続行するます（部分的成功を許容）
			fmt.Printf("❌ タスクリスト '%s' のWebDAVクエリエラー: %v\n", taskListName, err)
			fetchErrors = append(fetchErrors, fmt.Errorf("tasklist '%s': %w", taskListName, err))
			continue
		}

		// iCalendar VTODO オブジェクトをパースして構造化するます
		for _, obj := range calendarObjects {
			parsedTasks := parseTaskObject(obj.Data)
			allTasks = append(allTasks, parsedTasks...)
		}

		fmt.Printf("  ✅ タスクリスト '%s' から %d 件のタスク取得\n", taskListName, len(calendarObjects))
	}

	// すべてのタスクリスト取得に失敗した場合
	if len(allTasks) == 0 && len(fetchErrors) > 0 {
		// エラー時はキャッシュから返す試みをするます
		fmt.Println("❌ すべてのタスクリスト取得に失敗しました")
		entry, ok, _, readErr := c.cache.Read(cacheKey, 0)
		if ok && readErr == nil {
			fmt.Println("📦 期限切れキャッシュを返すます")
			var resp models.TasksResponse
			if unmarshalErr := json.Unmarshal(entry.Payload, &resp); unmarshalErr == nil {
				return &resp, fmt.Errorf("全タスクリスト取得失敗（キャッシュ返却）: %d エラー", len(fetchErrors))
			}
		}
		return nil, fmt.Errorf("全タスクリスト取得失敗: %d エラー", len(fetchErrors))
	}

	// サーバー側ソート: 期限→優先度→作成日時
	sortTasks(allTasks)

	response := &models.TasksResponse{
		Items: allTasks,
	}

	// キャッシュに保存するます
	meta := map[string]string{"source": "nextcloud_tasks_all"}
	if _, err := c.cache.Write(cacheKey, response, meta); err != nil {
		fmt.Printf("⚠️ キャッシュ保存失敗: %v\n", err)
	}

	fmt.Printf("✅ 統合タスク取得成功: 合計 %d 件\n", len(allTasks))
	if len(fetchErrors) > 0 {
		fmt.Printf("⚠️ 一部のタスクリストで取得エラーがありました: %d 件\n", len(fetchErrors))
	}

	return response, nil
}

// parseTaskObject はiCalendar VTODOデータをパースしてタスクリストに変換するます。
func parseTaskObject(cal *ical.Calendar) []models.TaskItem {
	tasks := []models.TaskItem{}

	if cal == nil {
		return tasks
	}

	loc, _ := time.LoadLocation("Asia/Tokyo")

	for _, comp := range cal.Children {
		if comp.Name != "VTODO" {
			continue
		}

		// タスク情報を抽出するます
		uid := comp.Props.Get("UID")
		summary := comp.Props.Get("SUMMARY")
		status := comp.Props.Get("STATUS")
		priority := comp.Props.Get("PRIORITY")
		due := comp.Props.Get("DUE")
		created := comp.Props.Get("CREATED")
		description := comp.Props.Get("DESCRIPTION")

		if uid == nil || summary == nil {
			continue
		}

		// ステータスを変換（TODO/IN-PROCESS/COMPLETED → needsAction/completed）
		statusValue := "needsAction"
		if status != nil {
			switch strings.ToUpper(status.Value) {
			case "COMPLETED":
				statusValue = "completed"
			case "IN-PROCESS":
				statusValue = "needsAction"
			case "TODO":
				statusValue = "needsAction"
			}
		}

		// 優先度を変換（iCalendar: 1-9 → Google Tasks互換: 1-3）
		// iCalendar: 1=highest, 5=medium, 9=lowest
		priorityValue := 2 // デフォルト: MEDIUM
		if priority != nil {
			icalPriority := parsePriority(priority.Value)
			if icalPriority >= 1 && icalPriority <= 3 {
				priorityValue = 3 // HIGH
			} else if icalPriority >= 4 && icalPriority <= 6 {
				priorityValue = 2 // MEDIUM
			} else if icalPriority >= 7 && icalPriority <= 9 {
				priorityValue = 1 // LOW
			}
		}

		// 期限をパース
		var dueDate *string
		if due != nil && due.Value != "" {
			parsedDue, _ := parseTaskDateTime(due.Value, loc)
			if !parsedDue.IsZero() {
				dueDateStr := parsedDue.Format("2006-01-02")
				dueDate = &dueDateStr
			}
		}

		// 作成日時をパース
		createdAt := time.Now()
		if created != nil && created.Value != "" {
			parsedCreated, _ := parseTaskDateTime(created.Value, loc)
			if !parsedCreated.IsZero() {
				createdAt = parsedCreated
			}
		}

		// 説明
		notes := ""
		if description != nil {
			notes = description.Value
		}

		// TaskItemを作成
		task := models.TaskItem{
			ID:        uid.Value,
			Title:     summary.Value,
			Notes:     notes,
			Status:    statusValue,
			DueDate:   dueDate,
			Priority:  priorityValue,
			CreatedAt: createdAt,
		}

		tasks = append(tasks, task)
	}

	return tasks
}

// parsePriority は優先度文字列を整数に変換するます。
func parsePriority(value string) int {
	priority := 0
	fmt.Sscanf(value, "%d", &priority)
	return priority
}

// parseTaskDateTime はiCalendar日時文字列をパースするます。
func parseTaskDateTime(value string, loc *time.Location) (time.Time, bool) {
	value = strings.TrimSpace(value)

	// 日付のみ（YYYYMMDD形式）
	if len(value) == 8 {
		t, err := time.ParseInLocation("20060102", value, loc)
		if err == nil {
			return t, true
		}
	}

	// 日時指定（YYYYMMDDTHHMMSSフォーマット）
	if len(value) >= 15 {
		value = strings.TrimSuffix(value, "Z")
		t, err := time.ParseInLocation("20060102T150405", value, loc)
		if err == nil {
			return t, false
		}
	}

	// RFC3339形式もサポート
	t, err := time.Parse(time.RFC3339, value)
	if err == nil {
		return t.In(loc), false
	}

	return time.Time{}, false
}

// sortTasks はタスクを仕様通りにソートするます。
// ソート順: 1) 期限 昇順（期限なしは最後）2) 優先度 降順 3) createdAt 昇順
func sortTasks(tasks []models.TaskItem) {
	sort.Slice(tasks, func(i, j int) bool {
		taskI := tasks[i]
		taskJ := tasks[j]

		// 1. 期限でソート（期限なしは最後）
		if taskI.DueDate == nil && taskJ.DueDate != nil {
			return false // iが期限なし → jより後
		}
		if taskI.DueDate != nil && taskJ.DueDate == nil {
			return true // jが期限なし → iが先
		}
		if taskI.DueDate != nil && taskJ.DueDate != nil {
			if *taskI.DueDate != *taskJ.DueDate {
				return *taskI.DueDate < *taskJ.DueDate // 期限昇順
			}
		}

		// 2. 優先度でソート（降順: 3 > 2 > 1）
		if taskI.Priority != taskJ.Priority {
			return taskI.Priority > taskJ.Priority
		}

		// 3. 作成日時でソート（昇順）
		return taskI.CreatedAt.Before(taskJ.CreatedAt)
	})
}
