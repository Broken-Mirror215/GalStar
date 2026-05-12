# Gal finder

## 项目定位

`Gal finder` 是一个 Galgame 个人检索与收藏系统。项目基于 `gin-vue-admin` 的工程结构进行轻量化改造，重点放在 Gin 后端开发、请求链路设计、中间件实践、第三方 API 代理和前后端接口约定上。

真正的训练重点是理解 Gin 的请求处理流程，以及中间件如何在真实接口中承担鉴权、日志、跨域、异常恢复、限流、超时控制等职责。

参考项目：

- `gin-vue-admin`：https://github.com/flipped-aurora/gin-vue-admin
- VNDB Kana API：https://api.vndb.org/kana

## 后端核心目标

后端负责统一承接前端请求，并代理 VNDB 第三方 API，避免前端直接依赖外部接口。

核心职责：

- 用户注册、登录、登录态校验。
- 通过 Gin 后端代理 VNDB 搜索接口。
- 对 VNDB 返回数据做字段裁剪和响应格式转换。
- 维护用户收藏数据。
- 通过中间件统一处理鉴权、日志、跨域、异常恢复、限流和请求超时。

第一版不追求复杂权限系统和微服务拆分，重点是把一个单体后端服务做完整、清晰、可维护。

## 项目训练重点

这个项目的主线不是“做几个增删改查接口”，而是围绕一次 HTTP 请求在 Gin 中的完整生命周期来设计后端。

实现时可以把 VNDB 搜索接口当作主线接口：先跑通请求，再逐步接入 JWT、Logger、Recovery、Rate Limit、Timeout。每加一个中间件，都观察它在请求进入、拦截、放行、失败返回这些环节中的作用。

一次搜索请求的目标链路：

```text
前端发起搜索
-> CORS 处理跨域
-> Request ID 标记请求
-> Logger 记录请求入口
-> Recovery 兜底 panic
-> JWT Auth 校验登录态
-> Rate Limit 限制搜索频率
-> api 层绑定和校验参数
-> service 层调用 VNDB
-> HTTP Client timeout 控制外部请求
-> response DTO 转换
-> Logger 记录状态码和耗时
-> 前端按统一响应格式渲染结果
```

CRUD 在这里的定位：

- 用户表用于支撑 JWT 登录和用户身份识别。
- 收藏表用于支撑登录用户的数据归属和接口权限判断。
- 搜索历史可以作为后续扩展，用来练习异步记录、统计和缓存。
- 普通增删改查不作为项目亮点，只作为中间件、接口设计和后端分层的业务载体。

## 技术选型

后端：

- Go
- Gin
- GORM
- MySQL
- JWT
- bcrypt
- Redis，第二阶段用于限流和缓存

前端：

- Vue 3
- Vite
- Element Plus 或 gin-vue-admin 原有组件体系

## 后端模块设计

项目结构参考 `gin-vue-admin`，但第一版保持轻量。

```text
server
├── main.go
├── config
├── global
├── initialize
├── router
├── api
│   └── v1
├── service
├── model
│   ├── request
│   └── response
├── middleware
└── utils
```

推荐职责划分：

```text
router       路由分组，区分公开接口和登录后接口
api/v1       Controller 层，负责参数绑定、校验、调用 service、返回响应
service      业务层，处理注册登录、VNDB 搜索、收藏逻辑
model        GORM 模型、请求 DTO、响应 DTO
middleware   JWT、CORS、Logger、Recovery、Rate Limit、Timeout
utils        JWT 生成解析、密码加密、HTTP Client、统一响应工具
```

后端分层约定：

- `api` 不直接写数据库逻辑，只处理 HTTP 相关内容。
- `service` 不依赖 Gin Context，方便后续测试和复用。
- `model/request` 保存前端入参结构。
- `model/response` 保存返回给前端的数据结构。
- 第三方 VNDB 返回结构不直接透传，统一转换成项目自己的 response。

## Gin 实践重点

这个项目要重点体会 Gin 在后端服务中的几个核心能力。

### 路由分组

通过路由分组区分公开接口和需要登录的接口。

