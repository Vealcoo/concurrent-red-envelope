#!/bin/bash
set -e

# 檢查是否提供了 campaign_id 和 cache_mode 參數
if [ -z "$1" ] || [ -z "$2" ]; then
    echo "Error: Please provide campaign_id and cache_mode parameters"
    echo "Usage: $0 <campaign_id> <cache_mode>"
    exit 1
fi

# 將參數賦值給變量
CAMPAIGN_ID="$1"
CACHE_MODE="$2"

# 驗證 cache_mode 是否為 mongodb 或 redis
if [ "$CACHE_MODE" != "mongodb" ] && [ "$CACHE_MODE" != "redis" ]; then
    echo "Error: cache_mode must be 'mongodb' or 'redis'"
    exit 1
fi

./script/start_campaign.sh $CAMPAIGN_ID $CACHE_MODE
./script/fortio_testing.sh $CAMPAIGN_ID $CACHE_MODE
