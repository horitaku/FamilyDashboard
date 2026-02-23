package google

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/rihow/FamilyDashboard/internal/models"
)

// GetTaskItems は共有タスクリストからタスクを取得するのです。
// サーバー側でソート（期限→優先度→createdAt）を適用して返すます。
// キャッシュが有効な場合はキャッシュから返し、無い場合はGoogle Tasks APIから取得するのです。
func (c *Client) GetTaskItems(ctx context.Context) (*models.TasksResponse, error) {
	// キャッシュキーを生成するのです
	cacheKey := "google_tasks_items"
	ttl := c.config.GetRefreshInterval("tasks")

	// キャッシュが有効か確認するのです
	var cachedData []byte
	entry, exists, stale, err := c.cache.Read(cacheKey, ttl)
	if err == nil && exists {
		if !stale {
			// キャッシュがヒットしたので、パースして返すのです
			var resp models.TasksResponse
			if err := parseJSONResponse(entry.Payload, &resp); err == nil {
				return &resp, nil
			}
			// パースエラーの場合は、APIから新たに取得するのです
		} else {
			// 期限切れでもフォールバック用に保持するのです
			cachedData = entry.Payload
		}
	}

	// トークンを確認して、必要に応じて自動更新するのです
	if err := c.EnsureTokenValid(ctx); err != nil {
		// トークンが無い または更新失敗の場合
		fmt.Printf("❌ トークン確認エラー: %v\n", err)
		// キャッシュがあれば使用するのです
		if cachedData != nil {
			var resp models.TasksResponse
			if err := parseJSONResponse(cachedData, &resp); err == nil {
				fmt.Printf("⚠️ キャッシュ（期限切れ）を使用するのです\n")
				return &resp, nil
			}
		}
		return nil, err
	}

	// Google Tasks APIから取得するのです
	// API URL例: https://www.googleapis.com/tasks/v1/lists/{taskListId}/tasks
	// 本実装では、settings.jsonの taskListId を使用するます（設定がない場合は @default）
	listID := "@default"
	if c.config.Google.TaskListID != "" {
		listID = c.config.Google.TaskListID
	}

	url := fmt.Sprintf(
		"https://www.googleapis.com/tasks/v1/lists/%s/tasks?showCompleted=false&maxResults=100",
		listID,
	)

	// HTTPリクエストを実行するのです
	body, err := c.doRequest(ctx, "GET", url)
	if err != nil {
		// APIエラーの場合は、キャッシュがあれば使用するのです
		if cachedData != nil {
			var resp models.TasksResponse
			if err := parseJSONResponse(cachedData, &resp); err == nil {
				fmt.Printf("⚠️ Google Tasks API error, cache を使用するのです: %v\n", err)
				return &resp, err
			}
		}
		return nil, fmt.Errorf("Google Tasks API error: %w", err)
	}

	// Google Tasks APIのレスポンス形式をパースするのです
	var gTasksResp GoogleTasksResponse
	if err := parseJSONResponse(body, &gTasksResp); err != nil {
		return nil, fmt.Errorf("tasks response parse error: %w", err)
	}

	// レスポンスを models.TasksResponse に変換するのです
	resp, err := c.convertTasksResponse(gTasksResp)
	if err != nil {
		return nil, fmt.Errorf("tasks conversion error: %w", err)
	}

	// サーバー側ソート: 期限 → 優先度 → createdAt
	c.sortTaskItems(resp)

	// キャッシュに保存するのです
	if data, err := models.ToJSON(resp); err == nil {
		_ = c.saveCache(cacheKey, data)
	}

	return resp, nil
}

// GoogleTasksResponse は Google Tasks API のレスポンス形式です。
type GoogleTasksResponse struct {
	Items []GoogleTask `json:"items"`
}

// GoogleTask は Google Tasks API のタスク形式です。
type GoogleTask struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Notes    string `json:"notes"`
	Status   string `json:"status"`   // "needsAction" or "completed"
	DueDate  string `json:"due"`      // ISO 8601 format (YYYY-MM-DD)
	Updated  string `json:"updated"`  // RFC 3339 format
	Position string `json:"position"` // Google Tasks の内部順序キー
}

