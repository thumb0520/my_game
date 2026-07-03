package llm

import (
	"context"
	"fmt"
	"math/rand"
)

// MockProvider 是一个不依赖真实 API 的模拟实现
// 用于开发测试和无网络环境下的游玩
type MockProvider struct {
	promptBuilder *PromptBuilder
}

// NewMockProvider 创建模拟 LLM 提供商
func NewMockProvider() *MockProvider {
	return &MockProvider{
		promptBuilder: NewPromptBuilder(),
	}
}

var mockMonsters = []MonsterResponse{
	{Name: "哥布林", Description: "一个矮小狡猾的绿皮生物", Traits: []string{"狡猾", "胆小"}, Dialogue: "嘎嘎嘎！把你的金币交出来！", Strategy: "喜欢偷袭，血量低时会逃跑"},
	{Name: "骷髅战士", Description: "被黑暗魔法驱动的不死战士", Traits: []string{"不死", "坚韧"}, Dialogue: "（骨头嘎吱作响）", Strategy: "不会逃跑，持续进攻"},
	{Name: "暗影刺客", Description: "隐匿在黑暗中的致命杀手", Traits: []string{"潜行", "暴击"}, Dialogue: "你什么都看不见……", Strategy: "高暴击，优先攻击弱点"},
	{Name: "石头魔像", Description: "由古代魔法赋予生命的岩石巨人", Traits: []string{"高防", "缓慢"}, Dialogue: "（沉重的脚步声）", Strategy: "防御极高，但行动缓慢"},
	{Name: "毒蛛", Description: "一只巨大的蜘蛛，毒液能腐蚀钢铁", Traits: []string{"毒素", "缠绕"}, Dialogue: "嘶嘶嘶……", Strategy: "施加中毒状态，持续伤害"},
	{Name: "火焰元素", Description: "燃烧的火焰凝聚成的元素生物", Traits: []string{"火焰", "灼烧"}, Dialogue: "（火焰噼啪作响）", Strategy: "火焰攻击，有几率灼烧"},
	{Name: "冰霜巨魔", Description: "来自极北冰原的巨大魔物", Traits: []string{"冰冻", "强壮"}, Dialogue: "吼——！入侵者！", Strategy: "冰冻攻击，降低玩家速度"},
	{Name: "幽灵", Description: "飘忽不定的亡灵，物理攻击难以命中", Traits: []string{"虚体", "恐惧"}, Dialogue: "（低沉的哭泣声）", Strategy: "闪避极高，有几率恐惧玩家"},
}

var mockBosses = []MonsterResponse{
	{Name: "地牢守卫者", Description: "守护这片地下城的远古骑士，铠甲上刻满了封印符文", Traits: []string{"精英", "格挡", "反击"}, Dialogue: "凡人……你不该来到这里。", Strategy: "会格挡攻击并反击，需要找准破绽"},
	{Name: "蛛后", Description: "统领所有毒蛛的巨大蛛后，巢穴中到处是蛛网", Traits: []string{"蛛网", "召唤", "毒素"}, Dialogue: "嘶嘶……我的孩子们，开饭了。", Strategy: "召唤小蜘蛛，用蛛网限制玩家行动"},
	{Name: "堕落法师", Description: "曾经的王国法师，被黑暗力量腐蚀后占据了这片地牢", Traits: []string{"魔法", "诅咒", "瞬移"}, Dialogue: "你来得正好，我需要新的实验材料。", Strategy: "使用各种诅咒和魔法攻击，会瞬移"},
}

var mockRooms = map[string][]string{
	"combat": {
		"走廊尽头传来金属碰撞的回声，空气中弥漫着铁锈和血液的味道。",
		"这间房间的地面上散落着破碎的骨头和生锈的武器残骸。",
		"火把的光芒在墙壁上投下诡异的影子，有什么东西在暗处窥视着你。",
		"一扇半开的门后传来低沉的咆哮声，地板上有新鲜的爪痕。",
	},
	"treasure": {
		"房间中央放着一个古朴的宝箱，上面覆盖着厚厚的灰尘。",
		"墙壁上的壁画似乎在暗示什么，角落里有一个不起眼的箱子。",
		"一束光线从天花板的裂缝中照下，照亮了一个镶金的宝箱。",
	},
	"shop": {
		"一个戴着兜帽的身影坐在角落的柜台后，面前摆满了各种奇怪的瓶瓶罐罐。",
		"灯光昏暗的小房间里，一个老矮人正在擦拭他的商品。",
	},
	"rest": {
		"这间房间异常安静，中央有一个已经熄灭的篝火，似乎还算安全。",
		"一个天然形成的石室，泉水从墙壁中渗出，形成了一个小水池。",
	},
	"event": {
		"你走进一间奇怪的房间，墙上刻满了你无法理解的符文。",
		"地板上有一个巨大的魔法阵，微弱的光芒在其中流转。",
		"房间中央有一面古老的铜镜，镜面中似乎映出了不同的景象。",
	},
	"entrance": {
		"你站在地牢的入口处。潮湿的石阶向下延伸，消失在黑暗中。空气中弥漫着霉味和……某种更古老的气息。",
	},
	"boss": {
		"推开沉重的石门，你进入了一个巨大的圆形大厅。空气中充满了压迫感，某种强大的存在正在这里等待着你。",
	},
}

