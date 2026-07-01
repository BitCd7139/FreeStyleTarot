## 关于本目录

公开仓库里 **`storage/` 仅保留本说明**；提示词、牌义库等资源在本地/私有构建时放入此目录，由 Go `embed` 打进后端二进制。下文描述的是**完整部署时**的设计思路，不逐一展开各文件内容。

---

## 项目目录（简览）

| 目录 | 做什么 |
|------|--------|
| `web/` | 前端：抽牌、牌阵、提问，把牌名/正逆位/阵位含义/坐标等 POST 给后端 |
| `api/` | HTTP 入口（流式占卜、调试 prompt 等） |
| `service/` | 牌义检索、提示词拼装、多 Agent 编排、调用 DeepSeek |
| `config/` | 服务与 Agent 的 `max_tokens`、温度等 |
| `model/` | 请求/响应结构体 |
| `storage/` | 嵌入用的提示词与知识库（公开版仅 README） |

`web/android/` 等为 Capacitor 打包产物，与占卜逻辑无关，可忽略。

---

## 占卜流水线（多 Agent）

一次解读大致分三步：

```text
视觉分析 Agent — 只提取元素内容和关系、不分析
    ↓
牌阵解读 Agent — 只分析、不给建议
    ↓
人设建议 Agent — 按角色卡给决策与行动建议（流式输出）
```

### 视觉分析 Agent

后端会接收每张牌的 **坐标**（`x`, `y`）及牌阵类型。设想能力是：根据牌阵模板提取关键元素，结合位置分析**视觉流动关系**（视线、张力、聚散），再与用户问题里的关键要素做对照，输出一段结构化「阵面视觉笔记」，供后续 Agent 引用。  

### 牌阵解读 Agent 

- 从知识库按抽到的牌做**精确检索**（非向量 RAG）：结构化关键词 + 长文牌义，按正逆位选取对应段落，拼成 Markdown「牌阵素材」。
- 结合用户问题做能量格局、逐牌解读、综合叙事。
- **禁止**替用户做决定或给行动清单。

### 人设建议 Agent 

- 输入：用户问题 + 牌阵**摘要**（digest，控制长度）+ 上一 Agent 的分析全文。
- 按所选角色卡（`prompt_*.md` 对应的人设）输出带风格的建议；用户看到的流式正文主要来自这一步。

每个 Agent 各自加载独立的 **System 提示**（`{agent}_system.md`），User 消息开头用 `[ACTIVE AGENT: …]` 标识当前职责。早期版本曾让多 Agent 共用一份长 System 以命中前缀缓存，随着职责拆分变复杂、命中率下降，已改为按 Agent 拆分独立 System。

Agent 的 `max_tokens`、温度、摘要截断长度等在根目录 `config/config.yaml` 的 `agents` 段配置，启动时加载。

---

## 提示词怎么拼

完整部署时，本目录大致两类 Markdown：

1. **Agent System**（`{agent}_system.md`）  
   每个 Agent 一份独立 System：`spread_analyst_system.md`（牌阵解读）、`advisor_system.md`（综合建议）、`persona_system.md`（人设转写）、`visual_flow_system.md`（视觉流动，Freestyle 模式）、`intent_clarifier_system.md`（意图澄清）。各 Agent 只加载自己的 System，互不共用。

2. **角色卡**（`prompt_{model}.md`）  
   与前端 `model` 字段对应，如 `default` / `mate` / `neko` / `zako`。仅由 `persona` Agent 拼在 `persona_system.md` 之后；解读与建议阶段忽略人设。

调试：对 `POST /prompt` 提交与占卜相同的 JSON，可查看拼装后的 System/User，不扣模型额度。

---

## 知识库与摘要（概念）

不展开具体文件名与字段，只说明分工：

| 层次 | 作用 |
|------|------|
| 结构化短库 | 关键词、元素/占星等元数据；与前端牌 `name` 键对齐 |
| 长文牌义库 | 按牌分块的长解读（正/逆/描述）；启动时建索引，按牌名命中 |
| 全文素材 | 给 **spread_analyst**，含当前阵位与完整牌义 |
| digest 摘要 | 给 **persona_advisor**，每张牌保留关键词与截断后的要点，避免重复塞入全文 |

正逆位：只注入与抽牌朝向一致的那一侧解读。

---

## 使用与维护提示

- **想要更中性的牌理**：前端选「默认」角色；强人设（如雌小鬼）会明显影响建议语气，解读阶段虽要求中性，终稿仍受人设牵引。
- **公开克隆后本地跑通**：需自行补全 `storage` 下嵌入资源再编译后端；仅有本 README 时无法提供牌义与 prompt。
- **改 Agent 行为**：改对应 Agent 的 `{agent}_system.md` 里的职责说明；改语气改对应 `prompt_*.md`。
- **视觉分析**：待接入后，计划在解读 Agent 之前增加一环，User 消息中多一节「阵面视觉笔记」即可，与现有素材拼装并列。

相关实现入口：`service/orchestrator.go`、`service/tarot_material.go`、`service/prompt_agents.go`（细节以代码为准）。
