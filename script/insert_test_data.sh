#!/bin/bash
set -e

# 發送 curl 請求
curl -X POST "http://localhost:8089/insert_test_data" \
    -H "Content-Type: application/json"
