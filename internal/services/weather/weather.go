package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/rihow/FamilyDashboard/internal/cache"
	"github.com/rihow/FamilyDashboard/internal/models"
)

// Client は天気APIクライアントなのです。
// 気象庁データまたはOpen-Meteo APIを使用して天気情報を取得するます。
type Client struct {
	// 気象庁パラメーター機能
	baseURL    string
	httpClient *http.Client
	fc         *cache.FileCache // キャッシュ機能
	geocodeURL string           // ジオコーディング情報を取得するバックエンド URL (例: http://localhost:8080)
}

// OpenMeteoWeatherResponse は Open-Meteo API のレスポンス構造体なのです。
// 気象庁データをラップしているため、日本の天気予報データを取得できるます。
type OpenMeteoWeatherResponse struct {
	Latitude  float64              `json:"latitude"`
	Longitude float64              `json:"longitude"`
	Current   OpenMeteoWeatherData `json:"current"`
	Daily     OpenMeteoDailyData   `json:"daily"`
	Hourly    OpenMeteoHourlyData  `json:"hourly"`
}

// OpenMeteoWeatherData は Open-Meteo API の現在データなのです。
type OpenMeteoWeatherData struct {
	Temperature      float64 `json:"temperature"`
	RelativeHumidity int     `json:"relative_humidity"`
	WindSpeed        float64 `json:"wind_speed_10m"`
	WeatherCode      int     `json:"weather_code"`
	Time             string  `json:"time"`
}

// OpenMeteoDailyData は Open-Meteo API の日別予報なのです。
type OpenMeteoDailyData struct {
	Time              []string  `json:"time"`
	MaxTemperature    []float64 `json:"temperature_2m_max"`
	MinTemperature    []float64 `json:"temperature_2m_min"`
	PrecipitationProb []int     `json:"precipitation_probability_max"`
	WeatherCode       []int     `json:"weather_code"`
}

// OpenMeteoHourlyData は Open-Meteo API の時間別降水確率なのです。
type OpenMeteoHourlyData struct {
	Time              []string `json:"time"`
	PrecipitationProb []int    `json:"precipitation_probability"`
}

// NewClient は天気APIクライアントを作成するます。
// geocodeURL はこのサーバー自身の URL (例: http://localhost:8080) です。
// 緯度経度を取得するために使用するます。
func NewClient(fc *cache.FileCache, geocodeURL string) *Client {
	return &Client{
		baseURL: "https://api.open-meteo.com/v1/forecast",
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		fc:         fc,
		geocodeURL: geocodeURL,
	}
}

// GetWeather は 指定都市の天気情報を取得するます。
// キャッシュをリスク判定して、有効な場合はそれを返します。
// 無効な場合は Open-Meteo API から取得して保存するます。
func (c *Client) GetWeather(ctx context.Context, cityName, country string) (*models.WeatherResponse, error) {
	cacheKey := fmt.Sprintf("weather:%s:%s", country, cityName)
	ttl := 5 * time.Minute // デフォルト5分

	// キャッシュを読み込もうするます
	var cachedWeather models.WeatherResponse
	_, found, stale, err := c.fc.ReadPayload(cacheKey, ttl, &cachedWeather)
	if found && !stale {
		// キャッシュが有効な場合は返すます
		return &cachedWeather, nil
	}

	// 緯度経度を取得するます（キャッシュ済み含む）
	coords, err := c.getCoordinates(ctx, cityName, country)
	if err != nil {
		return nil, fmt.Errorf("緯度経度取得失敗するます: %w", err)
	}

	// Open-Meteo APIから天気データを取得するます
	weatherRsp, err := c.fetchFromOpenMeteo(ctx, coords.Latitude, coords.Longitude, cityName)
	if err != nil {
		return nil, err
	}

	// キャッシュに保存するます
	_, _ = c.fc.Write(cacheKey, weatherRsp, map[string]string{
		"city":    cityName,
		"country": country,
		"source":  "open-meteo",
	})

	return weatherRsp, nil
}

