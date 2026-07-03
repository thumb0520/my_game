package game

import (
	"fmt"
	"math/rand"
)

// CombatResult 战斗结果
type CombatResult struct {
	PlayerDamage  int
	MonsterDamage int
	Messages      []string
	PlayerCrit    bool
	MonsterCrit   bool
	PlayerDodge   bool
	MonsterDodge  bool
	IsOver        bool
	PlayerWon     bool
}

// CombatAction 战斗行动类型
type CombatAction int

const (
	ActionAttack CombatAction = iota
	ActionSkill
	ActionItem
	ActionFlee
)

// CombatState 战斗状态
type CombatState struct {
	Player    *Character
	Monster   *Monster
	Round     int
	Log       []string
	IsOver    bool
	Won       bool
	Fled      bool
}

// NewCombatState 创建战斗状态
func NewCombatState(player *Character, monster *Monster) *CombatState {
	return &CombatState{
		Player:  player,
		Monster: monster,
		Round:   0,
		Log:     make([]string, 0),
	}
}

// ExecutePlayerAttack 玩家普通攻击
func (cs *CombatState) ExecutePlayerAttack() *CombatResult {
	cs.Round++
	result := &CombatResult{Messages: make([]string, 0)}

	// 检查闪避
	if rand.Intn(100) < cs.Monster.Dodge {
		result.MonsterDodge = true
		result.Messages = append(result.Messages, fmt.Sprintf("%s 闪避了你的攻击！", cs.Monster.Name))
		return result
	}

	// 计算伤害
	damage := cs.Player.ATK
	isCrit := rand.Intn(100) < cs.Player.CRIT
	if isCrit {
		damage = int(float64(damage) * 1.8)
		result.PlayerCrit = true
	}

	actualDamage := cs.Monster.TakeDamage(damage)
	result.PlayerDamage = actualDamage

	if isCrit {
		result.Messages = append(result.Messages, fmt.Sprintf("暴击！你对 %s 造成了 %d 点伤害！", cs.Monster.Name, actualDamage))
	} else {
		result.Messages = append(result.Messages, fmt.Sprintf("你攻击了 %s，造成 %d 点伤害。", cs.Monster.Name, actualDamage))
	}

	if !cs.Monster.IsAlive() {
		result.Messages = append(result.Messages, fmt.Sprintf("%s 被击败了！", cs.Monster.Name))
		result.IsOver = true
		result.PlayerWon = true
	}

	return result
}

