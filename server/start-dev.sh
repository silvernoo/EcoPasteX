#!/bin/bash

echo "🚀 启动 EcoPaste Webhook 服务端..."
echo ""

# 检查 Docker 是否安装
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安装，请先安装 Docker"
    exit 1
fi

# 检查 Docker Compose 是否安装
if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose 未安装，请先安装 Docker Compose"
    exit 1
fi

# 启动服务
echo "📦 启动 MongoDB..."
docker-compose up -d mongodb

echo "⏳ 等待 MongoDB 启动..."
sleep 5

echo "🔧 启动后端服务..."
cd backend
go mod download
go run main.go &
BACKEND_PID=$!
cd ..

echo "⏳ 等待后端启动..."
sleep 3

echo "🎨 启动前端服务..."
cd frontend
npm install
npm run dev &
FRONTEND_PID=$!
cd ..

echo ""
echo "✅ 服务启动成功！"
echo ""
echo "📍 访问地址:"
echo "   前端: http://localhost:5173"
echo "   后端: http://localhost:8080"
echo "   MongoDB: localhost:27017"
echo ""
echo "⚙️  配置 EcoPaste Webhook URL: http://localhost:8080/api/webhook"
echo ""
echo "按 Ctrl+C 停止所有服务"

# 等待用户中断
trap "echo ''; echo '🛑 停止服务...'; kill $BACKEND_PID $FRONTEND_PID; docker-compose down; echo '✅ 服务已停止'; exit" INT

wait