var mockEvents = []EventResponse{
	{
		Title:       "神秘商人",
		Description: "一个幽灵般的商人从墙壁中浮现，他的身体半透明，但手中的商品却异常真实。",
		Choices: []EventChoice{
			{Text: "查看商品", Description: "商人展示了一瓶散发着蓝光的药水，你用50金币买下了它。", Outcome: "good"},
			{Text: "无视并离开", Description: "你转身离开，商人的身影渐渐消散在空气中。", Outcome: "neutral"},
			{Text: "试图抢劫", Description: "商人发出刺耳的笑声，一道闪电击中了你。", Outcome: "bad"},
		},
	},
	{
		Title:       "陷阱房间",
		Description: "你踩到了一块松动的石板，墙壁中射出了几支毒箭！",
		Choices: []EventChoice{
			{Text: "翻滚躲避", Description: "你敏捷地翻滚躲过，只受了轻伤。", Outcome: "neutral"},
			{Text: "用盾牌格挡", Description: "毒箭叮叮当当地弹开，你的盾牌上多了几个凹痕。", Outcome: "good"},
			{Text: "站着不动", Description: "你被毒箭射中了！毒素开始在体内蔓延。", Outcome: "bad"},
		},
	},
	{
		Title:       "古老祭坛",
		Description: "一个散发着微光的祭坛出现在你面前，上面放着一把生锈的匕首和一个发光的宝石。",
		Choices: []EventChoice{
			{Text: "拿起宝石", Description: "宝石的光芒涌入你的身体，你感到力量增强了。", Outcome: "good"},
			{Text: "拿起匕首", Description: "匕首虽然生锈，但你感到它蕴含着某种古老的力量。", Outcome: "good"},
			{Text: "两者都拿", Description: "贪婪触发了祭坛的诅咒，你被一道黑光击中。", Outcome: "bad"},
			{Text: "不碰任何东西", Description: "你谨慎地离开了祭坛。也许下次吧。", Outcome: "neutral"},
		},
	},
}

func (m *MockProvider) GenerateNarrative(_ context.Context, req NarrativeRequest) (string, error) {
	rooms, ok := mockRooms[req.RoomType]
	if !ok {
		return "你来到了一个普通的房间。", nil
	}
	idx := rand.Intn(len(rooms))
	return rooms[idx], nil
}

func (m *MockProvider) GenerateMonster(_ context.Context, req MonsterRequest) (*MonsterResponse, error) {
	if req.IsBoss {
		idx := rand.Intn(len(mockBosses))
		monster := mockBosses[idx]
		return &monster, nil
	}
	idx := rand.Intn(len(mockMonsters))
	monster := mockMonsters[idx]

	// 根据深度调整怪物名称（加上前缀）
	if req.Depth > 3 {
		prefixes := []string{"精英", "强化", "变异", "暗黑"}
		monster.Name = prefixes[rand.Intn(len(prefixes))] + monster.Name
	}
	return &monster, nil
}

func (m *MockProvider) GenerateEvent(_ context.Context, req EventRequest) (*EventResponse, error) {
	idx := rand.Intn(len(mockEvents))
	event := mockEvents[idx]
	return &event, nil
}

func (m *MockProvider) GenerateDialogue(_ context.Context, req DialogueRequest) (string, error) {
	dialogues := map[string][]string{
		"shop": {
			"欢迎，冒险者。看看我的货物吧，保证物美价廉……大概。",
			"哦？又一个不怕死的。随便看看吧，别碰坏了。",
		},
		"mystic": {
			"命运之线在你身上交织……前方的路并不好走。",
			"我能感受到黑暗的力量在增长。小心行事。",
		},
	}
	lines, ok := dialogues[req.NPCType]
	if !ok {
		return "……", nil
	}
	return lines[rand.Intn(len(lines))], nil
}

func (m *MockProvider) GenerateCombatAction(_ context.Context, req CombatActionRequest) (string, error) {
	actions := []string{
		fmt.Sprintf("%s 发动了猛烈的攻击！", req.MonsterName),
		fmt.Sprintf("%s 发出怒吼，向你扑来！", req.MonsterName),
		fmt.Sprintf("%s 悄悄绕到你的身后，发动偷袭！", req.MonsterName),
		fmt.Sprintf("%s 蓄力后释放了一记重击！", req.MonsterName),
		fmt.Sprintf("%s 张开血盆大口咬向你！", req.MonsterName),
	}
	if req.MonsterHP < req.MonsterMaxHP/3 {
		actions = append(actions,
			fmt.Sprintf("%s 在绝境中爆发，发动了疯狂的攻击！", req.MonsterName),
			fmt.Sprintf("%s 感受到了死亡的威胁，变得异常凶猛！", req.MonsterName),
		)
	}
	return actions[rand.Intn(len(actions))], nil
}