```text
/api/auth          公开接口，注册和登录
/api/user          用户相关接口，需要 JWT
/api/v1/vndb       VNDB 搜索接口，需要 JWT 和限流
/api/v1/favorites  收藏接口，需要 JWT
```

这样可以把中间件挂在合适的路由组上，而不是每个接口里手动判断登录态。

### Context 传递

JWT 中间件解析 token 后，把 `userId` 写入 Gin Context。

后续接口只从 Context 中读取当前用户，不从前端请求体中接收 `userId`。

```text
JWT Auth middleware -> c.Set("userId", userId)
api handler -> c.Get("userId")
```

这样可以避免前端伪造用户 ID，也能更清楚地体会中间件和业务 handler 之间的协作。

### 参数绑定和校验

Gin 的 `ShouldBindJSON`、`ShouldBindQuery` 用来处理前端入参。

重点不是会不会接参数，而是形成统一约定：

- 请求参数绑定失败，统一返回参数错误。
- 分页参数有默认值和最大值限制。
- 搜索关键词不能为空。
- 收藏接口不允许前端传 `userId`。

### 统一响应和错误处理

所有 handler 都返回统一结构，避免前端针对不同接口写不同解析逻辑。

业务错误、参数错误、鉴权错误、第三方 API 错误都应该转换成稳定的响应格式。

## 前后端接口设计

接口设计的重点不是把数据库表直接暴露成 CRUD，而是明确前端能做什么、后端负责兜住什么。

### 接口职责划分

第一版可以按四类接口来设计：

- 用户主页：如果只做静态欢迎页，可以暂时不接后端；如果要展示当前用户昵称、收藏数量、收藏列表或搜索历史，就需要调用 `GET /api/user/profile` 等接口。
- 登录与注册：前端提交账号密码，后端负责参数校验、密码加密、数据库读写和 token 生成；登录成功后的页面跳转由前端完成。
- 查找接口：前端只传关键词和分页参数，后端负责调用 VNDB API、处理超时和错误，并把结果整理成稳定的响应结构返回给前端。
- 收藏接口：前端点击收藏后，把作品信息或 `vndbId` 发给后端；后端从 JWT 中读取当前用户 ID，再把作品写入当前用户的收藏夹。

这里要注意：前端不应该传 `userId` 来决定操作哪个用户的数据，用户身份统一由 JWT 中间件解析后写入 Gin Context。

### 通用响应格式

后端统一返回固定结构，前端只需要按 `code` 判断业务结果。

```json
{
  "code": 0,
  "message": "success",
  "data": {}
}
```

错误示例：

```json
{
  "code": 401,
  "message": "unauthorized",
  "data": null
}
```

约定：

- HTTP 状态码用于表示请求级别结果。
- `code` 用于表示业务级别结果。
- `message` 可直接用于前端提示。
- `data` 永远存在，无数据时返回 `null` 或空对象。

### 认证接口

#### 注册

```text
POST /api/auth/register
```

请求：

```json
{
  "username": "test",
  "password": "123456",
  "nickname": "tester"
}
```

响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "userId": 1,
    "username": "test",
    "nickname": "tester"
  }
}
```

#### 登录

```text
POST /api/auth/login
```

请求：

```json
{
  "username": "test",
  "password": "123456"
}
```

响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "token": "jwt_token",
    "expiresAt": 1710000000,
    "user": {
      "userId": 1,
      "username": "test",
      "nickname": "tester"
    }
  }
}
```

前端约定：

- 登录成功后保存 `token`。
- 后续请求在 Header 中携带：

```text
Authorization: Bearer jwt_token
```

### 用户接口

```text
GET /api/user/profile
```

说明：

- 需要登录。
- 后端从 JWT 中解析用户 ID。
- 前端不需要传 `userId`，避免越权风险。

