package llm

import (
	"fmt"
	"strings"
)

// PromptBuilder 用于构建发送给 LLM 的 prompt
type PromptBuilder struct {
	templates map[string]string
}

// NewPromptBuilder 创建 prompt 构建器
func NewPromptBuilder() *PromptBuilder {
	pb := &PromptBuilder{
		templates: make(map[string]string),
	}
	pb.loadDefaults()
	return pb
}

func (pb *PromptBuilder) loadDefaults() {
	pb.templates["narrative"] = `你是一个地下城探险游戏的叙事者。请为以下场景生成一段简短但生动的描述（2-3句话）。
风格：暗黑奇幻，略带幽默，像在写日志条目。

场景类型：{{.RoomType}}
地牢深度：第{{.Depth}}层
房间序号：#{{.RoomIndex}}
上下文：{{.Context}}
玩家状态：{{.PlayerInfo}}

要求：
- 用中文描述
- 简洁有力，不超过100字
- 可以包含隐藏线索或氛围暗示
- 语气像在记录系统日志`

	pb.templates["monster"] = `你是一个地下城探险游戏的怪物设计师。请生成一个有趣的怪物。

地牢深度：第{{.Depth}}层（越深越强）
场景类型：{{.RoomType}}
是否 Boss：{{.IsBoss}}
玩家状态：{{.PlayerInfo}}

请用 JSON 格式返回：
{
  "name": "怪物名称（中文，2-4个字）",
  "description": "简短描述（1句话）",
  "traits": ["特征1", "特征2"],
  "dialogue": "怪物的战吼或开场白（1句话）",
  "strategy": "怪物的战斗策略（简短描述）"
}`

	pb.templates["event"] = `你是一个地下城探险游戏的事件设计师。请生成一个有趣的随机事件。

地牢深度：第{{.Depth}}层
玩家状态：{{.PlayerInfo}}
之前的事件：{{.Previous}}

请用 JSON 格式返回：
{
  "title": "事件标题",
  "description": "事件描述（2-3句话）",
  "choices": [
    {"text": "选项1", "description": "结果描述", "outcome": "good/bad/neutral"},
    {"text": "选项2", "description": "结果描述", "outcome": "good/bad/neutral"}
  ]
}`

	pb.templates["combat_action"] = `你是一个地下城战斗的旁白者。请描述怪物的下一步行动。
简短有力，1-2句话，像战斗日志。

怪物：{{.MonsterName}}（{{.Traits}}）
策略：{{.Strategy}}
怪物 HP：{{.MonsterHP}}/{{.MonsterMaxHP}}
玩家 HP：{{.PlayerHP}}/{{.PlayerMaxHP}}
回合数：{{.Round}}
玩家上一行动：{{.LastAction}}

直接输出行动描述，不要加引号或前缀。`

	pb.templates["dialogue"] = `你是一个地下城探险游戏的 NPC 扮演者。
用简短的对话回应玩家（1-2句话），符合角色性格。

NPC 类型：{{.NPCType}}
对话主题：{{.Topic}}
上下文：{{.Context}}
玩家状态：{{.PlayerInfo}}

直接输出 NPC 的对话内容。`
}

// Build 构建 prompt，替换模板变量
func (pb *PromptBuilder) Build(templateName string, vars map[string]string) (string, error) {
	tmpl, ok := pb.templates[templateName]
	if !ok {
		return "", fmt.Errorf("template %q not found", templateName)
	}

	result := tmpl
	for key, val := range vars {
		result = strings.ReplaceAll(result, "{{."+key+"}}", val)
	}
	return result, nil
}

// SetTemplate 设置自定义模板
func (pb *PromptBuilder) SetTemplate(name, template string) {
	pb.templates[name] = template
}
