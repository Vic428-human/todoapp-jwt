// responsible for running database
package main

import (
	"log"
	"time"

	"todo_api/internal/config"
	"todo_api/internal/database"
	"todo_api/internal/handlers"
	"todo_api/internal/middleware"

	// "todo_api/internal/repository"
	// "todo_api/internal/service"

	// "cloud.google.com/go/storage"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool" // PostgreSQL 驅動程式的 connection pool 版本，提供高效連線管理
)

func main() {
	var cfg *config.Config
	var err error

	cfg, err = config.Load()
	if err != nil {
		log.Fatal(err)
	}

	var pool *pgxpool.Pool

	// 1. 應用啟動時：只建立一次連線池（生命週期 = 整個應用）
	pool, err = database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	// =========================
	// 暫時停用 GCS 相關初始化
	// 等 Render 上的 GCS credentials 設定好後再打開
	// =========================

	// 建立 GCS client
	// ctx := context.Background()
	//
	// storageClient, err := storage.NewClient(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer storageClient.Close()
	//
	// 建立 image repository
	// imageRepo := repository.NewGCImageRepository(
	// 	storageClient,
	// 	cfg.GCSBucketName,
	// )
	//
	// 建立 user service
	// userService := service.NewUserService(pool, imageRepo)

	// create server
	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// router.Use(middleware.CORSMiddleware())

	// health check
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message":  "todo api running successfully",
			"status":   "success",
			"database": "connected",
			"gcs":      "disabled",
		})
	})

	// Todo routes
	router.POST("/todos", handlers.CreateTodoHandler(pool))
	router.GET("/todos", handlers.GetTodosHandler(pool))
	router.GET("/todos/:id", handlers.GetTodoByIDHandler(pool))
	router.PUT("/todos/:id", handlers.UpdateToDoHandler(pool))

	// Auth routes
	router.POST("/auth/register", handlers.CreateUserHandler(pool))
	router.POST("/auth/login", handlers.LoginHandler(pool, cfg))

	// User routes
	// 這條就是之後用 Postman / 前端測試頭像上傳的 API
	// 先暫時註解，等 GCS credentials 設定好再打開
	// router.PUT("/users/:id/profile-image", handlers.SetProfileImageHandler(userService))

	// Middleware test route
	router.GET("/protected-test", middleware.AuthMiddleware(cfg), handlers.TestProtectedHandler())

	// Product routes
	router.POST("/products", handlers.CreatteProductHandler(pool))
	router.GET("/products", handlers.GetAllProductsHandler(pool))
	router.PUT("/products/:id", handlers.UpdateProductHandler(pool))
	router.GET("/products/:id", handlers.GetProductByIDHandler(pool))
	router.GET("/products/search", handlers.ListProductsHandler(pool))

	log.Printf("server starting on port %s\n", cfg.Port)
	log.Printf("GCS bucket in use: %s\n", cfg.GCSBucketName)

	router.Run(":" + cfg.Port)
}
