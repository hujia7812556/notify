#!/bin/bash

# 基础URL
BASE_URL="http://localhost:8080/api/v1"

# 颜色输出
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

# 测试健康检查接口
echo "Testing health check endpoint..."
health_response=$(curl -s -w "\n%{http_code}" $BASE_URL/health)
status_code=$(echo "$health_response" | tail -n1)
response_body=$(echo "$health_response" | head -n1)

if [ $status_code -eq 200 ]; then
    echo -e "${GREEN}Health check passed${NC}"
else
    echo -e "${RED}Health check failed with status $status_code${NC}"
    echo "Response: $response_body"
fi

# 测试WxPusher消息发送
echo -e "\nTesting WxPusher notification..."
wxpusher_response=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "platform": "wechat",
        "content": "Test message from API test",
        "extra": {
            "user_id": "UID_xxx"
        }
    }' \
    $BASE_URL/notify)
status_code=$(echo "$wxpusher_response" | tail -n1)
response_body=$(echo "$wxpusher_response" | head -n1)

if [ $status_code -eq 202 ]; then
    echo -e "${GREEN}WxPusher notification sent successfully${NC}"
else
    echo -e "${RED}WxPusher notification failed with status $status_code${NC}"
    echo "Response: $response_body"
fi

# 测试钉钉消息发送
echo -e "\nTesting DingTalk notification..."
dingtalk_response=$(curl -s -w "\n%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -d '{
        "platform": "dingtalk",
        "content": "Test message from API test",
        "extra": {}
    }' \
    $BASE_URL/notify)
status_code=$(echo "$dingtalk_response" | tail -n1)
response_body=$(echo "$dingtalk_response" | head -n1)

if [ $status_code -eq 202 ]; then
    echo -e "${GREEN}DingTalk notification sent successfully${NC}"
else
    echo -e "${RED}DingTalk notification failed with status $status_code${NC}"
    echo "Response: $response_body"
fi 