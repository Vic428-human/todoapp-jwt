// reponsible forrunning database
package main

import (
	"log"
	"time"
	"todo_api/internal/config"
	"todo_api/internal/database"
	"todo_api/internal/handlers"
	"todo_api/internal/middleware"

	"github.com/gin-contrib/cors"
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
	router.SetTrustedProxies(nil)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 將「同一個」pool 實例傳給所有 handler
	router.GET("/", func(c *gin.Context) {
		// gin.H is a shortcut for map[string]interface{} or map[string]any
		c.JSON(200, gin.H{
			"message":  "todo api running successfully",
			"status":   "success",
			"database": "connected",
		})
	})

	// 當前專案會用到
	router.POST("/todos", handlers.CreateTodoHandler(pool))
	// 0315的時候已經把 users + pagination 放進去交易所前台專案了
	router.GET("/todos", handlers.GetTodosHandler(pool)) // 有分頁
	router.GET("/todos/:id", handlers.GetTodoByIDHandler(pool))
	router.PUT("/todos/:id", handlers.UpdateToDoHandler(pool))
	router.POST("/auth/register", handlers.CreateUserHandler(pool))
	router.POST("/auth/login", handlers.LoginHandler(pool, cfg))

	// Middleware Test Route
	router.GET("/protected-test", middleware.AuthMiddleware(cfg), handlers.TestProtectedHandler())

	// 交易所才會用到，只是在這進行測試
	router.POST("/products", handlers.CreatteProductHandler(pool))
	router.GET("/products", handlers.GetAllProductsHandler(pool)) // 無 keyword：全拿
	router.PUT("/products/:id", handlers.UpdateProductHandler(pool))
	router.GET("/products/:id", handlers.GetProductByIDHandler(pool))
	// router 加這行（不碰現有）已經實驗過搜尋 "太陽神" 關鍵字會只拿到 太陽神有關的商品列表 => http://localhost:3000/products/search?keyword=太陽神
	router.GET("/products/search", handlers.ListProductsHandler(pool))

	router.Run(":" + cfg.Port) // listens on 0.0.0.0:8080 by default

}