// ExecutePlayerSkill 玩家使用技能
func (cs *CombatState) ExecutePlayerSkill(skillIdx int) *CombatResult {
	cs.Round++
	result := &CombatResult{Messages: make([]string, 0)}

	if skillIdx < 0 || skillIdx >= len(cs.Player.Skills) {
		result.Messages = append(result.Messages, "无效的技能！")
		return result
	}

	skill := cs.Player.Skills[skillIdx]

	// 检查冷却
	if skill.CurCD > 0 {
		result.Messages = append(result.Messages, fmt.Sprintf("%s 还在冷却中（剩余 %d 回合）！", skill.Name, skill.CurCD))
		cs.Round-- // 不消耗回合
		return result
	}

	// 设置冷却
	skill.CurCD = skill.Cooldown

	// 处理技能效果
	switch skill.Effect {
	case "buff_def":
		cs.Player.Buffs = append(cs.Player.Buffs, Buff{Name: "防御强化", Duration: 2, Stat: "def", Value: cs.Player.DEF / 2})
		result.Messages = append(result.Messages, fmt.Sprintf("你使用了 %s！防御力大幅提升！", skill.Name))

	case "buff_atk":
		cs.Player.Buffs = append(cs.Player.Buffs, Buff{Name: "攻击强化", Duration: 3, Stat: "atk", Value: cs.Player.ATK / 3})
		result.Messages = append(result.Messages, fmt.Sprintf("你使用了 %s！攻击力提升！", skill.Name))

	case "buff_dodge":
		cs.Player.Buffs = append(cs.Player.Buffs, Buff{Name: "闪避强化", Duration: 2, Stat: "dodge", Value: 50})
		result.Messages = append(result.Messages, fmt.Sprintf("你使用了 %s！闪避率大幅提升！", skill.Name))

	case "shield":
		cs.Player.Buffs = append(cs.Player.Buffs, Buff{Name: "护盾", Duration: 3, Stat: "shield", Value: cs.Player.MaxHP / 4})
		result.Messages = append(result.Messages, fmt.Sprintf("你使用了 %s！获得了一个护盾！", skill.Name))

	case "poison":
		damage := cs.Monster.TakeDamage(int(float64(cs.Player.ATK) * float64(skill.Damage) / 100))
		cs.Monster.Debuffs = append(cs.Monster.Debuffs, Debuff{Name: "中毒", Duration: 3, Damage: cs.Player.ATK / 4})
		result.PlayerDamage = damage
		result.Messages = append(result.Messages, fmt.Sprintf("你使用了 %s！造成 %d 点伤害，并使敌人中毒！", skill.Name, damage))

	case "burn":
		damage := cs.Monster.TakeDamage(int(float64(cs.Player.ATK) * float64(skill.Damage) / 100))
		cs.Monster.Debuffs = append(cs.Monster.Debuffs, Debuff{Name: "灼烧", Duration: 2, Damage: cs.Player.ATK / 3})
		result.PlayerDamage = damage
		result.Messages = append(result.Messages, fmt.Sprintf("你使用了 %s！造成 %d 点伤害，敌人被灼烧！", skill.Name, damage))

	case "freeze":
		damage := cs.Monster.TakeDamage(int(float64(cs.Player.ATK) * float64(skill.Damage) / 100))
		cs.Monster.Debuffs = append(cs.Monster.Debuffs, Debuff{Name: "冰冻", Duration: 1})
		result.PlayerDamage = damage
		result.Messages = append(result.Messages, fmt.Sprintf("你使用了 %s！造成 %d 点伤害，敌人被冻结！", skill.Name, damage))

	case "stun":
		damage := cs.Monster.TakeDamage(int(float64(cs.Player.ATK) * float64(skill.Damage) / 100))
		cs.Monster.Debuffs = append(cs.Monster.Debuffs, Debuff{Name: "眩晕", Duration: 1})
		result.PlayerDamage = damage
		result.Messages = append(result.Messages, fmt.Sprintf("你使用了 %s！造成 %d 点伤害，敌人被眩晕！", skill.Name, damage))

	default:
		// 纯伤害技能
		baseDamage := int(float64(cs.Player.ATK) * float64(skill.Damage) / 100)

		// 检查暴击
		isCrit := rand.Intn(100) < cs.Player.CRIT
		if skill.Effect == "guaranteed_crit" {
			isCrit = true
		}
		if isCrit {
			baseDamage = int(float64(baseDamage) * 1.8)
			result.PlayerCrit = true
		}

		if skill.Target == "all_enemies" {
			// AOE - 单目标战斗中就是普通伤害
			damage := cs.Monster.TakeDamage(baseDamage)
			result.PlayerDamage = damage
			if isCrit {
				result.Messages = append(result.Messages, fmt.Sprintf("暴击！%s 对 %s 造成了 %d 点伤害！", skill.Name, cs.Monster.Name, damage))
			} else {
				result.Messages = append(result.Messages, fmt.Sprintf("你使用了 %s！对 %s 造成 %d 点伤害。", skill.Name, cs.Monster.Name, damage))
			}
		} else {
			damage := cs.Monster.TakeDamage(baseDamage)
			result.PlayerDamage = damage
			if isCrit {
				result.Messages = append(result.Messages, fmt.Sprintf("暴击！%s 对 %s 造成了 %d 点伤害！", skill.Name, cs.Monster.Name, damage))
			} else {
				result.Messages = append(result.Messages, fmt.Sprintf("你使用了 %s！对 %s 造成 %d 点伤害。", skill.Name, cs.Monster.Name, damage))
			}
		}
	}

	if !cs.Monster.IsAlive() {
		result.Messages = append(result.Messages, fmt.Sprintf("%s 被击败了！", cs.Monster.Name))
		result.IsOver = true
		result.PlayerWon = true
	}

	return result
}

