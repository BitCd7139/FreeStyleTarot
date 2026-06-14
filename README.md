# FreeStyleTarot

## ⭐️ 项目概述

**FreeStyleTarot** 是一个基于塔罗牌的 AI 占卜系统：前端支持自由抽牌、自定义牌阵与牌阵自动识别；后端结合嵌入的塔罗知识库与大模型 API，通过多 Agent 流水线对用户问题与牌阵进行分层解读，并以流式 Markdown 输出个性化建议。

公开仓库中 **`storage/` 仅保留说明文档**，提示词与牌义库需在本地补全后再编译后端。AI 解析的设计细节见 [`storage/README.md`](storage/README.md)。

## ✨ 功能和亮点

- 🖥️ **多端适配**：PC 与移动端均可使用，支持 Capacitor 打包 Android（`web/android/`）。
- 🔐 **用户账号**：邮箱密码登录、验证码登录/注册；注册与找回密码均需邮箱验证码（30 分钟有效）；JWT 会话；个人资料页查看与修改昵称、会员状态与占卜配额。
- 🔮 **自由牌阵**：随意摆牌组成阵面，前端可自动识别常见牌阵，也可手动填写阵位含义。
- 📚 **知识库驱动**：按抽到的牌名与正逆位精确检索牌义，控制上下文长度（非向量 RAG）。
- 🤖 **多 Agent 解读**：视觉笔记 → 牌阵分析 → 人设建议，职责拆分、流式终稿。
- 👁️ **阵面坐标**：前端上报每张牌的坐标与阵位，供视觉分析环节理解牌阵空间关系（设计见 `storage/README.md`）。
- 💘 **多角色人设**：默认、女友、猫娘、雌小鬼等风格，主要影响建议阶段的语气。
- 🎯 **交互能力**：卡牌缩放与清除、阵位含义自定义、结果复制与存图、流式 Markdown 展示、系统公告。

## 👀 效果演示

|   ![](assets/g1.gif)    | ![](assets/g2.gif) | ![](assets/g3.gif) |
|:------------------:|:------------------:|:------------------:|

## 📁 项目结构

前后端分离：**Vue 3 + Vite**（`web/`）+ **Golang + Gin**（根目录 Go 模块）。

```
FreeStyleTarot/
├── main.go                 # 入口：初始化配置、数据库、路由
├── api/                    # HTTP 处理器与中间件
│   ├── auth_handler.go     # 注册、登录、验证码、个人资料
│   ├── announcement_handler.go
│   ├── HandlePredictStream.go
│   ├── HandlePrompt.go
│   └── middleware/         # CORS、JWT 鉴权、占卜配额
├── config/                 # config.yaml + 环境变量加载
├── db/                     # PostgreSQL 连接、Redis、SQL migration
│   └── migrations/
├── model/                  # 请求/响应、用户、提示词相关结构体
│   ├── request/
│   ├── response/
│   ├── user/
│   └── prompt/
├── repository/             # 数据访问层
│   ├── user/
│   └── verify/             # 验证码 Redis 存储
├── service/                # 业务逻辑
│   ├── auth/               # 用户认证、密码、会话与配额
│   ├── email/              # Resend 邮件发送
│   └── ...                 # 知识库检索、多 Agent 编排、大模型调用
├── storage/                # 嵌入用的提示词与知识库（公开版仅 README）
├── web/                    # 前端
│   ├── src/
│   │   ├── components/     # TarotMain、AuthDrawer、ProfileDrawer 等
│   │   ├── composables/    # useAuth、useAnnouncement
│   │   ├── spread/         # 牌阵识别模板（凯尔特十字、圣三角等）
│   │   └── utils/          # authApi、predictStream、cardGrid 等
│   └── android/            # Capacitor Android 工程
├── assets/                 # README 演示 GIF
├── .do/                    # DigitalOcean App Platform 部署配置
├── Dockerfile
└── .env-example            # 环境变量模板
```

