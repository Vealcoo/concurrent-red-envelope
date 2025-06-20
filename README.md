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

## 環境設置

### 1. 安裝 Docker 和 Docker Compose

確保已安裝 Docker 和 Docker Compose：

- **Docker**：請參考 [官方安裝指南](https://docs.docker.com/get-docker/)。
- **Docker Compose**：通常隨 Docker Desktop 一起安裝，或參考 [官方安裝指南](https://docs.docker.com/compose/install/)。

### 2. 安裝 Go

專案需要 Go 環境來運行後端服務：

- 下載並安裝 Go（建議版本 1.18 或更高）：[Go 官方下載](https://golang.org/dl/)。
- 驗證安裝：
  ```bash
  go version
  ```

### 3. 安裝 Fortio

Fortio 是用於高併發負載測試的工具，需安裝以執行測試腳本：

- 下載並安裝 Fortio：[Fortio 官方下載](https://github.com/fortio/fortio)。

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

1. **啟動服務**

   ```bash
   ./start_server.sh
   ```

2. **執行測試腳本**

   - 寫入測試資料

   ```bash
   ./script/insert_test_data.sh
   ```

   為模擬真實資料量，所以會寫入 1000 個活動，共 2000000 筆紅包資料，比較久一點

   - 執行測試 參數 1: 活動 ID 參數 2: 緩存模式

   ```bash
   ./exec_testing.sh 1 mongodb
   ```

   每個活動只能跑一次，每次測試請換不同活動 ID

   - 腳本路徑：script/fortio_testing.sh

   ```bash
   fortio load -qps 0 -c 50 -n 2000 -H "Content-Type: application/json" \
    -payload "{\"user_id\": \"test_user\", \"cache_mode\": \"${CACHE_MODE}\"}" \
    http://localhost:8089/campaign/${CAMPAIGN_ID}/claim
   ```

   可依照自己需求調整腳本
