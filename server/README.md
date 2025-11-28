# EcoPaste Webhook 服务端

这是 EcoPaste 剪贴板管理工具的 Webhook 服务端，使用 Gin + Vue + MongoDB 构建。

## 功能特性

- 📥 **Webhook 接收**: 接收来自 EcoPaste 客户端的剪贴板数据
- 🔍 **搜索功能**: 支持全文搜索剪贴板内容
- 🏷️ **类型筛选**: 可以筛选文本和图片类型
- 📄 **分页展示**: 每页展示 100 条记录，支持翻页
- 🖼️ **图片预览**: 支持 base64 图片的直接预览
- 🗑️ **删除功能**: 可以删除不需要的历史记录
- 🎨 **现代化 UI**: 使用 Tailwind CSS 构建的美观界面

## 技术栈

### 后端
- **Go 1.21+**: 高性能后端语言
- **Gin**: Web 框架
- **MongoDB**: NoSQL 数据库
- **MongoDB Driver**: Go 官方驱动

### 前端
- **Vue 3**: 渐进式 JavaScript 框架
- **Vite**: 下一代前端构建工具
- **Tailwind CSS**: 实用优先的 CSS 框架
- **Axios**: HTTP 客户端

## 项目结构

```
server/
├── backend/                # Go 后端
│   ├── main.go            # 主程序
│   ├── go.mod             # Go 模块文件
│   └── Dockerfile         # 后端 Docker 配置
├── frontend/              # Vue 前端
│   ├── src/
│   │   ├── App.vue        # 主组件
│   │   ├── main.js        # 入口文件
│   │   └── style.css      # 样式文件
│   ├── package.json       # NPM 依赖
│   ├── Dockerfile         # 前端 Docker 配置
│   └── nginx.conf         # Nginx 配置
└── docker-compose.yml     # Docker Compose 配置
```

## 快速开始

### 方式一：使用 Docker Compose（推荐）

1. 确保已安装 Docker 和 Docker Compose

2. 启动所有服务：
```bash
cd server
docker-compose up -d
```

3. 访问应用：
   - 前端: http://localhost:3000
   - 后端 API: http://localhost:8080
   - MongoDB: localhost:27017

4. 停止服务：
```bash
docker-compose down
```

### 方式二：本地开发

#### 启动 MongoDB

```bash
# 使用 Docker 启动 MongoDB
docker run -d -p 27017:27017 --name mongodb mongo:7.0

# 或使用本地安装的 MongoDB
mongod --dbpath /path/to/data
```

#### 启动后端

```bash
cd server/backend

# 安装依赖
go mod download

# 运行服务
go run main.go
```

后端将在 http://localhost:8080 启动

#### 启动前端

```bash
cd server/frontend

# 安装依赖
npm install

# 启动开发服务器
npm run dev
```

前端将在 http://localhost:5173 启动

## API 文档

### 1. Webhook 接收

接收来自 EcoPaste 客户端的剪贴板数据。

**端点**: `POST /api/webhook`

**请求体**:
```json
{
  "type": "text",
  "value": "剪贴板内容",
  "timestamp": "2023-11-28T10:00:00Z",
  "subtype": "plain"
}
```

**响应**:
```json
{
  "message": "Data received successfully"
}
```

### 2. 获取剪贴板列表

获取分页的剪贴板历史记录。

**端点**: `GET /api/clipboard`

**查询参数**:
- `page`: 页码（默认: 1）
- `pageSize`: 每页数量（默认: 100，最大: 100）
- `type`: 类型筛选（all/text/image）
- `search`: 搜索关键词

**响应**:
```json
{
  "items": [
    {
      "id": "507f1f77bcf86cd799439011",
      "type": "text",
      "value": "内容",
      "timestamp": "2023-11-28T10:00:00Z",
      "isImage": false,
      "preview": "内容预览..."
    }
  ],
  "total": 150,
  "page": 1,
  "pageSize": 100,
  "totalPages": 2
}
```

### 3. 删除记录

删除指定的剪贴板记录。

**端点**: `DELETE /api/clipboard/:id`

**响应**:
```json
{
  "message": "Item deleted successfully"
}
```

## EcoPaste 客户端配置

在 EcoPaste 偏好设置中配置 Webhook：

1. 打开 EcoPaste
2. 进入 **偏好设置** → **Webhook**
3. 启用 Webhook
4. 设置 URL: `http://localhost:8080/api/webhook`
5. 保存设置

现在，每次复制内容时，数据都会自动发送到服务端。

## 数据示例

### 文本数据

```json
{
  "type": "text",
  "value": "<html><head>...</head><body>CollectOneListAdapter</body></html>",
  "timestamp": "2023-11-28T10:00:00Z",
  "subtype": "html"
}
```

在前端会显示为：
- **类型**: 文本
- **预览**: CollectOneListAdapter（已去除 HTML 标签）

### 图片数据

```json
{
  "type": "image",
  "value": "data:image/png;base64,iVBORw0KGgo...",
  "timestamp": "2023-11-28T10:00:00Z"
}
```

在前端会显示为：
- **类型**: 图片
- **预览**: 直接显示图片

## 环境变量

### 后端

可以通过环境变量配置 MongoDB 连接：

```bash
MONGODB_URI=mongodb://localhost:27017
```

### 前端

在生产环境中，需要修改 API 地址。编辑 `frontend/src/App.vue`:

```javascript
const API_BASE_URL = 'http://your-backend-url/api'
```

## 生产部署

### 构建前端

```bash
cd server/frontend
npm run build
```

构建产物在 `dist` 目录。

### 构建后端

```bash
cd server/backend
go build -o ecopaste-server main.go
```

### 使用 Docker 部署

```bash
cd server
docker-compose up -d --build
```

## 开发说明

### 添加新功能

1. **后端**: 在 `backend/main.go` 中添加新的路由和处理函数
2. **前端**: 在 `frontend/src/App.vue` 中添加新的 UI 和逻辑

### 数据库索引

系统会自动在 `timestamp` 字段上创建索引，以优化查询性能。

### 自动刷新

前端每 30 秒自动刷新一次数据，确保显示最新内容。

## 故障排查

### MongoDB 连接失败

检查 MongoDB 是否正在运行：
```bash
docker ps | grep mongodb
```

### 后端启动失败

检查端口 8080 是否被占用：
```bash
lsof -i :8080
```

### 前端无法连接后端

1. 检查后端是否正在运行
2. 检查 CORS 配置
3. 确认 API_BASE_URL 设置正确

## 性能优化

- MongoDB 使用索引加速查询
- 前端使用防抖优化搜索
- 图片使用懒加载
- Nginx 启用 gzip 压缩

## 安全建议

在生产环境中：

1. 使用 HTTPS
2. 添加身份验证
3. 限制 API 访问频率
4. 使用环境变量管理敏感信息
5. 定期备份 MongoDB 数据

## 许可证

本项目采用与 EcoPaste 相同的许可证。

## 贡献

欢迎提交 Issue 和 Pull Request！

## 联系方式

如有问题，请在 EcoPaste 主仓库提交 Issue。
