package google

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/rihow/FamilyDashboard/internal/cache"
	"github.com/rihow/FamilyDashboard/internal/config"
)

// Client は Google API（Calendar/Tasks）のクライアントです。
// OAuth認証・データ取得・キャッシュを管理するのです。
type Client struct {
	cache          *cache.FileCache
	config         *config.Config
	accessToken    string
	refreshToken   string
	tokenExpiresAt time.Time
	httpClient     *http.Client
}

// NewClient は Google APIクライアントを初期化します。
// configからトークンを読み込み、キャッシュを設定するます。
func NewClient(fc *cache.FileCache, cfg *config.Config) *Client {
	return &Client{
		cache:      fc,
		config:     cfg,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

// SetAccessToken はアクセストークンを設定します。
// 本番環境ではOAuth認可コードフローで取得するのです。
func (c *Client) SetAccessToken(token string, expiresIn int) {
	c.accessToken = token
	c.tokenExpiresAt = time.Now().Add(time.Duration(expiresIn) * time.Second)
}

// SetRefreshToken はリフレッシュトークンを設定します。
// トークン失効時に再発行するのです。
func (c *Client) SetRefreshToken(token string) {
	c.refreshToken = token
}

// IsTokenValid はアクセストークンが有効かどうかをチェックします。
func (c *Client) IsTokenValid() bool {
	if c.accessToken == "" {
		return false
	}
	return time.Now().Before(c.tokenExpiresAt)
}

// doRequest は認可ヘッダー付きのHTTPリクエストを実行します。
// エラー時はキャッシュを返す仕様で対応するのです。
func (c *Client) doRequest(ctx context.Context, method, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成エラー: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))
	req.Header.Set("User-Agent", "FamilyDashboard/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP実行エラー: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンス読込エラー: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("APIエラー（ステータス %d）: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// parseJSONResponse はJSONレスポンスをパースして、構造体に変換します。
func parseJSONResponse(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return fmt.Errorf("JSON parse error: %w", err)
	}
	return nil
}

// OAuthAuthorizationCodeFlow は OAuth 2.0 認可コードフローを実行するのです。
// ユーザーから受け取った認可コードを使用して、Googleからアクセストークンを取得するます。
func (c *Client) OAuthAuthorizationCodeFlow(ctx context.Context, authCode string) (*TokenResponse, error) {
	// 設定の検証
	if c.config.Google.ClientID == "" || c.config.Google.ClientSecret == "" {
		return nil, fmt.Errorf("Google OAuth設定がまだ設定されていないのです。settings.json を確認するます")
	}

	// Google OAuth Token Endpoint にPOST
	tokenURL := "https://oauth2.googleapis.com/token"

	// リクエストボディを作成
	data := url.Values{}
	data.Set("code", authCode)
	data.Set("client_id", c.config.Google.ClientID)
	data.Set("client_secret", c.config.Google.ClientSecret)
	data.Set("redirect_uri", c.config.Google.RedirectUri)
	data.Set("grant_type", "authorization_code")

	// HTTPリクエスト実行
	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成エラー: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "FamilyDashboard/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP実行エラー: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("レスポンス読込エラー: %w", err)
	}

	// ステータスコード確認
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Google OAuth エラー（ステータス %d）: %s", resp.StatusCode, string(body))
	}

	// JSONレスポンスをパース
	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("レスポンス解析エラー: %w", err)
	}

	// トークンをクライアントに設定
	c.SetAccessToken(tokenResp.AccessToken, tokenResp.ExpiresIn)
	c.SetRefreshToken(tokenResp.RefreshToken)

	// トークンをファイルに保存
	if err := c.SaveTokens("./data/tokens.json"); err != nil {
		return nil, fmt.Errorf("トークン保存エラー: %w", err)
	}

	return &tokenResp, nil
}

// TokenResponse はGoogleOAuthのトークンレスポンスです。
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