响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "userId": 1,
    "username": "test",
    "nickname": "tester"
  }
}
```

### VNDB 搜索接口

```text
GET /api/v1/vndb/search?q=关键词&page=1&pageSize=20
```

说明：

- 需要登录。
- 前端只关心搜索关键词和分页。
- 后端负责拼接 VNDB 请求体、调用外部 API、处理超时和错误。
- 搜索接口后续可以单独加限流和缓存。

响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "id": "v97",
        "title": "Saya no Uta",
        "altTitle": "沙耶之歌",
        "released": "2003-12-26",
        "rating": 80,
        "imageUrl": "https://...",
        "thumbnailUrl": "https://...",
        "sexual": 1.3,
        "violence": 1.8
      }
    ],
    "page": 1,
    "pageSize": 20,
    "hasMore": true
  }
}
```

后端处理要点：

- VNDB 的 `image` 可能为 `null`，后端返回空字符串或 `null`，前端展示默认图。
- VNDB 的评分字段需要明确单位，避免前端误解。
- 外部 API 调用必须设置 timeout。
- VNDB 错误不直接暴露给前端，统一转换成项目错误码和提示。

### 作品详情接口

```text
GET /api/v1/vndb/:id
```

说明：

- 需要登录。
- 第一版可以先不做，搜索列表足够支撑收藏功能。
- 如果实现，后端仍然代理 VNDB API，不让前端直接请求 VNDB。

### 收藏接口

#### 添加收藏

```text
POST /api/v1/favorites
```

请求：

```json
{
  "vndbId": "v97",
  "title": "Saya no Uta",
  "imageUrl": "https://...",
  "thumbnailUrl": "https://...",
  "rating": 80,
  "released": "2003-12-26"
}
```

响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "favoriteId": 10
  }
}
```

后端处理要点：

- 用户 ID 从 JWT 中获取，不允许前端传。
- `user_id + vndb_id` 建唯一索引。
- 重复收藏返回明确业务错误，例如 `already favorited`。

#### 收藏列表

```text
GET /api/v1/favorites?page=1&pageSize=20
```

响应：

```json
{
  "code": 0,
  "message": "success",
  "data": {
    "list": [
      {
        "favoriteId": 10,
        "vndbId": "v97",
        "title": "Saya no Uta",
        "imageUrl": "https://...",
        "thumbnailUrl": "https://...",
        "rating": 80,
        "released": "2003-12-26",
        "createdAt": "2026-05-11T12:00:00Z"
      }
    ],
    "page": 1,
    "pageSize": 20,
    "total": 1
  }
}
```

#### 取消收藏

```text
DELETE /api/v1/favorites/:vndb_id
```

说明：

- 需要登录。
- 只允许删除当前登录用户自己的收藏。

## 中间件设计重点

这个项目可以重点写中间件，因为这正好补齐目前相对缺少的经验。

中间件不要只停留在 `router.Use()`，还要关注挂载位置、执行顺序、失败时是否 `Abort` 请求，以及向 Gin Context 写入哪些信息供后续 handler 使用。

### Recovery

目标：

- 捕获 panic。
- 返回统一错误响应。
- 记录 panic 信息，避免服务直接崩溃。

位置：

```text
router.Use(middleware.Recovery())
```

### Logger

目标：

- 记录请求方法、路径、状态码、耗时、客户端 IP。
- 后续结合 Request ID 定位一次完整请求。

日志字段建议：

```text
method
path
status
latency
client_ip
user_id
request_id
```

### CORS

目标：

- 允许本地 Vue 开发服务器访问 Gin 后端。
- 控制允许的 Method、Header 和 Origin。

注意：

- 开发环境可以允许 `localhost`。
- 生产环境不要直接放开所有 Origin。

### JWT Auth

目标：

- 从 `Authorization` Header 中读取 token。
- 校验 token 是否有效和过期。
- 从 token claims 中取出 `userId`。
- 将 `userId` 写入 Gin Context，后续接口直接读取。

后端接口约定：

```text
public group:
POST /api/auth/register
POST /api/auth/login

