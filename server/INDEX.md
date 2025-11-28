# 📚 EcoPaste Webhook 服务端 - 文档索引

欢迎使用 EcoPaste Webhook 服务端！这里是所有文档的导航页面。

## 🚀 快速开始

如果你是第一次使用，请按以下顺序阅读：

1. **[快速启动指南](QUICKSTART.md)** ⭐ 推荐首先阅读
   - 最快的启动方式
   - 5 分钟内运行起来
   - 包含测试步骤

2. **[项目总结](PROJECT_SUMMARY.md)**
   - 功能完整清单
   - 使用示例
   - 数据展示效果

3. **[完整文档](README.md)**
   - 详细的功能介绍
   - API 文档
   - 配置说明

## 📖 详细文档

### 核心文档

- **[README.md](README.md)** - 完整项目文档
  - 功能特性
  - 技术栈
  - 安装部署
  - API 文档
  - 故障排查

- **[QUICKSTART.md](QUICKSTART.md)** - 快速启动指南
  - 3 种启动方式
  - 配置步骤
  - 测试方法
  - 常见问题

- **[PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)** - 项目总结
  - 已完成功能
  - 文件清单
  - 使用方法
  - 测试建议

### 技术文档

- **[ARCHITECTURE.md](ARCHITECTURE.md)** - 架构说明
  - 项目结构
  - 核心文件说明
  - 数据流
  - 技术栈详情
  - 数据库结构
  - 性能指标
  - 扩展建议

- **[SYSTEM_ARCHITECTURE.md](SYSTEM_ARCHITECTURE.md)** - 系统架构图
  - 完整系统架构
  - 数据流向图
  - 部署架构
  - 技术选型理由

## 🛠️ 实用工具

### 脚本文件

- **[start-dev.sh](start-dev.sh)** - 开发环境启动脚本
  ```bash
  ./start-dev.sh
  ```
  自动启动 MongoDB、后端和前端服务

- **[test-webhook.sh](test-webhook.sh)** - Webhook 测试脚本
  ```bash
  ./test-webhook.sh
  ```
  发送测试数据到服务端

### 配置文件

- **[docker-compose.yml](docker-compose.yml)** - Docker Compose 配置
  - MongoDB 服务
  - 后端服务
  - 前端服务

- **[.gitignore](.gitignore)** - Git 忽略配置

## 📁 代码文件

### 后端代码

- **[backend/main.go](backend/main.go)** - Go 后端主程序
  - Webhook 处理
  - 数据库操作
  - API 路由

- **[backend/go.mod](backend/go.mod)** - Go 模块配置
- **[backend/Dockerfile](backend/Dockerfile)** - 后端容器配置

### 前端代码

- **[frontend/src/App.vue](frontend/src/App.vue)** - Vue 主组件
  - UI 界面
  - 搜索筛选
  - 分页逻辑

- **[frontend/src/main.js](frontend/src/main.js)** - 入口文件
- **[frontend/src/style.css](frontend/src/style.css)** - 样式文件
- **[frontend/package.json](frontend/package.json)** - NPM 依赖
- **[frontend/tailwind.config.js](frontend/tailwind.config.js)** - Tailwind 配置
- **[frontend/Dockerfile](frontend/Dockerfile)** - 前端容器配置
- **[frontend/nginx.conf](frontend/nginx.conf)** - Nginx 配置

## 🎯 按需查找

### 我想...

#### 快速运行项目
→ 阅读 [QUICKSTART.md](QUICKSTART.md)

#### 了解所有功能
→ 阅读 [PROJECT_SUMMARY.md](PROJECT_SUMMARY.md)

#### 查看 API 文档
→ 阅读 [README.md](README.md) 的 API 文档部分

#### 理解项目结构
→ 阅读 [ARCHITECTURE.md](ARCHITECTURE.md)

#### 查看系统架构
→ 阅读 [SYSTEM_ARCHITECTURE.md](SYSTEM_ARCHITECTURE.md)

#### 部署到生产环境
→ 阅读 [README.md](README.md) 的生产部署部分

#### 修改代码
→ 查看 [ARCHITECTURE.md](ARCHITECTURE.md) 的开发说明

#### 解决问题
→ 阅读 [README.md](README.md) 的故障排查部分
→ 阅读 [QUICKSTART.md](QUICKSTART.md) 的常见问题

## 📊 文档统计

- 📄 总文档数: 5 个
- 📝 总代码文件: 10+ 个
- 🛠️ 脚本工具: 2 个
- ⚙️ 配置文件: 6+ 个

## 🔗 相关链接

- **EcoPaste 主项目**: [GitHub](https://github.com/EcoPasteHub/EcoPaste)
- **EcoPaste 官网**: [ecopaste.cn](https://ecopaste.cn)

## 💡 提示

- 所有文档都使用 Markdown 格式
- 可以使用任何 Markdown 阅读器查看
- 建议使用支持目录的编辑器（如 VS Code）
- 代码示例可以直接复制使用

## 🆘 需要帮助？

1. 查看对应的文档
2. 运行测试脚本
3. 查看代码注释
4. 提交 Issue

## 📝 文档更新

最后更新: 2023-11-28

---

**祝使用愉快！** 🎉

如果觉得这个项目有用，请给 EcoPaste 主项目一个 ⭐ Star！