// convertTasksResponse は GoogleTasksResponse を models.TasksResponse に変換するのです。
func (c *Client) convertTasksResponse(gTasksResp GoogleTasksResponse) (*models.TasksResponse, error) {
	items := make([]models.TaskItem, 0)

	for _, gTask := range gTasksResp.Items {
		// 優先度はGoogle Tasksでは直接サポートされていないため、1（最高）で統一するのです
		// 後で別の仕組み（例: タスクのカスタムフィール）で優先度を管理することを検討するのです
		priority := 1

		// createdAt は updated フィールドから取得するのです
		createdAt := time.Now()
		if gTask.Updated != "" {
			if t, err := time.Parse(time.RFC3339, gTask.Updated); err == nil {
				createdAt = t
			}
		}

		// dueDate をパースするのです（ISO 8601フォーマット YYYY-MM-DD）
		var dueDatePtr *string
		if gTask.DueDate != "" {
			dueDatePtr = &gTask.DueDate
		}

		item := models.TaskItem{
			ID:        gTask.ID,
			Title:     gTask.Title,
			Notes:     gTask.Notes,
			Status:    gTask.Status,
			DueDate:   dueDatePtr,
			Priority:  priority,
			CreatedAt: createdAt,
		}

		items = append(items, item)
	}

	return &models.TasksResponse{Items: items}, nil
}

// sortTaskItems はタスクをサーバー側ソート規則に従ってソートするのです。
// ソート順序: 1) 期限 昇順（期限なしは最後） 2) 優先度 降順 3) createdAt 昇順
func (c *Client) sortTaskItems(resp *models.TasksResponse) {
	sort.SliceStable(resp.Items, func(i, j int) bool {
		a := resp.Items[i]
		b := resp.Items[j]

		// 1) 期限で比較。期限がある方が先。
		if (a.DueDate == nil) != (b.DueDate == nil) {
			return a.DueDate != nil // a が期限を持つ場合は true（先に並ぶ）
		}

		// 両方が期限を持つ場合は、期限の日付で比較
		if a.DueDate != nil && b.DueDate != nil {
			if *a.DueDate != *b.DueDate {
				return *a.DueDate < *b.DueDate // 期限 昇順
			}
		}

		// 2) 優先度で比較。優先度が高い（値が小さい）方が先。
		if a.Priority != b.Priority {
			return a.Priority < b.Priority // 優先度 降順（値が小さい=優先度高）
		}

		// 3) createdAt で比較
		return a.CreatedAt.Before(b.CreatedAt) // createdAt 昇順
	})
}

// generateDummyTasks はテスト用のダミータスクを生成するのです。
func (c *Client) generateDummyTasks() *models.TasksResponse {
	location, _ := time.LoadLocation("Asia/Tokyo")
	now := time.Now().In(location)

	items := []models.TaskItem{
		{
			ID:        "dummy-task-1",
			Title:     "買い物に行く",
			Notes:     "牛乳とパンを買う",
			Status:    "needsAction",
			DueDate:   stringPtr(now.AddDate(0, 0, 1).Format("2006-01-02")),
			Priority:  1,
			CreatedAt: now.Add(-24 * time.Hour),
		},
		{
			ID:        "dummy-task-2",
			Title:     "家の掃除",
			Notes:     "リビングとキッチン",
			Status:    "needsAction",
			DueDate:   stringPtr(now.AddDate(0, 0, 3).Format("2006-01-02")),
			Priority:  2,
			CreatedAt: now.Add(-48 * time.Hour),
		},
		{
			ID:        "dummy-task-3",
			Title:     "メールを返信する",
			Notes:     "",
			Status:    "needsAction",
			DueDate:   nil, // 期限なし
			Priority:  1,
			CreatedAt: now.Add(-12 * time.Hour),
		},
		{
			ID:        "dummy-task-4",
			Title:     "ガスの検針",
			Notes:     "月初に実施する",
			Status:    "needsAction",
			DueDate:   stringPtr(now.AddDate(0, 0, 2).Format("2006-01-02")),
			Priority:  3,
			CreatedAt: now.Add(-72 * time.Hour),
		},
	}

	// ソート規則に従ってソートするのです
	resp := &models.TasksResponse{Items: items}
	c.sortTaskItems(resp)

	return resp
}

// stringPtr は string へのポインタを返すヘルパーなのです。
func stringPtr(s string) *string {
	return &s
}
