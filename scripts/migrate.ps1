# 透過 PowerShell script 來執行冗長的 migration 指令
# 用途一句話總結：
# 這是一支 PowerShell wrapper，用來：
# - 自動載入 .env
# - 統一管理 golang-migrate 指令
# - 避免每次手打一長串參數

# 套用 migration
# .\scripts\migrate.ps1 up 

# 回滾 1 個（預設）
# .\scripts\migrate.ps1 down

# 回滾 3 個
# .\scripts\migrate.ps1 down 3

# 建立 migration
# .\scripts\migrate.ps1 create add_users_table

# 強制指定版本
# .\scripts\migrate.ps1 force 10

# 讀取 .env 檔案內容，一行一行處理
Get-Content .env | ForEach-Object {
    if ($_ -match '^([^#][^=]+)=(.+)$') {
        # 把 .env 裡的 KEY=VALUE 設成 PowerShell 的環境變數
        Set-Item -Path "env:$($matches[1])" -Value $matches[2]
    }
}

# 取得第一個 CLI 參數：指令 (up / down / create / force)
$command = $args[0]
# 取得第二個 CLI 參數：名稱或數量
$name = $args[1]

# 根據 command 決定要執行哪個 migrate 行為
switch ($command) {
    # ===== migrate up =====
    "up" { migrate -database $env:DATABASE_URL -path migrations up }

    # ===== migrate down =====
    "down" { 
        $count = if ($name) { $name } else { "1" }
        Write-Host "Rolling back $count migration(s). Continue? [y/N]"
        $confirm = Read-Host
        if ($confirm -eq 'y') {
            migrate -database $env:DATABASE_URL -path migrations down $count
        }
    }

    # ===== migrate create =====
    "create" { migrate create -ext sql -dir migrations -seq $name }

    # ===== migrate force =====
    "force" { migrate -database $env:DATABASE_URL -path migrations force $name }
}