// ExecutePlayerUseItem 玩家使用物品
func (cs *CombatState) ExecutePlayerUseItem(itemIdx int) *CombatResult {
	cs.Round++
	result := &CombatResult{Messages: make([]string, 0)}

	if itemIdx < 0 || itemIdx >= len(cs.Player.Inventory) {
		result.Messages = append(result.Messages, "无效的物品！")
		cs.Round--
		return result
	}

	item := cs.Player.Inventory[itemIdx]
	if item.Type != Consumable {
		result.Messages = append(result.Messages, "该物品无法在战斗中使用！")
		cs.Round--
		return result
	}

	// 处理消耗品效果
	if heal, ok := item.Effects["heal"]; ok {
		cs.Player.Heal(heal)
		result.Messages = append(result.Messages, fmt.Sprintf("你使用了 %s，恢复了 %d 点生命值！", item.Name, heal))
	}
	if _, ok := item.Effects["cure_poison"]; ok {
		newDebuffs := make([]Debuff, 0)
		for _, d := range cs.Player.Debuffs {
			if d.Name != "中毒" {
				newDebuffs = append(newDebuffs, d)
			}
		}
		cs.Player.Debuffs = newDebuffs
		result.Messages = append(result.Messages, fmt.Sprintf("你使用了 %s，中毒状态被清除！", item.Name))
	}
	if atk, ok := item.Effects["temp_atk"]; ok {
		cs.Player.Buffs = append(cs.Player.Buffs, Buff{Name: "力量增强", Duration: 3, Stat: "atk", Value: atk})
		result.Messages = append(result.Messages, fmt.Sprintf("你使用了 %s，攻击力临时提升 %d！", item.Name, atk))
	}
	if def, ok := item.Effects["temp_def"]; ok {
		cs.Player.Buffs = append(cs.Player.Buffs, Buff{Name: "防御增强", Duration: 3, Stat: "def", Value: def})
		result.Messages = append(result.Messages, fmt.Sprintf("你使用了 %s，防御力临时提升 %d！", item.Name, def))
	}

	// 移除物品
	cs.Player.Inventory = append(cs.Player.Inventory[:itemIdx], cs.Player.Inventory[itemIdx+1:]...)

	return result
}

// ExecutePlayerFlee 玩家逃跑
func (cs *CombatState) ExecutePlayerFlee() *CombatResult {
	result := &CombatResult{Messages: make([]string, 0)}

	// 逃跑成功率：基础 50% + 速度差
	fleeChance := 50 + (cs.Player.SPD - cs.Monster.SPD) * 2
	if cs.Monster.IsBoss {
		fleeChance -= 30 // Boss 战逃跑更难
	}

	if rand.Intn(100) < fleeChance {
		result.Messages = append(result.Messages, "你成功逃离了战斗！")
		result.IsOver = true
		cs.Fled = true
	} else {
		result.Messages = append(result.Messages, "逃跑失败！")
		// 逃跑失败，怪物免费攻击一次
		monsterDmg := cs.Monster.ATK
		actualDmg := cs.Player.TakeDamage(monsterDmg)
		result.MonsterDamage = actualDmg
		result.Messages = append(result.Messages, fmt.Sprintf("%s 趁机攻击了你，造成 %d 点伤害！", cs.Monster.Name, actualDmg))
		if !cs.Player.IsAlive() {
			result.Messages = append(result.Messages, "你倒下了……")
			result.IsOver = true
		}
	}

	return result
}

