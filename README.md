# ⚔ DungeonLog — 终端地牢探险 RPG

> 一款伪装成日志输出的终端地牢 RPG，适合上班摸鱼时游玩。

## 🎮 特性

- **伪装模式**：所有输出带 `[时间戳] INFO/WARN/ERROR` 前缀，看起来像后台服务日志
- **LLM 驱动**：怪物行为、场景描述、随机事件由大模型动态生成
- **战斗策略**：回合制战斗 + 技能 + 装备词缀 + Buff/Debuff 系统
- **4 种职业**：战士、游侠、法师、盗贼，各有独特技能
- **程序化地牢**：每层 8-12 个房间，类型随机生成
- **装备词缀系统**：普通/优秀/稀有/史诗/传说 5 种品质

## 🚀 快速开始

```bash
# 编译
go build -o dungeonlog ./cmd/dungeonlog

# 运行
./dungeonlog
```

## 📖 基本命令

| 命令 | 说明 |
|------|------|
| `start <名字> [职业]` | 开始新游戏 |
| `info` / `status` | 查看角色信息 |
| `bag` / `inventory` | 查看物品栏 |
| `skills` | 查看技能列表 |
| `map` | 查看地牢地图 |
| `go <编号>` | 前往指定房间 |
| `look` / `interact` | 与当前房间交互 |
| `equip <编号>` | 装备物品 |
| `save` | 保存游戏 |
| `load` | 加载游戏 |
| `stealth` | 切换伪装模式 |
| `config` | 查看 LLM 配置 |
| `help` | 查看帮助 |
| `quit` | 退出游戏 |

## ⚔ 战斗命令

| 命令 | 说明 |
|------|------|
| `1` / `attack` | 普通攻击 |
| `2` / `skill <编号>` | 使用技能 |
| `3` / `use <编号>` | 使用物品 |
| `4` / `flee` | 尝试逃跑 |

## 🎯 职业介绍

| 职业 | 特点 | 技能 |
|------|------|------|
| 🛡 战士 | 高防高血，擅长持久战 | 重击、盾墙、战吼、旋风斩 |
| 🏹 游侠 | 高暴击高闪避 | 精准射击、毒箭、闪避、致命一击 |
| 🔮 法师 | AOE 伤害 + 控制 | 火球术、冰冻术、雷电链、法力护盾 |
| 🗡 盗贼 | 高暴击 + 连击 | 背刺、毒刃、影遁、连击 |

## 🏗 项目结构

```
dungeonlog/
├── cmd/dungeonlog/      # 主程序入口
├── internal/
│   ├── engine/          # 游戏引擎（状态机、命令处理）
│   ├── display/         # 终端渲染（日志伪装格式）
│   ├── game/            # 游戏核心（角色、战斗、地牢、物品）
│   └── llm/             # LLM 接口（可切换不同提供商）
└── data/
    └── saves/           # 存档目录
```

## 🔌 LLM 配置

编辑 `data/config.yaml` 启用 LLM：

```yaml
llm:
  enabled: true                    # 启用 LLM
  base_url: "https://api.openai.com/v1"  # API 地址
  api_key: "sk-your-api-key"       # API Key
  model: "gpt-4o-mini"             # 模型名称
  temperature: 0.8                 # 温度
  max_tokens: 512                  # 最大 token
```

### 支持的 API 服务商

| 服务商 | base_url |
|--------|----------|
| OpenAI | `https://api.openai.com/v1` |
| DeepSeek | `https://api.deepseek.com/v1` |
| Moonshot | `https://api.moonshot.cn/v1` |
| 本地 Ollama | `http://localhost:11434/v1` |
| 兼容 OpenAI 格式的任意服务 | 自定义 URL |

也支持环境变量覆盖：
```bash
export DUNGEONLOG_API_KEY="sk-your-key"
export DUNGEONLOG_BASE_URL="https://api.deepseek.com/v1"
```

### 自定义 LLM Provider

实现 `internal/llm/provider.go` 中的 `LLMProvider` 接口即可接入其他 LLM：
- `GenerateNarrative()` — 生成场景描述
- `GenerateMonster()` — 生成怪物
- `GenerateEvent()` — 生成随机事件
- `GenerateDialogue()` — 生成 NPC 对话
- `GenerateCombatAction()` — 生成战斗行为描述

## 🎲 房间类型

| 图标 | 类型 | 说明 |
|------|------|------|
| 🚪 | 入口 | 地牢入口 |
| ⚔️ | 战斗 | 遭遇怪物 |
| 💎 | 宝箱 | 获得战利品 |
| 🏪 | 商店 | 购买装备物品 |
| 🏕 | 休息 | 恢复生命值 |
| ❓ | 事件 | 随机事件 |
| 💀 | Boss | Boss 战 |

## 📝 开发

```bash
# 安装依赖
go mod tidy

# 运行测试
go test ./...

# 编译
go build -o dungeonlog ./cmd/dungeonlog
```

## 📜 License

MIT
