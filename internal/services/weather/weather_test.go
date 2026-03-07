package weather

import (
	"testing"
	"time"

	"github.com/rihow/FamilyDashboard/internal/cache"
)

// TestConvertToWeatherResponse は気象コード変換テストなのです。
func TestConvertToWeatherResponse(t *testing.T) {
	fc := cache.New("./data/cache")
	c := NewClient(fc, "http://localhost:8080")

	jst := time.FixedZone("Asia/Tokyo", 9*3600)
	now := time.Now().In(jst)
	today := now.Format("2006-01-02")

	hourlyTimes := make([]string, 0, 8)
	hourlyPrecip := make([]int, 0, 8)
	for i := 1; i <= 8; i++ {
		t := now.Add(time.Duration(i*3) * time.Hour)
		hourlyTimes = append(hourlyTimes, t.Format("2006-01-02T15:04"))
		hourlyPrecip = append(hourlyPrecip, 10)
	}

	// ダミー Open-Meteo レスポンスを作成するます
	dummyResp := &OpenMeteoWeatherResponse{
		Latitude:  34.815,
		Longitude: 134.685,
		Current: OpenMeteoWeatherData{
			Temperature:      15.5,
			RelativeHumidity: 60,
			WindSpeed:        3.2,
			WeatherCode:      2,
			Time:             time.Now().Format(time.RFC3339),
		},
		Daily: OpenMeteoDailyData{
			Time:              []string{today},
			MaxTemperature:    []float64{18.0},
			MinTemperature:    []float64{10.0},
			PrecipitationProb: []int{15},
			WeatherCode:       []int{2},
		},
		Hourly: OpenMeteoHourlyData{
			Time:              hourlyTimes,
			PrecipitationProb: hourlyPrecip,
		},
	}

	// 変換を実行するます
	result := c.convertToWeatherResponse(dummyResp, "姫路市")

	// アサーションをするます
	if result.Location != "姫路市" {
		t.Errorf("Location: 期待: 姫路市, 実際: %s", result.Location)
	}

	if result.Current.Temperature != 15.5 {
		t.Errorf("Temperature: 期待: 15.5, 実際: %f", result.Current.Temperature)
	}

	if result.Current.Condition != "曇" {
		t.Errorf("Condition: 期待: 曇, 実際: %s", result.Current.Condition)
	}

	if result.Today.MaxTemp != 18.0 {
		t.Errorf("MaxTemp: 期待: 18.0, 実際: %f", result.Today.MaxTemp)
	}

	if result.Today.MinTemp != 10.0 {
		t.Errorf("MinTemp: 期待: 10.0, 実際: %f", result.Today.MinTemp)
	}

	if len(result.PrecipSlots) == 0 {
		t.Errorf("PrecipSlots: 期待: > 0, 実際: %d", len(result.PrecipSlots))
	}
}

// TestWeatherCodeToCondition は天気コード変換テストなのです。
func TestWeatherCodeToCondition(t *testing.T) {
	tests := []struct {
		code     int
		expected string
	}{
		{0, "晴"},
		{1, "曇"},
		{2, "曇"},
		{3, "曇"},
		{45, "霧"},
		{48, "霧"},
		{51, "小雨"},
		{61, "雨"},
		{71, "雪"},
		{95, "雷雨"},
	}

	for _, test := range tests {
		result := weatherCodeToCondition(test.code)
		if result != test.expected {
			t.Errorf("code %d: 期待: %s, 実際: %s", test.code, test.expected, result)
		}
	}
}

// TestWeatherCodeToIcon はアイコンコード変換テストなのです。
func TestWeatherCodeToIcon(t *testing.T) {
	tests := []struct {
		code     int
		expected string
	}{
		{0, "01d"},
		{1, "02d"},
		{2, "03d"},
		{45, "50d"},
		{61, "10d"},
		{71, "13d"},
		{95, "15d"},
	}

	for _, test := range tests {
		result := weatherCodeToIcon(test.code)
		if result != test.expected {
			t.Errorf("code %d: 期待アイコン: %s, 実際: %s", test.code, test.expected, result)
		}
	}
}
