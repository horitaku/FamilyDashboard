package geocode

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/rihow/FamilyDashboard/internal/cache"
)

// Location は都市のジオコーディング結果を表す構造体なのです。
type Location struct {
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	CityName  string    `json:"city_name"`
	Country   string    `json:"country"`
	FetchedAt time.Time `json:"fetched_at"`
}

// NominatimResponse は Nominatim APIのレスポンス構造体なのです。
type NominatimResponse struct {
	Lat         string `json:"lat"`
	Lon         string `json:"lon"`
	DisplayName string `json:"display_name"`
}

// Client はジオコーディングクライアントなのです。
type Client struct {
	baseURL    string
	userAgent  string
	httpClient *http.Client
	fc         *cache.FileCache
}

// NewClient はジオコーディングクライアントを作成するます。
// fcはジオコーディング結果のキャッシュを管理するためのFileCacheなのです。
// TTLは5分（ジオコーディング結果は頻繁に変わらないため、長めのTTL）。
func NewClient(fc *cache.FileCache) *Client {
	return &Client{
		baseURL:   "https://nominatim.openstreetmap.org",
		userAgent: "FamilyDashboard/1.0 (https://github.com/rihow/FamilyDashboard; personal-use; @rihow)",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		fc: fc,
	}
}

// GetCoordinates は都市名（と国コード）から緯度経度を取得するます。
// キャッシュがあればそれを返し、ない場合はNominatim APIを呼ぶのです。
func (c *Client) GetCoordinates(ctx context.Context, cityName, country string) (*Location, error) {
	// キャッシュキーを生成するます。
	cacheKey := fmt.Sprintf("geocode_%s_%s", cityName, country)

	// キャッシュをチェックするます。（TTL: 5分）
	if entry, exists, stale, err := c.fc.Read(cacheKey, 5*time.Minute); err == nil && exists && !stale {
		var loc Location
		if err := json.Unmarshal(entry.Payload, &loc); err == nil {
			return &loc, nil
		}
		// キャッシュのデシリアライズに失敗した場合は、APIを呼ぶます。
	}

	// Nominatim APIを呼ぶます。
	location, err := c.queryNominatim(ctx, cityName, country)
	if err != nil {
		return nil, err
	}

	// 結果をキャッシュに保存するます。
	if data, err := json.Marshal(location); err == nil {
		_, _ = c.fc.Write(cacheKey, json.RawMessage(data), map[string]string{
			"cityName": cityName,
			"country":  country,
		})
	}

	return location, nil
}

// queryNominatim はNominatim APIにクエリーを送り、座標を取得するます。
// Nominatim利用規約に従い、1秒あたり最大1リクエストのレート制限しぶりを想定するます。
// （実装アプリ側でリクエストを1秒以上間隔を空ける必要があります）
func (c *Client) queryNominatim(ctx context.Context, cityName, country string) (*Location, error) {
	// クエリーを構築するます。
	// format=json: JSON形式で返す
	// limit=1: 最初の1件のみ返す
	// countrycodes: 国コードでフィルタ（JP など）
	baseURL := c.baseURL + "/search"
	params := url.Values{}
	params.Set("q", cityName)
	params.Set("countrycodes", country)
	params.Set("format", "json")
	params.Set("limit", "1")

	url := baseURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("request creation failed: %w", err)
	}

	// User-Agent と Referer を設定するます。（Nominatim利用規約対応）
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Referer", "https://familydashboard.local")

	// リクエストを送るます。
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("nominatim error: status %d, body: %s", resp.StatusCode, string(body))
	}

	// レスポンスをパースするます。
	var results []NominatimResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("json decode failed: %w", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no results found for %s, %s", cityName, country)
	}

	// 最初の結果を使用するます。
	result := results[0]

	// 文字列から float64 に変換するます。
	var lat, lon float64
	if _, err := fmt.Sscanf(result.Lat, "%f", &lat); err != nil {
		return nil, fmt.Errorf("latitude parse failed: %w", err)
	}
	if _, err := fmt.Sscanf(result.Lon, "%f", &lon); err != nil {
		return nil, fmt.Errorf("longitude parse failed: %w", err)
	}

	location := &Location{
		Latitude:  lat,
		Longitude: lon,
		CityName:  cityName,
		Country:   country,
		FetchedAt: time.Now(),
	}

	return location, nil
}
