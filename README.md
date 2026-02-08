
> 這個專案主要用於理解 JWT + GIN + GOLANG + PostgreSQL 規劃登入註冊功能，完成後會把它遷移過去跟交易所的登入功能做整合

```
<!-- 快速生成檔案 -->
mkdir cmd\api internal\config,internal\database,internal\handlers,internal\middleware,internal\models,internal\repository,migrations
```


yourproject/
├── go.mod
├── cmd/  reponsible forrunning database
│
├── internal/   
│ └── config/    // set enviroment
│ └── database/    // database connection setup and pooling 
│ └── config/    // handling our http request
│ └── middleware/    // use for authentication
└── frontend/

│