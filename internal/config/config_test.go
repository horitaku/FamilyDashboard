package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestLoadConfigSuccess は正しい設定ファイルの読み込みをテストするます。
func TestLoadConfigSuccess(t *testing.T) {
	// テンポラリディレクトリ作成
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "settings.json")

	// テスト用の正しい設定ファイルを作成
	content := `{
		"refreshIntervals": {
			"weatherSec": 300,
			"calendarSec": 300,
			"tasksSec": 300
		},
		"location": {
			"cityName": "姫路市",
			"country": "JP"
		},
		"nextcloud": {
			"serverUrl": "https://nextcloud.example.com",
			"username": "testuser",
			"password": "testpass",
			"calendarNames": ["family"],
			"taskListNames": ["tasks"]
		},
		"weather": {
			"provider": "",
			"apiKey": "",
			"baseUrl": ""
		}
	}`

	err := os.WriteFile(configFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("テスト用設定ファイルの作成に失敗しました: %v", err)
	}

	// 設定ファイルを読み込む
	cfg, err := LoadConfig(configFile)
	if err != nil {
		t.Fatalf("設定ファイルの読み込みに失敗しました: %v", err)
	}

	// 値の検証
	if cfg.RefreshIntervals.WeatherSec != 300 {
		t.Errorf("weatherSec の値が一致しません。期待値：300、実際：%d", cfg.RefreshIntervals.WeatherSec)
	}
	if cfg.Location.CityName != "姫路市" {
		t.Errorf("cityName の値が一致しません。期待値：姫路市、実際：%s", cfg.Location.CityName)
	}
	if cfg.Location.Country != "JP" {
		t.Errorf("country の値が一致しません。期待値：JP、実際：%s", cfg.Location.Country)
	}

	// LoadedAt が設定されているか確認
	if cfg.LoadedAt().IsZero() {
		t.Error("LoadedAt が設定されていません")
	}
}

// TestLoadConfigFileNotFound はファイルが存在しない場合のエラーをテストするます。
func TestLoadConfigFileNotFound(t *testing.T) {
	cfg, err := LoadConfig("/nonexistent/path/settings.json")
	if err == nil {
		t.Errorf("エラーが発生すべきですが、成功しました。cfg=%v", cfg)
	}
}

// TestLoadConfigInvalidJSON は無効なJSONファイルのエラーをテストするます。
func TestLoadConfigInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "settings.json")

	// 不正なJSON
	content := `{ invalid json }`
	err := os.WriteFile(configFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("テスト用設定ファイルの作成に失敗しました: %v", err)
	}

	cfg, err := LoadConfig(configFile)
	if err == nil {
		t.Errorf("JSONパースエラーが発生すべきですが、成功しました。cfg=%v", cfg)
	}
}

// TestValidateRefreshIntervalsInvalid は無効な更新間隔のバリデーションテストです。
func TestValidateRefreshIntervalsInvalid(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "weatherSec が0以下",
			config: &Config{
				RefreshIntervals: RefreshIntervals{WeatherSec: 0, CalendarSec: 300, TasksSec: 300},
				Location:         Location{CityName: "姫路市", Country: "JP"},
			},
			wantErr: true,
		},
		{
			name: "calendarSec が負数",
			config: &Config{
				RefreshIntervals: RefreshIntervals{WeatherSec: 300, CalendarSec: -1, TasksSec: 300},
				Location:         Location{CityName: "姫路市", Country: "JP"},
			},
			wantErr: true,
		},
		{
			name: "tasksSec が0以下",
			config: &Config{
				RefreshIntervals: RefreshIntervals{WeatherSec: 300, CalendarSec: 300, TasksSec: 0},
				Location:         Location{CityName: "姫路市", Country: "JP"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("バリデーション結果が一致しません。期待エラー：%v、実際エラー：%v", tt.wantErr, err)
			}
		})
	}
}

// TestValidateLocationInvalid は無効なロケーション設定のバリデーションテストです。
func TestValidateLocationInvalid(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "cityName が空",
			config: &Config{
				RefreshIntervals: RefreshIntervals{WeatherSec: 300, CalendarSec: 300, TasksSec: 300},
				Location:         Location{CityName: "", Country: "JP"},
			},
			wantErr: true,
		},
		{
			name: "country が空",
			config: &Config{
				RefreshIntervals: RefreshIntervals{WeatherSec: 300, CalendarSec: 300, TasksSec: 300},
				Location:         Location{CityName: "姫路市", Country: ""},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("バリデーション結果が一致しません。期待エラー：%v、実際エラー：%v", tt.wantErr, err)
			}
		})
	}
}

// TestGetRefreshInterval は GetRefreshInterval メソッドのテストです。
func TestGetRefreshInterval(t *testing.T) {
	cfg := &Config{
		RefreshIntervals: RefreshIntervals{
			WeatherSec:  300,
			CalendarSec: 600,
			TasksSec:    450,
		},
	}

	tests := []struct {
		source   string
		expected time.Duration
	}{
		{"weather", 300 * time.Second},
		{"calendar", 600 * time.Second},
		{"tasks", 450 * time.Second},
		{"unknown", 5 * time.Minute}, // 既定値
	}

	for _, tt := range tests {
		got := cfg.GetRefreshInterval(tt.source)
		if got != tt.expected {
			t.Errorf("GetRefreshInterval(%q) = %v、期待値：%v", tt.source, got, tt.expected)
		}
	}
}

// TestGetLocationString は GetLocationString メソッドのテストです。
func TestGetLocationString(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		expected string
	}{
		{
			name:     "正常系",
			config:   &Config{Location: Location{CityName: "姫路市", Country: "JP"}},
			expected: "姫路市, JP",
		},
		{
			name:     "cityName が空",
			config:   &Config{Location: Location{CityName: "", Country: "JP"}},
			expected: "Unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.config.GetLocationString()
			if got != tt.expected {
				t.Errorf("GetLocationString() = %q、期待値：%q", got, tt.expected)
			}
		})
	}
}
