# CurrencyExchangeApp

外汇信息展示平台，支持汇率查询、文章浏览与点赞功能。

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + TypeScript + Vite + Element Plus + Vant + Axios + Pinia |
| 后端 | Go + Gin + GORM + Redis + MySQL |
| 认证 | JWT |

## 项目结构

```
exchangeapp/
├── backend/                 # 后端 Go 服务
│   ├── config/              # 配置（数据库、Redis）
│   ├── controllers/         # 控制器
│   ├── global/              # 全局变量
│   ├── middlewares/         # 中间件（JWT 认证）
│   ├── models/              # 数据模型
│   ├── router/              # 路由
│   ├── utils/               # 工具函数（密码加密、JWT）
│   ├── main.go              # 入口
│   ├── go.mod
│   └── go.sum
└── Exchangeapp_frontend/    # 前端 Vue 项目
    └── src/
        ├── views/           # 页面组件
        ├── types/           # TypeScript 类型定义
        ├── axios.ts         # Axios 实例（请求拦截、代理配置）
        └── ...
```

## 快速开始

### 环境要求

- Go 1.26+
- Node.js 18+
- MySQL 8.0+
- Redis 6.0+

### 1. 配置数据库与 Redis

修改 `backend/config/config.yml`：

```yaml
database:
  dsn: root:123456@tcp(127.0.0.1:3306)/exchange?charset=utf8mb4&parseTime=True&loc=Local

redis:
  addr: "localhost:6379"
  password: "123456"
  db: 1
```

### 2. 启动后端

```bash
cd backend
go mod tidy
go run main.go
```

服务默认运行在 `http://localhost:8081`。

### 3. 启动前端

```bash
cd Exchangeapp_frontend
npm install
npm run dev
```

开发服务器运行在 `http://localhost:5173`，API 请求通过 Vite 代理转发到后端 `8081` 端口。

## API 接口

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/auth/register` | 用户注册 | 否 |
| POST | `/api/auth/login` | 用户登录 | 否 |
| GET | `/api/exchangeRates` | 获取汇率列表 | 否 |
| POST | `/api/exchangeRates` | 创建汇率 | 是 |
| GET | `/api/articles` | 获取文章列表（Redis 缓存） | 是 |
| POST | `/api/articles` | 创建文章 | 是 |
| GET | `/api/articles/:id` | 获取文章详情 | 是 |
| POST | `/api/articles/:id/likes` | 点赞文章 | 是 |
| GET | `/api/articles/:id/likes` | 获取点赞数 | 是 |