private group with JWT:
GET    /api/user/profile
GET    /api/v1/vndb/search
POST   /api/v1/favorites
GET    /api/v1/favorites
DELETE /api/v1/favorites/:vndb_id
```

### Rate Limit

目标：

- 主要限制 VNDB 搜索接口。
- 避免单个用户或 IP 高频调用第三方 API。

第一版可以用内存限流：

```text
key = user_id 或 client_ip
limit = 每分钟 N 次
```

第二版再切换到 Redis：

```text
INCR rate_limit:{user_id}:{minute}
EXPIRE 60s
```

### Timeout

目标：

- VNDB 外部请求必须有超时。
- 避免第三方接口慢导致后端请求堆积。

建议：

```text
VNDB search timeout: 3s
overall request timeout: 5s
```

### Request ID

目标：

- 每个请求生成唯一 ID。
- Logger、错误响应和后续排查使用同一个 ID。

响应 Header：

```text
X-Request-ID: request_id
```

## 数据库设计

数据库设计只作为后端业务落地的基础，不作为项目核心亮点。表结构保持简单，重点放在登录用户的数据归属、唯一约束和接口权限判断。

### users

```text
id
username
password_hash
nickname
created_at
updated_at
deleted_at
```

索引：

```text
unique(username)
```

### favorites

```text
id
user_id
vndb_id
title
image_url
thumbnail_url
rating
released
created_at
updated_at
deleted_at
```

索引：

```text
index(user_id)
unique(user_id, vndb_id)
```

说明：

- 收藏表冗余保存标题、封面、评分，减少收藏列表对 VNDB 的依赖。
- 如果后续需要刷新详情，可以根据 `vndb_id` 再请求 VNDB。

## VNDB API 代理设计

后端调用：

```text
POST https://api.vndb.org/kana/vn
```

请求体：

```json
{
  "filters": ["search", "=", "saya no uta"],
  "fields": "title,alttitle,released,rating,image{url,thumbnail,sexual,violence}",
  "sort": "searchrank",
  "results": 20,
  "page": 1
}
```

建议封装：

```text
service/vndb_service.go
utils/http_client.go
model/response/vndb.go
```

处理流程：

```text
前端 search 请求
-> JWT Auth 中间件校验登录
-> Rate Limit 中间件检查频率
-> api/v1 绑定 query 参数
-> service 组装 VNDB 请求
-> HTTP Client 带 timeout 调用 VNDB
-> 转换为项目 response DTO
-> 返回统一响应
```

## 实现优先级

第一阶段：

1. 搭建 Gin 路由分组，区分 public group 和 private group。
2. 实现统一响应和统一错误处理。
3. 接入 Logger、Recovery、CORS，先把基础请求链路跑通。
4. 实现注册、登录和 JWT Auth，让登录态通过中间件进入业务接口。
5. 实现 VNDB 搜索代理接口，把第三方 API 调用接入 Gin 请求链路。
6. 实现收藏增删查，把 CRUD 作为登录鉴权和用户数据归属的业务载体。

第二阶段：

1. 给 VNDB 搜索接口接入 Rate Limit。
2. 给 VNDB HTTP Client 增加 Timeout。
3. 引入 Request ID，串联日志、错误响应和排查链路。
4. 细化错误码，区分参数错误、鉴权错误、业务错误和第三方服务错误。
5. Redis 缓存搜索结果，减少重复请求 VNDB。

第三阶段：

1. 搜索历史记录。
2. 收藏列表筛选和排序。
3. 作品详情页。
4. 针对中间件和接口补充测试，例如鉴权失败、限流命中、VNDB 超时。

## 简历表达方向

不建议写成“Gin 学习项目”，可以写成：

> 基于 Gin + GORM + MySQL 实现 Galgame 检索与收藏系统，负责后端接口设计、JWT 鉴权、VNDB 第三方 API 代理、收藏管理，以及 Logger、Recovery、CORS、Rate Limit、Timeout 等中间件建设。

可突出：

- 设计统一前后端接口响应格式，封装认证、搜索、收藏等 RESTful API。
- 基于 Gin 中间件实现 JWT 鉴权，将用户身份注入请求上下文，保护搜索和收藏接口。
- 基于 Gin 路由分组组织 public/private API，在请求链路中组合 Logger、Recovery、CORS、Rate Limit、Timeout 等中间件。
- 封装 VNDB API 客户端，统一处理请求超时、错误转换和响应 DTO 映射。
- 使用 GORM 设计用户表和收藏表，通过唯一索引保证同一用户不可重复收藏，使 CRUD 服务于认证和用户数据归属。
