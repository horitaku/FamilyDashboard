# weather

天気APIクライアントを実装するパッケージなのです。

## 機能

- **Open-Meteo API インテグレーション**: 気象庁データを含む天気予報を取得するます
- **ジオコーディング**: 都市名を緯度経度に変換（Open-Meteo geocoding API使用）
- **データ変換**: WMO天気コード→日本語条件・アイコン変換
- **キャッシュ管理**: 天気データのキャッシュ保存・有効期限管理（TTL: 5分）
- **エラーハンドリング**: ネットワーク障害時のエラー処理

## 実装詳細

### Client 構造体

```go
type Client struct {
    baseURL    string           // Open-Meteo API ベースURL
    httpClient *http.Client     // HTTP クライアント（タイムアウト10秒）
    fc         *cache.FileCache // キャッシュ管理
    geocodeURL string           // ジオコーディング用バックエンド URL
}
```

### GetWeather(ctx context.Context, cityName, country string) 関数

1. キャッシュをチェック（有効期限内なら返す）
2. 緯度経度を取得（getCoordinates）
3. Open-Meteo API から天気データ取得（fetchFromOpenMeteo）
4. models.WeatherResponse に変換（convertToWeatherResponse）
5. キャッシュに保存

### WMO 天気コード変換

- 0: 晴
- 1, 2, 3: 曇
- 45, 48: 霧
- 51, 53, 55: 小雨
- 61, 63, 65: 雨
- 71, 73, 75: 雪
- 95, 96, 99: 雷雨

## API レスポンス例

```json
{
  "location": "姫路市",
  "current": {
    "temperature": 15.5,
    "condition": "曇",
    "icon": "03d",
    "humidity": 65,
    "windSpeed": 3.2
  },
  "today": {
    "maxTemp": 20.0,
    "minTemp": 10.5,
    "summary": "曇"
  },
  "precipSlots": [
    {"time": "09:00", "precip": 10},
    {"time": "12:00", "precip": 5},
    {"time": "15:00", "precip": 0}
  ],
  "alerts": []
}
```

## テスト

- `TestConvertToWeatherResponse`: Open-Meteo レスポンス→モデル変換テスト
- `TestWeatherCodeToCondition`: 天気コード→日本語変換テスト
- `TestWeatherCodeToIcon`: 天気コード→アイコン変換テスト

## 注意事項
- Open-Meteo API は登録不要で無料
- レート制限: 推奨は毎秒10リクエスト以下
- キャッシュ検証のため、データ取得時は Asia/Tokyo タイムゾーン を使用
