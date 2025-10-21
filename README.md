# Bluebell 社区论坛项目

<div align="center">

![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)
![Gin](https://img.shields.io/badge/Gin-1.10+-00ADD8?style=flat&logo=go)
![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?style=flat&logo=mysql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-6.2-DC382D?style=flat&logo=redis&logoColor=white)
![Vue](https://img.shields.io/badge/Vue-2.6-4FC08D?style=flat&logo=vue.js&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Supported-2496ED?style=flat&logo=docker&logoColor=white)

一个高性能、易扩展的类 Reddit 社区论坛系统

[功能特性](#功能特性) • [技术栈](#技术栈) • [快速开始](#快速开始) • [架构设计](#架构设计) • [性能优化](#性能优化)

</div>

---

## 项目简介

Bluebell 是一个使用 Go 语言开发的现代化社区论坛系统，采用前后端分离架构。项目灵感来源于 Reddit，实现了完整的社区交互功能，包括用户系统、帖子管理、社区分类、投票排序等核心功能。

**项目亮点：**
-  **高性能优化**：N+1 查询优化、并发查询、多级缓存，性能提升 46.8%
-  **经典架构**：清晰的三层架构设计（Router-Controller-Logic-DAO）
-  **安全可靠**：JWT 认证、参数验证、SQL 注入防护
-  **读写分离**：MySQL 主从架构，轮询负载均衡
-  **易于部署**：Docker Compose 一键部署，完整的 CI/CD 支持
-  **完善文档**：Swagger API 文档，性能测试报告，架构设计文档

---

## 功能特性

### 核心功能

####  用户系统
- 用户注册与登录
- JWT Token 认证
- 密码加密存储
- 参数验证与中文错误提示

####  帖子管理
- 创建帖子（支持 Markdown）
- 帖子列表（支持分页）
- 帖子详情查看
- 多种排序方式（时间、热度）

#### 社区管理
- 社区分类
- 社区列表与详情
- 按社区筛选帖子

#### 投票系统
- 点赞/点踩功能
- 投票数统计
- 热度算法排序
- 防刷票机制（限流保护）

### 高级特性

-  **JWT 认证中间件**：无状态认证，易于扩展
-  **限流保护**：令牌桶算法，防止恶意刷票/刷帖
-  **性能监控**：数据库连接池统计、缓存命中率统计
-  **缓存优化**：Redis 多级缓存、分布式锁防击穿
-  **并发查询**：goroutine 并发优化，响应时间减少 30%
-  **读写分离**：MySQL 主从架构，吞吐量提升 50%
-  **结构化日志**：Zap 日志库，支持日志切割和归档
-  **API 文档**：Swagger 自动生成文档

---

##  技术栈

### 后端技术

| 技术 | 版本 | 用途 |
|------|------|------|
| Go | 1.23+ | 核心开发语言 |
| Gin | 1.10+ | Web 框架 |
| MySQL | 8.0 | 主数据库 |
| Redis | 6.2+ | 缓存与投票数据 |
| Viper | 1.20+ | 配置管理 |
| Zap | 1.27+ | 结构化日志 |
| JWT | 3.2+ | 用户认证 |
| Swagger | 1.16+ | API 文档 |
| Sqlx | 1.4+ | 数据库操作 |
| Snowflake | 0.3+ | 分布式 ID 生成 |

### 前端技术

| 技术 | 版本 | 用途 |
|------|------|------|
| Vue.js | 2.6 | 前端框架 |
| Vue Router | 3.1+ | 路由管理 |
| Vuex | 3.1+ | 状态管理 |
| Axios | 0.19+ | HTTP 请求 |

### 开发工具

- **Docker & Docker Compose**：容器化部署
- **Git**：版本控制
- **Makefile**：构建自动化
- **pprof**：性能分析

---

##  快速开始

### 前置要求

- Go 1.23+
- MySQL 8.0+
- Redis 6.2+
- Node.js 16+ (前端开发)
- Docker & Docker Compose (可选)

### 方式一：Docker Compose 部署（推荐）

```bash
# 克隆项目
git clone https://github.com/yourusername/Bluebell.git
cd Bluebell

# 使用 Docker Compose 启动
cd Bluebell
docker-compose up -d

# 访问服务
# 后端 API: http://localhost:8081
# Swagger 文档: http://localhost:8081/swagger/index.html
```

### 方式二：本地开发部署

#### 1. 后端部署

```bash
# 进入后端目录
cd Bluebell

# 安装依赖
go mod download

# 配置数据库（编辑 conf/config.yaml）
vim conf/config.yaml

# 导入数据库
mysql -u root -p < init.sql

# 编译运行
make build
./web-app ./conf/config.yaml

# 或直接运行
go run main.go ./conf/config.yaml
```

#### 2. 前端部署

```bash
# 进入前端目录
cd Bluebell_frontend

# 安装依赖
npm install

# 开发模式运行
npm run serve

# 生产构建
npm run build
```

### 访问项目

- **前端页面**: http://localhost:8080
- **后端 API**: http://localhost:8081
- **Swagger 文档**: http://localhost:8081/swagger/index.html

---

##  架构设计

### 整体架构

```
┌─────────────────────────────────────────────────────────┐
│                        Client                            │
│                   (Vue.js Frontend)                      │
└────────────────────┬────────────────────────────────────┘
                     │ HTTP/HTTPS
                     ↓
┌─────────────────────────────────────────────────────────┐
│                    Gin Web Server                        │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Middlewares (JWT / RateLimit / Logger / CORS)   │  │
│  └──────────────────────────────────────────────────┘  │
│                         ↓                                │
│  ┌──────────────────────────────────────────────────┐  │
│  │         Router (路由层)                           │  │
│  └──────────────────────────────────────────────────┘  │
│                         ↓                                │
│  ┌──────────────────────────────────────────────────┐  │
│  │     Controller (参数验证 + 响应格式)              │  │
│  └──────────────────────────────────────────────────┘  │
│                         ↓                                │
│  ┌──────────────────────────────────────────────────┐  │
│  │      Logic (业务逻辑 + 并发控制)                  │  │
│  └──────────────────────────────────────────────────┘  │
│                         ↓                                │
│  ┌──────────────────────────────────────────────────┐  │
│  │         DAO (数据访问层)                          │  │
│  │    ┌──────────────┐      ┌──────────────┐       │  │
│  │    │  MySQL DAO   │      │  Redis DAO   │       │  │
│  │    └──────────────┘      └──────────────┘       │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
                         ↓
        ┌────────────────┴────────────────┐
        ↓                                  ↓
┌──────────────┐                  ┌──────────────┐
│    MySQL     │                  │    Redis     │
│  (Master)    │                  │   (Cache)    │
│      ↓       │                  └──────────────┘
│   Slave-1    │
│   Slave-2    │
└──────────────┘
```

### 分层职责

**Router 层**
- 路由注册和分组
- 中间件挂载
- API 版本管理

**Controller 层**
- 参数绑定与验证
- 统一响应格式
- 错误码管理

**Logic 层**
- 核心业务逻辑
- 数据组装
- 并发控制
- 事务管理

**DAO 层**
- MySQL CRUD 操作
- Redis 缓存操作
- 批量查询优化
- 读写分离实现

---

##  性能优化

### 1. N+1 查询优化

**问题**：帖子列表接口存在 N+1 查询问题，获取 10 个帖子需要执行 21 次数据库查询。

**解决方案**：批量 IN 查询优化

```go
// 优化前：21 次查询
posts := GetPostList(10)          // 1 次
for _, post := range posts {
    user := GetUserByID(...)       // 10 次
    community := GetCommunityByID(...) // 10 次
}

// 优化后：3 次查询
posts := GetPostList(10)           // 1 次
userMap := BatchGetUsers(userIDs)  // 1 次 (IN 查询)
commMap := BatchGetComm(commIDs)   // 1 次 (IN 查询)
```

**效果**：
- 查询次数：21 次 → 3 次
- 性能提升：46.8%
- 响应时间：200ms → 106ms

### 2. 并发查询优化

**问题**：帖子详情需要查询帖子、用户、社区三张表，串行查询耗时长。

**解决方案**：goroutine 并发查询

```go
// 优化前：串行查询 150ms
post := GetPost(id)         // 50ms
user := GetUser(...)        // 50ms
community := GetComm(...)   // 50ms

// 优化后：并发查询 100ms
var wg sync.WaitGroup
wg.Add(2)
go func() { user = GetUser(...) }()
go func() { community = GetComm(...) }()
wg.Wait()
```

**效果**：
- 响应时间：150ms → 100ms（减少 33%）
- 并发能力提升

### 3. Redis 多级缓存

**问题**：高频查询给数据库带来压力，热点数据缓存失效导致缓存击穿。

**解决方案**：多级缓存 + 分布式锁

```go
// 查询流程
1. 查询 Redis 缓存
2. 命中 → 直接返回
3. 未命中 → 获取分布式锁
4. 双重检查缓存
5. 查询数据库
6. 异步更新缓存
```

**效果**：
- 缓存命中率：60-80%
- 缓存命中响应时间：< 10ms
- 防止缓存击穿

### 4. MySQL 读写分离

**架构**：1 主 2 从，轮询负载均衡

**实现**：
- 写操作 → 主库
- 读操作 → 从库（轮询）
- 主从同步 → Binlog 复制

**效果**：
- 读 QPS 提升：50%
- 主库压力减少：50%

### 5. 限流保护

**算法**：令牌桶算法

**实现**：
```go
// 投票接口：2 秒/次
RateLimitMiddleware(2*time.Second, 1)

// 发帖接口：10 秒/次
RateLimitMiddleware(10*time.Second, 1)
```

**效果**：
- 防止恶意刷票/刷帖
- 保护系统稳定性

### 性能指标汇总

| 优化项 | 优化前 | 优化后 | 提升 |
|--------|--------|--------|------|
| 帖子列表查询次数 | 21 次 | 3 次 | 85.7% |
| 帖子列表响应时间 | 200ms | 106ms | 46.8% |
| 帖子详情响应时间 | 150ms | 100ms | 33.3% |
| 缓存命中响应时间 | - | < 10ms | - |
| 读 QPS | 1000 | 1500 | 50% |
| 缓存命中率 | 0% | 60-80% | - |

---

##  项目结构

```
Bluebell/
├── Bluebell/                    # 后端项目
│   ├── main.go                  # 入口文件
│   ├── conf/                    # 配置文件
│   │   ├── config.yaml          # 本地配置
│   │   └── config.docker.yaml   # Docker 配置
│   ├── router/                  # 路由层
│   │   └── router.go
│   ├── controller/              # 控制器层
│   │   ├── user.go              # 用户相关
│   │   ├── post.go              # 帖子相关
│   │   ├── community.go         # 社区相关
│   │   ├── vote.go              # 投票相关
│   │   ├── response.go          # 统一响应
│   │   └── code.go              # 错误码定义
│   ├── logic/                   # 业务逻辑层
│   │   ├── user.go
│   │   ├── post.go
│   │   ├── community.go
│   │   └── vote.go
│   ├── dao/                     # 数据访问层
│   │   ├── mysql/               # MySQL 操作
│   │   │   ├── mysql.go         # 连接池
│   │   │   ├── user.go
│   │   │   ├── post.go
│   │   │   └── community.go
│   │   └── redis/               # Redis 操作
│   │       ├── redis.go
│   │       ├── keys.go
│   │       └── vote.go
│   ├── middlewares/             # 中间件
│   │   ├── auth.go              # JWT 认证
│   │   └── ratelimit.go         # 限流
│   ├── models/                  # 数据模型
│   │   └── params.go
│   ├── pkg/                     # 工具包
│   │   ├── jwt/                 # JWT 工具
│   │   └── snowflake/           # ID 生成
│   ├── logger/                  # 日志模块
│   ├── settings/                # 配置解析
│   ├── docs/                    # Swagger 文档
│   ├── docker-compose.yml       # Docker 编排
│   ├── Dockerfile              # Docker 镜像
│   ├── Makefile                # 构建脚本
│   ├── init.sql                # 数据库初始化
│   └── go.mod                  # Go 依赖
│
├── Bluebell_frontend/          # 前端项目
│   ├── src/
│   │   ├── views/              # 页面组件
│   │   ├── components/         # 通用组件
│   │   ├── router/             # 路由配置
│   │   ├── store/              # Vuex 状态
│   │   └── api/                # API 接口
│   ├── public/
│   ├── package.json
│   └── vue.config.js
│
└── README.md                   # 项目说明文档
```

---

## 📊 数据库设计

### 核心表结构

**用户表 (user)**
```sql
CREATE TABLE `user` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `user_id` bigint NOT NULL COMMENT '用户ID',
  `username` varchar(64) NOT NULL COMMENT '用户名',
  `password` varchar(64) NOT NULL COMMENT '密码',
  `email` varchar(64) DEFAULT NULL,
  `gender` tinyint DEFAULT 0,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`),
  UNIQUE KEY `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**社区表 (community)**
```sql
CREATE TABLE `community` (
  `id` int NOT NULL AUTO_INCREMENT,
  `community_id` int NOT NULL,
  `community_name` varchar(128) NOT NULL,
  `introduction` varchar(256) NOT NULL,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_community_id` (`community_id`),
  UNIQUE KEY `idx_community_name` (`community_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

**帖子表 (post)**
```sql
CREATE TABLE `post` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `post_id` bigint NOT NULL COMMENT '帖子ID',
  `title` varchar(128) NOT NULL COMMENT '标题',
  `content` text NOT NULL COMMENT '内容',
  `author_id` bigint NOT NULL COMMENT '作者ID',
  `community_id` bigint NOT NULL COMMENT '社区ID',
  `status` tinyint NOT NULL DEFAULT 1,
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_post_id` (`post_id`),
  KEY `idx_author_id` (`author_id`),
  KEY `idx_community_id` (`community_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

---

## 🔧 配置说明

### config.yaml 配置文件

```yaml
name: "bluebell"
mode: "dev"
port: 8081
version: "1.0.0"
start_time: "2020-01-01"
machine_id: 1

log:
  level: "debug"
  filename: "web-app.log"
  max_size: 200
  max_age: 30
  max_backups: 7

mysql:
  host: "127.0.0.1"
  port: 3306
  user: "root"
  password: "your_password"
  dbname: "bluebell"
  max_open_conns: 200
  max_idle_conns: 50
  enable_read_write_split: true  # 读写分离
  read_hosts:
    - "127.0.0.1:3307"
    - "127.0.0.1:3308"

redis:
  host: "127.0.0.1"
  port: 6379
  password: "your_password"
  db: 0
  pool_size: 100
```

---

##  API 文档

### Swagger 文档

项目集成了 Swagger 自动生成 API 文档，启动后访问：

```
http://localhost:8081/swagger/index.html
```

### 核心 API 列表

#### 用户相关

| 方法 | 路径 | 说明 | 是否需要认证 |
|------|------|------|-------------|
| POST | /api/v1/signup | 用户注册 | ❌ |
| POST | /api/v1/login | 用户登录 | ❌ |

#### 帖子相关

| 方法 | 路径 | 说明 | 是否需要认证 |
|------|------|------|-------------|
| GET | /api/v1/posts | 帖子列表（原版） | ❌ |
| GET | /api/v1/posts/optimized | 帖子列表（N+1优化） | ❌ |
| GET | /api/v1/posts/cached | 帖子列表（缓存版） | ❌ |
| GET | /api/v1/post/:id | 帖子详情 | ❌ |
| GET | /api/v1/post/:id/concurrent | 帖子详情（并发优化） | ❌ |
| POST | /api/v1/post | 创建帖子 | ✅ |

#### 社区相关

| 方法 | 路径 | 说明 | 是否需要认证 |
|------|------|------|-------------|
| GET | /api/v1/community | 社区列表 | ❌ |
| GET | /api/v1/community/:id | 社区详情 | ❌ |

#### 投票相关

| 方法 | 路径 | 说明 | 是否需要认证 |
|------|------|------|-------------|
| POST | /api/v1/vote | 帖子投票 | ✅ |

#### 监控相关

| 方法 | 路径 | 说明 | 是否需要认证 |
|------|------|------|-------------|
| GET | /api/v1/db/stats | 数据库连接池统计 | ❌ |
| GET | /api/v1/db/health | 健康检查 | ❌ |
| GET | /api/v1/cache/stats | 缓存统计 | ❌ |

---

##  测试

### 性能测试

项目包含完整的性能测试报告：

- `PERFORMANCE_TEST_REPORT.md` - 综合性能测试报告
- `CONCURRENT_OPTIMIZATION.md` - 并发优化测试报告
- `DATABASE_POOL_OPTIMIZATION_REPORT.md` - 数据库连接池优化报告
- `REDIS_CACHE_OPTIMIZATION_REPORT.md` - Redis 缓存优化报告
- `FINAL_OPTIMIZATION_REPORT.md` - 最终优化总结报告

### 使用 wrk 进行压测

```bash
# 测试帖子列表接口
wrk -t12 -c400 -d30s http://localhost:8081/api/v1/posts/optimized

# 测试帖子详情接口
wrk -t12 -c400 -d30s http://localhost:8081/api/v1/post/1/cached
```

---

## 📝 开发文档

项目提供了完整的开发文档：

- **架构笔记-思维导图版.md** - 项目架构思维导图
- **架构笔记-速查版.md** - 快速查阅手册
- **架构笔记-面试复习版.md** - 面试准备材料
- **面试题库-Bluebell项目.md** - 常见面试问题及答案

---

##  Docker 部署

### 使用 Docker Compose

```bash
# 启动所有服务（MySQL + Redis + App）
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down

# 停止并删除数据卷
docker-compose down -v
```

### 单独构建镜像

```bash
# 构建镜像
docker build -t bluebell:latest .

# 运行容器
docker run -d -p 8081:8081 bluebell:latest
```

---

## 监控与调试

### 性能分析 (pprof)

项目集成了 Go pprof 性能分析工具：

```bash
# 访问 pprof Web 界面
http://localhost:8081/debug/pprof/

# CPU 分析
go tool pprof http://localhost:8081/debug/pprof/profile

# 内存分析
go tool pprof http://localhost:8081/debug/pprof/heap

# goroutine 分析
go tool pprof http://localhost:8081/debug/pprof/goroutine
```

### 日志查看

```bash
# 实时查看日志
tail -f web-app.log

# 查看最近 100 行
tail -n 100 web-app.log
```

### 数据库监控

访问数据库连接池统计接口：

```bash
curl http://localhost:8081/api/v1/db/stats
```

返回示例：
```json
{
  "code": 1000,
  "msg": "success",
  "data": {
    "max_open_connections": 200,
    "open_connections": 45,
    "in_use": 12,
    "idle": 33,
    "wait_count": 0,
    "wait_duration": "0s"
  }
}
```

---

## 🤝贡献指南

欢迎贡献代码、提出问题和建议！

### 开发流程

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 提交 Pull Request

### 代码规范

- 遵循 Go 官方代码规范
- 使用 `gofmt` 格式化代码
- 添加必要的注释和文档
- 编写单元测试

---

## 📜 许可证

本项目采用 Apache 2.0 许可证 - 详见 [LICENSE](LICENSE) 文件

---


## 🙏 致谢

- [Gin](https://github.com/gin-gonic/gin) - 优秀的 Go Web 框架
- [go-redis](https://github.com/go-redis/redis) - Redis 客户端
- [sqlx](https://github.com/jmoiron/sqlx) - 数据库工具
- [Viper](https://github.com/spf13/viper) - 配置管理
- [Zap](https://github.com/uber-go/zap) - 高性能日志库

---



---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给一个 Star！⭐**


</div>

