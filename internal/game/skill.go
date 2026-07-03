package game

// Skill 技能
type Skill struct {
	Name     string
	Desc     string
	Damage   int    // 伤害倍率 (百分比，100 = 1x)
	Cost     int    // 消耗（目前无 MP，改为冷却制）
	Cooldown int    // 冷却回合数
	CurCD    int    // 当前剩余冷却
	Effect   string // 特殊效果：stun, poison, burn, freeze, heal, buff_def, buff_atk, aoe
	Target   string // single / self / all_enemies
}

// SkillDB 技能数据库
var SkillDB = map[string]*Skill{
	// === 战士技能 ===
	"重击": {
		Name: "重击", Desc: "蓄力后发动一记猛击，造成 180% 伤害",
		Damage: 180, Cooldown: 2, Effect: "", Target: "single",
	},
	"盾墙": {
		Name: "盾墙", Desc: "举起盾牌防御，本回合受到的伤害减半",
		Damage: 0, Cooldown: 3, Effect: "buff_def", Target: "self",
	},
	"战吼": {
		Name: "战吼", Desc: "发出震天怒吼，提升攻击力并震慑敌人",
		Damage: 0, Cooldown: 4, Effect: "buff_atk", Target: "self",
	},
	"旋风斩": {
		Name: "旋风斩", Desc: "旋转武器攻击所有敌人，造成 120% 伤害",
		Damage: 120, Cooldown: 3, Effect: "", Target: "all_enemies",
	},

	// === 游侠技能 ===
	"精准射击": {
		Name: "精准射击", Desc: "瞄准要害射击，必定暴击",
		Damage: 150, Cooldown: 2, Effect: "guaranteed_crit", Target: "single",
	},
	"毒箭": {
		Name: "毒箭", Desc: "射出毒箭，造成伤害并使敌人中毒",
		Damage: 100, Cooldown: 3, Effect: "poison", Target: "single",
	},
	"闪避": {
		Name: "闪避", Desc: "进入闪避状态，下次攻击必定闪避",
		Damage: 0, Cooldown: 3, Effect: "buff_dodge", Target: "self",
	},
	"致命一击": {
		Name: "致命一击", Desc: "瞄准弱点发动致命一击，造成 250% 伤害",
		Damage: 250, Cooldown: 5, Effect: "", Target: "single",
	},

	// === 法师技能 ===
	"火球术": {
		Name: "火球术", Desc: "发射火球，造成伤害并灼烧敌人",
		Damage: 160, Cooldown: 2, Effect: "burn", Target: "single",
	},
	"冰冻术": {
		Name: "冰冻术", Desc: "释放寒冰，造成伤害并冻结敌人一回合",
		Damage: 120, Cooldown: 3, Effect: "freeze", Target: "single",
	},
	"雷电链": {
		Name: "雷电链", Desc: "释放雷电攻击所有敌人，造成 100% 伤害",
		Damage: 100, Cooldown: 3, Effect: "", Target: "all_enemies",
	},
	"法力护盾": {
		Name: "法力护盾", Desc: "生成魔法护盾，吸收伤害",
		Damage: 0, Cooldown: 4, Effect: "shield", Target: "self",
	},

	// === 盗贼技能 ===
	"背刺": {
		Name: "背刺", Desc: "从背后攻击，造成 200% 伤害",
		Damage: 200, Cooldown: 2, Effect: "", Target: "single",
	},
	"毒刃": {
		Name: "毒刃", Desc: "在武器上涂毒，下次攻击附带中毒效果",
		Damage: 80, Cooldown: 2, Effect: "poison", Target: "single",
	},
	"影遁": {
		Name: "影遁", Desc: "隐入暗影，大幅提升闪避率",
		Damage: 0, Cooldown: 4, Effect: "buff_dodge", Target: "self",
	},
	"连击": {
		Name: "连击", Desc: "快速连续攻击 3 次，每次 60% 伤害",
		Damage: 180, Cooldown: 3, Effect: "multi_hit", Target: "single",
	},
}

// Buff 正面状态效果
type Buff struct {
	Name     string
	Duration int // 剩余回合
	Stat     string
	Value    int
}

// Debuff 负面状态效果
type Debuff struct {
	Name     string
	Duration int // 剩余回合
	Stat     string
	Value    int
	Damage   int // 每回合伤害（如中毒、灼烧）
}
