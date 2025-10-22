# Bluebell Community Forum Project

<div align="center">
![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)
![Gin](https://img.shields.io/badge/Gin-1.10+-00ADD8?style=flat&logo=go)
![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?style=flat&logo=mysql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-6.2-DC382D?style=flat&logo=redis&logoColor=white)
![Vue](https://img.shields.io/badge/Vue-2.6-4FC08D?style=flat&logo=vue.js&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Supported-2496ED?style=flat&logo=docker&logoColor=white)

A high-performance, scalable Reddit-like community forum system

[Features](#Features) • [Technology Stack](#Technology Stack) • [Quick Start](#Quick Start) • [Architecture Design](#Architecture Design) • [Performance optimization](#Performance optimization)

</div>

---

## Project Introduction

Bluebell is a modern community forum system developed in Go, featuring a front-end and back-end separation architecture. Inspired by Reddit, the project implements comprehensive community interaction functionality, including user systems, post management, community classification, and voting ranking.

**Project Highlights：**

-  **High-Performance Optimization**: N+1 query optimization, concurrent queries, and multi-level caching improve performance by 46.8%
- **Classic Architecture**: A clear three-tier architecture (Router-Controller-Logic-DAO)
- **Security and Reliability**: JWT authentication, parameter validation, and SQL injection protection
- **Read-Write Separation**: MySQL master-slave architecture with round-robin load balancing
- **Easy Deployment**: One-click deployment with Docker Compose and comprehensive CI/CD support
- **Complete Documentation**: Swagger API documentation, performance test reports, and architectural design documentation

---

- ## Features

  ### Core Functionality

  ####  User System

  - User Registration and Login
  - JWT Token Authentication
  - Encrypted Password Storage
  - Parameter Validation and Chinese Error Messages

  ####  Post Management

  - Create Post (Markdown Support)
  - Post List (Pagination Support)
  - View Post Details
  - Multiple Sorting Methods (Time, Popularity)

  #### Community Management

  - Community Categories

  - Community List and Details
  - Filter Posts by Community

  #### Voting System

  - Like/Dislike Functionality

  - Vote Statistics
  - Sorting by Popularity Algorithm
  - Anti-Vote Spam Mechanism (Current Limiting Protection)

  ### Advanced Features

  -  **JWT Authentication Middleware:** Stateless Authentication, Easily Scalable
  -  **Current Limiting Protection:** Token Bucket Algorithm to Prevent Malicious Voting/Post Spam
  -  **Performance Monitoring:** Database Connection Pool Statistics, Cache Hit Rate Statistics
  -  **Cache Optimization:** Redis Multi-Level Caching, Distributed Locking for Anti-Breach
  -  **Concurrent Queries: **Goroutine Concurrency Optimization, Response Time Reduction by 30%
  -  **Read/Write Separation:** MySQL Master-Slave Architecture, Throughput Increased by 50%
  -  **Structured Logging:** Zap Logging library, supporting log segmentation and archiving
  -  **API documentation:** Swagger automatically generates documentation

---

##  Technology Stack

### Backend Technology

| Technology | Version | Use |
|------|------|------|
| Go | 1.23+ | Core development language |
| Gin | 1.10+ | Web frame |
| MySQL | 8.0 | Master database |
| Redis | 6.2+ | Caching and voting data |
| Viper | 1.20+ | Configuration Management |
| Zap | 1.27+ | Structured logs |
| JWT | 3.2+ | User authentication |
| Swagger | 1.16+ | API document |
| Sqlx | 1.4+ | Database operations |
| Snowflake | 0.3+ | Distributed ID generation |

### Front-end technology

| Technology | Version | Use |
|------|------|------|
| Vue.js | 2.6 | Front-end framework |
| Vue Router | 3.1+ | Route Management |
| Vuex | 3.1+ | State Management |
| Axios | 0.19+ | HTTP Request |

### Development Tools

- **Docker & Docker Compose**：Containerized deployment
- **Git**：Version Control
- **Makefile**：Build Automation
- **pprof**：Performance Analysis

---

##  Quick Start

### Prerequisites

- Go 1.23+
- MySQL 8.0+
- Redis 6.2+
- Node.js 16+ (Front-end development)
- Docker & Docker Compose (Optional)

### Method 1: Docker Compose deployment (recommended)

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

### Method 2: Local development and deployment

#### 1. Backend deployment

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

#### 2. Front-end deployment

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

### Visit Project

- **Front-end page**: http://localhost:8080
- **Backend API**: http://localhost:8081
- **Swagger document**: http://localhost:8081/swagger/index.html

---

##  Architecture Design

### Overall architecture

```
┌─────────────────────────────────────────────────────────┐
│                        Client                           │
│                   (Vue.js Frontend)                     │
└────────────────────┬────────────────────────────────────┘
                     │ HTTP/HTTPS
                     ↓
┌─────────────────────────────────────────────────────────┐
│                    Gin Web Server                       │
│  ┌──────────────────────────────────────────────────┐   │
│  │  Middlewares (JWT / RateLimit / Logger / CORS)   │   │
│  └──────────────────────────────────────────────────┘   │
│                         ↓                               │
│  ┌──────────────────────────────────────────────────┐   │
│  │         Router (路由层)                           │   │
│  └──────────────────────────────────────────────────┘   │
│                         ↓                               │
│  ┌──────────────────────────────────────────────────┐   │
│  │     Controller (参数验证 + 响应格式)                │   │
│  └──────────────────────────────────────────────────┘   │
│                         ↓                               │
│  ┌──────────────────────────────────────────────────┐   │
│  │      Logic (业务逻辑 + 并发控制)                    │   │
│  └──────────────────────────────────────────────────┘   │
│                         ↓                               │
│  ┌──────────────────────────────────────────────────┐   │
│  │         DAO (数据访问层)                           │   │
│  │    ┌──────────────┐      ┌──────────────┐        │   │
│  │    │  MySQL DAO   │      │  Redis DAO   │        │   │
│  │    └──────────────┘      └──────────────┘        │   │
│  └──────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────┘
                         ↓
        ┌────────────────┴────────────────┐
        ↓                                 ↓
┌──────────────┐                  ┌──────────────┐
│    MySQL     │                  │    Redis     │
│  (Master)    │                  │   (Cache)    │
│      ↓       │                  └──────────────┘
│   Slave-1    │
│   Slave-2    │
└──────────────┘
```

### Layered responsibilities

**Router layer**

- Route registration and grouping

- Middleware mounting
- API version management

**Controller layer**

- Parameter binding and validation

- Unified response format
- Error code management

**Logic layer**

- Core business logic

- Data assembly
- Concurrency control
- Transaction management

**DAO layer**

- MySQL CRUD Operations

- Redis Cache Operations
- Batch Query Optimization
- Read-Write Splitting

---

##  Performance optimization

### 1. N+1 Query optimization

**Issue**: The post list API has an N+1 query problem, where getting 10 posts requires 21 database queries.

**Solution**: Bulk IN query optimization

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

- **Result**:

  Number of queries: 21 → 3

  Performance improvement: 46.8%

  Response time: 200ms → 106ms

### 2. Concurrent query optimization

**Problem**: Post details require querying three tables: posts, users, and communities. Serial queries take a long time.

**Solution**: Goroutine concurrent query

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

**Accomplish**：

- Response time: 150ms → 100ms (33% reduction)

- Improved concurrency

### 3. Redis Multi-level cache

**Problem**: High-frequency queries put pressure on the database, and cache failure of hot data leads to cache breakdown.

**Solution**: Multi-level caching + distributed locks

```go
// Query Process
1. Query the Redis cache
2. Hit → Return directly
3. Miss → Acquire the distributed lock
4. Double-check the cache
5. Query the database
6. Asynchronously update the cache
```

- **Effect**:

  Cache hit rate: 60-80%

  Cache hit response time: < 10ms

  Prevents cache overflow

### 4. MySQL Read-write separation

**Architecture**: 1 master and 2 slaves, round-robin load balancing

**Accomplish**：

- Write operations → Master

- Read operations → Slave (round-robin)
- Master-slave synchronization → Binlog replication

**Accomplish**：

- Read QPS increased by 50%

- Master database pressure reduced by 50%

### 5. Current limiting protection

**Algorithm**: Token Bucket Algorithm

**Accomplish**：

```go
// 投票接口：2 秒/次
RateLimitMiddleware(2*time.Second, 1)

// 发帖接口：10 秒/次
RateLimitMiddleware(10*time.Second, 1)
```

**Accomplish**：

- Preventing malicious vote manipulation/posting

- Protecting system stability

### Performance indicator summary

| Optimization | Before optimization | After optimization | Promote |
|--------|--------|--------|------|
| Post list query count | 21 次 | 3 次 | 85.7% |
| Post List Response Time | 200ms | 106ms | 46.8% |
| Post details response time | 150ms | 100ms | 33.3% |
| Cache hit response time | - | < 10ms | - |
| Read QPS | 1000 | 1500 | 50% |
| Cache hit rate | 0% | 60-80% | - |

---

##  Project Structure

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

## 📊 Database design

### Core table structure

**User Table**

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

**Community Table**

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

**Posts Table**

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

## 🔧 Configuration Instructions

### config.yaml Configuration File

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

##  API document

### Swagger document

The project integrates Swagger to automatically generate API documentation, which can be accessed after startup：

```
http://localhost:8081/swagger/index.html
```

### Core API List

#### User related

| Method | Path | Illustrate | Is certification required |
|------|------|------|-------------|
| POST | /api/v1/signup | User Registration | ❌ |
| POST | /api/v1/login | User login | ❌ |

#### Post related

| Method | Path | Illustrate | Is certification required |
|------|------|------|-------------|
| GET | /api/v1/posts | Post list (original version) | ❌ |
| GET | /api/v1/posts/optimized | Post list(N+1 optimization) | ❌ |
| GET | /api/v1/posts/cached | Post list(Cached version) | ❌ |
| GET | /api/v1/post/:id | Post Details | ❌ |
| GET | /api/v1/post/:id/concurrent | Post Details(Concurrency optimization) | ❌ |
| POST | /api/v1/post | Create a post | ✅ |

#### Community Related

| Method | Path | Illustrate | Is certification required |
|------|------|------|-------------|
| GET | /api/v1/community | Community List | ❌ |
| GET | /api/v1/community/:id | Community Details | ❌ |

#### Voting related

| Method | Path | Illustrate | Is certification required |
|------|------|------|-------------|
| POST | /api/v1/vote | Post Voting | ✅ |

#### Monitoring related

| Method | Path | Illustrate | Is certification required |
|------|------|------|-------------|
| GET | /api/v1/db/stats | Database connection pool statistics | ❌ |
| GET | /api/v1/db/health | Health Check | ❌ |
| GET | /api/v1/cache/stats | Cache Statistics | ❌ |

---

##  TEST

### Performance Testing

The project includes a complete performance test report:

- `PERFORMANCE_TEST_REPORT.md` - Comprehensive performance test report
- `CONCURRENT_OPTIMIZATION.md` - Concurrency Optimization Test Report
- `DATABASE_POOL_OPTIMIZATION_REPORT.md` - Database connection pool optimization report
- `REDIS_CACHE_OPTIMIZATION_REPORT.md` - Redis Cache Optimization Report
- `FINAL_OPTIMIZATION_REPORT.md` - Final optimization summary report

### Using wrk for stress testing

```bash
# 测试帖子列表接口
wrk -t12 -c400 -d30s http://localhost:8081/api/v1/posts/optimized

# 测试帖子详情接口
wrk -t12 -c400 -d30s http://localhost:8081/api/v1/post/1/cached
```

---

## 📝 Development Documentation

The project provides complete development documentation:

- **架构笔记-思维导图版.md** - Project Architecture Mind Map
- **架构笔记-速查版.md** - Quick Reference Manual
- **架构笔记-面试复习版.md** - - Interview preparation materials
- **面试题库-Bluebell项目.md** - Common Interview Questions and Answers

---

##  Docker deploy

### Using Docker Compose

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

### Build the image separately

```bash
# 构建镜像
docker build -t bluebell:latest .

# 运行容器
docker run -d -p 8081:8081 bluebell:latest
```

---

## Monitoring and debugging

### Performance Analysis (pprof)

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

### Log View

```bash
# 实时查看日志
tail -f web-app.log

# 查看最近 100 行
tail -n 100 web-app.log
```

### Database monitoring

Access the database connection pool statistics interface:

```bash
curl http://localhost:8081/api/v1/db/stats
```

Return example:
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

## 🤝Contribution Guidelines

Contributions, questions, and suggestions are welcome!

### Development Process

1. Fork this repository
2. Creating a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to branch (`git push origin feature/AmazingFeature`)
5. Pull Request

### Coding standards

- Follow the Go official coding standards
- Formatting code with gofmt
- Add necessary comments and documentation
- Writing unit tests

---

## 📜 License

This project is licensed under the Apache 2.0 license - see [LICENSE](LICENSE) file.

---


## 🙏 Acknowledgements

- [Gin](https://github.com/gin-gonic/gin) - Excellent Go Web Framework
- [go-redis](https://github.com/go-redis/redis) - Redis Client
- [sqlx](https://github.com/jmoiron/sqlx) - Database Tools
- [Viper](https://github.com/spf13/viper) - Configuration Management
- [Zap](https://github.com/uber-go/zap) - High-performance logging library

---



---

<div align="center">
**⭐ If this project is helpful to you, please give it a Star!⭐**

</div>