// getCoordinates はバックエンドのジオコーディング機能を使って緯度経度を取得するます。
func (c *Client) getCoordinates(ctx context.Context, cityName, country string) (*geocodeResult, error) {
	// バックエンドの /api/geocode エンドポイントを呼ぶます
	// ここではシンプルに直接 Open-Meteo の geocoding API を使うこともできるます
	url := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&country=%s&count=1&language=ja",
		cityName, country)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成失敗するます: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("APIリクエスト失敗するます: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("APIレスポンスエラー: code=%d, body=%s", resp.StatusCode, string(body))
	}

	var result struct {
		Results []struct {
			Latitude  float64 `json:"latitude"`
			Longitude float64 `json:"longitude"`
			Name      string  `json:"name"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("レスポンスパース失敗するます: %w", err)
	}

	if len(result.Results) == 0 {
		return nil, fmt.Errorf("都市 %s が見つからないます", cityName)
	}

	return &geocodeResult{
		Latitude:  result.Results[0].Latitude,
		Longitude: result.Results[0].Longitude,
	}, nil
}

// fetchFromOpenMeteo は Open-Meteo API から天気データを取得するます。
// 気象庁データベースが統合されているため、日本の天気データも取得できるます。
func (c *Client) fetchFromOpenMeteo(ctx context.Context, lat, lon float64, cityName string) (*models.WeatherResponse, error) {
	// Open-Meteo API リクエストを構築するます
	url := fmt.Sprintf(
		"%s?latitude=%.2f&longitude=%.2f&current=temperature_2m,relative_humidity_2m,weather_code,wind_speed_10m&daily=weather_code,temperature_2m_max,temperature_2m_min,precipitation_probability_max&hourly=precipitation_probability&timezone=Asia/Tokyo&forecast_days=7",
		c.baseURL, lat, lon,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("リクエスト作成失敗するます: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Open-Meteo APIリクエスト失敗するます: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("APIエラー: code=%d, body=%s", resp.StatusCode, string(body))
	}

	var omResp OpenMeteoWeatherResponse
	if err := json.NewDecoder(resp.Body).Decode(&omResp); err != nil {
		return nil, fmt.Errorf("レスポンスパース失敗するます: %w", err)
	}

	// Open-Meteo のレスポンスを models.WeatherResponse に変換するます
	weatherRsp := c.convertToWeatherResponse(&omResp, cityName)
	return weatherRsp, nil
}

// convertToWeatherResponse は Open-Meteo レスポンスを models.WeatherResponse に変換するます。
func (c *Client) convertToWeatherResponse(omResp *OpenMeteoWeatherResponse, cityName string) *models.WeatherResponse {
	condition := weatherCodeToCondition(omResp.Current.WeatherCode)
	icon := weatherCodeToIcon(omResp.Current.WeatherCode)

	// 現在の天候
	current := models.CurrentWeather{
		Temperature: omResp.Current.Temperature,
		Condition:   condition,
		Icon:        icon,
		Humidity:    omResp.Current.RelativeHumidity,
		WindSpeed:   omResp.Current.WindSpeed,
	}

	// 今日の天況
	today := models.TodayWeather{
		MaxTemp: omResp.Daily.MaxTemperature[0],
		MinTemp: omResp.Daily.MinTemperature[0],
		Summary: condition,
	}

	// 時間帯ごとの降水確率を取得するます（最大6時間分）
	precipSlots := []models.PrecipSlot{}
	for i := 0; i < len(omResp.Hourly.Time) && i < 24; i++ {
		// 3時間ごとにサンプリングするます
		if i%3 == 0 {
			precipSlots = append(precipSlots, models.PrecipSlot{
				Time:   fmt.Sprintf("%02d:00", i),
				Precip: omResp.Hourly.PrecipitationProb[i],
			})
		}
	}

	// 注意報・警報はここでは空にするます
	// （Open-Meteo では警報提供がないため、後で別途実装するます）
	alerts := []models.WeatherAlert{}

	return &models.WeatherResponse{
		Location:    cityName,
		Current:     current,
		Today:       today,
		PrecipSlots: precipSlots,
		Alerts:      alerts,
	}
}

// weatherCodeToCondition は WMO天気コードを日本語の気象情報に変換するます。
func weatherCodeToCondition(code int) string {
	switch code {
	case 0:
		return "晴"
	case 1, 2:
		return "曇"
	case 3:
		return "曇"
	case 45, 48:
		return "霧"
	case 51, 53, 55:
		return "小雨"
	case 61, 63, 65:
		return "雨"
	case 71, 73, 75:
		return "雪"
	case 77:
		return "吹雪"
	case 80, 81, 82:
		return "激しい雨"
	case 85, 86:
		return "にわか雨"
	case 95, 96, 99:
		return "雷雨"
	default:
		return "天候不明"
	}
}

// weatherCodeToIcon は WMO天気コードをアイコンコードに変換するます。
func weatherCodeToIcon(code int) string {
	switch code {
	case 0:
		return "01d" // 晴
	case 1:
		return "02d" // ほぼ晴
	case 2, 3:
		return "03d" // 曇
	case 45, 48:
		return "50d" // 霧
	case 51, 53, 55:
		return "09d" // 小雨
	case 61, 63, 65:
		return "10d" // 雨
	case 71, 73, 75:
		return "13d" // 雪
	case 77:
		return "14d" // 吹雪
	case 80, 81, 82:
		return "11d" // 激しい雨
	case 85, 86:
		return "12d" // にわか雨
	case 95, 96, 99:
		return "15d" // 雷雨
	default:
		return "04u" // 不明
	}
}

// geocodeResult はジオコーディング結果なのです。
type geocodeResult struct {
	Latitude  float64
	Longitude float64
}
