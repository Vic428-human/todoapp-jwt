// reponsible forrunning database
package main

import (
	"log"
	"todo_api/internal/config"
	"todo_api/internal/database"

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

	// 建立連線字串
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
	router.GET("/", func(c *gin.Context) {
		router.SetTrustedProxies(nil) // if you don't use any proxy, you can disable this feature by using nil, then Context.ClientIP() will return the remote address directly to avoid some unnecessary computation

		// gin.H is a shortcut for map[string]interface{} or map[string]any
		c.JSON(200, gin.H{
			"message":  "!!!todo api running successfully~~~",
			"status":   "success",
			"database": "connected",
		})
	})
	router.Run(":" + cfg.Port) // listens on 0.0.0.0:8080 by default

}
