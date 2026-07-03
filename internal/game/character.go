package game

import "fmt"

// Class 职业类型
type Class string

const (
	Warrior Class = "战士"
	Ranger  Class = "游侠"
	Mage    Class = "法师"
	Rogue   Class = "盗贼"
)

// ClassInfo 职业基础信息
type ClassInfo struct {
	Name        Class
	Description string
	BaseHP      int
	BaseATK     int
	BaseDEF     int
	BaseSPD     int
	BaseCRIT    int // 暴击率百分比
	BaseLuck    int
	Skills      []string
}

// 职业基础数据
var ClassData = map[Class]ClassInfo{
	Warrior: {
		Name: Warrior, Description: "近战之王，高防高血，擅长持久战",
		BaseHP: 120, BaseATK: 12, BaseDEF: 15, BaseSPD: 8, BaseCRIT: 5, BaseLuck: 5,
		Skills: []string{"重击", "盾墙", "战吼", "旋风斩"},
	},
	Ranger: {
		Name: Ranger, Description: "远程射手，高暴击高闪避",
		BaseHP: 80, BaseATK: 14, BaseDEF: 8, BaseSPD: 14, BaseCRIT: 20, BaseLuck: 10,
		Skills: []string{"精准射击", "毒箭", "闪避", "致命一击"},
	},
	Mage: {
		Name: Mage, Description: "元素法师，AOE伤害+控制",
		BaseHP: 70, BaseATK: 18, BaseDEF: 5, BaseSPD: 10, BaseCRIT: 10, BaseLuck: 8,
		Skills: []string{"火球术", "冰冻术", "雷电链", "法力护盾"},
	},
	Rogue: {
		Name: Rogue, Description: "暗影刺客，高暴击+连击",
		BaseHP: 85, BaseATK: 16, BaseDEF: 7, BaseSPD: 16, BaseCRIT: 25, BaseLuck: 15,
		Skills: []string{"背刺", "毒刃", "影遁", "连击"},
	},
}

// Character 玩家角色
type Character struct {
	Name      string
	Class     Class
	Level     int
	Exp       int
	ExpToNext int

	// 基础属性
	STR int // 力量 - 影响物理攻击
	DEX int // 敏捷 - 影响暴击和闪避
	INT int // 智力 - 影响魔法攻击
	VIT int // 体质 - 影响生命值
	LUK int // 幸运 - 影响掉落和暴击

	// 战斗属性（含装备加成）
	MaxHP int
	HP    int
	ATK   int
	DEF   int
	SPD   int
	CRIT  int // 暴击率 %
	Dodge int // 闪避率 %

	// 装备
	Equipment map[Slot]*Item

	// 技能
	Skills []*Skill

	// 物品栏
	Inventory []*Item

	// 金币
	Gold int

	// 状态效果
	Buffs   []Buff
	Debuffs []Debuff
}

// NewCharacter 创建新角色
func NewCharacter(name string, class Class) *Character {
	info := ClassData[class]
	c := &Character{
		Name:      name,
		Class:     class,
		Level:     1,
		Exp:       0,
		ExpToNext: 100,
		STR:       info.BaseATK,
		DEX:       info.BaseSPD,
		INT:       info.BaseATK / 2,
		VIT:       info.BaseHP / 10,
		LUK:       info.BaseLuck,
		Equipment: make(map[Slot]*Item),
		Skills:    make([]*Skill, 0),
		Inventory: make([]*Item, 0),
		Gold:      50,
	}

	c.RecalcStats()
	c.HP = c.MaxHP

	// 初始化职业技能
	for _, skillName := range info.Skills {
		if skill, ok := SkillDB[skillName]; ok {
			c.Skills = append(c.Skills, &Skill{
				Name:     skill.Name,
				Desc:     skill.Desc,
				Damage:   skill.Damage,
				Cost:     skill.Cost,
				Cooldown: skill.Cooldown,
				CurCD:    0,
				Effect:   skill.Effect,
				Target:   skill.Target,
			})
		}
	}

	// 初始物品
	c.Inventory = append(c.Inventory, &Item{
		Name: "小型药水", Type: Consumable, Rarity: Common,
		Description: "恢复 30 点生命值",
		Effects:     map[string]int{"heal": 30},
	})
	c.Inventory = append(c.Inventory, &Item{
		Name: "小型药水", Type: Consumable, Rarity: Common,
		Description: "恢复 30 点生命值",
		Effects:     map[string]int{"heal": 30},
	})

	return c
}

// GainExp 获得经验值
func (c *Character) GainExp(amount int) bool {
	c.Exp += amount
	if c.Exp >= c.ExpToNext {
		c.LevelUp()
		return true
	}
	return false
}

// LevelUp 升级
func (c *Character) LevelUp() {
	c.Level++
	c.Exp -= c.ExpToNext
	c.ExpToNext = int(float64(c.ExpToNext) * 1.5)

	// 属性成长
	c.STR += 2
	c.DEX += 1
	c.INT += 1
	c.VIT += 2
	c.LUK += 1

	c.RecalcStats()
	c.HP = c.MaxHP // 升级回满血
}

// RecalcStats 重新计算战斗属性（基础 + 装备）
func (c *Character) RecalcStats() {
	info := ClassData[c.Class]

	c.MaxHP = info.BaseHP + c.VIT*10
	c.ATK = info.BaseATK + c.STR*2
	c.DEF = info.BaseDEF + c.VIT/2
	spdExtra := info.BaseSPD + (c.DEX - info.BaseSPD)
	if spdExtra < 1 {
		spdExtra = 1
	}
	c.SPD = spdExtra
	c.CRIT = info.BaseCRIT + c.LUK/2
	c.Dodge = c.DEX / 3

	// 装备加成
	for _, item := range c.Equipment {
		if item != nil {
			for stat, val := range item.Effects {
				switch stat {
				case "atk":
					c.ATK += val
				case "def":
					c.DEF += val
				case "hp":
					c.MaxHP += val
				case "crit":
					c.CRIT += val
				case "dodge":
					c.Dodge += val
				}
			}
		}
	}

	if c.HP > c.MaxHP {
		c.HP = c.MaxHP
	}
}

// IsAlive 角色是否存活
func (c *Character) IsAlive() bool {
	return c.HP > 0
}

// TakeDamage 受到伤害
func (c *Character) TakeDamage(damage int) int {
	actualDamage := damage - c.DEF/2
	if actualDamage < 1 {
		actualDamage = 1
	}
	c.HP -= actualDamage
	if c.HP < 0 {
		c.HP = 0
	}
	return actualDamage
}

// Heal 恢复生命
func (c *Character) Heal(amount int) {
	c.HP += amount
	if c.HP > c.MaxHP {
		c.HP = c.MaxHP
	}
}

// String 返回角色状态摘要
func (c *Character) String() string {
	return fmt.Sprintf("[%s Lv.%d %s] HP:%d/%d ATK:%d DEF:%d SPD:%d CRIT:%d%%",
		c.Name, c.Level, c.Class, c.HP, c.MaxHP, c.ATK, c.DEF, c.SPD, c.CRIT)
}
