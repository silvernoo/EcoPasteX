#!/bin/bash

echo "🧪 测试 EcoPaste Webhook..."
echo ""

# 测试文本数据
echo "📝 发送文本数据..."
curl -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "type": "text",
    "value": "这是一条测试文本消息",
    "timestamp": "'$(date -u +"%Y-%m-%dT%H:%M:%SZ")'",
    "subtype": "plain"
  }'

echo ""
echo ""

# 测试 HTML 数据
echo "🌐 发送 HTML 数据..."
curl -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "type": "html",
    "value": "<html><head></head><body><h1>测试标题</h1><p>这是HTML内容</p></body></html>",
    "timestamp": "'$(date -u +"%Y-%m-%dT%H:%M:%SZ")'",
    "subtype": "html"
  }'

echo ""
echo ""

# 测试图片数据（小的 base64 图片）
echo "🖼️  发送图片数据..."
curl -X POST http://localhost:8080/api/webhook \
  -H "Content-Type: application/json" \
  -d '{
    "type": "image",
    "value": "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==",
    "timestamp": "'$(date -u +"%Y-%m-%dT%H:%M:%SZ")'"
  }'

echo ""
echo ""
echo "✅ 测试完成！"
echo "📍 访问 http://localhost:5173 查看结果"
