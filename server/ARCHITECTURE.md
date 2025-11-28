# EcoPaste Webhook 服务端 - 项目结构

```
server/
│
├── 📁 backend/                      # Go 后端服务
│   ├── main.go                      # 主程序文件
│   ├── go.mod                       # Go 模块依赖
│   ├── go.sum                       # 依赖校验文件（自动生成）
│   └── Dockerfile                   # 后端 Docker 镜像配置
│
├── 📁 frontend/                     # Vue 前端应用
│   ├── 📁 src/
│   │   ├── App.vue                  # 主组件（包含所有 UI 和逻辑）
│   │   ├── main.js                  # Vue 应用入口
│   │   └── style.css                # 全局样式（Tailwind CSS）
│   ├── 📁 public/                   # 静态资源
│   ├── index.html                   # HTML 模板
│   ├── package.json                 # NPM 依赖配置
│   ├── vite.config.js               # Vite 构建配置
│   ├── tailwind.config.js           # Tailwind CSS 配置
│   ├── postcss.config.js            # PostCSS 配置
│   ├── Dockerfile                   # 前端 Docker 镜像配置
│   └── nginx.conf                   # Nginx 服务器配置
│
├── docker-compose.yml               # Docker Compose 编排文件
├── README.md                        # 完整项目文档
├── QUICKSTART.md                    # 快速启动指南
├── start-dev.sh                     # 开发环境启动脚本
├── test-webhook.sh                  # Webhook 测试脚本
└── .gitignore                       # Git 忽略文件配置

```

## 核心文件说明

### 后端 (backend/)

#### `main.go` - 主程序
- **数据模型**: `ClipboardItem`, `WebhookPayload`
- **API 端点**:
  - `POST /api/webhook` - 接收 webhook 数据
  - `GET /api/clipboard` - 获取剪贴板列表（支持分页、搜索、筛选）
  - `DELETE /api/clipboard/:id` - 删除记录
- **数据库**: MongoDB 连接和操作
- **CORS**: 跨域配置
- **功能**:
  - HTML 标签清理
  - 图片预览提取
  - 文本预览生成
  - 自动索引创建

### 前端 (frontend/)

#### `src/App.vue` - 主组件
- **UI 组件**:
  - 搜索栏
  - 类型筛选按钮（全部/文本/图片）
  - 数据卡片列表
  - 分页控件
- **功能**:
  - 防抖搜索
  - 实时筛选
  - 图片预览
  - 删除确认
  - 自动刷新（30 秒）
  - 响应式布局
- **样式**: Tailwind CSS 实用类

## 数据流

```
EcoPaste 客户端
    ↓ (复制内容)
    ↓ POST /api/webhook
后端 API (Go + Gin)
    ↓ (保存数据)
MongoDB 数据库
    ↑ (查询数据)
    ↑ GET /api/clipboard
前端 UI (Vue + Tailwind)
    ↓ (显示给用户)
浏览器
```

## 技术栈详情

### 后端技术
- **语言**: Go 1.21+
- **框架**: Gin (Web 框架)
- **数据库**: MongoDB 7.0
- **驱动**: mongo-driver (官方)
- **CORS**: gin-contrib/cors

### 前端技术
- **框架**: Vue 3 (Composition API)
- **构建**: Vite 7
- **样式**: Tailwind CSS 3
- **HTTP**: Axios
- **图标**: SVG 内联

### 部署技术
- **容器**: Docker
- **编排**: Docker Compose
- **Web 服务器**: Nginx (生产环境)

## 端口分配

| 服务 | 开发端口 | 生产端口 |
|------|---------|---------|
| MongoDB | 27017 | 27017 |
| 后端 API | 8080 | 8080 |
| 前端开发 | 5173 | - |
| 前端生产 | - | 3000 (Nginx) |

## 环境要求

### 开发环境
- Go 1.21+
- Node.js 20+
- MongoDB 7.0+
- Docker (可选)

### 生产环境
- Docker
- Docker Compose

## 数据库结构

### Collection: `clipboard`

```javascript
{
  _id: ObjectId("..."),           // MongoDB 自动生成
  type: "text" | "image" | ...,   // 内容类型
  value: "...",                    // 实际内容
  timestamp: ISODate("..."),       // 时间戳
  subtype: "plain" | "html" | ..., // 子类型（可选）
  isImage: true | false,           // 是否为图片
  preview: "..."                   // 预览文本
}
```

### 索引
- `timestamp: -1` (降序) - 用于快速查询最新记录

## API 性能

- **分页**: 默认 100 条/页，最大 100 条
- **搜索**: 使用 MongoDB 正则表达式，不区分大小写
- **索引**: 时间戳字段已建立索引
- **响应时间**: 通常 < 100ms

## 安全考虑

### 当前实现
- CORS 配置（仅允许本地开发）
- 输入验证
- HTML 标签清理

### 生产环境建议
- [ ] 添加身份验证（JWT）
- [ ] API 限流
- [ ] HTTPS 加密
- [ ] 环境变量管理
- [ ] 日志记录
- [ ] 错误监控

## 扩展建议

### 功能扩展
- [ ] 用户系统
- [ ] 标签分类
- [ ] 收藏功能
- [ ] 导出功能
- [ ] 批量操作
- [ ] 数据统计

### 技术优化
- [ ] Redis 缓存
- [ ] 全文搜索引擎（Elasticsearch）
- [ ] CDN 加速
- [ ] 负载均衡
- [ ] 数据备份

## 开发工作流

1. **本地开发**: 使用 `start-dev.sh` 或手动启动
2. **测试**: 使用 `test-webhook.sh` 发送测试数据
3. **调试**: 查看控制台日志
4. **构建**: `docker-compose build`
5. **部署**: `docker-compose up -d`

## 维护建议

- 定期清理旧数据
- 监控 MongoDB 性能
- 备份数据库
- 更新依赖包
- 查看错误日志
