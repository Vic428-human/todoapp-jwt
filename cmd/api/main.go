// reponsible forrunning database
package main

import (
	"log"
	"todo_api/internal/config"
	"todo_api/internal/database"
	"todo_api/internal/handlers"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool" // PostgreSQL驅動程式的connection pool版本，提供高效連線管理
)

func main() {
	var cfg *config.Config
	var err error

	cfg, err = config.Load()
	if err != nil {
		log.Fatal(err)
	}
	var pool *pgxpool.Pool

	// 1️⃣ 應用啟動時：只建立「一次」連線池（生命週期 = 整個應用）
	pool, err = database.Connect(cfg.DatabaseURL)
	if err != nil {
		// 連線失敗時立即終止程式
		log.Fatal(err)
	}

	defer pool.Close() // 確保程式結束時關閉連線池

	// create server, take a look at routes, want api fast, use instance from the memory, pointer variable
	// * is a pointer, reference something in the memory
	// pointer refers to the address or instance in memory, and not copy entire thing
	var router *gin.Engine = gin.Default() // gin => do client request and response

	// 2️⃣ 將「同一個」pool 實例傳給所有 handler
	router.GET("/", func(c *gin.Context) {
		router.SetTrustedProxies(nil) // if you don't use any proxy, you can disable this feature by using nil, then Context.ClientIP() will return the remote address directly to avoid some unnecessary computation
		// gin.H is a shortcut for map[string]interface{} or map[string]any
		c.JSON(200, gin.H{
			"message":  "!!!todo api running successfully~~~",
			"status":   "success",
			"database": "connected",
		})
	})

	// 當前專案會用到
	router.POST("/todos", handlers.CreateTodoHandler(pool))
	router.GET("/todos", handlers.GetAllTodosHandler(pool))
	router.GET("/todos/:id", handlers.GetTodoByIDHandler(pool))

	// 交易所才會用到，只是在這進行測試
	router.POST("/products", handlers.CreatteProductHandler(pool))
	router.GET("/products", handlers.GetAllProductsHandler(pool))

	router.Run(":" + cfg.Port) // listens on 0.0.0.0:8080 by default

}
