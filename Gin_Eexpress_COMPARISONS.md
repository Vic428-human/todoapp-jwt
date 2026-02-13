


# Equivalent of `Express.js` in Go

```bash
go get github.com/gin-gonic/gin
```

# Equivalent of `Express.js Middleware` in Go

```go
func main() {
    cfg := loadConfig()
    // gin.Default()是對 gin.new() 的封装，加入了日誌和錯誤恢复中間件
    router := gin.Default()  

    admin := router.Group("/admin")
	admin.Use(h.AuthMiddleware())
    {
		admin.GET("", h.ServeAdminDashboard)
    }
    
    router.Run(":" + cfg.Port) //8080
}

type Config struct {
	Port             string
	DBPath           string
	SessionSecretKey string
}

func loadConfig() Config {
	return Config{
		Port:             getEnv("PORT", "8080"), // 定義key 跟 value
		DBPath:           getEnv("DATABASE_URL", "./data/orders.db"),
		SessionSecretKey: getEnv("SESSION_SECRET_KEY", "pizza-order-secret-key"),
	}
}

func getEnv(key, defaultValue string) string {
	// 環境變數存在時，用環境變數設定的值
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	// 環境變不存在時
	return defaultValue
}


```