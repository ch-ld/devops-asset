#!/bin/bash

# DNS提供商API测试脚本
BASE_URL="http://localhost:8080/api/v1/dns/providers"

echo "=== DNS提供商API测试 ==="

# 1. 获取提供商列表
echo "1. 测试获取提供商列表..."
curl -X GET "$BASE_URL" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  | jq '.' || echo "请求失败"

echo -e "\n"

# 2. 创建测试提供商
echo "2. 测试创建提供商..."
curl -X POST "$BASE_URL" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{
    "name": "测试阿里云DNS",
    "type": "aliyun",
    "credentials": {
      "access_key_id": "test-key-id",
      "access_key_secret": "test-key-secret",
      "region": "cn-hangzhou"
    },
    "remark": "测试用提供商"
  }' | jq '.' || echo "请求失败"

echo -e "\n"

# 3. 测试连接（需要先获取提供商ID）
echo "3. 测试连接功能..."
echo "需要先创建真实的提供商后才能测试连接功能"

echo -e "\n=== 测试完成 ===" 
