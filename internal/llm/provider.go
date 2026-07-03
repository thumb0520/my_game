package llm

import "context"

// LLMProvider 定义了大模型交互的抽象接口
// 不同的 LLM 提供商（Claude、OpenAI、Ollama）实现此接口即可
type LLMProvider interface {
	// GenerateNarrative 生成房间/场景描述
	GenerateNarrative(ctx context.Context, req NarrativeRequest) (string, error)

	// GenerateMonster 生成怪物的名称、描述和行为模式
	GenerateMonster(ctx context.Context, req MonsterRequest) (*MonsterResponse, error)

	// GenerateEvent 生成随机事件
	GenerateEvent(ctx context.Context, req EventRequest) (*EventResponse, error)

	// GenerateDialogue 生成 NPC 对话
	GenerateDialogue(ctx context.Context, req DialogueRequest) (string, error)

	// GenerateCombatAction 生成怪物的战斗行为描述
	GenerateCombatAction(ctx context.Context, req CombatActionRequest) (string, error)
}

// NarrativeRequest 叙事生成请求
type NarrativeRequest struct {
	RoomType   string            // 房间类型：combat, treasure, shop, boss, event, rest
	Depth      int               // 地牢深度
	RoomIndex  int               // 房间序号
	Context    string            // 上下文（前一个房间的信息）
	PlayerInfo string            // 玩家当前状态摘要
	Extra      map[string]string // 额外上下文
}

// MonsterRequest 怪物生成请求
type MonsterRequest struct {
	Depth      int    // 地牢深度（影响怪物强度）
	RoomType   string // 房间类型
	IsBoss     bool   // 是否是 Boss
	PlayerInfo string // 玩家当前状态
}

// MonsterResponse LLM 生成的怪物响应
type MonsterResponse struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Traits      []string `json:"traits"`     // 特征：凶猛、狡猾、胆小等
	Dialogue    string   `json:"dialogue"`   // 怪物的开场白/战吼
	Strategy    string   `json:"strategy"`   // 怪物的战斗策略描述
}

// EventRequest 随机事件请求
type EventRequest struct {
	Depth      int
	RoomType   string
	PlayerInfo string
	Previous   string // 之前的事件，避免重复
}

// EventResponse LLM 生成的事件响应
type EventResponse struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Choices     []EventChoice   `json:"choices"`
}

// EventChoice 事件选项
type EventChoice struct {
	Text        string `json:"text"`
	Description string `json:"description"` // 选择后的结果描述（LLM 生成）
	Outcome     string `json:"outcome"`     // good / neutral / bad
}

// DialogueRequest NPC 对话请求
type DialogueRequest struct {
	NPCType    string // 商店老板、神秘旅人、受伤冒险者等
	Context    string // 当前上下文
	PlayerInfo string
	Topic      string // 对话主题
}

// CombatActionRequest 战斗行为请求
type CombatActionRequest struct {
	MonsterName string
	Traits      []string
	Strategy    string
	MonsterHP   int
	MonsterMaxHP int
	PlayerHP    int
	PlayerMaxHP int
	Round       int
	LastAction  string // 玩家上一轮的行动
}
