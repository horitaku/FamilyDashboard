package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// main はGinサーバーのエントリーポイントです。
// いまは雛形として最低限のルーティングと静的配信だけを用意しています。
func main() {
    router := gin.Default()

    // ヘルスチェック用の最小エンドポイントです。
    router.GET("/api/health", func(ctx *gin.Context) {
        ctx.JSON(http.StatusOK, gin.H{
            "ok": true,
        })
    })

    // Svelteの静的ビルドを配信する前提のパスです。
    // 実際のビルドは後のステップで実行します。
    router.Static("/assets", "./frontend/dist/assets")
    router.NoRoute(func(ctx *gin.Context) {
        ctx.File("./frontend/dist/index.html")
    })

    // 既定ポートで起動します（必要に応じて設定で差し替えます）。
    _ = router.Run(":8080")
}
