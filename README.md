
> 這個專案主要用於理解 JWT + GIN + GOLANG + PostgreSQL 規劃登入註冊功能，完成後會把它遷移過去跟交易所的登入功能做整合

#### 0214 更新
> 已經確定 [分支](https://github.com/Vic428-human/todoapp-jwt/tree/connect-react-and-golang-postgres) 實驗過，React + ts + react-query + postgres + go + gin，本地端已經串起來。
> 所以接下來到交易所版本2 加入本專案的後端，然後實驗post API，補充，在交易所版本2中，會拔掉clerk。


- [本專案學習過程中的筆記整理](https://www.notion.so/2-2f6a54651e3e80d888ede6403ad3bf6a)

### 專案規劃流程 
- 下載所需套件
- 新專案要加資料表 (基於 mirgration 的方式添加)
- Handler
- Repository

#### 下載所需套件
> 套件、DB安裝

```
<!-- 快速生成檔案 -->
mkdir cmd\api internal\config,internal\database,internal\handlers,internal\middleware,internal\models,internal\repository,migrations

<!-- initialize the project -->
go mod init todo_api // go.mod會出現 module todo_api

<!-- run this project -->
go run cmd\api\main.go // 主程式入口點

<!-- setup dependencies :　本專案會需要用的的套件及其用途 -->

// https://ithelp.ithome.com.tw/articles/10387192
go get -u github.com/gin-gonic/gin // 處理 http requests 

<!-- 下載 postgres's driver -->
// https://github.com/jackc/pgx/blob/f56ca73076f3fc935a2a049cf78993bfcbba8f68/examples/url_shortener/main.go#L11
go get -u github.com/jackc/pgx/v5 
go get -u github.com/jackc/pgx/v5/pgxpool // PostgreSQL驅動程式的connection pool版本，提供高效連線管理

/*
為何選擇Migrate庫 https://github.com/golang-migrate/migrate/tree/master/database/postgres
版本控制和可追溯性: Migrate庫提供了一種簡潔的方式來版本化資料庫結構的改變。每個遷移都被保存為一個單獨的檔案，檔名通常包含時間戳和描述，這使得追蹤和審計資料庫結構的變更變得簡單直觀。這與手動管理一系列.sql腳本檔案相比，更加系統化和易於維護。
自動化操作: 使用Migrate庫可以實現遷移操作的自動化，如自動執行下一個未應用的遷移或回滾到特定版本。這種自動化大大降低了人為錯誤的風險，並提高了開發和部署的效率。
跨資料庫相容性: Migrate支援廣泛的資料庫技術，這意味著同一套遷移腳本可以用於不同的資料庫系統，從而簡化了多環境或多資料庫系統的遷移策略。
易於整合和擴展: 作為一個Go庫，Migrate可以輕鬆整合到Go應用程式中。它也支援透過插件來擴展更多的資料庫類型或自訂遷移邏輯。
*/
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

<!-- generates a bcrypt hash of a given password -->
go get -u golang.org/x/crypto/bcrypt  // https://pkg.go.dev/golang.org/x/crypto/bcrypt

<!--  read environment variables in golang -->
go get -u github.com/joho/godotenv  // https://github.com/joho/godotenv

<!-- 熱重載工具，作用是在你修改代碼文件後自動重新編譯並重啟程序 -->
go install github.com/cosmtrek/air@latest // 路徑已經改成下方這個
go install github.com/air-verse/air@latest // https://www.bilibili.com/opus/1068145464453365769 => 有詳細步驟

<!-- install postgres  5432 port , 密碼:提示結婚 (安裝時的設定)-->
psql --version // 確認是否下載
psql -U postgres // 指定以 postgres 用户身份連接數據庫，如果失敗改用 psql -U z0983 ，電腦名稱是 z0983。
```

#### 新專案要加資料表
> migration 操作

```
<!-- 新專案要加資料表、修改 schema（加欄位、改型別）、團隊合作時，讓大家資料庫版本一致時都會需要用的 -->
migrate create -ext sql -dir migrations -seq create_todos_table //  在 migrations 資料夾中建立一組用 SQL 撰寫、用遞增編號命名的「建立 todos 資料表」 分別為 1.up migration  2. down migration  檔案。

igrate.ps1 up //  先在 XXX_create_todos_table 裡面寫好sql語法創建table，接者操作 .\scripts\migrate.ps1 up 建立資料表

<!-- 下方是常見的db操作 -->
\l // look at database, show us all of db, 正常會出現 postgres、template0、template1

CREATE DARABSE XXX; // write sql query in uppercase

rmdir /S /Q "C:\Program Files\PostgreSQL" // 安裝失敗的時候強制刪除 

C:\Program Files\PostgreSQL\18\bin // 加到系統環境變數

介紹 pgAdmin UI介面用法 // https://www.youtube.com/watch?v=T1PrXly6kOs

\c todo_api // connected to database

\q // quit connect 

psql -U postgres -d todo_api // connet to certain db

\dt // find table 

NEW-ITEM -Path internal\config\config.go -ItemType File // create config.go file 
```

#### Handler 
- Handler 層用途 
> Handler 負責處理 HTTP 請求，執行業務邏輯（不涉及資料庫操作），接收請求、驗證參數、調用 Repository，然後回傳結果

#### Repository 
- Repository 層用途 
> Repository 封裝 SQL 查詢，負責與資料庫互動（新增、修改、查詢等），讓 Handler 只需調用簡單介面，避免直接寫 SQL。


#### auth middleware jwt validation
> API都規劃完成之後才處理,

### 補充
#### v5/pgxpool
[Creating the Connection Pool](https://resources.hexacluster.ai/blog/postgresql/postgresql-client-side-connection-pooling-in-golang-using-pgxpool/
)

#### golang air
```
// 因為本專案的main.go在cmd路徑裡
.air.toml 裡面的 cmd 加上 ./cmd/api // cmd = "go build -o ./tmp/main.exe ./cmd/api" 
```


### 專案架構 
```
yourproject/
├── go.mod
├── cmd/  reponsible forrunning database
│
├── internal/ 
│ └── config/  引用環境變數 
│ └── database/    // Creating the Connection and Pooling for postgres
│ └── config/    // handling our http request
│ └── middleware/    // use for authentication
│ └── models/    // 定義資料結構
│ └── repository/    // 封裝 SQL 查詢，負責與資料庫互動（新增、修改、查詢等，存取資料
│ └── handlers/ // HTTP API
├── scripts/ 
│ └── migrate.ps1/ //  把 powershell 寫成腳本，避免每次運行都寫一堆指令
│
├── migrations/  // migrate create -ext sql -dir migrations -seq create_todos_table 用途: 新專案要加資料表、修改 schema（加欄位、改型別）、團隊合作時，讓大家資料庫版本一致
└── frontend/ 預期放登入註冊畫面
```

### 知識點

- [Basic Timeout with Context](https://go-cookbook.com/snippets/context/using-context-for-timeouts)

### 比較這三個 Go Web 框架

| 框架 | 特點 | 效能 | Side Project 01 | Side Project 02|下載|
|------|------|------|--------|------------------|------|
| [Fiber](https://github.com/Vic428-human/todoapp-jwt/blob/main/Fiber_Eexpress_COMPARISONS.md) | 最像 Express.js，語法最接近 | 超高 |  | |go get github.com/gofiber/fiber/v2|
| [Gin](https://github.com/Vic428-human/todoapp-jwt/blob/main/Gin_Eexpress_COMPARISONS.md) | 成熟穩定，生態系豐富 | 高 | https://github.com/Vic428-human/todoapp-jwt | https://github.com/Vic428-human/go-pizza-order-tracker |	"github.com/gin-gonic/gin"|
| **Echo** | 功能完整，效能佳 | 高 | |  |尚未用過|
