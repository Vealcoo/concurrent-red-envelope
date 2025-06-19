# 紅包領取系統

## 專案概述

本專案是一個基於 Go 語言開發的紅包領取系統，旨在支援高併發場景。系統使用 MySQL、MongoDB 和 Redis 作為資料庫與緩存，並透過 RESTful API 提供服務。專案採用 Docker Compose 進行資料庫環境管理，提供創建活動、啟動活動、插入測試數據以及領取紅包的功能。

## 系統架構

- **後端**：使用 Go 語言，基於 Gin 框架構建 RESTful API。
- **資料庫**：
  - **MySQL**：儲存結構化數據，如活動和用戶信息。
  - **MongoDB**：儲存靈活的非結構化數據。
  - **Redis**：用於緩存，提升高併發場景下的存取效率。
- **容器化**：透過 Docker Compose 管理 MySQL、MongoDB 和 Redis 服務。
- **負載測試**：使用 Fortio 工具進行高併發負載測試。

## 路由與端點

系統提供以下 RESTful API 端點：

- **POST /insert_test_data**  
  插入測試數據，用於模擬初始化環境。
- **POST /campaign**  
  創建一個新的紅包活動。
- **POST /campaign/:campaign_id/start**  
  啟動指定的紅包活動，需指定 `cache_mode`（`mongodb` 或 `redis`）。
- **POST /campaign/:campaign_id/claim**  
  領取指定活動的紅包，需提供用戶 ID 和 `cache_mode`。

## 環境配置

專案使用 Docker Compose 管理資料庫和緩存服務，配置文件 `docker-compose.yml` 包含以下服務：

### MySQL

- **映像**：`mysql:9.0`
- **端口**：`3306:3306`
- **環境變數**：
  - `MYSQL_ROOT_PASSWORD=admin`
  - `MYSQL_DATABASE=test`
  - `MYSQL_USER=admin`
  - `MYSQL_PASSWORD=admin`
- **資源限制**：1 CPU，1024MB 記憶體

### MongoDB

- **映像**：`mongo:7.0`
- **端口**：`27017:27017`
- **環境變數**：
  - `MONGO_INITDB_ROOT_USERNAME=admin`
  - `MONGO_INITDB_ROOT_PASSWORD=admin`
- **資源限制**：0.5 CPU，512MB 記憶體

### Redis

- **映像**：`redis:7.4`
- **端口**：`6379:6379`
- **資源限制**：0.25 CPU，256MB 記憶體

**資源限制說明**：  
資源限制確保每個服務在 Docker 容器中佔用的 CPU 和記憶體不超過指定值，防止單一服務過度消耗系統資源。根據您的硬體配置，可在 `docker-compose.yml` 中調整 `cpus` 和 `memory` 參數。例如，若伺服器性能較高，可將 MySQL 的記憶體限制提升至 2048MB。

## 運行專案

按照以下步驟啟動專案：

1. **啟動 Docker 環境**

   ```bash
   docker-compose up -d
   ```

   此命令會啟動 MySQL、MongoDB 和 Redis 容器。

2. **啟動後端服務**

   ```bash
   go run main.go
   ```

3. **執行測試腳本**

   - 寫入測試資料

   ```bash
   ./script/insert_test_data.sh
   ```

   為模擬真實資料量，所以會寫入 2000000 筆 資料

   - 執行測試 參數 1: 活動 ID 參數 2: 緩存模式

   ```bash
   ./exec_testing.sh 1 mongodb
   ```
