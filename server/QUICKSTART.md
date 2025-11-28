# 快速启动指南

## 🚀 最快启动方式

### 1. 启动 MongoDB

```bash
# 使用 Docker 启动 MongoDB（推荐）
docker run -d -p 27017:27017 --name ecopaste-mongodb mongo:7.0
```

### 2. 启动后端

```bash
cd server/backend

# 首次运行需要下载依赖
go mod download

# 启动后端服务
go run main.go
```

后端将在 **http://localhost:8080** 启动

### 3. 启动前端

```bash
cd server/frontend

# 首次运行需要安装依赖
npm install

# 启动前端开发服务器
npm run dev
```

前端将在 **http://localhost:5173** 启动

### 4. 配置 EcoPaste

1. 打开 EcoPaste 应用
2. 进入 **偏好设置** → **Webhook**
3. 启用 Webhook
4. 设置 URL: `http://localhost:8080/api/webhook`
5. 保存设置

### 5. 测试

复制一些文本或图片，然后访问 http://localhost:5173 查看结果！

---

## 🧪 使用测试脚本

如果想快速测试 webhook 功能：

```bash
cd server
./test-webhook.sh
```

这将发送几条测试数据到服务端。

---

## 🐳 使用 Docker Compose（生产环境）

```bash
cd server
docker-compose up -d
```

这将启动：
- MongoDB (端口 27017)
- 后端服务 (端口 8080)
- 前端服务 (端口 3000)

访问 http://localhost:3000 即可使用。

---

## 📊 数据展示说明

### 文本内容

HTML 格式的文本会自动去除标签，只显示纯文本内容。例如：

**原始数据**:
```html
<html><head></head><body>CollectOneListAdapter</body></html>
```

**显示效果**:
```
CollectOneListAdapter
```

### 图片内容

支持两种格式：
1. **Base64 图片**: 直接在页面中预览
2. **HTML 中的图片**: 自动提取 `<img>` 标签中的 src 并显示

---

## 🔧 常见问题

### Q: MongoDB 连接失败？
A: 确保 MongoDB 正在运行：
```bash
docker ps | grep mongo
```

### Q: 端口被占用？
A: 修改端口配置：
- 后端: 在 `backend/main.go` 中修改 `:8080`
- 前端: 在 `frontend/vite.config.js` 中配置端口

### Q: 前端无法连接后端？
A: 检查 `frontend/src/App.vue` 中的 `API_BASE_URL` 是否正确。

---

## 📝 功能清单

- ✅ Webhook 接收剪贴板数据
- ✅ 文本和图片类型筛选
- ✅ 全文搜索功能
- ✅ 分页展示（每页 100 条）
- ✅ 图片预览
- ✅ 删除记录
- ✅ 自动刷新（30 秒）
- ✅ 响应式设计
- ✅ 暗色主题

---

## 🎯 下一步

1. 尝试复制不同类型的内容
2. 使用搜索功能查找历史记录
3. 测试图片预览功能
4. 体验分页和筛选功能

祝使用愉快！ 🎉
