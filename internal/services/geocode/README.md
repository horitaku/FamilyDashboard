# geocode パッケージ

都市名から座標（緯度・経度）を取得するジオコーディング機能を提供するパッケージなのです。

## 機能

- **Nominatim API統合**: OpenStreetMapのNominatim APIを使用した無料ジオコーディング
- **キャッシング**: ジオコーディング結果をJSONファイルにキャッシュ（5分TTL）
- **User-Agent/Referer対応**: Nominatim利用規約に準拠した適切なヘッダー送信
- **エラーハンドリング**: ネットワークエラーや不正なレスポンスへの対応

## 使用例

```go
package main

import (
    "context"
    "log"
    "github.com/rihow/FamilyDashboard/internal/cache"
    "github.com/rihow/FamilyDashboard/internal/services/geocode"
)

func main() {
    // キャッシュマネージャーを初期化するます
    cacheMgr := cache.NewManager("./data/cache")
    
    // ジオコーディングクライアントを作成するます
    client := geocode.NewClient(cacheMgr)
    
    // 都市名から座標を取得するます
    ctx := context.Background()
    location, err := client.GetCoordinates(ctx, "姫路市", "JP")
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("緯度: %f, 経度: %f\n", location.Latitude, location.Longitude)
}
```

## Nominatim 利用規約

- **1秒あたり最大1リクエスト**: アプリ側でリクエストを制限すること
- **User-Agent送信**: 識別可能なUser-Agentを必須で送信
- **Referer送信**: リクエスト元のURLを送信
- **帰属表示**: UIに「© OpenStreetMap contributors」と表示

## キャッシュキー

ジオコーディング結果は以下の形式でキャッシュされます：
- キー: `geocode_{cityName}_{country}` （例: `geocode_姫路市_JP`）
- 値: `Location` 構造体をJSON形式で保存
- TTL: 5分

## テスト

```bash
# Nominatim APIを実際に呼ぶテストを実行
go test -v ./internal/services/geocode
```

**注意**: テストはネットワーク接続が必要です。また、Nominatim利用規約に従い、テスト間隔を十分に空けてください。

## 今後の改善

- セルフホストのジオコーダー検討（利用が増えた場合）
- キャッシュの永続化最適化
- 複数都市の同時クエリー対応
