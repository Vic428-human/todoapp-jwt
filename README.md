
> 這個專案主要用於理解 JWT + GIN + GOLANG + PostgreSQL 規劃登入註冊功能，完成後會把它遷移過去跟交易所的登入功能做整合


####  本專案用到的dependencies
> 來源跟用途說明

```
<!-- 快速生成檔案 -->
mkdir cmd\api internal\config,internal\database,internal\handlers,internal\middleware,internal\models,internal\repository,migrations

<!-- initialize the project -->
go mod init todo_api

<!-- run this project -->
go run cmd\api\main.go

<!-- setup dependencies -->
// handle our http requests 
// https://ithelp.ithome.com.tw/articles/10387192
go get -u github.com/gin-gonic/gin

<!-- install postgres's driver -->
// https://github.com/jackc/pgx/blob/f56ca73076f3fc935a2a049cf78993bfcbba8f68/examples/url_shortener/main.go#L11
go get -u github.com/jackc/pgx/v5 
go get -u github.com/jackc/pgx/v5/pgxpool //for a concurrency safe connection pool

<!-- install jwt for authentication -->
/*
為何選擇Migrate庫
版本控制和可追溯性: Migrate庫提供了一種簡潔的方式來版本化資料庫結構的改變。每個遷移都被保存為一個單獨的檔案，檔名通常包含時間戳和描述，這使得追蹤和審計資料庫結構的變更變得簡單直觀。這與手動管理一系列.sql腳本檔案相比，更加系統化和易於維護。
自動化操作: 使用Migrate庫可以實現遷移操作的自動化，如自動執行下一個未應用的遷移或回滾到特定版本。這種自動化大大降低了人為錯誤的風險，並提高了開發和部署的效率。
跨資料庫相容性: Migrate支援廣泛的資料庫技術，這意味著同一套遷移腳本可以用於不同的資料庫系統，從而簡化了多環境或多資料庫系統的遷移策略。
易於整合和擴展: 作為一個Go庫，Migrate可以輕鬆整合到Go應用程式中。它也支援透過插件來擴展更多的資料庫類型或自訂遷移邏輯。
*/
// https://github.com/golang-migrate/migrate/tree/master/database/postgres
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

<!-- generates a bcrypt hash of a given password -->
// https://pkg.go.dev/golang.org/x/crypto/bcrypt
go get -u golang.org/x/crypto/bcrypt  

<!--  read environment variables in golang -->
// https://github.com/joho/godotenv
go get -u github.com/joho/godotenv 

<!-- 熱重載工具，作用是在你修改代碼文件後自動重新編譯並重啟程序 -->
// https://www.bilibili.com/opus/1068145464453365769 => 有詳細步驟
go install github.com/cosmtrek/air@latest // 路徑已經改成下方這個
go install github.com/air-verse/air@latest

<!-- install postgres  5432 port , 密碼:提示結婚 (安裝時的設定)-->
psql --version // 確認是否下載
psql -U postgres // 指定以 postgres 用户身份連接數據庫，如果失敗改用 psql -U z0983 ，電腦名稱是 z0983。

\l // look at database, show us all of db, 正常會出現 postgres、template0、template1
CREATE DARABSE XXX; // write sql query in uppercase

rmdir /S /Q "C:\Program Files\PostgreSQL" // 安裝失敗的時候強制刪除 
C:\Program Files\PostgreSQL\18\bin // 加到系統環境變數
介紹 pgAdmin UI介面用法 // https://www.youtube.com/watch?v=T1PrXly6kOs
\c todo_api // connected to database
\q // quit connect 

<!-- create config.go file -->
NEW-ITEM -Path internal\config\config.go -ItemType File
```
### 補充
#### v5/pgxpool
[Creating the Connection Pool](https://resources.hexacluster.ai/blog/postgresql/postgresql-client-side-connection-pooling-in-golang-using-pgxpool/
)

#### golang air
```
// 因為本專案的main.go在cmd路徑裡
.air.toml 裡面的 cmd 加上 ./cmd/api // cmd = "go build -o ./tmp/main.exe ./cmd/api" 
```



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
│ └── models/    // this going to provide structure of data 
│ └── repository/    // all of database operations
├── internal/
└── frontend/ 預期放登入註冊畫面
```