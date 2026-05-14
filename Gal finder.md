# Gal finder

## 项目定位

`Gal finder` 是一个 Galgame 个人检索与收藏系统。项目基于 `gin-vue-admin` 的工程结构进行轻量化改造，重点放在 Gin 后端开发、请求链路设计、中间件实践、第三方 API 代理和前后端接口约定上。

真正的训练重点是理解 Gin 的请求处理流程，以及中间件如何在真实接口中承担鉴权、日志、跨域、异常恢复、限流、超时控制等职责。

参考项目与API资源
- `gin-vue-admin`：https://github.com/flipped-aurora/gin-vue-admin
- VNDB Kana API：https://api.vndb.org/kana

## 后端核心目标

后端负责统一承接前端请求，并代理 VNDB 第三方 API，并返回调用结果
核心职责：

- 用户注册、登录、登录态校验。
- 通过 Gin 后端代理 VNDB 搜索接口。
- 对 VNDB 返回数据做字段裁剪和响应格式转换。
- 维护用户收藏数据。
- 通过中间件统一处理鉴权、日志、跨域、异常恢复、限流和请求超时。

## 项目亮点
可以把 VNDB 搜索接口当作主线接口：先跑通请求，再逐步接入 JWT、Logger、Recovery、Rate Limit、Timeout。每加一个中间件，都观察它在请求进入、拦截、放行、失败返回这些环节中的作用。

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
- Redis（开发中）

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

初步接口职责划分:

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

## 目前拥有的功能
1.用户注册与登录

2.有初步的搜索功能

3.有收藏功能，可以管理。


### to dolist

- 当前状态：阶段 1 和阶段 2 的 Gin 后端骨架已经基本完成，包括路由分组、统一响应、RequestID、JWTAuth、RateLimit、自定义 Logger、VNDB 搜索接口，以及 MySQL/GORM 初始化。
- 当前暂停点：先不要继续堆新功能，需要先把项目结构重新理清，避免只是在复制代码。
- 阶段 3.1：修正当前 `auth.go` 的登录逻辑，重点是 `UserID := user.ID`、bcrypt 使用 `user.PasswordHash` 校验前端传入的明文密码，并保证 `go build ./...` 通过。
- 阶段 3.2：完成注册接口接入 MySQL，注册时写入 `users` 表，密码只保存 bcrypt hash，不返回密码给前端。
- 阶段 3.3：完成登录接口接入 MySQL，登录时查询用户、校验密码、生成 JWT，并返回 token 与用户基础信息。
- 阶段 3.4：用 Postman/Apifox 验证 `register`、`login`、`profile` 三个接口，确认 JWT 中间件能从 token 中解析出 `userID`。
- 阶段 3.5：把收藏接口接入 MySQL 的 `favourites` 表，实现新增收藏、收藏列表、删除收藏，并通过 `userID + vndbID` 避免重复收藏。
- 阶段 4：梳理后端分层，把现在堆在 `api` 里的业务逻辑逐步拆到 `service`，让 `api` 只负责 HTTP 参数、Context、响应。
- 阶段 5：整理模型和 DTO，把数据库模型放在 `model`，前端请求结构放在 `request`，返回结构放在 `response`，避免数据库字段和接口字段混在一起。
- 阶段 6：搭建前端项目，基于 Vue 3 / Vite / Element Plus 或 gin-vue-admin 现有结构，先完成基础布局、路由、请求封装和 token 保存逻辑。
- 前端协作方式：由于前端基础较少，前端页面、组件、路由和请求封装主要由 AI提供代码；需要重点学习和提醒的是前后端对接点，包括接口 URL、请求方法、请求体字段、响应结构、token 保存、请求头携带、401/429/500 错误处理。
- 阶段 6.1：完成登录页和注册页，前端提交账号密码到后端，登录成功后保存 token，并跳转到主页面。
- 阶段 6.2：完成搜索页，前端输入关键词和分页参数，请求后端 VNDB 搜索接口，并展示标题、封面、发售时间、评分等信息。
- 阶段 6.3：完成收藏交互，搜索结果中点击收藏，前端调用收藏接口，把作品写入当前用户收藏夹。
- 阶段 6.4：完成收藏列表页，前端请求当前用户收藏列表，并支持取消收藏。
- 阶段 6.5：完成前端登录态处理，请求自动携带 `Authorization: Bearer token`，遇到 401 时清理 token 并跳转登录页。
- 阶段 7：前后端联调，确认登录、注册、搜索、收藏、收藏列表、取消收藏等流程完整可用。
- 阶段 8：增强中间件，把 Logger、RequestID、RateLimit、JWTAuth 的行为在真实接口中验证清楚，必要时补充统一错误码。
- 阶段 9：引入 Redis 作为增强项，不作为收藏主存储；后续主要用于搜索结果缓存、限流计数、热门关键词缓存。
- 阶段 10：补充接口测试和项目收尾，确认注册、登录、搜索、收藏、删除收藏、限流、鉴权失败等场景都能正常返回。
