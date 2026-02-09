# https://github.com/golang-migrate/migrate/blob/257fa847d614efe3948c25e9033e92b930527dec/database/postgres/TUTORIAL.md

# 透過 powershell script 來執行冗長的 migration 指令
# 你可以決定有多少次migration you want to rollback
# 整體用途一句話總結
#這是一支 PowerShell wrapper，用來：
#自動載入 .env
#統一管理 golang-migrate 指令
#避免每次手打一長串參數

# 套用 migration
# .\scripts\migrate.ps1 up 

# 回滾 1 個（預設）
#.\migrate.ps1 down

# 回滾 3 個
#.\migrate.ps1 down 3

# 建立 migration
#.\migrate.ps1 create add_users_table

# 強制指定版本
# .\scripts\migrate.ps1 force up

# 讀取 .env 檔案內容，一行一行處理
Get-Content .env | ForEach-Object {
    if ($_ -match '^([^#][^=]+)=(.+)$') {
        # 把 .env 裡的 KEY=VALUE
        # 設成 PowerShell 的環境變數
        Set-Item -Path "env:$($matches[1])" -Value $matches[2]
    }
}

# 組合完整的 URL
# 取得第一個 CLI 參數：指令 (up / down / create / force)
$command = $args[0]
# 取得第二個 CLI 參數：名稱或數量
$name = $args[1]

# 根據 command 決定要執行哪個 migrate 行為
switch ($command) {
    # Run migrations
    "up" { migrate -database $env:POSTGRES_URL -path migrations up }
    # 回滾 migration
    "down" { 
        # 如果有指定參數就用它（回滾幾個）
        # 沒指定預設回滾 1 個
        $count = if ($name) { $name } else { "1" }

        # 提示使用者確認，避免誤回滾, prevent destroy press N
        Write-Host "Rolling back $count migration(s). Continue? [y/N]"

        # 只有輸入 y 才真的執行
        $confirm = Read-Host
        if ($confirm -eq 'y') {
            # check if running reverse migration also works
            migrate -database $env:POSTGRES_URL -path migrations down $count
        }
    }
    # ===== migrate create =====
    # 建立新的 migration 檔案
    "create" { migrate create -ext sql -dir migrations -seq $name }

    # ===== migrate force =====
    # 強制設定 migration version
    # 常用於：
    # - migration 壞掉
    # - dirty 狀態修復
    "force" { migrate -database $env:POSTGRES_URL -path migrations force $name }
}