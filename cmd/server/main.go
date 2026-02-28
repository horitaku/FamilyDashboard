package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rihow/FamilyDashboard/internal/cache"
	"github.com/rihow/FamilyDashboard/internal/config"
	httproutes "github.com/rihow/FamilyDashboard/internal/http"
	"github.com/rihow/FamilyDashboard/internal/services/nextcloud"
	"github.com/rihow/FamilyDashboard/internal/services/weather"
	"github.com/rihow/FamilyDashboard/internal/status"
)

// main はGinサーバーのエントリーポイントなのです。
// 設定読み込み → APIルーティング → 静的ファイル配信 → サーバー起動 の順で処理するます。
func main() {
	// 設定ファイルを読み込むます。
	configFilePath := "./data/settings.json"
	cfg, err := config.LoadConfig(configFilePath)
	if err != nil {
		log.Fatalf("設定ファイルの読み込みに失敗しました: %v", err)
	}

	fmt.Printf("✨ 設定を読み込みました: %s\n", cfg.GetLocationString())
	fmt.Printf("   天気更新間隔: %v\n", cfg.GetRefreshInterval("weather"))
	fmt.Printf("   カレンダー更新間隔: %v\n", cfg.GetRefreshInterval("calendar"))
	fmt.Printf("   タスク更新間隔: %v\n", cfg.GetRefreshInterval("tasks"))

	// キャッシュを初期化するます
	fc := cache.New("./data/cache")

	// エラー状態ストアを初期化するます
	errorStore := status.NewErrorStore()

	// 天気APIクライアントを初期化するます
	weatherClient := weather.NewClient(fc, "http://localhost:8080")

	// Nextcloud CalDAV/WebDAV クライアントを初期化するます
	nextcloudClient, err := nextcloud.NewClient(fc, cfg)
	if err != nil {
		fmt.Printf("⚠️ Nextcloud クライアント初期化エラー: %v\n", err)
		// エラーでも継続する（設定不足の場合はダミーデータで動作）
	} else {
		fmt.Printf("✨ Nextcloud クライアントの初期化成功\n")
	}

	// Ginルーターを初期化
	router := gin.Default()

	// グローバルミドルウェアで設定・クライアントをコンテキストに保存するます。
	router.Use(func(ctx *gin.Context) {
		ctx.Set("config", cfg)
		ctx.Set("cache", fc)
		ctx.Set("weather", weatherClient)
		ctx.Set("nextcloud", nextcloudClient)
		ctx.Set("errorStore", errorStore)
		ctx.Next()
	})

	// APIルートの設定（internal/httpで定義したルートを登録）
	httproutes.SetupRoutes(router)

	// 静的ファイル配信の設定（Svelte ビルド成果物を配信）
	// 環境変数 FRONTEND_DIST_PATH でディレクトリを指定可能（デフォルト: ./frontend/build、vite.config.js で指定）
	frontendDistPath := os.Getenv("FRONTEND_DIST_PATH")
	if frontendDistPath == "" {
		frontendDistPath = "./frontend/build"
	}

	// assetsディレクトリを配信（JS/CSS/画像など）
	router.Static("/assets", frontendDistPath+"/assets")

	// ルートへのアクセスはindex.htmlを返す（SPA対応）
	router.NoRoute(func(ctx *gin.Context) {
		indexFile := frontendDistPath + "/index.html"
		if _, err := os.Stat(indexFile); err == nil {
			ctx.File(indexFile)
		} else {
			ctx.JSON(404, gin.H{
				"error":   "index.html not found. Frontend build required.",
				"path":    indexFile,
				"message": "Please run 'npm run build' in the frontend directory.",
			})
		}
	})

	// 既定ポート8080で起動するます。
	port := ":8080"
	fmt.Printf("🚀 サーバー起動するます！ http://localhost%s\n", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("サーバー起動に失敗しました: %v", err)
	}
}
