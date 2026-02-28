package nextcloud

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/emersion/go-webdav/caldav"
	"github.com/rihow/FamilyDashboard/internal/cache"
	"github.com/rihow/FamilyDashboard/internal/config"
)

// Client は Nextcloud CalDAV/WebDAV のクライアントなのです。
// カレンダー・タスクの取得とキャッシュ管理を担当するます。
type Client struct {
	cache        *cache.FileCache
	config       *config.Config
	httpClient   *http.Client
	caldavClient *caldav.Client
}

// NewClient は Nextcloud クライアントを初期化するます。
// Basic認証でCalDAVサーバーに接続する準備をするのです。
func NewClient(fc *cache.FileCache, cfg *config.Config) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("設定が nil なのです")
	}

	// Nextcloud 設定の検証
	if cfg.Nextcloud.ServerURL == "" {
		return nil, fmt.Errorf("Nextcloud ServerURL が設定されていません")
	}
	if cfg.Nextcloud.Username == "" {
		return nil, fmt.Errorf("Nextcloud Username が設定されていません")
	}
	if cfg.Nextcloud.Password == "" {
		return nil, fmt.Errorf("Nextcloud Password が設定されていません")
	}

	// HTTPクライアントを作成（タイムアウト付き＋Basic認証）
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &basicAuthTransport{
			Username: cfg.Nextcloud.Username,
			Password: cfg.Nextcloud.Password,
		},
	}

	// CalDAV クライアントを作成（Basic認証付きHTTPクライアント）
	caldavClient, err := caldav.NewClient(httpClient, cfg.Nextcloud.ServerURL)
	if err != nil {
		return nil, fmt.Errorf("CalDAVクライアント初期化エラー: %w", err)
	}

	client := &Client{
		cache:        fc,
		config:       cfg,
		httpClient:   httpClient,
		caldavClient: caldavClient,
	}

	fmt.Printf("✅ Nextcloud クライアント初期化成功: %s (ユーザー: %s)\n",
		cfg.Nextcloud.ServerURL, cfg.Nextcloud.Username)

	return client, nil
}

// basicAuthTransport は Basic認証用のHTTPトランスポートなのです。
type basicAuthTransport struct {
	Username string
	Password string
}

// RoundTrip はHTTPリクエストにBasic認証ヘッダーを追加するます。
func (t *basicAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(t.Username, t.Password)
	return http.DefaultTransport.RoundTrip(req)
}

// getCalendarPath はカレンダーのCalDAVパスを返すます。
// Nextcloudの標準パス: /remote.php/dav/calendars/USERNAME/CALENDARNAME/
func (c *Client) getCalendarPath(calendarName string) string {
	username := c.config.Nextcloud.Username
	if calendarName == "" {
		calendarName = "personal" // デフォルトカレンダー
	}
	return fmt.Sprintf("/remote.php/dav/calendars/%s/%s/", username, calendarName)
}

// getTasksPath はタスクのCalDAVパスを返すます。
// Nextcloudの標準パス: /remote.php/dav/calendars/USERNAME/TASKLISTNAME/
func (c *Client) getTasksPath(taskListName string) string {
	username := c.config.Nextcloud.Username
	if taskListName == "" {
		taskListName = "tasks" // デフォルトタスクリスト
	}
	return fmt.Sprintf("/remote.php/dav/calendars/%s/%s/", username, taskListName)
}

// getContext はタイムアウト付きのコンテキストを返すます。
func (c *Client) getContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 30*time.Second)
}
