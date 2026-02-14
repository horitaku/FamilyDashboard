package main

import (
	"github.com/gin-gonic/gin"
	httproutes "github.com/rihow/FamilyDashboard/internal/http"
)

// main はGinサーバーのエントリーポイントなのです。
// APIルーティングと静的ファイル配信を設定して、起動するもなのです。
func main() {
	router := gin.Default()

	// APIルートの設定（internal/httpで定義したルートを登録）
	httproutes.SetupRoutes(router)
	router.Static("/assets", "./frontend/dist/assets")

	// ルートへのアクセスはindex.htmlを返す（SPA対応）
	router.NoRoute(func(ctx *gin.Context) {
		ctx.File("./frontend/dist/index.html")
	})

	// 既定ポート8080で起動するます（必要に応じて設定で差し替えます）。
	_ = router.Run(":8080")
}
