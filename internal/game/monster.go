package game

import (
	"fmt"
	"math/rand"
)

// Monster 怪物
type Monster struct {
	Name     string
	Desc     string
	Traits   []string
	Dialogue string
	Strategy string

	MaxHP int
	HP    int
	ATK   int
	DEF   int
	SPD   int
	CRIT  int
	Dodge int

	ExpReward  int
	GoldReward int

	// 状态效果
	Buffs   []Buff
	Debuffs []Debuff

	// 是否是 Boss
	IsBoss bool
}

// NewMonster 从 LLM 响应创建怪物
func NewMonster(resp interface{}, depth int, isBoss bool) *Monster {
	// type assert - this will be *llm.MonsterResponse but we avoid import cycle
	// We'll use a simplified approach
	return nil
}

// NewMonsterFromData 直接从数据创建怪物
func NewMonsterFromData(name, desc string, traits []string, dialogue, strategy string, depth int, isBoss bool) *Monster {
	baseHP := 60 + depth*25
	baseATK := 15 + depth*5
	baseDEF := 8 + depth*3
	baseSPD := 8 + depth
	baseCRIT := 5 + depth*2

	if isBoss {
		baseHP *= 3
		baseATK = int(float64(baseATK) * 1.8)
		baseDEF = int(float64(baseDEF) * 1.5)
		baseCRIT += 10
	}

	m := &Monster{
		Name:       name,
		Desc:       desc,
		Traits:     traits,
		Dialogue:   dialogue,
		Strategy:   strategy,
		MaxHP:      baseHP,
		HP:         baseHP,
		ATK:        baseATK,
		DEF:        baseDEF,
		SPD:        baseSPD,
		CRIT:       baseCRIT,
		Dodge:      baseSPD / 3,
		ExpReward:  20 + depth*10,
		GoldReward: 10 + depth*5 + rand.Intn(10),
		IsBoss:     isBoss,
		Buffs:      make([]Buff, 0),
		Debuffs:    make([]Debuff, 0),
	}

	if isBoss {
		m.ExpReward *= 3
		m.GoldReward *= 3
	}

	return m
}

// IsAlive 怪物是否存活
func (m *Monster) IsAlive() bool {
	return m.HP > 0
}

// TakeDamage 受到伤害
func (m *Monster) TakeDamage(damage int) int {
	actualDamage := damage - m.DEF/2
	if actualDamage < 1 {
		actualDamage = 1
	}
	m.HP -= actualDamage
	if m.HP < 0 {
		m.HP = 0
	}
	return actualDamage
}

// String 怪物状态
func (m *Monster) String() string {
	bossTag := ""
	if m.IsBoss {
		bossTag = " ★BOSS★"
	}
	return fmt.Sprintf("%s%s HP:%d/%d ATK:%d DEF:%d", m.Name, bossTag, m.HP, m.MaxHP, m.ATK, m.DEF)
}

// ProcessDebuffs 处理回合结束时的持续伤害效果
func (m *Monster) ProcessDebuffs() []string {
	var messages []string
	newDebuffs := make([]Debuff, 0)
	for _, d := range m.Debuffs {
		if d.Damage > 0 {
			m.HP -= d.Damage
			if m.HP < 0 {
				m.HP = 0
			}
			messages = append(messages, fmt.Sprintf("%s受到 %d 点%s伤害！", m.Name, d.Damage, d.Name))
		}
		d.Duration--
		if d.Duration > 0 {
			newDebuffs = append(newDebuffs, d)
		} else {
			messages = append(messages, fmt.Sprintf("%s的%s效果消失了。", m.Name, d.Name))
		}
	}
	m.Debuffs = newDebuffs
	return messages
}

// HasDebuff 检查是否有特定 debuff
func (m *Monster) HasDebuff(name string) bool {
	for _, d := range m.Debuffs {
		if d.Name == name {
			return true
		}
	}
	return false
}

// HasTrait 检查是否有特定特征
func (m *Monster) HasTrait(trait string) bool {
	for _, t := range m.Traits {
		if t == trait {
			return true
		}
	}
	return false
}