// useCache はキャッシュキーとTTLを使用して、キャッシュの有効性を判定するのです。
// キャッシュが有効な場合は、キャッシュされたペイロードをパースして返すます。
func (c *Client) useCache(cacheKey string) ([]byte, bool, error) {
	ttl := c.config.GetRefreshInterval("calendar") // 天気/カレンダー/タスクは同じTTLを使用
	entry, exists, stale, err := c.cache.Read(cacheKey, ttl)
	if err != nil {
		return nil, false, err
	}

	if !exists || stale {
		return nil, false, nil
	}

	return entry.Payload, true, nil
}

// saveCache はデータをキャッシュに保存するのです。
func (c *Client) saveCache(cacheKey string, data []byte) error {
	meta := map[string]string{
		"source": "google",
	}
	_, err := c.cache.Write(cacheKey, json.RawMessage(data), meta)
	return err
}

// TokenFile は保存されたトークン情報を表すのです。
type TokenFile struct {
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	TokenExpiresAt time.Time `json:"token_expires_at"`
	SavedAt        time.Time `json:"saved_at"`
}

// SaveTokens はアクセストークン・リフレッシュトークンをファイルに保存するのです。
func (c *Client) SaveTokens(tokenFilePath string) error {
	// ディレクトリが存在するか確認
	dir := filepath.Dir(tokenFilePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("ディレクトリ作成エラー: %w", err)
	}

	// TokenFile を作成
	tf := TokenFile{
		AccessToken:    c.accessToken,
		RefreshToken:   c.refreshToken,
		TokenExpiresAt: c.tokenExpiresAt,
		SavedAt:        time.Now(),
	}

	// JSONに変換
	data, err := json.MarshalIndent(tf, "", "  ")
	if err != nil {
		return fmt.Errorf("JSON変換エラー: %w", err)
	}

	// ファイルに書き込み（権限600）
	if err := os.WriteFile(tokenFilePath, data, 0600); err != nil {
		return fmt.Errorf("ファイル書き込みエラー: %w", err)
	}

	return nil
}

// LoadTokens はファイルからトークンを読み込むのです。
func (c *Client) LoadTokens(tokenFilePath string) error {
	// ファイルが存在するか確認
	data, err := os.ReadFile(tokenFilePath)
	if err != nil {
		// ファイルが存在しない場合は、エラーではなく単に何もしない
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("ファイル読込エラー: %w", err)
	}

	// JSONをパース
	var tf TokenFile
	if err := json.Unmarshal(data, &tf); err != nil {
		return fmt.Errorf("JSON解析エラー: %w", err)
	}

	// トークンをクライアントに設定
	c.accessToken = tf.AccessToken
	c.refreshToken = tf.RefreshToken
	c.tokenExpiresAt = tf.TokenExpiresAt

	return nil
}

// RefreshAccessToken はリフレッシュトークンを使用して新しいアクセストークンを取得するのです。
func (c *Client) RefreshAccessToken(ctx context.Context) error {
	if c.refreshToken == "" {
		return fmt.Errorf("リフレッシュトークンが設定されていないのです")
	}

	tokenURL := "https://oauth2.googleapis.com/token"

	// リクエストボディを作成
	data := url.Values{}
	data.Set("client_id", c.config.Google.ClientID)
	data.Set("client_secret", c.config.Google.ClientSecret)
	data.Set("refresh_token", c.refreshToken)
	data.Set("grant_type", "refresh_token")

	// HTTPリクエスト実行
	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return fmt.Errorf("リクエスト作成エラー: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "FamilyDashboard/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP実行エラー: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("レスポンス読込エラー: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("トークンリフレッシュエラー（ステータス %d）: %s", resp.StatusCode, string(body))
	}

	// JSONレスポンスをパース
	var tokenResp TokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return fmt.Errorf("レスポンス解析エラー: %w", err)
	}

	// トークンを更新
	c.SetAccessToken(tokenResp.AccessToken, tokenResp.ExpiresIn)

	// リフレッシュトークンが返された場合は更新
	if tokenResp.RefreshToken != "" {
		c.SetRefreshToken(tokenResp.RefreshToken)
	}

	// 更新したトークンをファイルに保存
	if err := c.SaveTokens("./data/tokens.json"); err != nil {
		fmt.Printf("⚠️ トークン保存エラー: %v\n", err)
		// エラーでも継続する（トークンはメモリに保持されている）
	}

	return nil
}
