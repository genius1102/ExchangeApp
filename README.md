# CurrencyExchangeApp

外汇信息展示平台，支持汇率查询、文章浏览与点赞功能。

## 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | Vue 3 + TypeScript + Vite + Element Plus + Vant + Axios + Pinia + Vue Router |
| 后端 | Go + Gin + GORM + Redis + MySQL + JWT |
| 认证 | JWT（Bearer Token，72 小时有效） |

## 项目结构

```
exchangeapp/
├── backend/                          # 后端 Go 服务
│   ├── config/                       # 配置加载（Viper）与数据库/Redis 初始化
│   │   ├── config.go                 # 配置结构体 & 初始化入口
│   │   ├── config.yml                # 配置文件（数据库、Redis、JWT）
│   │   ├── db.go                     # GORM 数据库连接
│   │   └── redis.go                  # Redis 连接
│   ├── controllers/                  # 控制器层（参数绑定 → 调用 Service → 返回 HTTP 响应）
│   │   ├── article_controller.go     # 文章相关接口
│   │   ├── auth_controller.go        # 注册/登录接口
│   │   ├── exchange_rate_controller.go  # 汇率相关接口
│   │   └── likes_controller.go       # 点赞相关接口
│   ├── services/                     # 业务逻辑层
│   │   ├── article_service.go        # 文章 CRUD + 分页缓存（Cache-Aside）
│   │   ├── auth_service.go           # 用户注册/登录
│   │   ├── exchange_rate_service.go  # 汇率查询/创建
│   │   └── like_service.go           # 点赞/取消点赞（Lua 脚本）
│   ├── repository/                   # 数据访问层（GORM）
│   │   ├── article_repository.go     # 文章数据库操作
│   │   ├── user_repository.go        # 用户数据库操作
│   │   └── exchange_rate_repository.go  # 汇率数据库操作
│   ├── models/                       # 数据模型
│   │   ├── article.go                # Article 模型
│   │   ├── exchange_rate.go          # ExchangeRate 模型
│   │   ├── user.go                   # User 模型
│   │   └── page.go                   # 分页响应模型
│   ├── middlewares/                   # 中间件
│   │   └── auth_middleware.go         # JWT 认证中间件
│   ├── router/                        # 路由定义
│   │   └── router.go                 # Gin 路由 & CORS 配置
│   ├── utils/                         # 工具函数
│   │   └── utils.go                  # JWT 生成/解析、bcrypt 密码加密、分页参数解析
│   ├── global/                        # 全局变量（依赖注入容器）
│   │   └── global.go                 # DB、Redis、Service 实例
│   ├── main.go                        # 应用入口
│   ├── go.mod
│   └── go.sum
└── Exchangeapp_frontend/              # 前端 Vue 3 项目
    └── src/
        ├── views/                     # 页面组件
        │   ├── HomeView.vue           # 首页
        │   ├── CurrencyExchangeView.vue  # 货币兑换计算器
        │   ├── NewsView.vue           # 文章列表
        │   ├── NewsDetailView.vue     # 文章详情 + 点赞
        ├── components/                # 可复用组件
        │   ├── Login.vue              # 登录表单
        │   └── Register.vue           # 注册表单
        ├── router/                    # 前端路由
        │   └── index.ts
        ├── store/                     # Pinia 状态管理
        │   └── auth.ts                # 认证状态
        ├── types/                     # TypeScript 类型定义
        │   └── Article.d.ts
        ├── axios.ts                   # Axios 实例（请求拦截、Token 注入）
        ├── App.vue                    # 根组件
        └── main.ts                    # 应用入口
```

## 架构设计

后端采用 **Controller → Service → Repository** 三层架构，通过构造函数注入依赖：

```
HTTP 请求 → Controller（参数绑定）→ Service（业务逻辑）→ Repository（数据访问）→ DB/Redis
```

- **Controller**：仅负责 HTTP 层面的参数绑定与响应返回，不包含业务逻辑
- **Service**：封装核心业务逻辑，通过构造函数注入 Repository 和 Redis 客户端
- **Repository**：封装 GORM 数据库操作，不感知业务
- **Global**：在 `config.InitConfig()` 中完成依赖链组装，Controller 通过 `global.*Svc` 直接使用

## 快速开始

### 环境要求

- Go 1.26+
- Node.js 18+
- MySQL 8.0+
- Redis 6.0+

### 1. 配置数据库与 Redis

修改 `backend/config/config.yml`：

```yaml
app:
  name: CurrencyExchangeApp
  port: 8081
  jwt_secret: "your-secret-key-change-in-production"

database:
  dsn: root:123456@tcp(127.0.0.1:3306)/exchange?charset=utf8mb4&parseTime=True&loc=Local
  MaxIdleConns: 11
  MaxOpenConns: 114

redis:
  addr: "localhost:6379"
  password: "123456"
  db: 1
```

也可以通过环境变量覆盖配置（优先级高于配置文件）：

| 环境变量 | 对应配置 |
|----------|----------|
| `DB_DSN` | 数据库连接串 |
| `REDIS_ADDR` | Redis 地址 |
| `REDIS_PASSWORD` | Redis 密码 |
| `REDIS_DB` | Redis 数据库编号 |
| `JWT_SECRET` | JWT 签名密钥 |

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

### 认证相关

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/auth/register` | 用户注册，返回 JWT Token | 否 |
| POST | `/api/auth/login` | 用户登录，返回 JWT Token | 否 |

### 汇率相关

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/exchangeRates` | 获取全部汇率列表 | 否 |
| POST | `/api/exchangeRates` | 创建汇率记录 | 是 |

### 文章相关

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| GET | `/api/articles` | 分页获取文章列表（支持 `?page=1&size=10`） | 是 |
| POST | `/api/articles` | 创建文章 | 是 |
| GET | `/api/articles/:id` | 获取文章详情 | 是 |

### 点赞相关

| 方法 | 路径 | 说明 | 认证 |
|------|------|------|------|
| POST | `/api/articles/:id/likes` | 切换点赞状态（已赞则取消，未赞则点赞） | 是 |
| GET | `/api/articles/:id/likes` | 获取点赞数及当前用户是否已赞 | 是 |

## 功能亮点

### 分页与缓存策略

- 文章列表支持 `page` / `size` 分页参数，默认每页 10 条，最大 100 条
- 前三页文章使用 **Cache-Aside 模式** 缓存到 Redis（TTL 10 分钟）
- 后续页面直接查询数据库，避免缓存占用过多内存
- 新建文章时主动失效前三页缓存，TTL 作为兜底

### 点赞功能（Lua 脚本）

- 使用 Redis Set 存储每篇文章的点赞用户（去重），使用 String 计数
- 点赞/取消点赞通过 **Lua 脚本** 保证「判断 → 添加/移除 → 更新计数」三步的原子性
- 返回当前用户是否已赞，UI 据此展示点赞状态

### 安全特性

- 密码使用 **bcrypt** 加密存储，不保存明文
- JWT Token 72 小时过期
- 认证中间件拦截未登录请求，返回 401
