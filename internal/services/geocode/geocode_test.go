package geocode

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/rihow/FamilyDashboard/internal/cache"
)

// TestGetCoordinates_Nominatim はNominatim APIへのクエリーテストです。
// 注意: このテストは実際にNominatim APIを呼び出します。
// テスト実行時にはネットワーク接続が必要です。
// また、Nominatim利用規約に従い、テスト間隔を十分に空けてください。
func TestGetCoordinates_Nominatim(t *testing.T) {
	tmpDir := t.TempDir()
	fc := cache.New(tmpDir)
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	client := NewClient(fc)
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// 姫路市で座標取得テスト
	t.Run("get_himeji_coordinates", func(t *testing.T) {
		location, err := client.GetCoordinates(ctx, "姫路市", "JP")
		if err != nil {
			t.Fatalf("failed to get coordinates: %v", err)
		}

		if location == nil {
			t.Fatal("location is nil")
		}

		// 姫路市の座標は約（北緯34.8度、東経134.7度）
		if location.Latitude < 34 || location.Latitude > 36 {
			t.Errorf("unexpected latitude: %f", location.Latitude)
		}
		if location.Longitude < 134 || location.Longitude > 136 {
			t.Errorf("unexpected longitude: %f", location.Longitude)
		}

		t.Logf("✨ Himeji coordinates: lat=%f, lon=%f", location.Latitude, location.Longitude)
	})

	// 2回目のクエリーはキャッシュから返ってくるはずです。
	t.Run("cache_hit", func(t *testing.T) {
		// 最初のクエリー（キャッシュミス）
		loc1, err := client.GetCoordinates(ctx, "京都市", "JP")
		if err != nil {
			t.Fatalf("first query failed: %v", err)
		}

		// 2回目のクエリー（キャッシュヒット）
		loc2, err := client.GetCoordinates(ctx, "京都市", "JP")
		if err != nil {
			t.Fatalf("second query failed: %v", err)
		}

		if loc1.Latitude != loc2.Latitude || loc1.Longitude != loc2.Longitude {
			t.Errorf("cached result mismatch: %f,%f vs %f,%f", loc1.Latitude, loc1.Longitude, loc2.Latitude, loc2.Longitude)
		}

		t.Logf("✨ Cache hit confirmed: lat=%f, lon=%f", loc2.Latitude, loc2.Longitude)
	})

	// 存在しない都市でのエラーハンドリング
	t.Run("invalid_city", func(t *testing.T) {
		_, err := client.GetCoordinates(ctx, "ぜったいないまち12345", "JP")
		if err == nil {
			t.Fatal("expected error for invalid city, but got nil")
		}

		t.Logf("✨ Error handling works correctly: %v", err)
	})
}

// TestLocation_Serialization はLocation構造体のJSONシリアライゼーションテストです
func TestLocation_Serialization(t *testing.T) {
	loc := &Location{
		Latitude:  34.8,
		Longitude: 134.7,
		CityName:  "姫路市",
		Country:   "JP",
		FetchedAt: time.Now(),
	}

	// Marshalテスト
	data, err := json.Marshal(loc)
	if err != nil {
		t.Fatalf("marshal failed: %v", err)
	}

	// Unmarshalテスト
	var unmarshalled Location
	if err := json.Unmarshal(data, &unmarshalled); err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if unmarshalled.Latitude != loc.Latitude || unmarshalled.Longitude != loc.Longitude {
		t.Errorf("coordinates mismatch after roundtrip")
	}

	t.Logf("✨ Serialization works correctly")
}
