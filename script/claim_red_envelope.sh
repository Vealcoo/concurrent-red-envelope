#!/bin/bash
set -e

# 檢查是否提供了 campaign_id、user_id 和 cache_mode 參數
if [ -z "$1" ] || [ -z "$2" ] || [ -z "$3" ]; then
    echo "Error: Please provide campaign_id, user_id, and cache_mode parameters"
    echo "Usage: $0 <campaign_id> <user_id> <cache_mode>"
    exit 1
fi

# 將參數賦值給變量
CAMPAIGN_ID="$1"
USER_ID="$2"
CACHE_MODE="$3"

# 驗證 cache_mode 是否為 mongodb 或 redis
if [ "$CACHE_MODE" != "mongodb" ] && [ "$CACHE_MODE" != "redis" ]; then
    echo "Error: cache_mode must be 'mongodb' or 'redis'"
    exit 1
fi

# 發送 curl 請求
curl -X POST "http://localhost:8089/campaign/${CAMPAIGN_ID}/claim" \
    -H "Content-Type: application/json" \
    -d "{\"user_id\": \"${USER_ID}\", \"cache_mode\": \"${CACHE_MODE}\"}"