| 目录 | 说明 |
|------|------|
| `web/` | 抽牌、牌阵编辑与识别、登录/个人资料、提问；将牌名、正逆位、阵位含义、坐标等提交后端 |
| `api/` | HTTP 入口：认证（`/auth/*`）、公告（`/announcement`）、占卜（`/predict`、`/prompt`） |
| `service/` | 用户认证与配额、邮件、知识库检索、牌阵素材拼装、多 Agent 编排、大模型调用 |
| `repository/` | 用户表与验证码的持久化访问 |
| `db/` | PostgreSQL（pgx）、Redis、启动时 migration |
| `config/` | 服务运行模式、鉴权开关、公告与各 Agent 的 token、温度等 |
| `model/` | 请求/响应、用户、提示词相关结构体 |
| `storage/` | 嵌入用的提示词与知识库（公开版仅 README，见该目录说明） |

入口：`main.go` 启动 Gin 并注册路由。环境变量见 [`.env-example`](.env-example)（大模型 API、数据库、JWT、Redis、邮件等）。`config/config.yaml` 可配置是否强制登录（`auth.force_login`）、是否跳过邮件验证（`auth.skip_verify`）及公告内容。

## 🔌 API 概览

| 路由 | 鉴权 | 说明 |
|------|------|------|
| `GET /announcement` | 无 | 系统公告 |
| `GET /auth/config` | 无 | 前端鉴权配置（如是否强制登录） |
| `POST /auth/send-code` | 无 | 发送邮箱验证码 |
| `POST /auth/verify` | 无 | 校验邮箱（注册前置） |
| `POST /auth/verify-code` | 无 | 验证码登录 |
| `POST /auth/complete-code-signup` | 无 | 验证码注册补全 |
| `POST /auth/login` | 无 | 密码登录 |
| `POST /auth/register` | 无 | 密码注册（需邮箱验证码） |
| `POST /auth/reset-password` | 无 | 邮箱验证码重置密码 |
| `GET /auth/me` | 需登录 | 获取当前用户 |
| `PATCH /auth/me` | 需登录 | 更新个人资料（如昵称） |
| `POST /predict` | 需登录 + 配额 | 流式占卜 |
| `POST /prompt` | 需登录 + 配额 | 调试提示词拼装 |

## 📚 AI 解析（概览）

一次占卜大致经过三步 Agent（详见 [`storage/README.md`](storage/README.md)）：

```text
视觉分析 (TODO) — 根据牌阵与坐标提取阵面元素与视觉流动关系，对照用户问题
    ↓
牌阵解读（spread_analyst）— 结合知识库做能量与逐牌分析，不给行动清单
    ↓
人设建议（persona_advisor）— 按角色卡输出建议（用户看到的流式正文）
```

- **素材**：结构化短库 + 长文牌义按牌名命中；解读 Agent 用全文素材，建议 Agent 用截断后的 digest + 上一步分析。
- **提示词**：共享 System（规则与各 Agent 职责）+ 角色卡；同次占卜多次调用共用 System，以 `[ACTIVE AGENT: …]` 切换职责，便于前缀缓存。
- **输出**：流式结果以 JSON 包装逐段推送，避免 SSE 与 Markdown 冲突。

本地完整运行需自行准备 `storage` 嵌入资源、PostgreSQL、Redis（验证码场景）及 `.env` 配置后编译；仅克隆公开仓库时请先阅读 `storage/README.md`。

## 🤗 想为项目做贡献? 欢迎新建 issue 讨论你的想法！

包括不限于 **Prompt 调整、牌阵建议、知识库扩展、交互优化、美术素材、Bug 修复** 等内容。

## 📜 免责声明

本项目仅供娱乐和学习使用，所有占卜结果仅供参考，不应被视为专业建议。开发者不对任何因使用本项目而产生的直接或间接损失负责。
