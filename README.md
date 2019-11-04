# GameServer

# 1. Server啟動處理流程

1-1.
  - Config檔案解析

1-2.
- Config資料初始化DB管理系統並連線

1-3.
- 從連線的DB取得最新的Setting表

1-4.
- 根據Setting表Key值與Config資料相同的變數將更新Config資料不相同的資料將解析為gamesetting格式

1-5.
- 存儲Config、gamesetting於Server Memory

1-6.
- 初始化各系統功能

1-7.
- 啟動各服務項目(使用goroutin啟動多個服務接口)