// ExecuteMonsterTurn 怪物回合
func (cs *CombatState) ExecuteMonsterTurn(monsterAction string) *CombatResult {
	result := &CombatResult{Messages: make([]string, 0)}

	if !cs.Monster.IsAlive() {
		return result
	}

	// 检查是否被冻结/眩晕
	if cs.Monster.HasDebuff("冰冻") || cs.Monster.HasDebuff("眩晕") {
		result.Messages = append(result.Messages, fmt.Sprintf("%s 处于无法行动的状态！", cs.Monster.Name))
		return result
	}

	// 检查玩家闪避
	totalDodge := cs.Player.Dodge
	for _, buff := range cs.Player.Buffs {
		if buff.Stat == "dodge" {
			totalDodge += buff.Value
		}
	}
	if rand.Intn(100) < totalDodge {
		result.PlayerDodge = true
		result.Messages = append(result.Messages, fmt.Sprintf("你闪避了 %s 的攻击！", cs.Monster.Name))
		return result
	}

	// 计算伤害
	damage := cs.Monster.ATK
	isCrit := rand.Intn(100) < cs.Monster.CRIT
	if isCrit {
		damage = int(float64(damage) * 1.5)
		result.MonsterCrit = true
	}

	// 检查玩家护盾
	shieldDmg := 0
	for _, buff := range cs.Player.Buffs {
		if buff.Stat == "shield" && buff.Value > 0 {
			shieldDmg = buff.Value
			break
		}
	}

	// 检查玩家防御buff
	defBonus := 0
	for _, buff := range cs.Player.Buffs {
		if buff.Stat == "def" {
			defBonus += buff.Value
		}
	}

	originalDEF := cs.Player.DEF
	cs.Player.DEF += defBonus
	actualDamage := cs.Player.TakeDamage(damage)
	cs.Player.DEF = originalDEF

	// 护盾吸收
	if shieldDmg > 0 {
	 absorbed := min(shieldDmg, actualDamage)
		actualDamage -= absorbed
		cs.Player.HP += absorbed // 补回护盾吸收的部分
		// 护盾值减少
		for i := range cs.Player.Buffs {
			if cs.Player.Buffs[i].Stat == "shield" {
				cs.Player.Buffs[i].Value -= absorbed
				if cs.Player.Buffs[i].Value < 0 {
					cs.Player.Buffs[i].Value = 0
				}
			}
		}
		result.Messages = append(result.Messages, fmt.Sprintf("护盾吸收了 %d 点伤害！", absorbed))
	}

	result.MonsterDamage = actualDamage

	if monsterAction != "" {
		result.Messages = append(result.Messages, monsterAction)
	}

	if isCrit {
		result.Messages = append(result.Messages, fmt.Sprintf("暴击！%s 对你造成了 %d 点伤害！", cs.Monster.Name, actualDamage))
	} else {
		result.Messages = append(result.Messages, fmt.Sprintf("%s 攻击了你，造成 %d 点伤害。", cs.Monster.Name, actualDamage))
	}

	if !cs.Player.IsAlive() {
		result.Messages = append(result.Messages, "你倒下了……")
		result.IsOver = true
	}

	return result
}

// ProcessTurnEnd 处理回合结束效果
func (cs *CombatState) ProcessTurnEnd() []string {
	var messages []string

	// 处理怪物持续伤害
	monsterDebuffMsgs := cs.Monster.ProcessDebuffs()
	messages = append(messages, monsterDebuffMsgs...)

	// 处理玩家 debuff
	newDebuffs := make([]Debuff, 0)
	for _, d := range cs.Player.Debuffs {
		if d.Damage > 0 {
			cs.Player.HP -= d.Damage
			if cs.Player.HP < 0 {
				cs.Player.HP = 0
			}
			messages = append(messages, fmt.Sprintf("你受到 %d 点%s伤害！", d.Damage, d.Name))
		}
		d.Duration--
		if d.Duration > 0 {
			newDebuffs = append(newDebuffs, d)
		}
	}
	cs.Player.Debuffs = newDebuffs

	// 处理玩家 buff 持续时间
	newBuffs := make([]Buff, 0)
	for _, b := range cs.Player.Buffs {
		b.Duration--
		if b.Duration > 0 {
			newBuffs = append(newBuffs, b)
		}
	}
	cs.Player.Buffs = newBuffs

	// 处理技能冷却
	for _, skill := range cs.Player.Skills {
		if skill.CurCD > 0 {
			skill.CurCD--
		}
	}

	return messages
}

// GetBuffDefBonus 获取防御 buff 加成
func (cs *CombatState) GetBuffDefBonus() int {
	bonus := 0
	for _, buff := range cs.Player.Buffs {
		if buff.Stat == "def" {
			bonus += buff.Value
		}
	}
	return bonus
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
